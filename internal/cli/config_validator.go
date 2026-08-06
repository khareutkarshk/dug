package cli

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
)

// ValidationIssue is a single configuration problem.
type ValidationIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (i ValidationIssue) String() string {
	if i.Field == "" {
		return i.Message
	}
	return fmt.Sprintf("%s: %s", i.Field, i.Message)
}

// ValidationResult holds all issues found while validating a config.
type ValidationResult struct {
	Issues []ValidationIssue `json:"issues,omitempty"`
}

func (r *ValidationResult) Ok() bool {
	return len(r.Issues) == 0
}

func (r *ValidationResult) Error() error {
	if r.Ok() {
		return nil
	}

	msgs := make([]string, 0, len(r.Issues))
	for _, issue := range r.Issues {
		msgs = append(msgs, issue.String())
	}
	return errors.New(strings.Join(msgs, "\n"))
}

func (r *ValidationResult) add(field, message string) {
	r.Issues = append(r.Issues, ValidationIssue{
		Field:   field,
		Message: message,
	})
}

func checkConfig(cfg *config.Config) *ValidationResult {
	result := &ValidationResult{}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		result.add("server.port", "must be between 1 and 65535")
	}

	if cfg.Server.Retries < 0 {
		result.add("server.retries", "must be >= 0")
	}

	if cfg.Server.RateLimit.RPS < 0 {
		result.add("server.rate_limit.rps", "must be >= 0")
	}

	if cfg.Server.RateLimit.Burst < 0 {
		result.add("server.rate_limit.burst", "must be >= 0")
	}

	if cfg.Server.ReadTimeout < 0 {
		result.add("server.read_timeout", "must be >= 0")
	}
	if cfg.Server.WriteTimeout < 0 {
		result.add("server.write_timeout", "must be >= 0")
	}
	if cfg.Server.IdleTimeout < 0 {
		result.add("server.idle_timeout", "must be >= 0")
	}

	if cfg.Server.Limits.BodySize < 0 {
		result.add("server.limits.body_size", "must be >= 0")
	}

	if cfg.Server.Compression.MinSize < 0 {
		result.add("server.compression.min_size", "must be >= 0")
	}

	validateSecurityHeaders(result, cfg.Server.Security.Headers)

	if cfg.Server.TLS.Enabled {
		if strings.TrimSpace(cfg.Server.TLS.CertFile) == "" {
			result.add("server.tls.cert_file", "required when tls.enabled is true")
		} else if _, err := os.Stat(cfg.Server.TLS.CertFile); err != nil {
			result.add("server.tls.cert_file", fmt.Sprintf("file not found: %s", cfg.Server.TLS.CertFile))
		}

		if strings.TrimSpace(cfg.Server.TLS.KeyFile) == "" {
			result.add("server.tls.key_file", "required when tls.enabled is true")
		} else if _, err := os.Stat(cfg.Server.TLS.KeyFile); err != nil {
			result.add("server.tls.key_file", fmt.Sprintf("file not found: %s", cfg.Server.TLS.KeyFile))
		}
	}

	if len(cfg.Routes) == 0 {
		result.add("routes", "at least one route is required")
		return result
	}

	seenPaths := make(map[string]int)

	for i, route := range cfg.Routes {
		prefix := fmt.Sprintf("routes[%d]", i)

		path := strings.TrimSpace(route.Path)
		if path == "" {
			result.add(prefix+".path", "must not be empty")
		} else {
			if prev, ok := seenPaths[path]; ok {
				result.add(prefix+".path", fmt.Sprintf("duplicate of routes[%d].path (%s)", prev, path))
			} else {
				seenPaths[path] = i
			}
			if !strings.HasPrefix(path, "/") {
				result.add(prefix+".path", "must start with '/'")
			}
		}

		switch route.Strategy {
		case "",
			upstream.StrategyLeastConnections,
			upstream.StrategySmoothWeighted:
		default:
			result.add(prefix+".strategy", fmt.Sprintf(
				"unknown value %q (allowed: %s, %s)",
				route.Strategy,
				upstream.StrategySmoothWeighted,
				upstream.StrategyLeastConnections,
			))
		}

		if route.Timeout < 0 {
			result.add(prefix+".timeout", "must be >= 0")
		}

		validateHeaderRules(result, prefix+".request_headers", route.RequestHeaders)
		validateHeaderRules(result, prefix+".response_headers", route.ResponseHeaders)

		if len(route.Upstreams) == 0 {
			result.add(prefix+".upstreams", "at least one upstream is required")
			continue
		}

		seenUpstreamURLs := make(map[string]int)

		for j, u := range route.Upstreams {
			upPrefix := fmt.Sprintf("%s.upstreams[%d]", prefix, j)

			if strings.TrimSpace(u.URL) == "" {
				result.add(upPrefix+".url", "must not be empty")
				continue
			}

			parsed, err := url.ParseRequestURI(u.URL)
			if err != nil {
				result.add(upPrefix+".url", fmt.Sprintf("invalid URL: %v", err))
				continue
			}

			if parsed.Scheme != "http" && parsed.Scheme != "https" {
				result.add(upPrefix+".url", "scheme must be http or https")
			}

			if parsed.Host == "" {
				result.add(upPrefix+".url", "host is required")
			}

			normalized := strings.TrimRight(u.URL, "/")
			if prev, ok := seenUpstreamURLs[normalized]; ok {
				result.add(upPrefix+".url", fmt.Sprintf("duplicate of %s.upstreams[%d].url", prefix, prev))
			} else {
				seenUpstreamURLs[normalized] = j
			}

			if u.Weight <= 0 {
				result.add(upPrefix+".weight", "must be > 0")
			}
		}

		if route.CORS.Enabled {
			if route.CORS.AllowCredentials {
				for _, origin := range route.CORS.AllowOrigins {
					if origin == "*" {
						result.add(prefix+".cors", "allow_credentials cannot be true when allow_origins includes '*'")
						break
					}
				}
			}
		}
	}

	return result
}

func validateHeaderRules(result *ValidationResult, prefix string, rules config.HeaderRules) {
	seenAdd := make(map[string]string)
	for name := range rules.Add {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			result.add(prefix+".add", "header name must not be empty")
			continue
		}
		if prev, ok := seenAdd[key]; ok {
			result.add(prefix+".add", fmt.Sprintf("duplicate header %q (also as %q)", name, prev))
			continue
		}
		seenAdd[key] = name
	}

	seenRemove := make(map[string]struct{})
	for _, name := range rules.Remove {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			result.add(prefix+".remove", "header name must not be empty")
			continue
		}
		if _, ok := seenRemove[key]; ok {
			result.add(prefix+".remove", fmt.Sprintf("duplicate header %q", name))
			continue
		}
		seenRemove[key] = struct{}{}

		if prev, ok := seenAdd[key]; ok {
			result.add(prefix, fmt.Sprintf("header %q cannot be both added (%q) and removed", name, prev))
		}
	}
}

func validateSecurityHeaders(result *ValidationResult, headers config.SecurityHeaders) {
	if v := strings.TrimSpace(headers.XFrameOptions); v != "" {
		upper := strings.ToUpper(v)
		if upper != "DENY" && upper != "SAMEORIGIN" && !strings.HasPrefix(upper, "ALLOW-FROM ") {
			result.add(
				"server.security.headers.x_frame_options",
				`must be "DENY", "SAMEORIGIN", or "ALLOW-FROM <uri>"`,
			)
		}
	}

	if v := strings.TrimSpace(headers.XContentTypeOptions); v != "" {
		if !strings.EqualFold(v, "nosniff") {
			result.add(
				"server.security.headers.x_content_type_options",
				`must be "nosniff"`,
			)
		}
	}

	if v := strings.TrimSpace(headers.StrictTransportSecurity); v != "" {
		if !strings.Contains(strings.ToLower(v), "max-age=") {
			result.add(
				"server.security.headers.strict_transport_security",
				`must include "max-age="`,
			)
		}
	}

	if v := strings.TrimSpace(headers.ReferrerPolicy); v != "" {
		if !validReferrerPolicy(v) {
			result.add(
				"server.security.headers.referrer_policy",
				"invalid referrer policy value",
			)
		}
	}
}

func validReferrerPolicy(value string) bool {
	switch strings.ToLower(value) {
	case "no-referrer",
		"no-referrer-when-downgrade",
		"origin",
		"origin-when-cross-origin",
		"same-origin",
		"strict-origin",
		"strict-origin-when-cross-origin",
		"unsafe-url":
		return true
	default:
		return false
	}
}

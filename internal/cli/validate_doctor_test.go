package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

func TestCheckConfigValid(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
		Routes: []config.Route{{
			Path:      "/",
			Strategy:  "least_connections",
			Upstreams: []config.Upstream{{URL: "http://localhost:3001", Weight: 1}},
		}},
	}

	result := checkConfig(cfg)
	if !result.Ok() {
		t.Fatalf("expected valid config, got: %v", result.Error())
	}
}

func TestCheckConfigCollectsIssues(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 0},
		Routes: []config.Route{
			{Path: "api", Strategy: "round_robin", Upstreams: nil},
			{Path: "/ok", Upstreams: []config.Upstream{{URL: "ftp://bad", Weight: -1}}},
		},
	}

	result := checkConfig(cfg)
	if result.Ok() {
		t.Fatal("expected validation issues")
	}

	joined := result.Error().Error()
	for _, want := range []string{
		"server.port",
		"routes[0].path",
		"routes[0].strategy",
		"routes[0].upstreams",
		"routes[1].upstreams[0].url",
		"routes[1].upstreams[0].weight",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing issue for %s in:\n%s", want, joined)
		}
	}
}

func TestCheckConfigV1Features(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Limits: config.LimitsConfig{
				BodySize: -1,
			},
			Compression: config.CompressionConfig{
				Enabled: true,
				MinSize: -5,
			},
			Security: config.SecurityConfig{
				Headers: config.SecurityHeaders{
					XFrameOptions:           "ALLOWALL",
					XContentTypeOptions:     "invalid",
					StrictTransportSecurity: "includeSubDomains",
					ReferrerPolicy:          "not-a-policy",
				},
			},
			TLS: config.TLSConfig{
				Enabled: true,
			},
		},
		Routes: []config.Route{
			{
				Path:    "/api",
				Timeout: -1,
				RequestHeaders: config.HeaderRules{
					Add:    map[string]string{"X-A": "1"},
					Remove: []string{"X-A", "X-B", "x-b"},
				},
				Upstreams: []config.Upstream{
					{URL: "http://localhost:3001", Weight: 0},
					{URL: "http://localhost:3001/", Weight: 2},
				},
			},
			{
				Path:      "/api",
				Upstreams: []config.Upstream{{URL: "http://localhost:3002", Weight: 1}},
			},
		},
	}

	result := checkConfig(cfg)
	if result.Ok() {
		t.Fatal("expected validation issues")
	}

	joined := result.Error().Error()
	for _, want := range []string{
		"server.limits.body_size",
		"server.compression.min_size",
		"server.security.headers.x_frame_options",
		"server.security.headers.x_content_type_options",
		"server.security.headers.strict_transport_security",
		"server.security.headers.referrer_policy",
		"server.tls.cert_file",
		"server.tls.key_file",
		"routes[0].timeout",
		"routes[0].upstreams[0].weight",
		"routes[0].upstreams[1].url",
		"routes[1].path",
		"routes[0].request_headers",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing issue for %s in:\n%s", want, joined)
		}
	}
}

func TestValidateCommandJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ok.yaml")
	content := []byte(`
server:
  port: 8080
routes:
  - path: /
    upstreams:
      - url: http://127.0.0.1:9
        weight: 1
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := validate(&stdout, &stderr, []string{"-config", path, "-json"}); err != nil {
		t.Fatal(err)
	}

	var rep validateReport
	if err := json.Unmarshal(stdout.Bytes(), &rep); err != nil {
		t.Fatal(err)
	}
	if !rep.Valid || rep.Routes != 1 || rep.Server.Port != 8080 {
		t.Fatalf("unexpected report: %+v", rep)
	}
}

func TestValidateCommandInvalidQuiet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	content := []byte(`
server:
  port: 8080
routes: []
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := validate(&stdout, &stderr, []string{"-config", path, "-quiet"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(stderr.String(), "invalid") {
		t.Fatalf("expected invalid message, got %q", stderr.String())
	}
}

func TestDoctorUnreachableUpstreamFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doctor.yaml")
	content := []byte(`
server:
  port: 18080
routes:
  - path: /
    upstreams:
      - url: http://127.0.0.1:1
        weight: 1
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := doctor(&stdout, &stderr, []string{
		"-config", path,
		"-timeout", (50 * time.Millisecond).String(),
		"-json",
	})
	if err == nil {
		t.Fatal("expected doctor failure for unreachable upstream")
	}

	var report DoctorReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Failed == 0 {
		t.Fatalf("expected failing checks, got %+v", report)
	}
}

package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Upstream struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

type HeaderRules struct {
	Add    map[string]string `yaml:"add"`
	Remove []string          `yaml:"remove"`
}

type Route struct {
	Path      string        `yaml:"path"`
	Upstreams []Upstream    `yaml:"upstreams"`
	Strategy  string        `yaml:"strategy"`
	Timeout   time.Duration `yaml:"timeout"`

	RequestHeaders  HeaderRules `yaml:"request_headers"`
	ResponseHeaders HeaderRules `yaml:"response_headers"`
	CORS            CORSConfig  `yaml:"cors"`
}

type ServerConfig struct {
	Port    int `yaml:"port"`
	Retries int `yaml:"retries"`

	RateLimit RateLimitConfig `yaml:"rate_limit"`

	TLS TLSConfig `yaml:"tls"`

	Compression CompressionConfig `yaml:"compression"`
	Limits      LimitsConfig      `yaml:"limits"`
	Security    SecurityConfig    `yaml:"security"`

	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type RateLimitConfig struct {
	RPS   float64 `yaml:"rps"`
	Burst int     `yaml:"burst"`
}

// CompressionConfig controls optional gzip response compression.
// Disabled by default for backward compatibility.
type CompressionConfig struct {
	Enabled bool `yaml:"enabled"`
	// MinSize is the minimum response body size in bytes before compression.
	MinSize int `yaml:"min_size"`
}

// LimitsConfig holds request limits. Zero body_size means unlimited.
type LimitsConfig struct {
	BodySize int64 `yaml:"body_size"`
}

// SecurityConfig holds optional HTTP security headers.
type SecurityConfig struct {
	Headers SecurityHeaders `yaml:"headers"`
}

// SecurityHeaders are emitted only when non-empty.
type SecurityHeaders struct {
	XFrameOptions           string `yaml:"x_frame_options"`
	XContentTypeOptions     string `yaml:"x_content_type_options"`
	StrictTransportSecurity string `yaml:"strict_transport_security"`
	ReferrerPolicy          string `yaml:"referrer_policy"`
	ContentSecurityPolicy   string `yaml:"content_security_policy"`
	PermissionsPolicy       string `yaml:"permissions_policy"`
}

type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

type TLSConfig struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	Enabled  bool   `yaml:"enabled"`
}

func Load(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

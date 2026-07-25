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

	RateLimit struct {
		RPS   float64 `yaml:"rps"`
		Burst int     `yaml:"burst"`
	} `yaml:"rate_limit"`
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

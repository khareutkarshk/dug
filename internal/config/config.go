package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Upstream struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

type Route struct {
	Path      string     `yaml:"path"`
	Upstreams []Upstream `yaml:"upstreams"`
	Strategy  string     `yaml:"strategy"`
}

type ServerConfig struct {
	Port    int `yaml:"port"`
	Retries int `yaml:"retries"`

	RateLimit struct {
		RPS   float64 `yaml:"rps"`
		Burst int     `yaml:"burst"`
	} `yaml:"rate_limit"`
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

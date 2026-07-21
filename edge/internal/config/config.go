package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path      string   `yaml:"path"`
	Upstreams []string `yaml:"upstreams"`
}

type Serve struct {
	Port int `yaml:"port"`
}

type Config struct {
	Server Serve   `yaml:"server"`
	Routes []Route `yaml:"routes"`
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

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Route struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Routes []Route      `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &cfg, nil
}

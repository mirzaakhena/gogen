package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  bool   `yaml:"ssl_mode"`
	LogMode  bool   `yaml:"log_mode"`
}

func ReadConfig() (*Config, error) {

	bytes, err := os.ReadFile("infrastructure/config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
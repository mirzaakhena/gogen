package config

import (
	"encoding/json"
	"os"
)

func ReadConfig() *Config {

	bytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}

	var cfg Config

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err.Error())
	}

	return &cfg
}

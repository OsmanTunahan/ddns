package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DNSProvider string `json:"dns_provider"`
	APIKey      string `json:"api_key"`
	Email       string `json:"email"`
	Domain      string `json:"domain"`
	LogLevel    string `json:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

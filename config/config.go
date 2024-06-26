package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DNSProvider string `json:"dns_provider"`
	APIKey      string `json:"api_key"`
	Email       string `json:"email"`
	ProjectID   string `json:"project_id"`
	Domain      string `json:"domain"`
	LogLevel    string `json:"log_level"`
}

func LoadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

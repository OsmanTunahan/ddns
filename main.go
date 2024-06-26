package main

import (
	"ddns/config"
	"ddns/logger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.LogLevel)
}

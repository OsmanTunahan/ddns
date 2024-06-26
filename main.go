package main

import (
	"ddns/config"
	"ddns/logger"
	"ddns/service"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.LogLevel)

	// Start DDNS service
	s := service.NewService(cfg, log)
	s.Start()
}

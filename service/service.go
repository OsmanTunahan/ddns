package service

import (
	"ddns/config"
	"ddns/logger"
	"ddns/updater"
	"time"
)

type Service struct {
	config  *config.Config
	logger  *logger.Logger
	updater *updater.Updater
}

func NewService(cfg *config.Config, log *logger.Logger) *Service {
	return &Service{
		config:  cfg,
		logger:  log,
		updater: updater.NewUpdater(cfg, log),
	}
}

func (s *Service) Start() {
	s.logger.Println("Starting DDNS service...")

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.updater.UpdateDNS()
	}
}

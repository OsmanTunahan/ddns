package updater

import (
	"ddns/config"
	"ddns/dns"
	"ddns/ip"
	"ddns/logger"
)

type Updater struct {
	config     *config.Config
	logger     *logger.Logger
	ipResolver ip.IPResolver
	dnsClient  dns.DNSClient
}

func NewUpdater(cfg *config.Config, log *logger.Logger) *Updater {
	return &Updater{
		config:     cfg,
		logger:     log,
		ipResolver: ip.NewIPResolver(),
		dnsClient:  dns.NewDNSClient(cfg.DNSProvider, cfg.APIKey, cfg.Email),
	}
}

func (u *Updater) UpdateDNS() {
	currentIP, err := u.ipResolver.GetCurrentIP()
	if err != nil {
		u.logger.Printf("Error resolving IP: %v", err)
		return
	}

	if err := u.dnsClient.UpdateDNSRecord(u.config.Domain, currentIP); err != nil {
		u.logger.Printf("Error updating DNS record: %v", err)
		return
	}

	u.logger.Println("DNS record updated successfully")
}

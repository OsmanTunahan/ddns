package dns

type DNSClient interface {
	UpdateDNSRecord(domain string, ip string) error
}

type GenericDNSClient struct {
	provider  string
	apiKey    string
	email     string // Cloudflare
	projectID string // Google Cloud DNS
}

func NewDNSClient(provider, apiKey, email, projectID string) DNSClient {
	return &GenericDNSClient{
		provider:  provider,
		apiKey:    apiKey,
		email:     email,
		projectID: projectID,
	}
}

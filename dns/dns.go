package dns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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

func (g *GenericDNSClient) UpdateDNSRecord(domain, ip string) error {
	switch g.provider {
	case "cloudflare":
		return g.updateCloudflareDNS(domain, ip)
	case "google":
		return g.updateGoogleDNS(domain, ip)
	default:
		return errors.New("unsupported DNS provider")
	}
}

// Example of updating a DNS record for Cloudflare
func (g *GenericDNSClient) updateCloudflareDNS(domain, ip string) error {
	zoneID, err := g.getCloudflareZoneID(domain)
	if err != nil {
		return err
	}

	recordID, err := g.getCloudflareRecordID(zoneID, domain)
	if err != nil {
		return err
	}

	// Construct the request to update Cloudflare DNS record
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	body := map[string]string{
		"type":    "A",
		"name":    domain,
		"content": ip,
		"ttl":     "1",
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", g.email)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to update DNS record, status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (g *GenericDNSClient) getCloudflareZoneID(domain string) (string, error) {
	// Extract the base domain
	domainParts := strings.Split(domain, ".")
	baseDomain := domainParts[len(domainParts)-2] + "." + domainParts[len(domainParts)-1]

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones?name=%s", baseDomain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", g.email)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get zone ID, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if result["success"].(bool) {
		zones := result["result"].([]interface{})
		if len(zones) > 0 {
			return zones[0].(map[string]interface{})["id"].(string), nil
		}
	}

	return "", errors.New("zone ID not found")
}

func (g *GenericDNSClient) getCloudflareRecordID(zoneID, domain string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", zoneID, domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", g.email)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get record ID, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if result["success"].(bool) {
		records := result["result"].([]interface{})
		if len(records) > 0 {
			return records[0].(map[string]interface{})["id"].(string), nil
		}
	}

	return "", errors.New("record ID not found")
}

// Example of updating a DNS record for Google Cloud DNS
func (g *GenericDNSClient) updateGoogleDNS(domain, ip string) error {
	zoneID, err := g.getGoogleDNSZoneID(domain)
	if err != nil {
		return err
	}

	recordName := domain + "."

	// Construct the request to update Google Cloud DNS record
	url := fmt.Sprintf("https://dns.googleapis.com/dns/v1/projects/%s/managedZones/%s/rrsets/%s/A", g.projectID, zoneID, recordName)
	body := map[string]interface{}{
		"kind":    "dns#resourceRecordSet",
		"name":    recordName,
		"type":    "A",
		"ttl":     300,
		"rrdatas": []string{ip},
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to update DNS record, status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (g *GenericDNSClient) getGoogleDNSZoneID(domain string) (string, error) {
	url := fmt.Sprintf("https://dns.googleapis.com/dns/v1/projects/%s/managedZones", g.projectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get managed zones, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if managedZones, ok := result["managedZones"].([]interface{}); ok {
		for _, zone := range managedZones {
			zoneMap := zone.(map[string]interface{})
			if strings.Contains(domain, zoneMap["dnsName"].(string)) {
				return zoneMap["id"].(string), nil
			}
		}
	}

	return "", errors.New("managed zone ID not found")
}

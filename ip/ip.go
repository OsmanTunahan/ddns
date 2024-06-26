package ip

import (
	"io/ioutil"
	"net/http"
)

type IPResolver interface {
	GetCurrentIP() (string, error)
}

type ExternalIPResolver struct{}

func NewIPResolver() IPResolver {
	return &ExternalIPResolver{}
}

func (e *ExternalIPResolver) GetCurrentIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

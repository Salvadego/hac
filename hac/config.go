package hac

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Config struct {
	BaseURL       string
	Username      string
	Password      string
	Timeout       time.Duration
	SkipTLSVerify bool
	UserAgent     string
}

func (cfg *Config) sanitize() {
	if cfg.BaseURL == "" {
		panic("BaseURL is required")
	}

	cfg.BaseURL = trimSuffixSlash(cfg.BaseURL)
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Second * 30
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = "hac-go-client/1.0"
	}

}

func (cfg *Config) newHttpClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.SkipTLSVerify,
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}
}

func trimSuffixSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}

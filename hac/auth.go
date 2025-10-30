package hac

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type AuthService struct {
	client *HACClient
}

var csrfRe = regexp.MustCompile(`name=["']_csrf["']\s+value=["'](.+?)["']`)

func (s *AuthService) Login(ctx context.Context) error {
	s.client.clearSession()

	csrf, err := s.fetchInitialCSRF(ctx)
	if err != nil {
		return fmt.Errorf("fetch initial csrf: %w", err)
	}

	form := url.Values{}
	form.Set("j_username", s.client.username)
	form.Set("j_password", s.client.password)
	form.Set("_csrf", csrf)
	form.Set("_spring_security_remember_me", "true")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp, err := s.client.doRequest(ctx,
		http.MethodPost,
		"j_spring_security_check",
		strings.NewReader(form.Encode()),
		headers,
	)
	if err != nil {
		return fmt.Errorf("login POST failed: %w", err)
	}

	bodyStr, err := readAllBody(resp)
	if err != nil {
		return fmt.Errorf("read login response: %w", err)
	}

	if !validateLoginResponse(bodyStr) {
		resp2, err := s.client.doRequest(ctx, http.MethodGet, "", nil, nil)
		if err != nil {
			return fmt.Errorf("verify login: %w", err)
		}
		bodyStr2, err := readAllBody(resp2)
		if err != nil {
			return fmt.Errorf("read verify response: %w", err)
		}
		if !validateLoginResponse(bodyStr2) {
			return fmt.Errorf("login failed: invalid credentials")
		}
		bodyStr = bodyStr2
	}

	newCsrf, err := extractCSRFFromString(bodyStr)
	if err != nil {
		return fmt.Errorf("post-login csrf: %w", err)
	}
	s.client.csrf = newCsrf

	return nil
}

func (s *AuthService) fetchInitialCSRF(ctx context.Context) (string, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, "", nil, nil)
	if err != nil {
		return "", fmt.Errorf("get login page: %w", err)
	}

	bodyStr, err := readAllBody(resp)
	if err != nil {
		return "", fmt.Errorf("read login page: %w", err)
	}

	if strings.Contains(bodyStr, "503: This service is down for maintenance") ||
		strings.Contains(bodyStr, "SAP Commerce Cloud - Maintenance") {
		return "", fmt.Errorf("service down for maintenance")
	}

	csrf, err := extractCSRFFromString(bodyStr)
	if err != nil {
		return "", err
	}

	s.client.csrf = csrf
	return csrf, nil
}

func extractCSRFFromString(body string) (string, error) {
	m := csrfRe.FindStringSubmatch(body)
	if len(m) < 2 {
		return "", fmt.Errorf("csrf token not found")
	}
	return m[1], nil
}

func validateLoginResponse(body string) bool {
	return strings.Contains(body, "You're")
}

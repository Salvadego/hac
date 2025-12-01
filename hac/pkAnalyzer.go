package hac

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type PKAnalyzerService struct {
	client *HACClient
}

func NewPKAnalyzerService(c *HACClient) *PKAnalyzerService {
	return &PKAnalyzerService{client: c}
}

func (s *PKAnalyzerService) Analyze(ctx context.Context, pkRequest PKAnalyzeRequest) (*PKAnalyzeResponse, error) {
	form := s.client.buildForm(pkRequest)

	resp, err := s.client.doRequest(
		ctx,
		"POST",
		"platform/pkanalyzer/analyze",
		strings.NewReader(form.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze PK: %w", err)
	}

	body, err := readAllBody(resp)
	if err != nil {
		return nil, err
	}

	var out PKAnalyzeResponse
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		return nil, fmt.Errorf("failed to decode PK response: %w body: %s", err, body)
	}

	return &out, nil
}

package hac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GroovyService struct {
	client *HACClient
}

func (s *GroovyService) Execute(
	ctx context.Context,
	q GroovyRequest,
) (*GroovyResponse, error) {

	q = applyDefaults(q)
	form := s.client.buildForm(q)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/scripting/execute",
		strings.NewReader(form.Encode()),
		headers,
	)
	if err != nil {
		return nil, err
	}

	body, err := readAllBody(resp)
	if err != nil {
		return nil, err
	}

	var out GroovyResponse
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		return nil, fmt.Errorf("decode groovy response: %w", err)
	}

	return &out, nil
}

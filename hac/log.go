package hac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LogService struct {
	client *HACClient
}

func NewLogService(c *HACClient) *LogService {
	return &LogService{client: c}
}

func (s *LogService) GetCurrentLoggers(ctx context.Context) ([]Logger, error) {
	response, err := s.ChangeLogLevel(ctx, "root", LogLevelInfo)
	if err != nil {
		return nil, err
	}

	return response.Loggers, nil
}

func (s *LogService) GetLogLevels(ctx context.Context) ([]Level, error) {
	response, err := s.ChangeLogLevel(ctx, "root", LogLevelInfo)
	if err != nil {
		return nil, err
	}

	return response.Levels, nil
}

func (s *LogService) ChangeLogLevel(ctx context.Context, loggerName string, level LogLevelName) (*ChangeLogLevelResponse, error) {
	req := ChangeLogRequest{
		LoggerName:   loggerName,
		LogLevelName: level,
	}

	form := s.client.buildForm(req)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"platform/log4j/changeLevel",
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

	var result ChangeLogLevelResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, fmt.Errorf("decode ChangeLogLevel response %w", err)
	}

	return &result, nil
}

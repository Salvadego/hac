package hac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FlexService struct {
	client *HACClient
}

func NewFlexService(c *HACClient) *FlexService {
	return &FlexService{client: c}
}

func (s *FlexService) Execute(
	ctx context.Context,
	q FlexQuery,
	opts *FlexExecuteOptions,
) (*FlexSearchResponse, error) {

	if opts == nil {
		opts = &FlexExecuteOptions{}
	}

	q = applyDefaults(q)
	form := s.client.buildForm(q)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/flexsearch/execute",
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

	var result FlexSearchResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, fmt.Errorf("decode flex response: %w", err)
	}

	if opts == nil {
		opts = &FlexExecuteOptions{}
	}

	if opts.NoBlacklist {
		return &result, nil
	}

	if len(result.ResultList) == 0 {
		return &result, nil
	}

	valid := selectColumns(result.Headers, result.ResultList, opts.ColumnBlacklist)
	result.Headers = valid.Headers
	result.ResultList = valid.Rows

	return &result, nil
}

type sliced struct {
	Headers []string
	Rows    [][]string
}

func selectColumns(headers []string, rows [][]string, blacklist []string) sliced {
	validIdx := make([]int, 0)

	for i, h := range headers {
		if isBlacklisted(h, blacklist) {
			continue
		}
		if columnHasValue(rows, i) {
			validIdx = append(validIdx, i)
		}
	}

	out := sliced{
		Headers: make([]string, len(validIdx)),
		Rows:    make([][]string, len(rows)),
	}

	for newIdx, oldIdx := range validIdx {
		out.Headers[newIdx] = headers[oldIdx]
	}

	for r, row := range rows {
		out.Rows[r] = make([]string, len(validIdx))
		for newIdx, oldIdx := range validIdx {
			out.Rows[r][newIdx] = row[oldIdx]
		}
	}
	return out
}

func columnHasValue(rows [][]string, col int) bool {
	for _, r := range rows {
		if col < len(r) && r[col] != "" && r[col] != "null" {
			return true
		}
	}
	return false
}

func isBlacklisted(header string, blacklist []string) bool {
	n := normalize(header)
	for _, b := range blacklist {
		if n == normalize(b) {
			return true
		}
	}
	return false
}

func normalize(s string) string {
	s = strings.TrimPrefix(s, "p_")
	s = strings.TrimPrefix(s, "P_")
	return strings.ToLower(s)
}

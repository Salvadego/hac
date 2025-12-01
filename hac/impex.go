package hac

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"
)

type ImpexService struct {
	client *HACClient
}

func NewImpexService(c *HACClient) *ImpexService {
	return &ImpexService{client: c}
}

func (s *ImpexService) Import(
	ctx context.Context,
	q ImpexImportRequest,
) (string, error) {

	q = applyDefaults(q)
	form := s.client.buildForm(q)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/impex/import",
		strings.NewReader(form.Encode()),
		headers,
	)
	if err != nil {
		return "", err
	}

	body, err := readAllBody(resp)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(string(body))
	if doc.Error != nil {
		return "", doc.Error
	}

	resultTag := doc.Find("span", "id", "impexResult")
	if resultTag.Error != nil {
		resultTag = doc.Find("div", "class", "impexResult")
	}

	var result string
	if resultTag.Error != nil {
		result = ""
	} else {
		result = resultTag.Attrs()["data-result"]
		if result == "" {
			result = resultTag.FullText()
		}
	}

	return result, nil
}

func (s *ImpexService) FetchTypeAndAttributes(t TypeAttributesRequest) (*TypeAttributesResponse, error) {
	form := s.client.buildForm(t)

	req, err := s.client.doRequest(
		context.Background(),
		http.MethodPost,
		"console/impex/typeAndAttributes",
		strings.NewReader(form.Encode()),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	if err != nil {
		return nil, err
	}

	body, err := readAllBody(req)
	if err != nil {
		return nil, err
	}

	var out TypeAttributesResponse
	if e := json.Unmarshal([]byte(body), &out); e != nil {
		return nil, fmt.Errorf("decode failed: %w body: %s", e, body)
	}
	return &out, nil
}

func (s *ImpexService) Export(ctx context.Context, q ImpexExportRequest) (string, string, error) {
	q = applyDefaults(q)
	form := s.client.buildForm(q)

	resp, err := s.client.doRequest(
		ctx,
		http.MethodPost,
		"console/impex/export",
		strings.NewReader(form.Encode()),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)

	if err != nil {
		return "", "", fmt.Errorf("failed to execute impex export: %w", err)
	}

	body, err := readAllBody(resp)
	if err != nil {
		return "", "", err
	}

	doc := soup.HTMLParse(string(body))
	if doc.Error != nil {
		return "", "", doc.Error
	}

	resultTag := doc.Find("span", "id", "impexResult")

	if resultTag.Error != nil {
		return "", "", fmt.Errorf("missing impexResult tag")
	}

	result := resultTag.Attrs()["data-result"]

	parent := doc.Find("div", "id", "downloadExportResultData")
	if parent.Error != nil {
		return result, "", fmt.Errorf("missing downloadExportResultData container")
	}

	link := parent.Find("a")
	if link.Error != nil {
		return result, "", fmt.Errorf("missing download link inside container")
	}

	href, ok := link.Attrs()["href"]

	if !ok || href == "" {
		return result, "", fmt.Errorf("download link missing href attribute")
	}

	return result, href, nil
}

func (s *ImpexService) DownloadExportZip(downloadPath string) ([]byte, error) {
	base := s.client.baseURL
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}

	full := base + "console/impex/" + downloadPath

	req, err := http.NewRequest(http.MethodGet, full, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", base+"console/impex/export")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", s.client.userAgent)
	if s.client.csrf != "" {
		req.Header.Set("X-CSRF-TOKEN", s.client.csrf)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "zip") && !strings.Contains(ct, "octet-stream") {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("not zip. status %d body:\n%s", resp.StatusCode, b)
	}

	return io.ReadAll(resp.Body)
}

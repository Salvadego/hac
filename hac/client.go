package hac

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type HACClient struct {
	baseURL    string
	username   string
	password   string
	userAgent  string
	httpClient *http.Client
	jar        *cookiejar.Jar
	csrf       string

	Auth   *AuthService
	Flex   *FlexService
	Groovy *GroovyService
	Impex  *ImpexService
	PKA    *PKAnalyzerService
	Log    *LogService
}

func NewClient(cfg *Config) *HACClient {
	cfg.sanitize()

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize cookie jar: %s", err.Error()))
	}

	httpClient := cfg.newHttpClient()
	httpClient.Jar = jar

	client := &HACClient{
		baseURL:    cfg.BaseURL,
		username:   cfg.Username,
		password:   cfg.Password,
		userAgent:  cfg.UserAgent,
		httpClient: httpClient,
		jar:        jar,
	}

	client.Auth = NewAuthService(client)
	client.Flex = NewFlexService(client)
	client.Groovy = NewGroovyService(client)
	client.Impex = NewImpexService(client)
	client.PKA = NewPKAnalyzerService(client)
	client.Log = NewLogService(client)
	return client
}

func (c *HACClient) doRequest(
	ctx context.Context,
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.csrf != "" {
		req.Header.Set("X-CSRF-TOKEN", c.csrf)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return c.doRequestRetry405(ctx, method, path, body, headers)
	}

	return resp, nil
}

func (c *HACClient) clearSession() {
	newJar, _ := cookiejar.New(nil)
	c.httpClient.Jar = newJar
	c.jar = newJar
}

func readAllBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func (c *HACClient) doRequestRetry405(
	ctx context.Context,
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*http.Response, error) {
	c.clearSession()

	if c.Auth != nil {
		_, _ = c.Auth.fetchInitialCSRF(ctx)
	}

	return c.doRequest(ctx, method, path, body, headers)
}

func applyDefaults[T any](q T) T {
	rv := reflect.ValueOf(&q).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		d := f.Tag.Get("default")
		if d == "" {
			continue
		}
		v := rv.Field(i)
		if !v.IsZero() {
			continue
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(d)
		case reflect.Int, reflect.Int64:
			n, _ := strconv.Atoi(d)
			v.SetInt(int64(n))
		case reflect.Bool:
			b, _ := strconv.ParseBool(d)
			v.SetBool(b)
		}
	}
	return q
}

func (c *HACClient) buildForm(v any) url.Values {
	form := url.Values{}

	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if rt.Kind() != reflect.Struct {
		panic("buildForm expects struct")
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("form")
		if tag == "" {
			continue
		}

		keepZero := false
		name := tag

		if strings.Contains(tag, ",") {
			parts := strings.Split(tag, ",")
			name = parts[0]
			for _, p := range parts[1:] {
				if p == "keepzero" {
					keepZero = true
				}
			}
		}

		val := rv.Field(i)

		if val.IsZero() && !keepZero {
			continue
		}

		switch val.Kind() {
		case reflect.Bool:
			form.Set(name, strconv.FormatBool(val.Bool()))
		case reflect.Int, reflect.Int64:
			form.Set(name, fmt.Sprintf("%d", val.Int()))
		case reflect.String:
			form.Set(name, val.String())
		}
	}

	return form
}

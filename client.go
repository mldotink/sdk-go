package ink

import (
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

const (
	DefaultBaseURL = "https://api.ml.ink/graphql"
	DefaultExecURL = "wss://exec-eu-central-1.ml.ink"
)

type Config struct {
	APIKey     string
	BaseURL    string
	ExecURL    string
	HTTPClient *http.Client
}
type Client struct {
	gql     graphql.Client
	execURL string
}

func NewClient(cfg Config) *Client {
	if cfg.APIKey == "" {
		panic("ink: APIKey is required")
	}
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	execURL := cfg.ExecURL
	if execURL == "" {
		execURL = DefaultExecURL
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = &authTransport{
		apiKey: cfg.APIKey,
		base:   transportOrDefault(httpClient.Transport),
	}
	return &Client{
		gql:     graphql.NewClient(baseURL, httpClient),
		execURL: execURL,
	}
}
func (c *Client) ExecBaseURL() string { return c.execURL }

type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.base.RoundTrip(req)
}

func transportOrDefault(t http.RoundTripper) http.RoundTripper {
	if t != nil {
		return t
	}
	return http.DefaultTransport
}
func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func optInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

type Error struct {
	Message    string         `json:"message"`
	Path       []string       `json:"path"`
	Extensions map[string]any `json:"extensions"`
}

func (e *Error) Error() string { return e.Message }

type Errors []*Error

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Message
	}
	return strings.Join(msgs, "; ")
}

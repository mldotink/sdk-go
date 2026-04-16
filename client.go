// Package ink provides a Go client for the Ink platform API (ml.ink).
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

// Config configures an Ink API client.
type Config struct {
	// APIKey is required. Create one at https://ml.ink/account/api-keys.
	APIKey string

	// BaseURL overrides the GraphQL endpoint. Default: https://api.ml.ink/graphql
	BaseURL string

	// ExecURL overrides the exec-proxy WebSocket endpoint.
	ExecURL string

	// HTTPClient overrides the underlying HTTP client. Auth headers are added
	// automatically; do not set Authorization in its default headers.
	HTTPClient *http.Client
}

// Client is an Ink platform API client.
type Client struct {
	gql     graphql.Client
	execURL string
}

// NewClient creates a new Ink API client.
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

// ExecBaseURL returns the configured exec-proxy WebSocket base URL.
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

// optStr converts an empty string to nil and a non-empty string to a pointer,
// matching the GraphQL convention where null and "" have different meanings.
func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// optInt converts 0 to nil and a non-zero int to a pointer, matching the
// GraphQL convention where null means "use server default".
func optInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// Error is a GraphQL error returned by the Ink API.
type Error struct {
	Message    string         `json:"message"`
	Path       []string       `json:"path"`
	Extensions map[string]any `json:"extensions"`
}

func (e *Error) Error() string { return e.Message }

// Errors is a list of GraphQL errors.
type Errors []*Error

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Message
	}
	return strings.Join(msgs, "; ")
}

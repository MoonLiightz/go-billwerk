package optimize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/moonliightz/go-billwerk/pkg/request"
	"io"
	"net/http"
	"time"
)

// BaseURL is the base URL for all Billwerk Optimize API requests.
var BaseURL = "https://api.reepay.com/v1"

// Billwerk represents the API client object.
type Billwerk struct {
	apiKey     string
	apiKeyB64  string
	httpClient *http.Client
}

// Option is a function that sets options for the Billwerk client configuration.
type Option func(billwerk *Billwerk)

// WithHTTPClient allows setting a custom HTTP client for the Billwerk client.
func WithHTTPClient(client *http.Client) Option {
	return func(billwerk *Billwerk) {
		billwerk.httpClient = client
	}
}

// New creates a new Billwerk client with an API key and optional configuration options.
func New(apiKey string, opts ...Option) *Billwerk {
	b := &Billwerk{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// newBillwerkRequest creates a new HTTP request with base URL and authentication.
func (b *Billwerk) newBillwerkRequest(ctx context.Context) request.Builder {
	return request.New(ctx).
		WithBaseURL(BaseURL).
		WithBasicAuth(b.apiKey, "").
		WithHeader("Accept", "application/json; charset=utf-8")
}

// Do executes an HTTP request and json decodes the response into v (if provided).
//
// The function checks the status code of the response and returns an error
// if the status code indicates a failure (4xx or 5xx).
// If v is not nil, the response body is json decoded into the provided value.
func (b *Billwerk) Do(req *http.Request, v interface{}) error {
	res, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errRes
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}

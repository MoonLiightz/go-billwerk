package request

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/moonliightz/go-billwerk"
	"io"
	"net/http"
	"net/url"
)

const (
	UserAgent = "go-billwerk/" + billwerk.Version
)

// Builder is an HTTP request builder.
type Builder interface {
	// WithBaseURL sets the base URL for the request.
	// Should include the scheme and hostname.
	//
	// Example: WithBaseURL("https://api.reepay.com")
	WithBaseURL(baseURL string) Builder

	// WithEndpoint sets the endpoint for the request.
	// The endpoint is appended to the base URL.
	//
	// Example: WithEndpoint("/v1/plans")
	WithEndpoint(endpoint string) Builder

	// WithParam sets a query parameter for the request.
	// If the key already exists, it is overwritten.
	//
	// Example: WithParam("page", "1") // Results in ?page=1
	WithParam(key, value string) Builder

	// AddParam adds a query parameter for the request.
	// If the key already exists, the value is appended.
	//
	// Example:
	//  AddParam("filter", "active") // Results in ?filter=active,
	//  AddParam("filter", "inactive") // Results in ?filter=active&filter=inactive
	AddParam(key, value string) Builder

	// WithHeader sets a header for the request.
	// If the key already exists, it is overwritten.
	WithHeader(key, value string) Builder

	// WithBasicAuth sets the Authorization header to a Basic Auth value.
	WithBasicAuth(username, password string) Builder

	// WithContentTypeJSON sets the Content-Type header to application/json; charset=utf-8.
	// Use this method for requests that send JSON payloads, which is a common standard for REST APIs.
	WithContentTypeJSON() Builder

	// WithBody sets the request body to the provided io.Reader.
	WithBody(body io.Reader) Builder

	// WithJSONBody sets the request body to the JSON encoding of v.
	// The Content-Type header is set to application/json; charset=utf-8.
	// Note: Any encoding errors are silently ignored. Ensure that v is JSON-serializable.
	WithJSONBody(v interface{}) Builder

	// GET builds an HTTP GET request.
	GET() (*http.Request, error)

	// POST builds an HTTP POST request.
	POST() (*http.Request, error)

	// PUT builds an HTTP PUT request.
	PUT() (*http.Request, error)

	// DELETE builds an HTTP DELETE request.
	DELETE() (*http.Request, error)

	// build builds the HTTP request.
	build() (*http.Request, error)
}

type request struct {
	ctx      context.Context
	method   string
	baseURL  string
	endpoint string
	header   http.Header
	params   url.Values
	body     io.Reader
}

func New(ctx context.Context) Builder {
	return &request{
		ctx:    ctx,
		method: http.MethodGet,
		header: make(http.Header),
		params: make(url.Values),
	}
}

func (r *request) WithBaseURL(baseURL string) Builder {
	r.baseURL = baseURL
	return r
}

func (r *request) WithEndpoint(endpoint string) Builder {
	r.endpoint = endpoint
	return r
}

func (r *request) WithParam(key, value string) Builder {
	r.params.Set(key, value)
	return r
}

func (r *request) AddParam(key, value string) Builder {
	r.params.Add(key, value)
	return r
}

func (r *request) WithHeader(key, value string) Builder {
	r.header.Set(key, value)
	return r
}

func (r *request) WithBasicAuth(username, password string) Builder {
	rawCredentials := fmt.Sprintf("%s:%s", username, password)
	credentials := base64.StdEncoding.EncodeToString([]byte(rawCredentials))
	value := fmt.Sprintf("Basic %s", credentials)
	return r.WithHeader("Authorization", value)
}

func (r *request) WithContentTypeJSON() Builder {
	return r.WithHeader("Content-Type", "application/json; charset=utf-8")
}

func (r *request) WithBody(body io.Reader) Builder {
	r.body = body
	return r
}

func (r *request) WithJSONBody(v interface{}) Builder {
	r.WithContentTypeJSON()
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(v)
	r.body = buf
	return r
}

func (r *request) GET() (*http.Request, error) {
	r.method = http.MethodGet
	return r.build()
}

func (r *request) POST() (*http.Request, error) {
	r.method = http.MethodPost
	return r.build()
}

func (r *request) PUT() (*http.Request, error) {
	r.method = http.MethodPut
	return r.build()
}

func (r *request) DELETE() (*http.Request, error) {
	r.method = http.MethodDelete
	return r.build()
}

func (r *request) build() (*http.Request, error) {
	fullURL := r.baseURL + r.endpoint
	req, err := http.NewRequestWithContext(r.ctx, r.method, fullURL, r.body)
	if err != nil {
		return nil, err
	}

	r.header.Set("User-Agent", UserAgent)

	req.Header = r.header
	req.URL.RawQuery = r.params.Encode()

	return req, nil
}

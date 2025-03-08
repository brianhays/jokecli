package testutils

import (
	"io"
	"net/http"
	"strings"
)

// RoundTripFunc is a type that implements http.RoundTripper interface
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip executes the mock round trip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns a mock http client
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// MockResponse creates a mock HTTP response
func MockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

package jokesapi

import "net/http"

// HTTPClient interface allows us to mock the HTTP client for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	// DefaultClient is the default HTTP client used for fetching jokes.
	// It's defined as a variable so it can be replaced during testing.
	DefaultClient HTTPClient = &http.Client{}
)

// SetCommonHeaders sets the common headers required by the joke APIs.
func SetCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	// It's good practice to identify your client via User-Agent
	req.Header.Set("User-Agent", "jokecli (https://github.com/brianhays/jokecli)")
}

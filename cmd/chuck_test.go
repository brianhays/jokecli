package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/brianhays/jokecli/internal/jokesapi"
	"github.com/brianhays/jokecli/internal/testutils"
)

// MockChuckClient is a mock HTTP client for testing Chuck Norris jokes
type MockChuckClient struct {
	statusCode int
	body       string
	err        error
}

// Do implements the jokesapi.HTTPClient interface for MockChuckClient
func (m *MockChuckClient) Do(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(m.body))),
		Header:     make(http.Header),
	}, nil
}

func TestGetChuckNorrisJoke(t *testing.T) {
	tests := []struct {
		name                string
		client              jokesapi.HTTPClient
		want                *jokesapi.ChuckNorrisJoke
		wantErr             bool
		expectedErrContains string
	}{
		{
			name:    "successful joke fetch",
			client:  &MockChuckClient{statusCode: http.StatusOK, body: `{"id": "1", "value": "Chuck Norris can divide by zero.", "url": ""}`},
			want:    &jokesapi.ChuckNorrisJoke{ID: "1", Value: "Chuck Norris can divide by zero.", URL: ""},
			wantErr: false,
		},
		{
			name:                "server error",
			client:              &MockChuckClient{statusCode: http.StatusInternalServerError, body: "Internal Server Error"},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "Chuck Norris API returned unexpected status code: 500",
		},
		{
			name:                "network error",
			client:              &MockChuckClient{err: fmt.Errorf("network error")},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "failed to fetch Chuck Norris joke: network error",
		},
		{
			name:                "invalid json response",
			client:              &MockChuckClient{statusCode: http.StatusOK, body: `invalid json`},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "failed to parse Chuck Norris joke:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jokesapi.GetChuckNorrisJoke(tt.client)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetChuckNorrisJoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErrContains) {
					t.Errorf("GetChuckNorrisJoke() error = %v, want error containing %q", err, tt.expectedErrContains)
				}
			}

			// Use basic comparison for non-error cases
			if !tt.wantErr && (got == nil || tt.want == nil || got.ID != tt.want.ID || got.Value != tt.want.Value || got.URL != tt.want.URL) {
				t.Errorf("GetChuckNorrisJoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChuckCommand(t *testing.T) {
	// Enable testing mode
	isTesting = true
	defer func() { isTesting = false }()

	// Save the original client and restore it after the test
	originalClient := jokesapi.DefaultClient
	defer func() { jokesapi.DefaultClient = originalClient }()

	// Use the correct struct type from the jokesapi package
	mockJoke := jokesapi.ChuckNorrisJoke{
		ID:    "abc",
		Value: "Test Chuck Norris Fact",
		URL:   "http://example.com",
	}

	mockJSON, err := json.Marshal(mockJoke)
	if err != nil {
		t.Fatal(err)
	}

	// Set the default client to our test client
	jokesapi.DefaultClient = testutils.NewTestClient(func(req *http.Request) *http.Response {
		return testutils.MockResponse(http.StatusOK, string(mockJSON))
	})

	// Test the command execution
	if err := chuckCmd.Execute(); err != nil {
		t.Errorf("chuckCmd.Execute() error = %v", err)
	}
}

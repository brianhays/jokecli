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

// MockDadClient is a mock HTTP client for testing Dad jokes
type MockDadClient struct {
	statusCode int
	body       string
	err        error
}

// Do implements the jokesapi.HTTPClient interface for MockDadClient
func (m *MockDadClient) Do(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(m.body))),
		Header:     make(http.Header),
	}, nil
}

func TestGetDadJoke(t *testing.T) {
	tests := []struct {
		name                string
		client              jokesapi.HTTPClient // Use the interface type
		want                *jokesapi.DadJoke
		wantErr             bool
		expectedErrContains string // Check for substring in error
	}{
		{
			name:    "successful joke fetch",
			client:  &MockDadClient{statusCode: http.StatusOK, body: `{"id": "1", "joke": "Test Joke", "status": 200}`},
			want:    &jokesapi.DadJoke{ID: "1", Joke: "Test Joke", Status: 200},
			wantErr: false,
		},
		{
			name:                "server error",
			client:              &MockDadClient{statusCode: http.StatusInternalServerError, body: "Internal Server Error"},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "dad joke API returned unexpected status code: 500",
		},
		{
			name:                "network error",
			client:              &MockDadClient{err: fmt.Errorf("network error")},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "failed to fetch dad joke: network error",
		},
		{
			name:                "invalid json response",
			client:              &MockDadClient{statusCode: http.StatusOK, body: `invalid json`},
			want:                nil,
			wantErr:             true,
			expectedErrContains: "failed to parse dad joke:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function directly with the mock client
			got, err := jokesapi.GetDadJoke(tt.client)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDadJoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.expectedErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErrContains) {
					t.Errorf("GetDadJoke() error = %v, want error containing %q", err, tt.expectedErrContains)
				}
			}

			// Use basic comparison for non-error cases, consider deep equal if structs get complex
			if !tt.wantErr && (got == nil || tt.want == nil || got.ID != tt.want.ID || got.Joke != tt.want.Joke || got.Status != tt.want.Status) {
				t.Errorf("GetDadJoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDadCommand(t *testing.T) {
	// Enable testing mode
	isTesting = true
	defer func() { isTesting = false }()

	// Save the original client and restore it after the test
	originalClient := jokesapi.DefaultClient
	defer func() { jokesapi.DefaultClient = originalClient }()

	// Use the correct struct type from the jokesapi package
	mockJoke := jokesapi.DadJoke{
		ID:     "123",
		Joke:   "Test joke",
		Status: 200,
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
	if err := dadCmd.Execute(); err != nil {
		t.Errorf("dadCmd.Execute() error = %v", err)
	}
}

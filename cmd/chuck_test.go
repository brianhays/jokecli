package cmd

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brianhays/jokecli/internal/testutils"
)

func TestGetChuckNorrisJoke(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantJoke       string
		wantErr        bool
	}{
		{
			name: "successful joke fetch",
			mockResponse: `{
				"id": "abc123",
				"value": "Chuck Norris can divide by zero.",
				"url": "https://api.chucknorris.io/jokes/abc123"
			}`,
			mockStatusCode: http.StatusOK,
			wantJoke:       "Chuck Norris can divide by zero.",
			wantErr:        false,
		},
		{
			name:           "invalid json response",
			mockResponse:   `{"invalid json"}`,
			mockStatusCode: http.StatusOK,
			wantJoke:       "",
			wantErr:        true,
		},
		{
			name:           "server error",
			mockResponse:   `{"message": "Internal Server Error"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantJoke:       "",
			wantErr:        false,
		},
		{
			name: "missing required fields",
			mockResponse: `{
				"id": "abc123",
				"url": "https://api.chucknorris.io/jokes/abc123"
			}`,
			mockStatusCode: http.StatusOK,
			wantJoke:       "",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testutils.NewTestClient(func(req *http.Request) *http.Response {
				// Test request headers
				if req.Header.Get("Accept") != "application/json" {
					t.Error("Accept header not set correctly")
				}
				if req.Header.Get("User-Agent") == "" {
					t.Error("User-Agent header not set")
				}

				return testutils.MockResponse(tt.mockStatusCode, tt.mockResponse)
			})

			got, err := getChuckNorrisJoke(client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChuckNorrisJoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil && got.Value != tt.wantJoke {
				t.Errorf("getChuckNorrisJoke() = %v, want %v", got.Value, tt.wantJoke)
			}
		})
	}
}

func TestChuckCommand(t *testing.T) {
	// Enable testing mode
	isTesting = true
	defer func() { isTesting = false }()

	// Save the original client and restore it after the test
	originalClient := defaultClient
	defer func() { defaultClient = originalClient }()

	mockJoke := ChuckNorrisJoke{
		ID:    "abc123",
		Value: "Test Chuck Norris fact",
		URL:   "https://api.chucknorris.io/jokes/abc123",
	}

	mockJSON, err := json.Marshal(mockJoke)
	if err != nil {
		t.Fatal(err)
	}

	defaultClient = testutils.NewTestClient(func(req *http.Request) *http.Response {
		return testutils.MockResponse(http.StatusOK, string(mockJSON))
	})

	// Test the command execution
	if err := chuckCmd.Execute(); err != nil {
		t.Errorf("chuckCmd.Execute() error = %v", err)
	}
}

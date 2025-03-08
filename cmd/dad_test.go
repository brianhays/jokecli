package cmd

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bhays/jokecli/internal/testutils"
)

func TestGetDadJoke(t *testing.T) {
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
				"id": "123",
				"joke": "Why don't eggs tell jokes? They'd crack up!",
				"status": 200
			}`,
			mockStatusCode: http.StatusOK,
			wantJoke:       "Why don't eggs tell jokes? They'd crack up!",
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

			got, err := getDadJoke(client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDadJoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Joke != tt.wantJoke {
				t.Errorf("getDadJoke() = %v, want %v", got.Joke, tt.wantJoke)
			}
		})
	}
}

func TestDadCommand(t *testing.T) {
	// Enable testing mode
	isTesting = true
	defer func() { isTesting = false }()

	// Save the original client and restore it after the test
	originalClient := defaultClient
	defer func() { defaultClient = originalClient }()

	mockJoke := DadJoke{
		ID:     "123",
		Joke:   "Test joke",
		Status: 200,
	}

	mockJSON, err := json.Marshal(mockJoke)
	if err != nil {
		t.Fatal(err)
	}

	defaultClient = testutils.NewTestClient(func(req *http.Request) *http.Response {
		return testutils.MockResponse(http.StatusOK, string(mockJSON))
	})

	// Test the command execution
	if err := dadCmd.Execute(); err != nil {
		t.Errorf("dadCmd.Execute() error = %v", err)
	}
}

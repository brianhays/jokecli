package jokesapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DadJoke represents the structure of a dad joke from the API.
type DadJoke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// GetDadJoke fetches a random dad joke using the provided HTTP client.
func GetDadJoke(client HTTPClient) (*DadJoke, error) {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dad joke request: %w", err)
	}

	SetCommonHeaders(req) // Use the helper from client.go

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dad joke: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dad joke API returned unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read dad joke response: %w", err)
	}

	var joke DadJoke
	if err := json.Unmarshal(body, &joke); err != nil {
		return nil, fmt.Errorf("failed to parse dad joke: %w", err)
	}

	return &joke, nil
}

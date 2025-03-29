package jokesapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChuckNorrisJoke represents the structure of a Chuck Norris joke from the API.
type ChuckNorrisJoke struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	URL   string `json:"url"`
}

// GetChuckNorrisJoke fetches a random Chuck Norris joke using the provided HTTP client.
func GetChuckNorrisJoke(client HTTPClient) (*ChuckNorrisJoke, error) {
	req, err := http.NewRequest("GET", "https://api.chucknorris.io/jokes/random", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Chuck Norris joke request: %w", err)
	}

	SetCommonHeaders(req) // Use the helper from client.go

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Chuck Norris joke: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Chuck Norris API returned unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Chuck Norris joke response: %w", err)
	}

	var joke ChuckNorrisJoke
	if err := json.Unmarshal(body, &joke); err != nil {
		return nil, fmt.Errorf("failed to parse Chuck Norris joke: %w", err)
	}

	return &joke, nil
}

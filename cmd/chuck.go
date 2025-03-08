package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type ChuckNorrisJoke struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	URL   string `json:"url"`
}

var chuckCmd = &cobra.Command{
	Use:   "chuck",
	Short: "Get a random Chuck Norris fact",
	Long:  `Fetches a random Chuck Norris fact from api.chucknorris.io`,
	RunE: func(cmd *cobra.Command, args []string) error {
		joke, err := getChuckNorrisJoke(defaultClient)
		if err != nil {
			return err
		}
		fmt.Println(joke.Value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chuckCmd)
}

func getChuckNorrisJoke(client httpClient) (*ChuckNorrisJoke, error) {
	req, err := http.NewRequest("GET", "https://api.chucknorris.io/jokes/random", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jokecli (https://github.com/bhays/jokecli)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch joke: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var joke ChuckNorrisJoke
	if err := json.Unmarshal(body, &joke); err != nil {
		return nil, fmt.Errorf("failed to parse joke: %w", err)
	}

	return &joke, nil
}

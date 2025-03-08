package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type DadJoke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// httpClient interface allows us to mock the HTTP client for testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	// defaultClient is the default HTTP client used for production
	defaultClient httpClient = &http.Client{}
)

var dadCmd = &cobra.Command{
	Use:   "dad",
	Short: "Get a random dad joke",
	Long:  `Fetches a random dad joke from icanhazdadjoke.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		joke, err := getDadJoke(defaultClient)
		if err != nil {
			return err
		}
		fmt.Println(joke.Joke)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dadCmd)
}

func getDadJoke(client httpClient) (*DadJoke, error) {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Required headers for icanhazdadjoke API
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jokecli (https://github.com/brianhays/jokecli)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch joke: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var joke DadJoke
	if err := json.Unmarshal(body, &joke); err != nil {
		return nil, fmt.Errorf("failed to parse joke: %w", err)
	}

	return &joke, nil
}

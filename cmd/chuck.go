package cmd

import (
	"fmt"

	"github.com/brianhays/jokecli/internal/jokesapi"
	"github.com/spf13/cobra"
)

var chuckCmd = &cobra.Command{
	Use:   "chuck",
	Short: "Get a random Chuck Norris fact",
	Long:  `Fetches a random Chuck Norris fact from api.chucknorris.io`,
	RunE: func(cmd *cobra.Command, args []string) error {
		joke, err := jokesapi.GetChuckNorrisJoke(jokesapi.DefaultClient)
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

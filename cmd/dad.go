package cmd

import (
	"fmt"

	"github.com/brianhays/jokecli/internal/jokesapi"
	"github.com/spf13/cobra"
)

var dadCmd = &cobra.Command{
	Use:   "dad",
	Short: "Get a random dad joke",
	Long:  `Fetches a random dad joke from icanhazdadjoke.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		joke, err := jokesapi.GetDadJoke(jokesapi.DefaultClient)
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

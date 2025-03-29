package cmd

import (
	"fmt"
	"os"

	"github.com/brianhays/jokecli/internal/jokesapi"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// testing flag to bypass interactive mode
var isTesting bool

var rootCmd = &cobra.Command{
	Use:   "jokecli",
	Short: "A CLI tool for fetching jokes from various sources",
	Long: `jokecli is a command line interface that provides jokes and funny facts
from various sources across the internet. You can get Chuck Norris facts,
dad jokes, and more!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isTesting {
			return nil
		}
		return runInteractiveMode()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runInteractiveMode() error {
	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What kind of joke would you like?").
				Options(
					huh.NewOption("Chuck Norris Fact", "chuck"),
					huh.NewOption("Dad Joke", "dad"),
				).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("failed to run interactive mode: %w", err)
	}

	switch selected {
	case "chuck":
		joke, err := jokesapi.GetChuckNorrisJoke(jokesapi.DefaultClient)
		if err != nil {
			return err
		}
		fmt.Println(joke.Value)
	case "dad":
		joke, err := jokesapi.GetDadJoke(jokesapi.DefaultClient)
		if err != nil {
			return err
		}
		fmt.Println(joke.Joke)
	default:
		return fmt.Errorf("invalid selection: %s", selected)
	}

	return nil
}

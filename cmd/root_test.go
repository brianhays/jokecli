package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test help output
	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Add root command as a subcommand to test command to capture output
	cmd.AddCommand(rootCmd)

	// Execute help command
	cmd.SetArgs([]string{"jokecli", "--help"})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if help output contains expected strings
	expectedStrings := []string{
		"jokecli",
		"Usage:",
		"Available Commands:",
		"chuck       Get a random Chuck Norris fact",
		"dad         Get a random dad joke",
		"help        Help about any command",
	}

	output := buf.String()
	for _, str := range expectedStrings {
		if !strings.Contains(output, str) {
			t.Errorf("Expected help output to contain %q", str)
		}
	}
}

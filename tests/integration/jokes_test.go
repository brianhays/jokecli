package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func getProjectRoot(t *testing.T) string {
	// Get the project root directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up until we find go.mod
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("Could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func buildTestBinary(t *testing.T, projectRoot string) {
	buildCmd := exec.Command("go", "build", "-o", "jokecli")
	buildCmd.Dir = projectRoot
	buildCmd.Env = append(os.Environ(), "GO111MODULE=on")
	var stderr bytes.Buffer
	buildCmd.Stderr = &stderr
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v\nStderr: %s", err, stderr.String())
	}
}

func TestJokeCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	projectRoot := getProjectRoot(t)

	tests := []struct {
		name    string
		command []string
		verify  func(t *testing.T, output string)
	}{
		{
			name:    "chuck norris joke",
			command: []string{"chuck"},
			verify: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty Chuck Norris joke")
				}
			},
		},
		{
			name:    "dad joke",
			command: []string{"dad"},
			verify: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty dad joke")
				}
			},
		},
		{
			name:    "help command",
			command: []string{"--help"},
			verify: func(t *testing.T, output string) {
				expectedStrings := []string{
					"jokecli",
					"Usage:",
					"Available Commands:",
					"chuck       Get a random Chuck Norris fact",
					"dad         Get a random dad joke",
				}
				for _, str := range expectedStrings {
					if !strings.Contains(output, str) {
						t.Errorf("Expected output to contain %q", str)
					}
				}
			},
		},
	}

	// Build the binary for testing
	buildTestBinary(t, projectRoot)
	defer os.Remove(filepath.Join(projectRoot, "jokecli"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add timeout for API calls
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create command with context
			cmd := exec.CommandContext(ctx, "./jokecli", tt.command...)
			cmd.Dir = projectRoot
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Run command
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}

			output := stdout.String()
			tt.verify(t, output)
		})
	}
}

// TestRateLimiting verifies we don't exceed API rate limits
func TestRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limit test in short mode")
	}

	projectRoot := getProjectRoot(t)

	// Build the binary for testing
	buildTestBinary(t, projectRoot)
	defer os.Remove(filepath.Join(projectRoot, "jokecli"))

	// Make multiple requests in quick succession
	for i := 0; i < 5; i++ {
		cmd := exec.Command("./jokecli", "chuck")
		cmd.Dir = projectRoot
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Request %d failed: %v\nStderr: %s", i+1, err, stderr.String())
		}
		// Add small delay between requests
		time.Sleep(500 * time.Millisecond)
	}
}

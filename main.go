package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/hamzabow/co/internal/apikeyinput"
	"github.com/hamzabow/co/internal/config"
	"github.com/hamzabow/co/internal/genmessage"
	"github.com/hamzabow/co/internal/messagetextarea"
)

// Provider constants for API keys
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
	ProviderOllama    = "ollama"
)

// Define error style
var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF5555")).
	MarginTop(1).
	MarginBottom(1).
	MarginLeft(1).
	MarginRight(1)

func displayError(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Println(errorStyle.Render("Error: " + errorMsg))
	os.Exit(1)
}

func main() {
	// Try to get API key from environment variable first
	key := os.Getenv("OPENAI_API_KEY")

	// If not in environment, try to load from config file
	if key == "" {
		var err error
		key, err = config.LoadAPIKey(ProviderOpenAI)
		if err != nil && err != config.ErrNoConfigFile && err != config.ErrProviderNotFound {
			displayError("Failed to load API key from config: %v", err)
		}

		// If still no key, prompt the user
		if key == "" {
			key, err = apikeyinput.PromptApiKeyWithRetries()
			if err != nil {
				if err == apikeyinput.ErrEmptyApiKey {
					fmt.Println("No API key provided. Exiting.")
					os.Exit(1)
				} else {
					displayError("%v", err)
				}
			}

			// Save the key to config for future use
			if key != "" {
				if err := config.SaveAPIKey(ProviderOpenAI, key); err != nil {
					fmt.Printf("Warning: Failed to save API key to config: %v\n", err)
				}
			}
		}
	}

	// from now on, key is not empty
	response, err := genmessage.GenerateCommitMessage(key)
	if err != nil {
		displayError("Failed to generate commit message: %v", err)
	}
	commitMessage, commitResult := messagetextarea.MessageTextArea(response)

	if commitMessage == "" {
		fmt.Println("No commit message provided")
		return
	}

	if commitResult == messagetextarea.ResultCommit {
		commit(commitMessage)
	} else {
		fmt.Println("Commit cancelled")
		return
	}
}

func commit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	err := cmd.Run()
	if err != nil {
		displayError("Failed to commit: %v", err)
	}
	fmt.Println("Commit successful")
}

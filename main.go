package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/hamzabow/co/internal/apikeyinput"
	"github.com/hamzabow/co/internal/genmessage"
	"github.com/hamzabow/co/internal/messagetextarea"
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
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		var err error
		key, err = apikeyinput.PromptApiKeyWithRetries()
		if err != nil {
			if err == apikeyinput.ErrEmptyApiKey {
				fmt.Println("No API key provided. Exiting.")
				os.Exit(1)
			} else {
				displayError("%v", err)
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

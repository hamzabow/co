package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/hamzabow/co/internal/apikeyinput"
	"github.com/hamzabow/co/internal/config"
	"github.com/hamzabow/co/internal/genmessage"
	"github.com/hamzabow/co/internal/messagetextarea"
	"github.com/spf13/cobra"
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

var (
	// Used for flags
	providerName string
	skipPrompt   bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "co",
		Short: "Generate AI-powered Git commit messages",
		Long: `Co is a CLI tool that generates Git commit messages using AI.
It analyzes your staged changes and suggests a meaningful commit message.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCommand()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().StringVarP(&providerName, "provider", "p", ProviderOpenAI, "AI provider to use (openai, anthropic, ollama)")
	rootCmd.Flags().BoolVarP(&skipPrompt, "yes", "y", false, "Skip the confirmation prompt and automatically commit")
}

func displayError(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Println(errorStyle.Render("Error: " + errorMsg))
	os.Exit(1)
}

func runRootCommand() error {
	// Initialize API key variable
	var key string
	var err error

	// Try to load API key from config file
	key, err = config.LoadAPIKey(providerName)
	if err != nil && err != config.ErrNoConfigFile && err != config.ErrProviderNotFound {
		displayError("Failed to load API key from config: %v", err)
	}

	// If no key found, prompt the user
	if key == "" {
		key, err = apikeyinput.PromptApiKeyWithRetries()
		if err != nil {
			if err == apikeyinput.ErrEmptyApiKey {
				fmt.Println("No API key provided. Exiting.")
				return fmt.Errorf("no API key provided")
			} else {
				return fmt.Errorf("%v", err)
			}
		}

		// Save the key to config for future use
		if key != "" {
			if err := config.SaveAPIKey(providerName, key); err != nil {
				fmt.Printf("Warning: Failed to save API key to config: %v\n", err)
			}
		}
	}

	// Generate commit message using the AI provider
	response, err := genmessage.GenerateCommitMessage(key)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %v", err)
	}

	if skipPrompt {
		// Skip message editing and directly commit
		return commit(response)
	}

	// Show text area for editing the message
	commitMessage, commitResult := messagetextarea.MessageTextArea(response)

	if commitMessage == "" {
		fmt.Println("No commit message provided")
		return nil
	}

	if commitResult == messagetextarea.ResultCommit {
		return commit(commitMessage)
	}

	fmt.Println("Commit cancelled")
	return nil
}

func commit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}
	fmt.Println("Commit successful")
	return nil
}

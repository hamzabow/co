package cmd

import (
	"fmt"

	"github.com/hamzabow/co/internal/config"
	"github.com/spf13/cobra"
)

var (
	apiKey     string
	provider   string
	showConfig bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure API key settings",
	Long: `Configure AI provider API keys for use with Co.
You can set, update, or view your API key configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showConfig {
			displayConfiguration()
			return
		}

		if apiKey != "" {
			if err := config.SaveAPIKey(provider, apiKey); err != nil {
				fmt.Printf("Error saving API key: %v\n", err)
				return
			}
			fmt.Printf("API key for %s saved successfully\n", provider)
		} else {
			fmt.Println("No API key provided. Use --key to set an API key.")
			fmt.Println("Use --show to display current configuration.")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVar(&apiKey, "key", "", "Set the API key for the selected provider")
	configCmd.Flags().StringVar(&provider, "provider", ProviderOpenAI, "Set the AI provider (openai, anthropic, ollama)")
	configCmd.Flags().BoolVar(&showConfig, "show", false, "Show current configuration")
}

func displayConfiguration() {
	fmt.Println("Current Configuration:")
	fmt.Println("----------------------")

	for _, p := range []string{ProviderOpenAI, ProviderAnthropic, ProviderOllama} {
		key, err := config.LoadAPIKey(p)
		if err != nil && err != config.ErrNoConfigFile && err != config.ErrProviderNotFound {
			fmt.Printf("%s: Error loading configuration: %v\n", p, err)
			continue
		}

		if key == "" {
			fmt.Printf("%s: Not configured\n", p)
		} else {
			// Show only first few characters for security
			maskedKey := key[:4] + "..." + key[len(key)-4:]
			fmt.Printf("%s: %s\n", p, maskedKey)
		}
	}
}

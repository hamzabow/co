package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

var (
	// ErrNoConfigFile is returned when the config file doesn't exist
	ErrNoConfigFile = errors.New("config file does not exist")
)

// GetConfigDir returns the directory where the config will be stored
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Use standard XDG config directory instead of custom .co directory
	configDir := filepath.Join(homeDir, ".config", "co")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return configDir, nil
}

// GetConfigFilePath returns the path to the config file
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// SaveAPIKey saves the API key to the config file
func SaveAPIKey(apiKey string) error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// Create config structure
	config := Config{
		OpenAIAPIKey: apiKey,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restrictive permissions (600 - user read/write only)
	return os.WriteFile(configPath, data, 0600)
}

// LoadAPIKey loads the API key from the config file
func LoadAPIKey() (string, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", ErrNoConfigFile
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	// Unmarshal JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}

	return config.OpenAIAPIKey, nil
}

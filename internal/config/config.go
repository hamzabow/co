package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
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
// Uses platform-specific directories:
// - Windows: %APPDATA%\Co\
// - macOS: ~/Library/Application Support/Co/
// - Linux: ~/.config/co/ (follows XDG spec)
func GetConfigDir() (string, error) {
	var configDir string

	// Determine config directory based on OS
	switch runtime.GOOS {
	case "windows":
		// Windows: use %APPDATA%\Co\
		appData := os.Getenv("APPDATA")
		if appData == "" {
			// Fallback if APPDATA is not set (very rare)
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "Co")

	case "darwin":
		// macOS: use ~/Library/Application Support/Co/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", "Co")

	default:
		// Linux/Unix: use ~/.config/co/ (XDG spec)
		// First check if XDG_CONFIG_HOME is set
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			configDir = filepath.Join(xdgConfigHome, "co")
		} else {
			// Fall back to default ~/.config/co/
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, ".config", "co")
		}
	}

	// Create directory if it doesn't exist
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

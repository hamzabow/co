package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config holds application configuration
type Config struct {
	OpenAIAPIKey string `toml:"openai_api_key"`
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
	// Note: On Windows, the 0700 permission is approximated as the default NTFS permissions
	// for the current user. Unix-style permissions don't directly apply on Windows.
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

	return filepath.Join(configDir, "config.toml"), nil
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

	// Create a new file with restrictive permissions
	// Note: On Windows, these permissions are approximated and don't directly map to
	// the Unix permission model. Files in %APPDATA% are typically only accessible
	// to the creating user by default on Windows.
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the TOML data
	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
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

	// Read and parse TOML
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return "", err
	}

	return config.OpenAIAPIKey, nil
}

// TODO: For enhanced Windows security, consider implementing Windows-specific
// ACL control using a library like github.com/hectane/go-acl or using syscalls
// to set Windows-specific file security attributes.

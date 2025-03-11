package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// ErrNoConfigFile is returned when the config file doesn't exist
	ErrNoConfigFile = errors.New("config file does not exist")
	// ErrDecryptionFailed is returned when credential decryption fails
	ErrDecryptionFailed = errors.New("failed to decrypt credentials")
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

// GetCredentialsFilePath returns the path to the credentials file
func GetCredentialsFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	// Using .credentials filename for the encrypted data
	filename := ".credentials"

	return filepath.Join(configDir, filename), nil
}

// deriveEncryptionKey generates an encryption key based on machine-specific information
// This provides a consistent key for the same machine without requiring a password
func deriveEncryptionKey() ([]byte, error) {
	// Get machine-specific identifiers
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Get user home directory as another machine-specific input
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Application-specific salt to ensure uniqueness
	appSalt := "co-credential-encryption-v1"

	// Combine all inputs
	data := hostname + homeDir + appSalt + runtime.GOOS

	// Use SHA-256 to derive a suitable key
	hash := sha256.Sum256([]byte(data))
	return hash[:], nil
}

// encrypt encrypts the plaintext using AES-GCM and returns the ciphertext
func encrypt(plaintext string) (string, error) {
	// Derive the encryption key
	key, err := deriveEncryptionKey()
	if err != nil {
		return "", err
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce (Number used ONCE)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and authenticate the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for safe storage
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// decrypt decrypts the ciphertext and returns the plaintext
func decrypt(encodedCiphertext string) (string, error) {
	// Derive the encryption key (should be the same as used for encryption)
	key, err := deriveEncryptionKey()
	if err != nil {
		return "", err
	}

	// Decode the base64 encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	// Create a new cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check if the ciphertext is long enough to contain a nonce
	if len(ciphertext) < gcm.NonceSize() {
		return "", ErrDecryptionFailed
	}

	// Extract the nonce from the beginning of the ciphertext
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	// Decrypt and authenticate
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// SaveAPIKey saves the API key to the encrypted credentials file
func SaveAPIKey(apiKey string) error {
	credentialsPath, err := GetCredentialsFilePath()
	if err != nil {
		return err
	}

	// Encrypt the API key
	encryptedKey, err := encrypt(apiKey)
	if err != nil {
		return err
	}

	// Create a new file with restrictive permissions
	file, err := os.OpenFile(credentialsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the encrypted data directly to file
	_, err = file.WriteString(encryptedKey)
	if err != nil {
		return err
	}

	return nil
}

// LoadAPIKey loads the API key from the encrypted credentials file
func LoadAPIKey() (string, error) {
	credentialsPath, err := GetCredentialsFilePath()
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return "", ErrNoConfigFile
	}

	// Read the encrypted content
	encryptedBytes, err := os.ReadFile(credentialsPath)
	if err != nil {
		return "", err
	}

	// Decrypt the content
	apiKey, err := decrypt(string(encryptedBytes))
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

// TODO: For enhanced Windows security, consider implementing Windows-specific
// ACL control using a library like github.com/hectane/go-acl or using syscalls
// to set Windows-specific file security attributes.

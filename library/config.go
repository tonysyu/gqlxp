package library

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// configDir returns the standard configuration directory for gqlxp.
// It follows platform conventions:
// - macOS/Linux: $HOME/.config/gqlxp/
// - Windows: %APPDATA%\gqlxp\
// - Fallback: $HOME/.gqlxp/
func configDir() (string, error) {
	var baseDir string

	// Try platform-specific config directory first
	if runtime.GOOS == "windows" {
		baseDir = os.Getenv("APPDATA")
		if baseDir != "" {
			return filepath.Join(baseDir, "gqlxp"), nil
		}
	} else {
		// macOS and Linux: use XDG config directory
		baseDir = os.Getenv("XDG_CONFIG_HOME")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			baseDir = filepath.Join(home, ".config")
		}
		return filepath.Join(baseDir, "gqlxp"), nil
	}

	// Fallback to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".gqlxp"), nil
}

// schemasDir returns the schemas directory within the config directory.
func schemasDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "schemas"), nil
}

// metadataFile returns the path to the metadata.json file.
func metadataFile() (string, error) {
	schemasDir, err := schemasDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(schemasDir, "metadata.json"), nil
}

// userConfigFile returns the path to the user config.json file.
func userConfigFile() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// InitConfigDir creates the configuration directory structure if it doesn't exist.
func InitConfigDir() error {
	schemasDir, err := schemasDir()
	if err != nil {
		return err
	}

	// Create the schemas directory (creates parent directories as needed)
	if err := os.MkdirAll(schemasDir, 0755); err != nil {
		return fmt.Errorf("failed to create schemas directory: %w", err)
	}

	// Create empty metadata.json if it doesn't exist
	metadataFile, err := metadataFile()
	if err != nil {
		return err
	}

	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		if err := os.WriteFile(metadataFile, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("failed to create metadata file: %w", err)
		}
	}

	return nil
}

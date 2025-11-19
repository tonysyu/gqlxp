package library

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

// White-box tests for internal config functions

func TestConfigDir(t *testing.T) {
	is := is.New(t)

	t.Run("returns non-empty config directory", func(t *testing.T) {
		configDir, err := configDir()
		is.NoErr(err)
		is.True(configDir != "")
	})
}

func TestInitConfigDir(t *testing.T) {
	is := is.New(t)

	// Create temporary directory for testing
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	t.Run("creates config directory structure", func(t *testing.T) {
		err := InitConfigDir()
		is.NoErr(err)

		// Verify schemas directory exists
		schemasDir, err := schemasDir()
		is.NoErr(err)

		info, err := os.Stat(schemasDir)
		is.NoErr(err)
		is.True(info.IsDir())

		// Verify metadata file exists
		metadataFile, err := metadataFile()
		is.NoErr(err)

		_, err = os.Stat(metadataFile)
		is.NoErr(err)
	})

	t.Run("calling InitConfigDir multiple times is safe", func(t *testing.T) {
		err := InitConfigDir()
		is.NoErr(err)

		err = InitConfigDir()
		is.NoErr(err)
	})
}

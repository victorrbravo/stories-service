// Package testenv loads the repository .env file for tests, so secrets such as
// API keys never have to be written in test code.
package testenv

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Load walks up from the current directory to the module root (where go.mod
// lives) and loads its .env file. Variables already set in the environment win.
func Load() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			_ = godotenv.Load(filepath.Join(dir, ".env"))
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

// Get returns the value of an environment variable or fallback when unset.
func Get(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"log"
)

var dirPerms fs.FileMode = 0750 // rwxr-x---
var subdirs = []string{"scriptbundles"}

func createDir(dir string) error {
	if err := os.MkdirAll(dir, dirPerms); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func ticliHomePath() (string, error) {
	var homeDir string
	if homeDir = os.Getenv("TICLI_HOME"); homeDir == "" {

		if runtime.GOOS == "windows" {
			// For Windows, use the LOCALAPPDATA environment variable
			homeDir = os.Getenv("LOCALAPPDATA")
			if homeDir == "" {
				return "", fmt.Errorf("LOCALAPPDATA environment variable is not set")
			}
			homeDir = filepath.Join(homeDir, "ticli")
		} else {
			// For *nix platforms, use the HOME environment variable
			homeDir = os.Getenv("HOME")
			if homeDir == "" {
				return "", fmt.Errorf("HOME environment variable is not set")
			}
			homeDir = filepath.Join(homeDir, ".ticli")
		}
	}
	return homeDir, nil
}

// SetupTicliEnv initializes the home directory for the Ticli tool.
// It determines the home directory path based on the "TICLI_HOME" environment variable, if set,
// or by using platform-specific defaults (e.g., %LOCALAPPDATA%/ticli on Windows or ~/.ticli on Unix-like systems).
// The function then creates the home directory, including any necessary parent directories, with specific permissions.
//
// Returns:
//
//	string: The absolute path to the created or identified home directory.
//	error: An error value if any error occurs while determining the home path or creating the directory.
//	       This includes situations where necessary environment variables are not set or if directory creation fails.
//
// Note: The function will log errors using the log package, providing detailed information on failures.
//
// Usage:
//
//	homeDir, err := SetupTicliEnv()
//	if err != nil {
//	    log.Println(err)
//	}
func SetupTicliEnv() (string, error) {
	homeDir, err := ticliHomePath()
	if err != nil {
		log.Println("Error getting ticli home path")
		return "", err
	}
	err = createDir(homeDir)
	if err != nil {
		log.Println("Error creating ticli home directory")
		return "", err
	}

	for _, subdir := range subdirs {
		err = createDir(filepath.Join(homeDir, subdir))
		if err != nil {
			log.Println("Error creating ticli home subdirectory ", subdir)
			return "", err
		}
	}
	log.Println("Ticli home directory:", homeDir)

	return homeDir, nil
}

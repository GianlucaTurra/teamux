package data

import (
	"os"
	"strings"
)

// GetPWDorPlaceholder returns the path for the PWD, checking if it exists or return home
// if the path is empty. Error is returned if the directory is not empty and
// does not exists
func GetPWDorPlaceholder(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "$HOME", nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, err
	}
	return path, nil
}

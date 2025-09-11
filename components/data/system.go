package data

import (
	"os"
	"strings"
)

// getPWD returns the path for the PWD, checking if it exists or return home
// if the path is empty. Error is returned if the directory is not empty and
// does not exists
func getPWD(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		if home, err := os.UserHomeDir(); err != nil {
			return "", err
		} else {
			return home, nil
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, err
	}
	return path, nil
}

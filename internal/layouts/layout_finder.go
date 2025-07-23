package layouts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadLayouts() []string {
	var layoutFiles []string
	f, err := os.Open("./")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	files, err := f.ReadDir(0)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".sh" {
			continue
		}
		if !strings.Contains(file.Name(), "teamux") {
			continue
		}
		layoutFiles = append(layoutFiles, file.Name())
	}
	return layoutFiles
}

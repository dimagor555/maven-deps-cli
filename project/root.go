package project

import (
	"os"
	"path/filepath"
)

var buildFileMarkers = []string{
	"settings.gradle.kts",
	"settings.gradle",
	"build.gradle.kts",
	"build.gradle",
}

func FindRoot(startDir string) string {
	current, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}

	for {
		for _, marker := range buildFileMarkers {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

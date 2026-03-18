package project

import (
	"os"
	"path/filepath"
)

var skipDirs = map[string]bool{
	"build":        true,
	".gradle":      true,
	".idea":        true,
	"node_modules": true,
	".git":         true,
}

func FindBuildFiles(root, fileName string) []string {
	var results []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && skipDirs[d.Name()] {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == fileName {
			results = append(results, path)
		}
		return nil
	})
	return results
}

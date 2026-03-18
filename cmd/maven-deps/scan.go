package main

import (
	"os"
	"path/filepath"

	"dimagor555.pro/maven-deps/gradle"
	"dimagor555.pro/maven-deps/project"
)

type scannedDep struct {
	GroupID       string `json:"groupId"`
	ArtifactID    string `json:"artifactId"`
	Version       string `json:"version,omitempty"`
	Configuration string `json:"configuration"`
	Source        string `json:"source"`
}

func scanProject(path string) ([]scannedDep, error) {
	root := project.FindRoot(path)
	if root == "" {
		return nil, nil
	}

	seen := make(map[string]bool)
	var deps []scannedDep

	catalogPath := filepath.Join(root, "gradle", "libs.versions.toml")
	if data, err := os.ReadFile(catalogPath); err == nil {
		catalog := gradle.ParseVersionCatalog(string(data))
		for _, entry := range catalog {
			key := entry.GroupID + ":" + entry.ArtifactID + ":" + entry.Version
			if !seen[key] {
				seen[key] = true
				deps = append(deps, scannedDep{
					GroupID:       entry.GroupID,
					ArtifactID:    entry.ArtifactID,
					Version:       entry.Version,
					Configuration: "catalog",
					Source:        "libs.versions.toml",
				})
			}
		}
	}

	for _, fileName := range []string{"build.gradle.kts", "build.gradle"} {
		for _, file := range project.FindBuildFiles(root, fileName) {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			rel, _ := filepath.Rel(root, file)
			for _, dep := range gradle.ParseDependencies(string(data), rel) {
				if dep.GroupID == "" || dep.ArtifactID == "" {
					continue
				}
				key := dep.GroupID + ":" + dep.ArtifactID + ":" + dep.Version
				if !seen[key] {
					seen[key] = true
					deps = append(deps, scannedDep{
						GroupID:       dep.GroupID,
						ArtifactID:    dep.ArtifactID,
						Version:       dep.Version,
						Configuration: dep.Configuration,
						Source:        rel,
					})
				}
			}
		}
	}

	return deps, nil
}

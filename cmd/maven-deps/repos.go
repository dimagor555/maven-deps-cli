package main

import (
	"os"
	"path/filepath"

	"dimagor555.pro/maven-deps/gradle"
	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/project"
)

var gradleFiles = []string{
	"settings.gradle.kts",
	"settings.gradle",
	"build.gradle.kts",
	"build.gradle",
}

func discoverRepos(path string) []maven.Repository {
	root := project.FindRoot(path)
	if root == "" {
		return defaultRepos()
	}

	var repos []maven.Repository
	seen := make(map[string]bool)

	for _, file := range gradleFiles {
		full := filepath.Join(root, file)
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		for _, rc := range gradle.ParseRepositories(string(data)) {
			if !seen[rc.URL] {
				seen[rc.URL] = true
				repos = append(repos, maven.NewRepository(rc.Name, rc.URL))
			}
		}
	}

	for _, fallback := range defaultRepos() {
		if !seen[fallback.URL] {
			seen[fallback.URL] = true
			repos = append(repos, fallback)
		}
	}

	return repos
}

func defaultRepos() []maven.Repository {
	return []maven.Repository{maven.Google, maven.GradlePlugin, maven.Central}
}

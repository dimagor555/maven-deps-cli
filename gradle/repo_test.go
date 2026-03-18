package gradle

import "testing"

func containsRepo(repos []RepoConfig, name, url string) bool {
	for _, r := range repos {
		if r.Name == name && r.URL == url {
			return true
		}
	}
	return false
}

func TestParseRepositories_WhenMavenCentral_Parses(t *testing.T) {
	content := `repositories { mavenCentral() }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "Maven Central", "https://repo1.maven.org/maven2") {
		t.Errorf("expected Maven Central, got %v", repos)
	}
}

func TestParseRepositories_WhenGoogle_Parses(t *testing.T) {
	content := `repositories { google() }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "Google", "https://maven.google.com") {
		t.Errorf("expected Google, got %v", repos)
	}
}

func TestParseRepositories_WhenGradlePluginPortal_Parses(t *testing.T) {
	content := `repositories { gradlePluginPortal() }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "Gradle Plugin Portal", "https://plugins.gradle.org/m2") {
		t.Errorf("expected Gradle Plugin Portal, got %v", repos)
	}
}

func TestParseRepositories_WhenMavenDirect_Parses(t *testing.T) {
	content := `repositories { maven("https://jitpack.io") }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "https://jitpack.io", "https://jitpack.io") {
		t.Errorf("expected jitpack, got %v", repos)
	}
}

func TestParseRepositories_WhenMavenUrlParam_Parses(t *testing.T) {
	content := `repositories { maven(url = "https://maven.pkg.jetbrains.space/public/p/compose/dev") }`
	repos := ParseRepositories(content)
	url := "https://maven.pkg.jetbrains.space/public/p/compose/dev"
	if !containsRepo(repos, url, url) {
		t.Errorf("expected compose dev, got %v", repos)
	}
}

func TestParseRepositories_WhenMavenBlockUri_Parses(t *testing.T) {
	content := `repositories { maven { url = uri("https://repo.spring.io/milestone") } }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "https://repo.spring.io/milestone", "https://repo.spring.io/milestone") {
		t.Errorf("expected spring milestone, got %v", repos)
	}
}

func TestParseRepositories_WhenGroovyMavenBlock_Parses(t *testing.T) {
	content := `repositories { maven { url 'https://repo.spring.io/milestone' } }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "https://repo.spring.io/milestone", "https://repo.spring.io/milestone") {
		t.Errorf("expected spring milestone, got %v", repos)
	}
}

func TestParseRepositories_WhenGoogleBlock_Parses(t *testing.T) {
	content := `repositories {
        google {
            mavenContent {
                includeGroupAndSubgroups("androidx")
            }
        }
    }`
	repos := ParseRepositories(content)
	if !containsRepo(repos, "Google", "https://maven.google.com") {
		t.Errorf("expected Google, got %v", repos)
	}
}

func TestParseRepositories_WhenNoRepos_ReturnsEmpty(t *testing.T) {
	content := `plugins { id("org.jetbrains.kotlin.jvm") }`
	repos := ParseRepositories(content)
	if len(repos) != 0 {
		t.Errorf("expected empty, got %v", repos)
	}
}

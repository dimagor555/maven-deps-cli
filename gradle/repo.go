package gradle

import "regexp"

type RepoConfig struct {
	Name string
	URL  string
}

var wellKnownRepos = map[string]RepoConfig{
	"mavenCentral":       {Name: "Maven Central", URL: "https://repo1.maven.org/maven2"},
	"google":             {Name: "Google", URL: "https://maven.google.com"},
	"gradlePluginPortal": {Name: "Gradle Plugin Portal", URL: "https://plugins.gradle.org/m2"},
}

var (
	mavenDirectRe      = regexp.MustCompile(`\bmaven\s*\(\s*["']([^"']+)["']\s*\)`)
	mavenUrlParamRe    = regexp.MustCompile(`\bmaven\s*\(\s*url\s*=\s*["']([^"']+)["']\s*\)`)
	mavenBlockUriRe    = regexp.MustCompile(`\bmaven\s*\{[^}]*url\s*=\s*uri\s*\(\s*["']([^"']+)["']\s*\)`)
	mavenBlockGroovyRe = regexp.MustCompile(`\bmaven\s*\{[^}]*url\s+["']([^"']+)["']`)
)

func ParseRepositories(content string) []RepoConfig {
	var repos []RepoConfig

	for funcName, config := range wellKnownRepos {
		re := regexp.MustCompile(`\b` + funcName + `\s*(\(\s*\)|\s*\{)`)
		if re.MatchString(content) {
			repos = append(repos, config)
		}
	}

	for _, re := range []*regexp.Regexp{mavenDirectRe, mavenUrlParamRe, mavenBlockUriRe, mavenBlockGroovyRe} {
		for _, m := range re.FindAllStringSubmatch(content, -1) {
			repos = append(repos, RepoConfig{Name: m[1], URL: m[1]})
		}
	}

	return repos
}

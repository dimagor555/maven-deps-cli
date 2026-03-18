package gradle

import (
	"regexp"
	"strings"
)

type Dependency struct {
	GroupID       string
	ArtifactID    string
	Version       string
	Configuration string
	Source        string
	CatalogRef    string
}

var configurations = []string{
	"implementation", "api", "compileOnly", "runtimeOnly",
	"testImplementation", "testCompileOnly", "testRuntimeOnly",
	"kapt", "ksp", "annotationProcessor",
}

var (
	stringDepRe  *regexp.Regexp
	catalogDepRe *regexp.Regexp
	pluginIDRe     *regexp.Regexp
	pluginKotlinRe *regexp.Regexp
)

func init() {
	pattern := strings.Join(configurations, "|")
	stringDepRe = regexp.MustCompile(
		`\b(` + pattern + `)\s*[(  ]\s*["']([^"':]+):([^"':]+)(?::([^"']+))?["']\s*\)?`,
	)
	catalogDepRe = regexp.MustCompile(
		`\b(` + pattern + `)\s*\(\s*libs\.([a-zA-Z0-9.]+)\s*\)`,
	)
	pluginIDRe = regexp.MustCompile(
		`\bid\s*\(\s*["']([^"']+)["']\s*\)\s*version\s*["']([^"']+)["']`,
	)
	pluginKotlinRe = regexp.MustCompile(
		`\bkotlin\s*\(\s*["']([^"']+)["']\s*\)\s*version\s*["']([^"']+)["']`,
	)
}

func ParseDependencies(content, source string) []Dependency {
	var deps []Dependency

	for _, m := range stringDepRe.FindAllStringSubmatch(content, -1) {
		deps = append(deps, Dependency{
			GroupID:       m[2],
			ArtifactID:    m[3],
			Version:       m[4],
			Configuration: m[1],
			Source:        source,
		})
	}

	for _, m := range catalogDepRe.FindAllStringSubmatch(content, -1) {
		deps = append(deps, Dependency{
			Configuration: m[1],
			Source:        source,
			CatalogRef:    m[2],
		})
	}

	for _, m := range pluginIDRe.FindAllStringSubmatch(content, -1) {
		pluginID := m[1]
		deps = append(deps, Dependency{
			GroupID:       pluginID,
			ArtifactID:    pluginID + ".gradle.plugin",
			Version:       m[2],
			Configuration: "plugin",
			Source:        source,
		})
	}

	for _, m := range pluginKotlinRe.FindAllStringSubmatch(content, -1) {
		pluginID := "org.jetbrains.kotlin." + m[1]
		deps = append(deps, Dependency{
			GroupID:       pluginID,
			ArtifactID:    pluginID + ".gradle.plugin",
			Version:       m[2],
			Configuration: "plugin",
			Source:        source,
		})
	}

	return deps
}

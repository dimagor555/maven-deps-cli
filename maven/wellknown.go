package maven

var (
	Central      = NewRepository("Maven Central", "https://repo1.maven.org/maven2")
	Google       = NewRepository("Google", "https://maven.google.com")
	GradlePlugin = NewRepository("Gradle Plugin Portal", "https://plugins.gradle.org/m2")
)

var ProxyTargetURLs = map[string]bool{
	Central.URL: true,
}

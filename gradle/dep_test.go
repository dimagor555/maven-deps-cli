package gradle

import "testing"

func TestParseDependencies_WhenStringNotation_ParsesAll(t *testing.T) {
	content := `implementation("io.ktor:ktor-client-core:3.1.1")`
	deps := ParseDependencies(content, "build.gradle.kts")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	d := deps[0]
	if d.GroupID != "io.ktor" {
		t.Errorf("GroupID = %q", d.GroupID)
	}
	if d.ArtifactID != "ktor-client-core" {
		t.Errorf("ArtifactID = %q", d.ArtifactID)
	}
	if d.Version != "3.1.1" {
		t.Errorf("Version = %q", d.Version)
	}
	if d.Configuration != "implementation" {
		t.Errorf("Configuration = %q", d.Configuration)
	}
}

func TestParseDependencies_WhenNoVersion_VersionEmpty(t *testing.T) {
	content := `api("com.google.code.gson:gson")`
	deps := ParseDependencies(content, "build.gradle.kts")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].Version != "" {
		t.Errorf("Version = %q, want empty", deps[0].Version)
	}
}

func TestParseDependencies_WhenCatalogRef_ParsesRef(t *testing.T) {
	content := `implementation(libs.ktor.client.core)`
	deps := ParseDependencies(content, "build.gradle.kts")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].CatalogRef != "ktor.client.core" {
		t.Errorf("CatalogRef = %q", deps[0].CatalogRef)
	}
}

func TestParseDependencies_WhenGroovyString_Parses(t *testing.T) {
	content := `testImplementation 'junit:junit:4.13.2'`
	deps := ParseDependencies(content, "build.gradle")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	if deps[0].GroupID != "junit" || deps[0].Version != "4.13.2" {
		t.Errorf("unexpected dep: %+v", deps[0])
	}
}

func TestParseDependencies_WhenMultipleConfigs_ParsesAll(t *testing.T) {
	content := `
implementation("io.ktor:ktor-client-core:3.1.1")
testImplementation("org.junit:junit:5.0.0")
kapt("com.google.dagger:dagger-compiler:2.50")
`
	deps := ParseDependencies(content, "build.gradle.kts")
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %d", len(deps))
	}
}

func TestParseDependencies_WhenPluginID_ParsesAsGradlePlugin(t *testing.T) {
	content := `
plugins {
    id("io.gitlab.arturbosch.detekt") version "1.23.7"
}
`
	deps := ParseDependencies(content, "build.gradle.kts")
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	d := deps[0]
	if d.GroupID != "io.gitlab.arturbosch.detekt" {
		t.Errorf("GroupID = %q", d.GroupID)
	}
	if d.ArtifactID != "io.gitlab.arturbosch.detekt.gradle.plugin" {
		t.Errorf("ArtifactID = %q", d.ArtifactID)
	}
	if d.Version != "1.23.7" {
		t.Errorf("Version = %q", d.Version)
	}
	if d.Configuration != "plugin" {
		t.Errorf("Configuration = %q", d.Configuration)
	}
}

func TestParseDependencies_WhenPluginIDNoVersion_NotParsed(t *testing.T) {
	content := `id("com.android.application")`
	deps := ParseDependencies(content, "build.gradle.kts")
	for _, d := range deps {
		if d.Configuration == "plugin" {
			t.Errorf("should not parse plugin without version, got %+v", d)
		}
	}
}

func TestParseDependencies_WhenKotlinShorthand_ParsesAsKotlinPlugin(t *testing.T) {
	content := `
plugins {
    kotlin("jvm") version "2.2.21"
    kotlin("plugin.serialization") version "2.2.21"
}
`
	deps := ParseDependencies(content, "build.gradle.kts")
	var plugins []Dependency
	for _, d := range deps {
		if d.Configuration == "plugin" {
			plugins = append(plugins, d)
		}
	}
	if len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(plugins))
	}
	if plugins[0].GroupID != "org.jetbrains.kotlin.jvm" {
		t.Errorf("GroupID = %q, want org.jetbrains.kotlin.jvm", plugins[0].GroupID)
	}
	if plugins[0].ArtifactID != "org.jetbrains.kotlin.jvm.gradle.plugin" {
		t.Errorf("ArtifactID = %q", plugins[0].ArtifactID)
	}
	if plugins[0].Version != "2.2.21" {
		t.Errorf("Version = %q", plugins[0].Version)
	}
	if plugins[1].GroupID != "org.jetbrains.kotlin.plugin.serialization" {
		t.Errorf("GroupID = %q", plugins[1].GroupID)
	}
}

func TestParseDependencies_WhenKotlinNoVersion_NotParsed(t *testing.T) {
	content := `kotlin("jvm")`
	deps := ParseDependencies(content, "build.gradle.kts")
	for _, d := range deps {
		if d.Configuration == "plugin" {
			t.Errorf("should not parse kotlin plugin without version, got %+v", d)
		}
	}
}

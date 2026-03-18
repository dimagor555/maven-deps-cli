package gradle

import "testing"

func TestParseVersionCatalog_WhenVersionRef_Resolves(t *testing.T) {
	content := `
[versions]
ktor = "3.1.1"
kotlin = "2.1.0"

[libraries]
ktor-client-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
kotlin-stdlib = { module = "org.jetbrains.kotlin:kotlin-stdlib", version.ref = "kotlin" }
`
	entries := ParseVersionCatalog(content)

	e, ok := entries["ktor-client-core"]
	if !ok {
		t.Fatal("missing ktor-client-core")
	}
	if e.GroupID != "io.ktor" || e.ArtifactID != "ktor-client-core" || e.Version != "3.1.1" {
		t.Errorf("unexpected: %+v", e)
	}

	e, ok = entries["kotlin-stdlib"]
	if !ok {
		t.Fatal("missing kotlin-stdlib")
	}
	if e.Version != "2.1.0" {
		t.Errorf("Version = %q, want 2.1.0", e.Version)
	}
}

func TestParseVersionCatalog_WhenInlineVersion_Parses(t *testing.T) {
	content := `
[libraries]
gson = { module = "com.google.code.gson:gson", version = "2.11.0" }
`
	entries := ParseVersionCatalog(content)
	e, ok := entries["gson"]
	if !ok {
		t.Fatal("missing gson")
	}
	if e.GroupID != "com.google.code.gson" || e.Version != "2.11.0" {
		t.Errorf("unexpected: %+v", e)
	}
}

func TestParseVersionCatalog_WhenGroupName_Parses(t *testing.T) {
	content := `
[versions]
ktor = "3.1.1"

[libraries]
ktor-core = { group = "io.ktor", name = "ktor-client-core", version.ref = "ktor" }
`
	entries := ParseVersionCatalog(content)
	e, ok := entries["ktor-core"]
	if !ok {
		t.Fatal("missing ktor-core")
	}
	if e.GroupID != "io.ktor" || e.ArtifactID != "ktor-client-core" || e.Version != "3.1.1" {
		t.Errorf("unexpected: %+v", e)
	}
}

func TestParseVersionCatalog_WhenNoVersion_VersionEmpty(t *testing.T) {
	content := `
[libraries]
bom-lib = { module = "io.ktor:ktor-bom" }
`
	entries := ParseVersionCatalog(content)
	e, ok := entries["bom-lib"]
	if !ok {
		t.Fatal("missing bom-lib")
	}
	if e.Version != "" {
		t.Errorf("Version = %q, want empty", e.Version)
	}
}

func TestParseVersionCatalog_WhenEmpty_ReturnsEmpty(t *testing.T) {
	entries := ParseVersionCatalog("")
	if entries != nil && len(entries) != 0 {
		t.Errorf("expected nil or empty, got %v", entries)
	}
}

func TestParseVersionCatalog_WhenPlugin_ParsesAsGradlePlugin(t *testing.T) {
	content := `
[versions]
kotlin = "2.1.0"

[plugins]
kotlin-jvm = { id = "org.jetbrains.kotlin.jvm", version.ref = "kotlin" }
`
	entries := ParseVersionCatalog(content)
	e, ok := entries["plugin-kotlin-jvm"]
	if !ok {
		t.Fatal("missing plugin-kotlin-jvm")
	}
	if e.GroupID != "org.jetbrains.kotlin.jvm" {
		t.Errorf("GroupID = %q", e.GroupID)
	}
	if e.ArtifactID != "org.jetbrains.kotlin.jvm.gradle.plugin" {
		t.Errorf("ArtifactID = %q", e.ArtifactID)
	}
	if e.Version != "2.1.0" {
		t.Errorf("Version = %q", e.Version)
	}
}

func TestParseVersionCatalog_WhenStringNotation_Parses(t *testing.T) {
	content := `
[libraries]
gson = "com.google.code.gson:gson:2.11.0"
`
	entries := ParseVersionCatalog(content)
	e, ok := entries["gson"]
	if !ok {
		t.Fatal("missing gson")
	}
	if e.GroupID != "com.google.code.gson" || e.ArtifactID != "gson" || e.Version != "2.11.0" {
		t.Errorf("unexpected: %+v", e)
	}
}

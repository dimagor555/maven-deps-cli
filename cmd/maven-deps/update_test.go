package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"dimagor555.pro/maven-deps/gradle"
	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/update"
)

type fakeResolver struct {
	versions map[string][]string
}

func (f *fakeResolver) Resolve(_ context.Context, groupID, artifactID string) (maven.Metadata, error) {
	key := groupID + ":" + artifactID
	return maven.Metadata{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Versions:   f.versions[key],
	}, nil
}

type failingResolver struct {
	err error
}

func (f *failingResolver) Resolve(_ context.Context, _, _ string) (maven.Metadata, error) {
	return maven.Metadata{}, f.err
}

func TestExecuteUpdate_ResolveFails_ReturnsErrorNotNothingToUpdate(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	resolver := &failingResolver{err: errors.New("connection reset by peer")}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader(""), root, catalog, resolver, update.LevelMinor, true, false)
	if err == nil {
		t.Fatal("expected error when resolve fails, got nil")
	}
	if strings.Contains(out.String(), "nothing to update") {
		t.Errorf("must not claim 'nothing to update' when resolve failed:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "WARNING") {
		t.Errorf("expected WARNING about failed resolves, got:\n%s", out.String())
	}
}

func TestResolveUpgrades_NotFound_NotTreatedAsFailure(t *testing.T) {
	entries := map[string]gradle.CatalogEntry{
		"a": {SourceAlias: "a", GroupID: "g", ArtifactID: "a", Version: "1.0.0", VersionRef: "a"},
	}
	resolver := &failingResolver{err: maven.NewNotFoundError("any repository")}
	_, failures := resolveUpgrades(context.Background(), entries, resolver)
	if len(failures) != 0 {
		t.Errorf("not-found must not be reported as failure, got %d", len(failures))
	}
}

func TestResolveUpgrades_TransientError_ReportedAsFailure(t *testing.T) {
	entries := map[string]gradle.CatalogEntry{
		"a": {SourceAlias: "a", GroupID: "g", ArtifactID: "a", Version: "1.0.0", VersionRef: "a"},
	}
	resolver := &failingResolver{err: errors.New("timeout")}
	_, failures := resolveUpgrades(context.Background(), entries, resolver)
	if len(failures) != 1 {
		t.Fatalf("transient error must be reported, got %d", len(failures))
	}
}

func TestExecuteUpdate_AppliesRefAndInlineChanges(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
gson = { module = "com.google.code.gson:gson", version = "2.11.0" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	resolver := &fakeResolver{versions: map[string][]string{
		"io.ktor:ktor-client-core":       {"3.0.0", "3.1.0"},
		"com.google.code.gson:gson":      {"2.11.0", "2.12.0"},
	}}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader(""), root, catalog, resolver, update.LevelMinor, true, false)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(catalog)
	got := string(data)
	if !strings.Contains(got, `ktor = "3.1.0"`) {
		t.Errorf("ktor not bumped:\n%s", got)
	}
	if !strings.Contains(got, `version = "2.12.0"`) {
		t.Errorf("gson not bumped:\n%s", got)
	}
}

func TestExecuteUpdate_NoUpgrades_PrintsNothing(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.1.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	resolver := &fakeResolver{versions: map[string][]string{
		"io.ktor:ktor-client-core": {"3.1.0"},
	}}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader(""), root, catalog, resolver, update.LevelMinor, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "nothing to update") {
		t.Errorf("expected 'nothing to update', got %q", out.String())
	}
}

func TestExecuteUpdate_DirtyRepo_Blocks(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	if err := os.WriteFile(catalog, []byte(`[versions]
ktor = "3.0.0"
# dirty
[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`), 0o644); err != nil {
		t.Fatal(err)
	}
	resolver := &fakeResolver{versions: map[string][]string{
		"io.ktor:ktor-client-core": {"3.0.0", "3.1.0"},
	}}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader(""), root, catalog, resolver, update.LevelMinor, true, false)
	if err == nil {
		t.Errorf("expected dirty error")
	}
}

func TestExecuteUpdate_DirtyRepo_ForceDirty_Proceeds(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	os.WriteFile(catalog, []byte(`[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`), 0o644)
	resolver := &fakeResolver{versions: map[string][]string{
		"io.ktor:ktor-client-core": {"3.0.0", "3.1.0"},
	}}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader(""), root, catalog, resolver, update.LevelMinor, true, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteUpdate_ConfirmN_DoesNotWrite(t *testing.T) {
	root := initUpdateRepo(t, `[versions]
ktor = "3.0.0"

[libraries]
ktor-core = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`)
	catalog := filepath.Join(root, "gradle", "libs.versions.toml")
	resolver := &fakeResolver{versions: map[string][]string{
		"io.ktor:ktor-client-core": {"3.0.0", "3.1.0"},
	}}
	var out bytes.Buffer
	err := executeUpdate(context.Background(), &out, strings.NewReader("n\n"), root, catalog, resolver, update.LevelMinor, false, false)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(catalog)
	if strings.Contains(string(data), "3.1.0") {
		t.Errorf("file was written despite decline")
	}
	if !strings.Contains(out.String(), "aborted") {
		t.Errorf("expected 'aborted'")
	}
}

func initUpdateRepo(t *testing.T, catalogContent string) string {
	t.Helper()
	dir := t.TempDir()
	writeFileTest(t, filepath.Join(dir, "settings.gradle.kts"), "")
	gradleDir := filepath.Join(dir, "gradle")
	if err := os.MkdirAll(gradleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFileTest(t, filepath.Join(gradleDir, "libs.versions.toml"), catalogContent)
	runCmd(t, dir, "git", "init")
	runCmd(t, dir, "git", "config", "user.email", "t@t")
	runCmd(t, dir, "git", "config", "user.name", "t")
	runCmd(t, dir, "git", "add", ".")
	runCmd(t, dir, "git", "commit", "-m", "init")
	return dir
}

func writeFileTest(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runCmd(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	c := exec.Command(name, args...)
	c.Dir = dir
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

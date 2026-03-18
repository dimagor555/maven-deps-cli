package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRoot_WhenSettingsGradleKts_ReturnsDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "settings.gradle.kts"), []byte(""), 0644)
	sub := filepath.Join(dir, "app", "src")
	os.MkdirAll(sub, 0755)

	got := FindRoot(sub)
	if got != dir {
		t.Errorf("got %q, want %q", got, dir)
	}
}

func TestFindRoot_WhenBuildGradle_ReturnsDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "build.gradle"), []byte(""), 0644)

	got := FindRoot(dir)
	if got != dir {
		t.Errorf("got %q, want %q", got, dir)
	}
}

func TestFindRoot_WhenNoMarker_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	got := FindRoot(dir)
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

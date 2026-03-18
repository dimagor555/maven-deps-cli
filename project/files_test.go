package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindBuildFiles_WhenMultipleModules_FindsAll(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "build.gradle.kts"), []byte(""), 0644)
	os.MkdirAll(filepath.Join(dir, "app"), 0755)
	os.WriteFile(filepath.Join(dir, "app", "build.gradle.kts"), []byte(""), 0644)

	files := FindBuildFiles(dir, "build.gradle.kts")
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestFindBuildFiles_WhenSkipDirs_IgnoresThem(t *testing.T) {
	dir := t.TempDir()
	buildDir := filepath.Join(dir, "build")
	os.MkdirAll(buildDir, 0755)
	os.WriteFile(filepath.Join(buildDir, "build.gradle.kts"), []byte(""), 0644)

	files := FindBuildFiles(dir, "build.gradle.kts")
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

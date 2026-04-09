package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIsDirty_CleanRepo_ReturnsFalse(t *testing.T) {
	dir := initTestRepo(t)
	file := filepath.Join(dir, "tracked.txt")
	writeFile(t, file, "hello")
	runGit(t, dir, "add", "tracked.txt")
	runGit(t, dir, "commit", "-m", "init")

	dirty, err := IsDirty(dir, []string{file})
	if err != nil {
		t.Fatal(err)
	}
	if dirty {
		t.Errorf("expected clean")
	}
}

func TestIsDirty_ModifiedFile_ReturnsTrue(t *testing.T) {
	dir := initTestRepo(t)
	file := filepath.Join(dir, "tracked.txt")
	writeFile(t, file, "hello")
	runGit(t, dir, "add", "tracked.txt")
	runGit(t, dir, "commit", "-m", "init")
	writeFile(t, file, "changed")

	dirty, err := IsDirty(dir, []string{file})
	if err != nil {
		t.Fatal(err)
	}
	if !dirty {
		t.Errorf("expected dirty")
	}
}

func TestIsDirty_NonGitDir_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := IsDirty(dir, []string{filepath.Join(dir, "x.txt")})
	if err == nil {
		t.Errorf("expected error for non-git dir")
	}
}

func TestIsDirty_OtherFileDirty_DoesNotAffect(t *testing.T) {
	dir := initTestRepo(t)
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	writeFile(t, a, "a")
	writeFile(t, b, "b")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init")
	writeFile(t, b, "changed")

	dirty, err := IsDirty(dir, []string{a})
	if err != nil {
		t.Fatal(err)
	}
	if dirty {
		t.Errorf("expected clean for a.txt only")
	}
}

func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "t@t")
	runGit(t, dir, "config", "user.name", "t")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

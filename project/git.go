package project

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsDirty(root string, files []string) (bool, error) {
	args := []string{"status", "--porcelain", "--"}
	args = append(args, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git status: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)) != "", nil
}

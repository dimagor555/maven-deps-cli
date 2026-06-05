package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"dimagor555.pro/maven-deps/gradle"
	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/project"
	"dimagor555.pro/maven-deps/update"
	"github.com/spf13/cobra"
)

var (
	updateYes        bool
	updateForceDirty bool
)

var updateCmd = &cobra.Command{
	Use:   "update <patch|minor|major>",
	Short: "Update catalog dependencies to newer versions",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&updateYes, "yes", false, "Skip confirmation prompt")
	updateCmd.Flags().BoolVar(&updateForceDirty, "force-dirty", false, "Proceed even if catalog is dirty in git")
}

type dependencyResolver interface {
	Resolve(ctx context.Context, groupID, artifactID string) (maven.Metadata, error)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	level, err := parseUpdateLevel(args[0])
	if err != nil {
		return err
	}
	root := project.FindRoot(projectPath)
	if root == "" {
		return fmt.Errorf("project root not found")
	}
	catalogPath := filepath.Join(root, "gradle", "libs.versions.toml")
	resolver := maven.NewResolver(discoverRepos(projectPath))
	return executeUpdate(cmd.Context(), cmd.OutOrStdout(), os.Stdin, root, catalogPath, resolver, level, updateYes, updateForceDirty)
}

func parseUpdateLevel(s string) (update.Level, error) {
	switch s {
	case "patch":
		return update.LevelPatch, nil
	case "minor":
		return update.LevelMinor, nil
	case "major":
		return update.LevelMajor, nil
	}
	return "", fmt.Errorf("invalid level %q (expected patch|minor|major)", s)
}

func executeUpdate(
	ctx context.Context,
	w io.Writer,
	in io.Reader,
	root, catalogPath string,
	resolver dependencyResolver,
	level update.Level,
	yes, forceDirty bool,
) error {
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		return fmt.Errorf("read catalog: %w", err)
	}
	content := string(data)
	entries := gradle.ParseVersionCatalog(content)
	upgrades, failures := resolveUpgrades(ctx, entries, resolver)
	plan := update.Plan(upgrades, level)
	if len(plan.VersionRefUpdates) == 0 && len(plan.InlineUpdates) == 0 {
		if len(failures) > 0 {
			printResolveFailures(w, failures)
			return errResolveFailures
		}
		fmt.Fprintln(w, "nothing to update")
		return nil
	}
	if len(failures) > 0 {
		printResolveFailures(w, failures)
	}
	if err := ensureGitClean(root, catalogPath, forceDirty); err != nil {
		return err
	}
	printUpdatePreview(w, plan)
	if !yes {
		ok, _ := Confirm(w, in, "Apply these updates?")
		if !ok {
			fmt.Fprintln(w, "aborted")
			return nil
		}
	}
	newContent, err := applyPlan(content, plan)
	if err != nil {
		return err
	}
	if err := os.WriteFile(catalogPath, []byte(newContent), 0o644); err != nil {
		return err
	}
	if len(failures) > 0 {
		return errResolveFailures
	}
	return nil
}

func ensureGitClean(root, catalogPath string, forceDirty bool) error {
	dirty, err := project.IsDirty(root, []string{catalogPath})
	if err != nil {
		return nil
	}
	if dirty && !forceDirty {
		return fmt.Errorf("catalog file is dirty in git; commit changes or use --force-dirty")
	}
	return nil
}

func printUpdatePreview(w io.Writer, plan update.Result) {
	fmt.Fprintln(w, "Upgrades to apply:")
	for _, ru := range plan.VersionRefUpdates {
		fmt.Fprintf(w, "  [versions] %s: %s → %s  (used by %s)\n",
			ru.Ref, ru.OldVersion, ru.NewVersion, strings.Join(ru.Aliases, ", "))
	}
	for _, iu := range plan.InlineUpdates {
		fmt.Fprintf(w, "  [%s] %s: %s → %s (%s)\n",
			iu.Section, iu.Alias, iu.OldVersion, iu.NewVersion, iu.Type)
	}
	total := len(plan.VersionRefUpdates) + len(plan.InlineUpdates)
	fmt.Fprintf(w, "\n%d changes\n", total)
}

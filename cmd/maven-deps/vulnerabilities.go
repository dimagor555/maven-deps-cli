package main

import (
	"context"
	"fmt"

	"dimagor555.pro/maven-deps/vulnerability"
	"github.com/spf13/cobra"
)

var vulnsCmd = &cobra.Command{
	Use:   "vulnerabilities",
	Short: "Check dependencies for known CVEs via OSV.dev",
	RunE:  runVulnerabilities,
}

func runVulnerabilities(cmd *cobra.Command, _ []string) error {
	deps, err := scanProject(projectPath)
	if err != nil {
		return err
	}

	var refs []vulnerability.DependencyRef
	for _, dep := range deps {
		if dep.Version == "" {
			continue
		}
		refs = append(refs, vulnerability.DependencyRef{
			GroupID:    dep.GroupID,
			ArtifactID: dep.ArtifactID,
			Version:    dep.Version,
		})
	}

	if len(refs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No dependencies with versions to check")
		return nil
	}

	ctx := context.Background()
	results, err := vulnerability.QueryBatch(ctx, refs)
	if err != nil {
		return fmt.Errorf("query OSV: %w", err)
	}

	var vulnCount int
	var vulnDeps []vulnerability.Result

	for _, r := range results {
		if len(r.Vulnerabilities) > 0 {
			vulnCount++
			vulnDeps = append(vulnDeps, r)
		}
	}

	if jsonOutput {
		printJSON(vulnDeps)
		return nil
	}

	if vulnCount == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No vulnerabilities found")
		return nil
	}

	for _, r := range vulnDeps {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s:%s\n", r.GroupID, r.ArtifactID, r.Version)
		for _, v := range r.Vulnerabilities {
			fixed := ""
			if v.FixedVersion != "" {
				fixed = fmt.Sprintf("  (fixed: %s)", v.FixedVersion)
			}
			severity := v.Severity
			if severity == "" {
				severity = "UNKNOWN"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  %s%s\n", v.ID, severity, v.Summary, fixed)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%d dependencies with vulnerabilities\n", vulnCount)
	return nil
}

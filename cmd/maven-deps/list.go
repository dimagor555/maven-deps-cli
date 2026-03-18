package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all dependencies from project build files",
	RunE:  runList,
}

func runList(cmd *cobra.Command, _ []string) error {
	deps, err := scanProject(projectPath)
	if err != nil {
		return err
	}

	if jsonOutput {
		printJSON(deps)
		return nil
	}

	for _, d := range deps {
		ver := d.Version
		if ver == "" {
			ver = "(no version)"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s  %s  %s  %s\n",
			d.GroupID, d.ArtifactID, ver, d.Configuration, d.Source)
	}
	return nil
}

package main

import (
	"context"
	"fmt"

	"dimagor555.pro/maven-deps/version"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check group:artifact version",
	Short: "Check if a specific version exists",
	Args:  cobra.ExactArgs(2),
	RunE:  runCheck,
}

type checkResult struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Exists     bool   `json:"exists"`
	Stability  string `json:"stability,omitempty"`
	Repository string `json:"repository,omitempty"`
}

func runCheck(cmd *cobra.Command, args []string) error {
	groupID, artifactID, err := parseGA(args[0])
	if err != nil {
		return err
	}
	ver := args[1]

	repos := discoverRepos(projectPath)
	ctx := context.Background()

	for _, repo := range repos {
		meta, err := repo.FetchMetadata(ctx, groupID, artifactID)
		if err != nil {
			continue
		}
		for _, v := range meta.Versions {
			if v == ver {
				r := checkResult{
					GroupID:    groupID,
					ArtifactID: artifactID,
					Version:    ver,
					Exists:     true,
					Stability:  string(version.Classify(ver)),
					Repository: repo.Name,
				}
				if jsonOutput {
					printJSON(r)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%s:%s:%s  exists  %s  %s\n",
						r.GroupID, r.ArtifactID, r.Version, r.Stability, r.Repository)
				}
				return nil
			}
		}
	}

	r := checkResult{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Version:    ver,
		Exists:     false,
	}
	if jsonOutput {
		printJSON(r)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s:%s  not found\n", r.GroupID, r.ArtifactID, r.Version)
	}
	return nil
}

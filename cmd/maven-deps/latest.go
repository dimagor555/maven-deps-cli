package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/version"
	"github.com/spf13/cobra"
)

var (
	stableOnly  bool
	allVersions bool
)

var latestCmd = &cobra.Command{
	Use:   "latest [group:artifact...]",
	Short: "Get latest versions of Maven artifacts",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runLatest,
}

func init() {
	latestCmd.Flags().BoolVar(&stableOnly, "stable", false, "Only stable versions")
	latestCmd.Flags().BoolVar(&allVersions, "all", false, "Include all pre-release versions")
}

type latestResult struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Stability  string `json:"stability"`
	Error      string `json:"error,omitempty"`
	failed     bool
}

func runLatest(cmd *cobra.Command, args []string) error {
	filter := version.PreferStable
	if stableOnly {
		filter = version.StableOnly
	} else if allVersions {
		filter = version.All
	}

	repos := discoverRepos(projectPath)
	ctx := context.Background()
	resolver := maven.NewResolver(repos)
	sem := make(chan struct{}, maxConcurrency)
	results := make([]latestResult, len(args))

	var wg sync.WaitGroup
	wg.Add(len(args))

	for i, arg := range args {
		go func(i int, arg string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			groupID, artifactID, err := parseGA(arg)
			if err != nil {
				results[i] = latestResult{Error: err.Error()}
				return
			}

			meta, err := resolver.Resolve(ctx, groupID, artifactID)
			if err != nil {
				results[i] = latestResult{GroupID: groupID, ArtifactID: artifactID, Error: err.Error(), failed: !maven.IsNotFound(err)}
				return
			}

			selected := version.FindLatest(meta.Versions, filter)
			if selected == "" {
				results[i] = latestResult{GroupID: groupID, ArtifactID: artifactID, Error: "no version found"}
				return
			}

			results[i] = latestResult{
				GroupID:    groupID,
				ArtifactID: artifactID,
				Version:    selected,
				Stability:  string(version.Classify(selected)),
			}
		}(i, arg)
	}
	wg.Wait()

	anyFailed := false
	for _, r := range results {
		if r.failed {
			anyFailed = true
			break
		}
	}

	if jsonOutput {
		printJSON(results)
		if anyFailed {
			return errResolveFailures
		}
		return nil
	}

	for _, r := range results {
		if r.Error != "" {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s:%s  error: %s\n", r.GroupID, r.ArtifactID, r.Error)
			continue
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s  %s  %s\n", r.GroupID, r.ArtifactID, r.Version, r.Stability)
	}
	if anyFailed {
		return errResolveFailures
	}
	return nil
}

func parseGA(s string) (string, string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid format %q, expected group:artifact", s)
	}
	return parts[0], parts[1], nil
}

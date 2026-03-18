package main

import (
	"context"
	"fmt"
	"sync"

	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/version"
	"github.com/spf13/cobra"
)

const maxConcurrency = 20

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Show dependencies that have newer versions available",
	RunE:  runOutdated,
}

type outdatedResult struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Current    string `json:"current"`
	Latest     string `json:"latest"`
	Upgrade    string `json:"upgrade"`
}

type outdatedSummary struct {
	Results []outdatedResult `json:"results"`
	Total   int              `json:"total"`
	Major   int              `json:"major"`
	Minor   int              `json:"minor"`
	Patch   int              `json:"patch"`
}

type indexedResult struct {
	idx    int
	result *outdatedResult
}

func runOutdated(cmd *cobra.Command, _ []string) error {
	deps, err := scanProject(projectPath)
	if err != nil {
		return err
	}

	repos := discoverRepos(projectPath)
	ctx := context.Background()
	resolver := maven.NewResolver(repos)
	sem := make(chan struct{}, maxConcurrency)

	var withVersion []struct {
		idx int
		dep scannedDep
	}
	for i, dep := range deps {
		if dep.Version != "" {
			withVersion = append(withVersion, struct {
				idx int
				dep scannedDep
			}{i, dep})
		}
	}

	ch := make(chan indexedResult, len(withVersion))
	var wg sync.WaitGroup
	wg.Add(len(withVersion))

	for _, item := range withVersion {
		go func(idx int, dep scannedDep) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			meta, err := resolver.Resolve(ctx, dep.GroupID, dep.ArtifactID)
			if err != nil {
				ch <- indexedResult{idx: idx}
				return
			}
			latest := version.FindLatestForCurrent(meta.Versions, dep.Version)
			if latest == "" || latest == dep.Version {
				ch <- indexedResult{idx: idx}
				return
			}
			upgrade := version.GetUpgradeType(dep.Version, latest)
			if upgrade == version.None {
				ch <- indexedResult{idx: idx}
				return
			}
			ch <- indexedResult{idx: idx, result: &outdatedResult{
				GroupID:    dep.GroupID,
				ArtifactID: dep.ArtifactID,
				Current:    dep.Version,
				Latest:     latest,
				Upgrade:    string(upgrade),
			}}
		}(item.idx, item.dep)
	}

	if jsonOutput {
		wg.Wait()
		close(ch)
		collected := make(map[int]*outdatedResult)
		for item := range ch {
			if item.result != nil {
				collected[item.idx] = item.result
			}
		}
		var results []outdatedResult
		for i := range deps {
			if r, ok := collected[i]; ok {
				results = append(results, *r)
			}
		}
		printJSON(buildSummary(results))
		return nil
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []outdatedResult
	for item := range ch {
		if item.result != nil {
			results = append(results, *item.result)
			r := item.result
			fmt.Fprintf(cmd.OutOrStdout(), "%s:%s  %s → %s  %s\n",
				r.GroupID, r.ArtifactID, r.Current, r.Latest, r.Upgrade)
		}
	}

	if len(results) > 0 {
		s := buildSummary(results)
		fmt.Fprintf(cmd.OutOrStdout(), "\n%d outdated (%d major, %d minor, %d patch)\n",
			s.Total, s.Major, s.Minor, s.Patch)
	}
	return nil
}

func buildSummary(results []outdatedResult) outdatedSummary {
	s := outdatedSummary{Results: results, Total: len(results)}
	for _, r := range results {
		switch version.UpgradeType(r.Upgrade) {
		case version.Major:
			s.Major++
		case version.Minor:
			s.Minor++
		case version.Patch:
			s.Patch++
		}
	}
	return s
}

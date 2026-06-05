package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"dimagor555.pro/maven-deps/maven"
	"dimagor555.pro/maven-deps/version"
	"github.com/spf13/cobra"
)

const maxConcurrency = 4

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
	Failed  []failedResult   `json:"failed,omitempty"`
}

type failedResult struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Error      string `json:"error"`
}

type indexedResult struct {
	idx    int
	result *outdatedResult
	failed *resolveFailure
}

type resolveFailure struct {
	GroupID    string
	ArtifactID string
	Err        error
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
				if maven.IsNotFound(err) {
					ch <- indexedResult{idx: idx}
					return
				}
				ch <- indexedResult{idx: idx, failed: &resolveFailure{
					GroupID:    dep.GroupID,
					ArtifactID: dep.ArtifactID,
					Err:        err,
				}}
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	total := len(withVersion)
	done := 0
	var results []outdatedResult
	var failures []resolveFailure
	for item := range ch {
		done++
		if !jsonOutput {
			fmt.Fprintf(os.Stderr, "\r%s", progressLine(done, total))
		}
		if item.result != nil {
			results = append(results, *item.result)
		}
		if item.failed != nil {
			failures = append(failures, *item.failed)
		}
	}
	if !jsonOutput && total > 0 {
		fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 30))
	}
	sortResults(results)

	if jsonOutput {
		summary := buildSummary(results)
		summary.Failed = toFailedResults(failures)
		printJSON(summary)
		if len(failures) > 0 {
			return errResolveFailures
		}
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s  %s → %s  %s\n",
			r.GroupID, r.ArtifactID, r.Current, r.Latest, r.Upgrade)
	}

	if len(results) > 0 {
		s := buildSummary(results)
		fmt.Fprintf(cmd.OutOrStdout(), "\n%d outdated (%d major, %d minor, %d patch)\n",
			s.Total, s.Major, s.Minor, s.Patch)
	}

	if len(failures) > 0 {
		printResolveFailures(cmd.ErrOrStderr(), failures)
		return errResolveFailures
	}
	return nil
}

var errResolveFailures = fmt.Errorf("some dependencies could not be resolved")

func toFailedResults(failures []resolveFailure) []failedResult {
	out := make([]failedResult, 0, len(failures))
	for _, f := range failures {
		out = append(out, failedResult{
			GroupID:    f.GroupID,
			ArtifactID: f.ArtifactID,
			Error:      f.Err.Error(),
		})
	}
	return out
}

func printResolveFailures(w io.Writer, failures []resolveFailure) {
	fmt.Fprintf(w, "\nWARNING: failed to resolve %d dependencies (results above may be incomplete):\n", len(failures))
	for _, f := range failures {
		fmt.Fprintf(w, "  %s:%s  %v\n", f.GroupID, f.ArtifactID, f.Err)
	}
}

func progressLine(done, total int) string {
	return fmt.Sprintf("checking %d/%d...", done, total)
}

func upgradeOrder(t string) int {
	switch version.UpgradeType(t) {
	case version.Major:
		return 0
	case version.Minor:
		return 1
	case version.Patch:
		return 2
	default:
		return 3
	}
}

func sortResults(results []outdatedResult) {
	sort.Slice(results, func(i, j int) bool {
		oi := upgradeOrder(results[i].Upgrade)
		oj := upgradeOrder(results[j].Upgrade)
		if oi != oj {
			return oi < oj
		}
		ki := results[i].GroupID + ":" + results[i].ArtifactID
		kj := results[j].GroupID + ":" + results[j].ArtifactID
		return ki < kj
	})
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

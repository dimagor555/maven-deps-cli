package main

import (
	"context"
	"fmt"
	"strings"

	"dimagor555.pro/maven-deps/search"
	"github.com/spf13/cobra"
)

var searchLimit int

var searchCmd = &cobra.Command{
	Use:   "search query",
	Short: "Search Maven Central by keyword",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "Max results (default: 10, max: 100)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	ctx := context.Background()

	results, err := search.Search(ctx, query, searchLimit)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if jsonOutput {
		printJSON(results)
		return nil
	}

	for _, r := range results {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%s  %s  (%d versions)\n",
			r.GroupID, r.ArtifactID, r.LatestVersion, r.VersionCount)
	}
	return nil
}

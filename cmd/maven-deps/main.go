package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	projectPath string
	jsonOutput  bool
)

var rootCmd = &cobra.Command{
	Use:     "maven-deps",
	Short:   "Maven/Gradle dependency intelligence",
	Version: "1.1.0",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "C", ".", "Project path for repository discovery")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "JSON output")

	rootCmd.AddCommand(latestCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(outdatedCmd)
	rootCmd.AddCommand(vulnsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "repos",
	Short: "A CLI tool for discovering and managing local Git repositories",
	Long: `Repos is a CLI tool that helps you discover, list, inspect, and open
Git repositories on your local machine.

Configure one or more directories to scan, and repos will find all
Git repositories within them.`,
}

func Execute() {
	rootCmd.Execute()
}

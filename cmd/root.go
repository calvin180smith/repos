package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "repos",
	Short: "cli tool for exploring repos",
}

func Execute(){
	rootCmd.Execute()
}
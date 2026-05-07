package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func Open(repoPath string) (string, error) {

	cmd := exec.Command("code", ".")
	cmd.Dir = repoPath
	out, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	outputStr := strings.TrimSpace(string(out))

	return outputStr, nil

}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open a repo in VS Code",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := readConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		latest, _ := cmd.Flags().GetBool("latest")

		if latest {
			listRepos, err := listRepos(config.Path)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			latestRepo := listRepos[0]
			Open(latestRepo.Path)

		} else if len(args) > 0 {
			path := args[0]
			repoPath := filepath.Join(config.Path, path)
			Open(repoPath)

		} else {
			fmt.Println("please provide a repo name or use --latest")
		}

	},
}

func init() {
	openCmd.Flags().Bool("latest", false, "open the most recently modified repo")
	rootCmd.AddCommand(openCmd)
}

package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)

func Open(repoPath string,editor string) (string, error) {

	cmd := exec.Command(editor, ".")
	cmd.Dir = repoPath
	out, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	outputStr := strings.TrimSpace(string(out))

	return outputStr, nil

}

var openCmd = &cobra.Command{
	Use:   "open [repository name]",
	Short: "Open a repository in VS Code",
	Long:  "Opens the specified repository in VS Code by name, or use --latest to open the most recently modified one.",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := readConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		ide, err := cmd.Flags().GetString("ide")
		if err != nil {
			panic(err)
		}

		latest, _ := cmd.Flags().GetBool("latest")

		if latest {
			listRepos, err := listRepos(config.Paths)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			latestRepo := listRepos[0]
			if ide != "" {
				Open(latestRepo.Path,ide)
			} else {
				Open(latestRepo.Path,config.Editor)
			}

		} else if len(args) > 0 {
			repoName := args[0]
			repoPath, err := findRepoPath(repoName, config.Paths)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			if ide != "" {
				Open(repoPath,ide)
			} else {
				Open(repoPath,config.Editor)
			}

		} else {
			fmt.Println("please provide a repo name or use --latest")
		}

	},
}

func init() {
	openCmd.Flags().Bool("latest", false, "open the most recently modified repository")
	openCmd.Flags().String("ide","","editor to use when opening repo")
	rootCmd.AddCommand(openCmd)
}

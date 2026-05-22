package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)

type RepoInfo struct {
	Name             string
	Path             string
	CurrentBranch    string
	RemoteUrl        string
	NumberofBranches int
	LastCommitInfo   string
}

func (r *RepoInfo) Print() {
	fmt.Printf("%-12s %s\n", "Name:", r.Name)
	fmt.Printf("%-12s %s\n", "Path:", r.Path)
	fmt.Printf("%-12s %s\n", "Remote:", r.RemoteUrl)
	fmt.Printf("%-12s %s\n", "Branch:", r.CurrentBranch)
	fmt.Printf("%-12s %v\n", "Branches:", r.NumberofBranches)
	fmt.Printf("%-12s %s\n", "Last Commit:", r.LastCommitInfo)

}

func runGitCommand(repoPath string, args ...string) (string, error) {

	cmdArgs := append([]string{"-C", repoPath}, args...)
	output, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return "N/A", nil
	}

	return outputStr, nil

}

func repoInfo(repoPath string) (RepoInfo, error) {

	remoteArgs := []string{"config", "--get", "remote.origin.url"}
	remote, err := runGitCommand(repoPath, remoteArgs...)
	if err != nil {
		remote = "No remote branch url found"
	}

	branchArgs := []string{"branch", "--show-current"}
	branch, err := runGitCommand(repoPath, branchArgs...)
	if err != nil {
		branch = fmt.Sprintf("Could not fetch current branch: %s", err)
	}

	nrBranchArgs := []string{"branch"}
	nrBranch, err := runGitCommand(repoPath, nrBranchArgs...)

	branches := 0
	for _, line := range strings.Split(nrBranch, "\n") {
		if strings.TrimSpace(line) != "" {
			branches++
		}
	}

	lastCommitArgs := []string{"log", "-1", "--format=%s (%cd) <%an>)"}
	lastCommitInfo, err := runGitCommand(repoPath, lastCommitArgs...)
	if err != nil {
		lastCommitInfo = "No commits found"
	}

	return RepoInfo{RemoteUrl: remote, CurrentBranch: branch, NumberofBranches: branches, LastCommitInfo: lastCommitInfo}, err

}

var infoCmd = &cobra.Command{
	Use:   "info [repository name]",
	Short: "Show detailed information about a repository",
	Long:  "Displays information about a Git repository including its remote URL, current branch, number of local branches, and last commit details.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		config, err := readConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		repoName := args[0]

		repoPath, err := findRepoPath(repoName, config.Paths)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		info, err := repoInfo(repoPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		info.Path = repoPath
		info.Name = repoName
		info.Print()

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

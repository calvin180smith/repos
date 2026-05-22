package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type Repo struct {
	Name         string
	Path         string
	LastModified time.Time
}

func listRepos(paths []string) ([]Repo, error) {

	var reposList []Repo

	for _, path := range paths {

		items, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("Error: %w", err)
		}

		for _, item := range items {

			if item.IsDir() {
				repoPath := filepath.Join(path, item.Name())
				gitPath := filepath.Join(repoPath, ".git")
				_, err := os.Stat(gitPath)
				if err != nil {
					continue
				}

				info, err := os.Stat(filepath.Join(gitPath, "COMMIT_EDITMSG"))
				if err != nil {
					info, err = os.Stat(gitPath)
					if err != nil {
						return nil, fmt.Errorf("Error: %w", err)
					}

				}

				lastModified := info.ModTime()

				reposList = append(reposList, Repo{Name: item.Name(), Path: repoPath, LastModified: lastModified})

			}
		}
	}

	sort.Slice(reposList, func(i, j int) bool {
		return reposList[i].LastModified.After(reposList[j].LastModified)
	})

	return reposList, nil

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered Git repositories, sorted by last modified",
	Long: `Lists all Git repositories found in configured directories, sorted by
most recently modified first. Use --limit to restrict the number of results.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetString("limit")

		config, err := readConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		repos, err := listRepos(config.Paths)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if limit != "" {
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				fmt.Printf("Please give a valid limit integer \n Error: %v\n", err)
				return
			}
			if limitInt > len(repos) {
				limitInt = len(repos)
			}
			repos = repos[0:limitInt]

		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"NAME", "PATH", "LAST MODIFIED"})
		t.AppendSeparator()

		var rows []table.Row

		for _, repo := range repos {

			timeSince := time.Since(repo.LastModified)

			hours := timeSince.Hours()
			minutes := timeSince.Minutes()
			seconds := timeSince.Seconds()

			var returnTimeString string

			if hours < 1 {
				if minutes == 0 {
					returnTimeString = fmt.Sprintf("%.0f seconds ago", seconds)
				} else if minutes == 1 {
					returnTimeString = fmt.Sprintf("%.0f minute ago", minutes)
				}
				returnTimeString = fmt.Sprintf("%.0f minutes ago", minutes)
			} else if hours == 1 {
				returnTimeString = fmt.Sprintf("%.0f hour and %.0f minutes ago", hours, minutes-60)
			} else if hours > 1 && hours < 24 {
				returnTimeString = fmt.Sprintf("%.0f hours ago", hours)
			} else if hours > 24 && hours < 25 {
				returnTimeString = fmt.Sprintf("%.0f day ago", math.Round(hours/24))
			} else if hours > 25 {
				returnTimeString = fmt.Sprintf("%.0f days ago", math.Round(hours/24))
			}

			rows = append(rows, table.Row{repo.Name, repo.Path, returnTimeString})
		}

		t.AppendRows(rows)
		t.Render()

	},
}

func init() {
	listCmd.Flags().String("limit", "", "maximum number of repositories to display")
	rootCmd.AddCommand(listCmd)

}

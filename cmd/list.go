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

func listRepos(path string) ([]Repo, error) {
	items, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	var repos []Repo

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

			repos = append(repos, Repo{Name: item.Name(), Path: repoPath, LastModified: lastModified})

		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].LastModified.After(repos[j].LastModified)
	})

	return repos, nil

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repos",
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetString("limit")

		config, err := readConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		repos, err := listRepos(config.Path)
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

			if hours == 0 {
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

			// if hours == 1 {
			// 	returnTimeString = fmt.Sprintf("%.1f hour and %.1f minutes ago",hours,minutes)
			// } else if hours > 1 {
			// 	returnTimeString = fmt.Sprintf("%.1f hours and %.1f minutes ago",hours,minutes)

			// } else if hours == 0 {
			// 	if minutes == 0 {
			// 		returnTimeString = fmt.Sprintf("%.1f seconds ago",seconds)
			// 	}
			// 	returnTimeString = fmt.Sprintf("%.1f minutes ago",minutes)
			// }

			rows = append(rows, table.Row{repo.Name, repo.Path, returnTimeString})
		}

		t.AppendRows(rows)
		t.Render()

	},
}

func init() {
	listCmd.Flags().String("limit", "", "limit the number of repos to list")
	rootCmd.AddCommand(listCmd)

}

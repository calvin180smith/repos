package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
			repoPath := filepath.Join(path,item.Name())
			gitPath := filepath.Join(repoPath,".git")
			_, err := os.Stat(gitPath)
			if err != nil {
				continue
			}

			info, err := os.Stat(repoPath)
			if err != nil {
				return nil, fmt.Errorf("Error: %w",err)
			}

			repos = append(repos,Repo{Name: item.Name(),Path: repoPath,LastModified: info.ModTime()})
			
			
			// info, err := os.Stat(item.Name())
			// if err != nil{
			// 	return nil, fmt.Errorf("Error: %w",err)
			// }
			// repoName := info.Name()
			// lastModified := info.ModTime()
			// _, err := os.Stat()




			
		}
	}
	return repos, nil

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repos",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	listCmd.Flags().String("limit", "", "limit the number of repos to list")
	rootCmd.AddCommand(listCmd)

}

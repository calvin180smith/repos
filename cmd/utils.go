package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Home directory not found: %w", err)
	}
	return filepath.Join(dir, ".repos.yaml"), nil
}


func readConfig() (*Config, error) {
	dir, err := getConfigFilePath()
	if err != nil {
		return &Config{},fmt.Errorf("Error %w",err)
	}
	file, err := os.ReadFile(dir)
	if err != nil {
		return &Config{},fmt.Errorf("Config file not found, please use repos config set --path <path> to set it.")
	}
	config := Config{}
	if err := yaml.Unmarshal(file,&config); err != nil {
		return &Config{},fmt.Errorf("Could not unmarshal config: %w",err)
	}


	return &config, nil

}

func checkDuplicatePath(paths []string) bool {
    seen := make(map[string]struct{}, len(paths))
    for _, path := range paths {
        if _, ok := seen[path]; ok {
            return true
        }
        seen[path] = struct{}{}
    }
    return false
}

func findRepoPath(repoName string,paths []string) (string, error) {
	repos, err := listRepos(paths)
	if err != nil {
		return "", fmt.Errorf("Error: %v\n", err)
	}

	for r := range repos {
		if repos[r].Name == repoName {
			return repos[r].Path, nil
			
		}
	}

	return "", fmt.Errorf("Could not find repo")
}
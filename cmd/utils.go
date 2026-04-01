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
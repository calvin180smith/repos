package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Path string `yaml:"path"`
}

func setConfig(path string, configFilePath string) error {

	f, err := os.ReadFile(configFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			out, err := yaml.Marshal(Config{Path: path})
			if err != nil {
				return fmt.Errorf("could not marshal config: %w", err)
			}
			err = os.WriteFile(configFilePath, out, 0644)
			if err != nil {
				return fmt.Errorf("could not write config: %w", err)
			}
			return nil

		}
	}
	config := Config{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		return fmt.Errorf("could not unmarshal config: %w", err)
	}

	config.Path = path
	out, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}
	err = os.WriteFile(configFilePath, out, 0644)
	if err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return err

}

var cfgCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize config",
}

var setCfgCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config for repos",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		configFilePath, err := getConfigFilePath()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		err = setConfig(path, configFilePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Printf("Config file successfully set at %v:", configFilePath)

	},
}

var showCfgCmd = &cobra.Command{
	Use:   "show",
	Short: "Show config",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readConfig()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(config)


	},
}

func init() {
	setCfgCmd.Flags().String("path", "", "path to config file")
	setCfgCmd.MarkFlagRequired("set")
	cfgCmd.AddCommand(showCfgCmd)
	cfgCmd.AddCommand(setCfgCmd)
	rootCmd.AddCommand(cfgCmd)

}

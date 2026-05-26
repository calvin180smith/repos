package cmd

import (
	"errors"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Paths []string `yaml:"paths"`
	Editor string   `yaml:"editor"`
}

func setConfig(paths []string,ide string ,configFilePath string) error {

	if ide == "" {
		ide = "code"
	}

	f, err := os.ReadFile(configFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			out, err := yaml.Marshal(Config{Paths: paths,Editor: ide})
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

	config.Paths = paths
	config.Editor = ide

	dupCheck := checkDuplicatePath(config.Paths)
	if dupCheck {
		return fmt.Errorf("Paths must be unique")
	}

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

func addPath(paths []string, configFilePath string) error {
	f, err := os.ReadFile(configFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			out, err := yaml.Marshal(Config{Paths: paths})
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

	config.Paths = append(config.Paths, paths...)

	dupCheck := checkDuplicatePath(config.Paths)
	if dupCheck {
		return fmt.Errorf("Paths must be unique")
	}

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
	Short: "Manage directory paths where repos scans for Git repositories",
	Long:  "Configure which directories repos should scan for Git repositories.\nSettings are stored in ~/.repos.yaml.",
}

var paths []string
var addPaths []string

var setCfgCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the list of directories to scan (replaces existing paths)",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := cmd.Flags().GetStringSlice("path")
		if err != nil {
			panic(err)
		}

		ide, err := cmd.Flags().GetString("ide")
		if err != nil {
			panic(err)
		}

		configFilePath, err := getConfigFilePath()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		err = setConfig(paths,ide,configFilePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Printf("Config file successfully set at %v:", configFilePath)

	},
}

var showCfgCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readConfig()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%-12s %s\n", "paths:", config.Paths)
		fmt.Printf("%-12s %s\n", "editor:", config.Editor)
	},
}

var addPathCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more directories to the scan list",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := cmd.Flags().GetStringSlice("path")
		if err != nil {
			panic(err)
		}

		dupCheck := checkDuplicatePath(paths)
		if dupCheck {
			fmt.Printf("Paths must be unique")
			return
		}

		configFilePath, err := getConfigFilePath()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		err = addPath(paths, configFilePath)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Printf("Config file successfully updated at %v:", configFilePath)

	},
}



func init() {
	setCfgCmd.Flags().StringSliceVar(&paths, "path", []string{}, "directories to scan for Git repositories (comma-separated or repeated)")
	setCfgCmd.MarkFlagRequired("path")

	setCfgCmd.Flags().String("ide","","default editor to open a repo in (defaults to VS Code)")

	addPathCmd.Flags().StringSliceVar(&addPaths, "path", []string{}, "directories to add (comma-separated or repeated)")
	addPathCmd.MarkFlagRequired("path")

	cfgCmd.AddCommand(showCfgCmd)
	cfgCmd.AddCommand(setCfgCmd)
	cfgCmd.AddCommand(addPathCmd)
	rootCmd.AddCommand(cfgCmd)

}

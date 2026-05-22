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
}

func setConfig(paths []string, configFilePath string) error {

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

	config.Paths = paths

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

		configFilePath, err := getConfigFilePath()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		err = setConfig(paths, configFilePath)
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
		fmt.Println(config)

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

	addPathCmd.Flags().StringSliceVar(&addPaths, "path", []string{}, "directories to add (comma-separated or repeated)")
	addPathCmd.MarkFlagRequired("path")

	cfgCmd.AddCommand(showCfgCmd)
	cfgCmd.AddCommand(setCfgCmd)
	cfgCmd.AddCommand(addPathCmd)
	rootCmd.AddCommand(cfgCmd)

}

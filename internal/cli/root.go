package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Persistent flags
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without subcommands
var rootCmd = &cobra.Command{
	Use:   "skill",
	Short: "Skills manager from Git repositories",
	Long: `skill is a CLI for managing skills stored in Git repositories.

Skills are structured as individual directories with multiple files
(config.yaml, scripts, README.md, etc.) and are stored in a shared Git repository.

This CLI allows you to initialize, install, list and manage skills from
a configured remote repository.

Available Commands:
  repository    Manage skill repositories (add, remove, list, set-current)
  config        Manage CLI configuration (show, set-repo, verify)
  list          List available skills from repositories
  install       Install a skill to a provider (Claude, Cursor)

Global Flags:
  --config string   Configuration file (default: ~/.config/skill/config.yaml)
  -v, --verbose     Detailed output
  --help            Show help for command
  --version         Show version information`,
	Version: "0.1.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file (default: ~/.config/skill/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "detailed output")

	// Add subcommands
	rootCmd.AddCommand(repositoryCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
}

// initConfig reads the configuration
func initConfig() {
	if cfgFile != "" {
		// If the user specified a config file, use it
		// For now we don't do anything special with this
	}

	if verbose {
		fmt.Println("Verbose mode activated")
	}
}

// GetVerbose returns whether verbose mode is enabled
func GetVerbose() bool {
	return verbose
}

// SetExitOnError sets whether Cobra should exit on error
func SetExitOnError(exit bool) {
	if !exit {
		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true
	}
}

// PrintError prints an error in a friendly format
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Println(msg)
}

// PrintInfo prints an informational message
func PrintInfo(msg string) {
	if verbose {
		fmt.Println(msg)
	}
}

package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Persistent flags
	cfgFile string
	verbose bool
)

// getColoredHelp returns the colored help text for the root command
func getColoredHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

Skills are structured as individual directories with multiple files
(config.yaml, scripts, README.md, etc.) and are stored in a shared Git repository.

This CLI allows you to initialize, install, list and manage skills from
a configured remote repository.

%s
  %s
  %s

%s
  %s              %s
    %s            %s
    %s               %s
    %s                        %s
    %s          %s

  %s                %s
    %s                        %s
    %s              %s
    %s                      %s

  %s                %s
  %s %s
  %s %s

  %s              %s
  %s                     %s

%s
  %s   %s
  %s     %s
  %s        %s

%s
  %s
  %s
  %s
  %s
  %s

%s`,
		white("skill is a CLI for managing skills stored in Git repositories."),

		cyan("USAGE:"),
		green("skill [command] [flags]"),
		green("skill [command] [subcommand] [arguments] [flags]"),

		cyan("AVAILABLE COMMANDS:"),
		yellow("repository <subcommand>"), dim("Manage skill repositories"),
		green("add <name> <url>"), dim("Add a new repository"),
		green("remove <name>"), dim("Remove a repository"),
		green("list"), dim("List all repositories"),
		green("set-current <name>"), dim("Set active repository"),

		yellow("config <subcommand>"), dim("Manage CLI configuration"),
		green("show"), dim("Show current configuration"),
		green("set-repo <url>"), dim("Change repository URL"),
		green("verify"), dim("Verify repository access"),

		yellow("list [flags]"), dim("List available or installed skills"),
		yellow("install <skill-name> [flags]"), dim("Install a skill to a provider"),
		yellow("share --path <path> [flags]"), dim("Share a skill via PR"),

		yellow("help [command]"), dim("Show help for any command"),
		yellow("version"), dim("Show version information"),

		cyan("GLOBAL FLAGS:"),
		green("--config string"), dim("Configuration file (default: ~/.config/skill/config.yaml)"),
		green("-v, --verbose"), dim("Detailed output"),
		green("-h, --help"), dim("Show help for command"),

		cyan("EXAMPLES:"),
		green("skill repository add myrepo https://github.com/org/skills.git"),
		green("skill list"),
		green("skill install explain-code --provider claude"),
		green("skill share --path ./my-skill"),
		green("skill config show"),

		dim("Use \"skill [command] --help\" for more information about a command."),
	)
}

// rootCmd represents the base command when called without subcommands
var rootCmd = &cobra.Command{
	Use:   "skill [command]",
	Short: "Skills manager from Git repositories",
	Long:  getColoredHelp(),
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
	rootCmd.AddCommand(shareCmd)
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

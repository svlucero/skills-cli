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
	yellow := color.New(color.FgYellow).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return NewHelpBuilder().
		Text(white("skills is a CLI for managing skills stored in Git repositories.")).
		EmptyLine().
		Text("Skills are structured as individual directories with multiple files").
		Text("(config.yaml, scripts, README.md, etc.) and are stored in a shared Git repository.").
		EmptyLine().
		Text("This CLI allows you to initialize, install, list and manage skills from").
		Text("a configured remote repository.").
		Section("USAGE:").
		Example("skills [command] [flags]", "").
		Example("skills [command] [subcommand] [arguments] [flags]", "").
		Section("AVAILABLE COMMANDS:").
		Item(yellow("repository <subcommand>"), "Manage skill repositories").
		SubItem("add <name> <url>", "Add a new repository").
		SubItem("remove <name>", "Remove a repository").
		SubItem("list", "List all repositories").
		SubItem("set-current <name>", "Set active repository").
		EmptyLine().
		Item(yellow("config <subcommand>"), "Manage CLI configuration").
		SubItem("show", "Show current configuration").
		SubItem("set-repo <url>", "Change repository URL").
		SubItem("verify", "Verify repository access").
		EmptyLine().
		Item(yellow("list [flags]"), "List available or installed skills").
		Item(yellow("install <skill-name> [flags]"), "Install a skill to a provider").
		Item(yellow("uninstall <skill-name> [flags]"), "Uninstall a skill from a provider").
		Item(yellow("share --path <path> [flags]"), "Share a skill via PR").
		EmptyLine().
		Item(yellow("help [command]"), "Show help for any command").
		Item(yellow("version"), "Show version information").
		Section("GLOBAL FLAGS:").
		Item("--config string", "Configuration file (default: ~/.config/skill/config.yaml)").
		Item("-v, --verbose", "Detailed output").
		Item("-h, --help", "Show help for command").
		Section("EXAMPLES:").
		Example("skills repository add myrepo https://github.com/org/skills.git", "").
		Example("skills list", "").
		Example("skills install explain-code --provider claude", "").
		Example("skills share --path ./my-skill", "").
		Example("skills config show", "").
		EmptyLine().
		Text(dim("Use \"skills [command] --help\" for more information about a command.")).
		Build()
}

// rootCmd represents the base command when called without subcommands
var rootCmd = &cobra.Command{
	Use:   "skills [command]",
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
	rootCmd.AddCommand(uninstallCmd)
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

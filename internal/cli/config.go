package cli

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
	"github.com/silvinalucero/skill_cli/internal/git"
	"github.com/spf13/cobra"
)

// getConfigHelp returns colored help text for config command
func getConfigHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s                  %s
                        %s

  %s        %s
                        %s

  %s                %s
                        %s

%s
  %s
  %s
  %s`,
		"Commands to view and modify skill CLI configuration.",

		cyan("AVAILABLE SUBCOMMANDS:"),
		yellow("show"), dim("Show current configuration"),
		dim("Shows all repositories, active repo, and status"),

		yellow("set-repo <url>"), dim("Change configured repository"),
		dim("Flags: --no-verify"),

		yellow("verify"), dim("Verify repository access"),
		dim("Tests connectivity to the configured repository"),

		cyan("EXAMPLES:"),
		green("skill config show"),
		green("skill config set-repo https://github.com/org/new-repo.git"),
		green("skill config verify"),
	)
}

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manages skill CLI configuration",
	Long:  getConfigHelp(),
}

// getConfigShowHelp returns colored help text
func getConfigShowHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	return fmt.Sprintf(`%s

%s
  - Configuration file path
  - Configuration version
  - Active repository name
  - All configured repositories with:
    * Repository name
    * URL (HTTPS or SSH)
    * Local path
    * Authentication type
    * Last verification time
    * Repository status (clean, pending changes, not cloned)
    * Active indicator (*)

%s
  %s`,
		"Shows current skill CLI configuration.",

		cyan("OUTPUT INCLUDES:"),

		cyan("EXAMPLE:"),
		green("skill config show"),
	)
}

// configShowCmd shows current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  getConfigShowHelp(),
	RunE:  runConfigShow,
}

// getConfigSetRepoHelp returns colored help text
func getConfigSetRepoHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s    %s

%s
  %s         %s

%s
  - Validates URL format
  - Verifies repository access (unless --no-verify)
  - Updates configuration with new URL
  - Does not automatically update local repository
  - Run 'skill update' to sync (coming soon)

%s
  %s
  %s`,
		"Changes the configured skills repository.",

		cyan("PARAMETERS:"),
		yellow("<repository-url>"), dim("New Git repository URL - HTTPS or SSH format (required)"),

		cyan("FLAGS:"),
		green("--no-verify"), dim("Skip repository verification before saving"),

		cyan("BEHAVIOR:"),

		cyan("EXAMPLES:"),
		green("skill config set-repo https://github.com/org/new-repo.git"),
		green("skill config set-repo git@github.com:org/new-repo.git --no-verify"),
	)
}

// configSetRepoCmd changes the configured repository
var configSetRepoCmd = &cobra.Command{
	Use:   "set-repo <repository-url> [flags]",
	Short: "Change skills repository",
	Long:  getConfigSetRepoHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigSetRepo,
}

// getConfigVerifyHelp returns colored help text
func getConfigVerifyHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	return fmt.Sprintf(`%s

%s
  - Tests connectivity to the configured repository
  - Checks authentication (SSH keys or HTTPS credentials)
  - Updates last verification timestamp on success
  - Provides troubleshooting help on failure

%s
  - For HTTPS: Tests repository read access
  - For SSH: Validates SSH key authentication

%s
  %s`,
		"Verifies that the configured repository is accessible.",

		cyan("BEHAVIOR:"),

		cyan("VERIFICATION METHODS:"),

		cyan("EXAMPLES:"),
		green("skill config verify"),
	)
}

// configVerifyCmd verifies repository access
var configVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify configured repository access",
	Long:  getConfigVerifyHelp(),
	RunE:  runConfigVerify,
}

var (
	noVerify bool
)

func init() {
	// Add subcommands to config
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetRepoCmd)
	configCmd.AddCommand(configVerifyCmd)

	// Flags for set-repo
	configSetRepoCmd.Flags().BoolVar(&noVerify, "no-verify", false, "don't verify repository before saving")
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	// Colors for output
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check if configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skill init <name> <repo-url>' to initialize")
		return errors.ErrConfigNotFound
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	// Get configuration path
	configPath, _ := config.GetConfigPath()

	// Show configuration
	fmt.Printf("%s %s\n\n", cyan("Configuration:"), configPath)
	fmt.Printf("Version: %s\n", cfg.Version)
	fmt.Printf("Active repository: %s\n\n", green(cfg.ActiveRepo))

	// Show all repositories
	fmt.Printf("%s\n", cyan("Repositories:"))

	if len(cfg.Repositories) == 0 {
		fmt.Printf("  %s No repositories configured\n", yellow("⚠"))
		return nil
	}

	for name, repo := range cfg.Repositories {
		// Active repository indicator
		activeMarker := "  "
		if name == cfg.ActiveRepo {
			activeMarker = green("* ")
		}

		fmt.Printf("\n%s%s\n", activeMarker, cyan(name))
		fmt.Printf("    URL:            %s\n", repo.URL)
		fmt.Printf("    Local Path:     %s\n", repo.LocalPath)
		fmt.Printf("    Auth Type:      %s\n", repo.AuthType)
		fmt.Printf("    Last Verified:  %s\n", repo.LastVerified.Format("2006-01-02 15:04:05"))

		// Check local repository status
		if git.RepoExists(repo.LocalPath) {
			isClean, err := git.IsClean(repo.LocalPath)
			if err != nil {
				fmt.Printf("    Status:         %s\n", yellow("Could not verify"))
			} else if isClean {
				fmt.Printf("    Status:         %s\n", green("Clean"))
			} else {
				fmt.Printf("    Status:         %s\n", yellow("Pending changes"))
			}
		} else {
			fmt.Printf("    Status:         %s\n", red("Not cloned"))
		}
	}

	fmt.Println()

	return nil
}

func runConfigSetRepo(cmd *cobra.Command, args []string) error {
	repoURL := args[0]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check if configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skill init <repo-url>' to initialize first")
		return errors.ErrConfigNotFound
	}

	// Validate URL format
	fmt.Printf("Validating repository URL...\n")
	if err := git.ValidateURL(repoURL); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Valid URL\n", green("✓"))

	// Verify access (unless --no-verify is used)
	if !noVerify {
		fmt.Printf("Verifying repository access...\n")
		if err := git.VerifyBasic(repoURL); err != nil {
			fmt.Printf("%s %v\n", red("✗"), err)
			return err
		}
		fmt.Printf("%s Repository accessible\n", green("✓"))
	}

	// Load existing configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	// Update URL and auth type
	cfg.Repository.URL = repoURL
	cfg.Repository.AuthType = string(git.DetectAuthType(repoURL))
	cfg.Repository.LastVerified = time.Now()

	// Save configuration
	fmt.Printf("Updating configuration...\n")
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Error saving configuration: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Repository updated to: %s\n", green("✓"), repoURL)
	fmt.Println("\nNote: Local repository was not automatically updated.")
	fmt.Println("      Run 'skill update' to sync with the new repository (coming soon)")

	return nil
}

func runConfigVerify(cmd *cobra.Command, args []string) error {
	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check if configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skill init <repo-url>' to initialize first")
		return errors.ErrConfigNotFound
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	// Verify repository access
	fmt.Printf("Verifying access to: %s\n", cfg.Repository.URL)
	if err := git.VerifyBasic(cfg.Repository.URL); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)

		// Help messages
		authType := git.DetectAuthType(cfg.Repository.URL)
		if errors.IsAuthenticationFailed(err) {
			fmt.Println("\nHelp:")
			if authType == config.AuthSSH {
				fmt.Println("  - Verify that your SSH key is configured")
				fmt.Println("  - Run: ssh -T git@github.com (for GitHub)")
			} else {
				fmt.Println("  - Verify your access credentials")
				fmt.Println("  - For private repos, configure an access token")
			}
		}

		return err
	}

	// Update last verification
	cfg.Repository.LastVerified = time.Now()
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Warning: could not update timestamp: %v\n", red("⚠"), err)
	}

	fmt.Printf("%s Repository is accessible\n", green("✓"))
	fmt.Printf("Last verification: %s\n", cfg.Repository.LastVerified.Format("2006-01-02 15:04:05"))

	return nil
}

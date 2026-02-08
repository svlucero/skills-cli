package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
	"github.com/silvinalucero/skill_cli/internal/git"
	"github.com/spf13/cobra"
)

var (
	// Flags for repository commands
	forceRepo  bool
	skipVerify bool
	setCurrent bool
	keepLocal  bool
)

// getRepositoryHelp returns colored help text for repository command
func getRepositoryHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s    %s
                      %s

  %s       %s
                      %s

  %s                %s

  %s  %s

%s
  %s
  %s
  %s
  %s`,
		"Commands to manage skill repositories.",

		cyan("AVAILABLE SUBCOMMANDS:"),
		yellow("add <name> <url>"), dim("Add a new repository"),
		dim("Flags: --force, --skip-verify, --set-current"),

		yellow("remove <name>"), dim("Remove a repository"),
		dim("Flags: --keep-local"),

		yellow("list"), dim("List all configured repositories"),

		yellow("set-current <name>"), dim("Set the current active repository"),

		cyan("EXAMPLES:"),
		green("skill repository add myrepo https://github.com/org/skills.git"),
		green("skill repository list"),
		green("skill repository set-current myrepo"),
		green("skill repository remove oldrepo"),
	)
}

// repositoryCmd represents the repository command group
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Manage skill repositories",
	Long:  getRepositoryHelp(),
}

// getRepositoryAddHelp returns colored help text
func getRepositoryAddHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s              %s
  %s    %s

%s
  %s         %s
  %s       %s
  %s       %s

%s
  - Validates repository URL format
  - Verifies repository access (unless --skip-verify)
  - Clones to ~/.local/share/skill/repos/<name>
  - Saves configuration to ~/.config/skill/config.yaml
  - First repository added automatically becomes active

%s
  %s
  %s
  %s
  %s`,
		"Add a new Git repository where skills are stored.",

		cyan("PARAMETERS:"),
		yellow("<name>"), dim("Name to identify this repository (required)"),
		yellow("<repository-url>"), dim("Git repository URL - HTTPS or SSH format (required)"),

		cyan("FLAGS:"),
		green("-f, --force"), dim("Overwrite existing repository with same name"),
		green("--skip-verify"), dim("Skip repository verification (not recommended)"),
		green("--set-current"), dim("Set this repository as the active one"),

		cyan("BEHAVIOR:"),

		cyan("EXAMPLES:"),
		green("skill repository add myrepo https://github.com/org/skills-repo.git"),
		green("skill repository add company git@github.com:company/internal-skills.git"),
		green("skill repository add myrepo https://github.com/org/skills-repo.git --force"),
		green("skill repository add secondary https://github.com/org/another.git --set-current"),
	)
}

// repositoryAddCmd adds a new repository
var repositoryAddCmd = &cobra.Command{
	Use:   "add <name> <repository-url> [flags]",
	Short: "Add a new skill repository",
	Long:  getRepositoryAddHelp(),
	Args:  cobra.ExactArgs(2),
	RunE:  runRepositoryAdd,
}

// getRepositoryRemoveHelp returns colored help text
func getRepositoryRemoveHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s              %s

%s
  %s        %s

%s
  - Removes repository from configuration
  - Deletes local repository files (unless --keep-local)
  - Cannot remove the active repository (switch first)

%s
  %s   %s
  %s   %s`,
		"Remove a repository from configuration.",

		cyan("PARAMETERS:"),
		yellow("<name>"), dim("Name of the repository to remove (required)"),

		cyan("FLAGS:"),
		green("--keep-local"), dim("Keep local repository files, only remove from config"),

		cyan("BEHAVIOR:"),

		cyan("EXAMPLES:"),
		green("skill repository remove oldrepo"), dim("# Remove repo and local files"),
		green("skill repository remove oldrepo --keep-local"), dim("# Remove from config only"),
	)
}

// repositoryRemoveCmd removes a repository
var repositoryRemoveCmd = &cobra.Command{
	Use:   "remove <name> [flags]",
	Short: "Remove a repository",
	Long:  getRepositoryRemoveHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runRepositoryRemove,
}

// getRepositoryListHelp returns colored help text
func getRepositoryListHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	return fmt.Sprintf(`%s

%s
  - Repository name
  - URL (HTTPS or SSH)
  - Local path
  - Authentication type
  - Last verification time
  - Clone status
  - Active repository indicator (*)

%s
  %s`,
		"Shows all configured repositories with basic information.",

		cyan("OUTPUT INCLUDES:"),

		cyan("EXAMPLE:"),
		green("skill repository list"),
	)
}

// repositoryListCmd lists all repositories
var repositoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured repositories",
	Long:  getRepositoryListHelp(),
	RunE:  runRepositoryList,
}

// getRepositorySetCurrentHelp returns colored help text
func getRepositorySetCurrentHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s              %s

%s
  - Sets the specified repository as the active one
  - Commands like 'skill list' and 'skill install' use the active repo by default
  - The repository must already be added to configuration

%s
  %s
  %s`,
		"Change the active repository.",

		cyan("PARAMETERS:"),
		yellow("<name>"), dim("Name of the repository to set as active (required)"),

		cyan("BEHAVIOR:"),

		cyan("EXAMPLES:"),
		green("skill repository set-current myrepo"),
		green("skill repository set-current company"),
	)
}

// repositorySetCurrentCmd sets the active repository
var repositorySetCurrentCmd = &cobra.Command{
	Use:   "set-current <name>",
	Short: "Set the current active repository",
	Long:  getRepositorySetCurrentHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runRepositorySetCurrent,
}

func init() {
	// Add subcommands to repository
	repositoryCmd.AddCommand(repositoryAddCmd)
	repositoryCmd.AddCommand(repositoryRemoveCmd)
	repositoryCmd.AddCommand(repositoryListCmd)
	repositoryCmd.AddCommand(repositorySetCurrentCmd)

	// Flags for add
	repositoryAddCmd.Flags().BoolVarP(&forceRepo, "force", "f", false, "overwrite existing repository with same name")
	repositoryAddCmd.Flags().BoolVar(&skipVerify, "skip-verify", false, "skip repository verification (not recommended)")
	repositoryAddCmd.Flags().BoolVar(&setCurrent, "set-current", false, "set this repository as active")

	// Flags for remove
	repositoryRemoveCmd.Flags().BoolVar(&keepLocal, "keep-local", false, "keep local repository (only remove from config)")
}

func runRepositoryAdd(cmd *cobra.Command, args []string) error {
	repoName := args[0]
	repoURL := args[1]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// 1. Validate repository name
	if err := config.ValidateRepoName(repoName); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return err
	}

	// 2. Load or create configuration
	var cfg *config.Config
	var isNewConfig bool

	if config.Exists() {
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
			return err
		}
		isNewConfig = false

		// Check if repository already exists
		if _, err := config.GetRepo(cfg, repoName); err == nil && !forceRepo {
			fmt.Printf("%s Repository '%s' already exists\n", red("✗"), repoName)
			fmt.Println("  Use --force to overwrite")
			return errors.ErrRepoAlreadyExists
		}
	} else {
		cfg = &config.Config{
			Version:      "2",
			Repositories: make(map[string]config.Repository),
		}
		isNewConfig = true
	}

	// 3. Validate URL format
	fmt.Printf("Validating repository URL...\n")
	if err := git.ValidateURL(repoURL); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		fmt.Println("\nValid formats:")
		fmt.Println("  - https://github.com/org/repo.git")
		fmt.Println("  - git@github.com:org/repo.git")
		fmt.Println("  - ssh://git@github.com/org/repo.git")
		return err
	}
	fmt.Printf("%s Valid URL\n", green("✓"))

	// 4. Detect authentication type
	authType := git.DetectAuthType(repoURL)
	fmt.Printf("Detected authentication type: %s\n", authType)

	// 5. Verify repository access (unless --skip-verify)
	if !skipVerify {
		fmt.Printf("Verifying repository access...\n")
		if err := git.VerifyBasic(repoURL); err != nil {
			fmt.Printf("%s %v\n", red("✗"), err)

			// Help messages based on error
			if errors.IsAuthenticationFailed(err) {
				fmt.Println("\nHelp:")
				if authType == config.AuthSSH {
					fmt.Println("  - Verify your SSH key is configured")
					fmt.Println("  - Run: ssh -T git@github.com (for GitHub)")
				} else {
					fmt.Println("  - Verify your access credentials")
					fmt.Println("  - For private repos, configure an access token")
				}
			} else if errors.IsNetworkError(err) {
				fmt.Println("\nHelp:")
				fmt.Println("  - Check your internet connection")
				fmt.Println("  - Verify you can reach the repository host")
			} else if errors.IsRepositoryNotFound(err) {
				fmt.Println("\nHelp:")
				fmt.Println("  - Verify the repository URL is correct")
				fmt.Println("  - Verify you have read permissions on the repository")
			}

			return err
		}
		fmt.Printf("%s Repository accessible\n", green("✓"))
	} else {
		fmt.Printf("%s Verification skipped (--skip-verify)\n", yellow("⚠"))
	}

	// 6. Get local path for repository
	repoPath := config.GetRepoPathForRepo(repoName)

	// 7. If local repo exists and using --force, remove it
	if git.RepoExists(repoPath) && forceRepo {
		fmt.Printf("Removing previous local repository...\n")
		if err := os.RemoveAll(repoPath); err != nil {
			fmt.Printf("%s Error removing previous repository: %v\n", red("✗"), err)
			return err
		}
	}

	// 8. Clone the repository
	if !git.RepoExists(repoPath) {
		fmt.Printf("Cloning repository to %s...\n", repoPath)
		if err := git.Clone(repoURL, repoPath); err != nil {
			fmt.Printf("%s Error cloning repository: %v\n", red("✗"), err)
			return err
		}
		fmt.Printf("%s Repository cloned successfully\n", green("✓"))
	} else {
		fmt.Printf("%s Local repository already exists, using existing\n", green("✓"))
	}

	// 9. Create or update repository in configuration
	repo := config.Repository{
		Name:         repoName,
		URL:          repoURL,
		LocalPath:    repoPath,
		LastVerified: time.Now(),
		AuthType:     string(authType),
	}

	// If repo exists (--force), remove it first
	if _, err := config.GetRepo(cfg, repoName); err == nil {
		delete(cfg.Repositories, repoName)
	}

	// Add the repository
	if err := config.AddRepository(cfg, repo); err != nil {
		fmt.Printf("%s Error adding repository: %v\n", red("✗"), err)
		return err
	}

	// If --set-current, change the active one
	if setCurrent {
		if err := config.SetActiveRepo(cfg, repoName); err != nil {
			fmt.Printf("%s Error setting active repository: %v\n", red("✗"), err)
			return err
		}
	}

	// 10. Validate configuration
	if err := config.Validate(cfg); err != nil {
		fmt.Printf("%s Invalid configuration: %v\n", red("✗"), err)
		return err
	}

	// 11. Save configuration
	fmt.Printf("Saving configuration...\n")
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Error saving configuration: %v\n", red("✗"), err)
		return err
	}

	configPath, _ := config.GetConfigPath()
	fmt.Printf("%s Configuration saved to: %s\n", green("✓"), configPath)

	// 12. Success message
	fmt.Println()
	if isNewConfig {
		fmt.Printf("%s Repository initialized successfully!\n\n", green("✓"))
	} else {
		fmt.Printf("%s Repository added successfully!\n\n", green("✓"))
	}

	fmt.Printf("Repository: %s\n", repoName)
	if cfg.ActiveRepo == repoName {
		fmt.Printf("Status: %s\n", green("active"))
	} else {
		fmt.Printf("Status: available\n")
		fmt.Printf("To activate: skill repository set-current %s\n", repoName)
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  - Run 'skill config show' to see configuration")
	fmt.Println("  - Run 'skill list' to see available skills")
	fmt.Println("  - Run 'skill repository list' to see all repositories")

	return nil
}

func runRepositoryList(cmd *cobra.Command, args []string) error {
	// Colors for output
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skill repository add <name> <repo-url>' to initialize")
		return errors.ErrConfigNotFound
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	if len(cfg.Repositories) == 0 {
		fmt.Printf("%s No repositories configured\n", yellow("⚠"))
		return nil
	}

	fmt.Printf("%s\n", cyan("Configured repositories:"))
	fmt.Println()

	for name, repo := range cfg.Repositories {
		// Active repository indicator
		activeMarker := "  "
		activeSuffix := ""
		if name == cfg.ActiveRepo {
			activeMarker = green("* ")
			activeSuffix = green(" (current)")
		}

		fmt.Printf("%s%s%s\n", activeMarker, name, activeSuffix)
		fmt.Printf("    URL:   %s\n", repo.URL)
		fmt.Printf("    Auth:  %s\n", repo.AuthType)

		// Local repo status
		if git.RepoExists(repo.LocalPath) {
			fmt.Printf("    Local: %s\n", green("✓ cloned"))
		} else {
			fmt.Printf("    Local: %s\n", red("✗ not cloned"))
		}

		fmt.Println()
	}

	return nil
}

func runRepositoryRemove(cmd *cobra.Command, args []string) error {
	repoName := args[0]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		return errors.ErrConfigNotFound
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	// Get repository before removing (to have the path)
	repo, err := config.GetRepo(cfg, repoName)
	if err != nil {
		fmt.Printf("%s Repository '%s' not found\n", red("✗"), repoName)
		return errors.ErrRepoNotFound
	}

	// Try to remove repository
	if err := config.RemoveRepository(cfg, repoName); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		if cfg.ActiveRepo == repoName {
			fmt.Println("\nTo remove the active repository:")
			fmt.Println("  1. Switch to another repository: skill repository set-current <other-name>")
			fmt.Printf("  2. Remove this repository: skill repository remove %s\n", repoName)
		}
		return err
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Error saving configuration: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Repository '%s' removed from configuration\n", green("✓"), repoName)

	// Remove local repository (unless --keep-local)
	if !keepLocal && git.RepoExists(repo.LocalPath) {
		fmt.Printf("Removing local repository...\n")
		if err := os.RemoveAll(repo.LocalPath); err != nil {
			fmt.Printf("%s Warning: could not remove local directory: %v\n", red("⚠"), err)
		} else {
			fmt.Printf("%s Local repository removed\n", green("✓"))
		}
	} else if keepLocal {
		fmt.Printf("Local repository kept at: %s\n", repo.LocalPath)
	}

	return nil
}

func runRepositorySetCurrent(cmd *cobra.Command, args []string) error {
	repoName := args[0]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skill repository add <name> <repo-url>' to initialize")
		return errors.ErrConfigNotFound
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
		return err
	}

	// Verify repository exists
	if _, err := config.GetRepo(cfg, repoName); err != nil {
		fmt.Printf("%s Repository '%s' not found\n", red("✗"), repoName)
		fmt.Println("\nAvailable repositories:")
		for name := range cfg.Repositories {
			fmt.Printf("  - %s\n", name)
		}
		return errors.ErrRepoNotFound
	}

	// Already active?
	if cfg.ActiveRepo == repoName {
		fmt.Printf("%s '%s' is already the active repository\n", green("✓"), repoName)
		return nil
	}

	// Change the active one
	if err := config.SetActiveRepo(cfg, repoName); err != nil {
		fmt.Printf("%s Error changing active repository: %v\n", red("✗"), err)
		return err
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Error saving configuration: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Active repository changed to '%s'\n", green("✓"), repoName)

	return nil
}

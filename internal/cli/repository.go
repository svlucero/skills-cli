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

// formatError formats an error message with red color
// Usage: fmt.Println(formatError("Error: %s", err))
func formatError(format string, args ...interface{}) string {
	red := color.New(color.FgRed).SprintFunc()
	msg := fmt.Sprintf(format, args...)
	return red(msg)
}

// formatSuccess formats a success message with green color
func formatSuccess(format string, args ...interface{}) string {
	green := color.New(color.FgGreen).SprintFunc()
	msg := fmt.Sprintf(format, args...)
	return green(msg)
}

// formatWarning formats a warning message with yellow color
func formatWarning(format string, args ...interface{}) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	msg := fmt.Sprintf(format, args...)
	return yellow(msg)
}

// getRepositoryHelp returns colored help text for repository command
func getRepositoryHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return NewHelpBuilder().
		Description("Commands to manage skill repositories.").
		Section("AVAILABLE SUBCOMMANDS:").
		Item(yellow("add <name> <url>"), "Add a new repository").
		Text("                      "+dim("Flags: --force, --skip-verify, --set-current, --skills-path")).
		EmptyLine().
		Item(yellow("update <name>"), "Update repository configuration").
		Text("                      "+dim("Flags: --skills-path")).
		EmptyLine().
		Item(yellow("remove <name>"), "Remove a repository").
		Text("                      "+dim("Flags: --keep-local")).
		EmptyLine().
		Item(yellow("list"), "List all configured repositories").
		EmptyLine().
		Item(yellow("set-current <name>"), "Set the current active repository").
		Section("EXAMPLES:").
		Example("skill repository add myrepo https://github.com/org/skills.git", "").
		Example("skill repository update myrepo --skills-path skills", "").
		Example("skill repository list", "").
		Example("skill repository set-current myrepo", "").
		Example("skill repository remove oldrepo", "").
		Build()
}

// repositoryCmd represents the repository command group
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Manage skill repositories",
	Long:  getRepositoryHelp(),
}

// getRepositoryAddHelp returns colored help text
func getRepositoryAddHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()

	return NewHelpBuilder().
		Description("Add a new Git repository where skills are stored.").
		Section("PARAMETERS:").
		Item(yellow("<name>"), "Name to identify this repository (required)").
		Item(yellow("<repository-url>"), "Git repository URL - HTTPS or SSH format (required)").
		Section("FLAGS:").
		Item("-f, --force", "Overwrite existing repository or allow duplicate URLs").
		Item("--skip-verify", "Skip repository verification (not recommended)").
		Item("--set-current", "Set this repository as the active one").
		Section("BEHAVIOR:").
		BulletList([]string{
			"Validates repository name and URL format",
			"Checks for duplicate repository URLs (prevents same repo with different names)",
			"Verifies repository access (unless --skip-verify)",
			"Clones to ~/.local/share/skill/repos/<name>",
			"Saves configuration to ~/.config/skill/config.yaml",
			"First repository added automatically becomes active",
		}).
		Section("EXAMPLES:").
		Example("skill repository add myrepo https://github.com/org/skills-repo.git", "").
		Example("skill repository add company git@github.com:company/internal-skills.git", "").
		Example("skill repository add myrepo https://github.com/org/skills-repo.git --force", "").
		Example("skill repository add secondary https://github.com/org/another.git --set-current", "").
		Build()
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
	yellow := color.New(color.FgYellow).SprintFunc()

	return NewHelpBuilder().
		Description("Remove a repository from configuration.").
		Section("PARAMETERS:").
		Item(yellow("<name>"), "Name of the repository to remove (required)").
		Section("FLAGS:").
		Item("--keep-local", "Keep local repository files, only remove from config").
		Section("BEHAVIOR:").
		BulletList([]string{
			"Removes repository from configuration",
			"Deletes local repository files (unless --keep-local)",
			"Cannot remove the active repository (switch first)",
		}).
		Section("EXAMPLES:").
		Example("skill repository remove oldrepo", "# Remove repo and local files").
		Example("skill repository remove oldrepo --keep-local", "# Remove from config only").
		Build()
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
	return NewHelpBuilder().
		Description("Shows all configured repositories with basic information.").
		Section("OUTPUT INCLUDES:").
		BulletList([]string{
			"Repository name",
			"URL (HTTPS or SSH)",
			"Local path",
			"Authentication type",
			"Last verification time",
			"Clone status",
			"Active repository indicator (*)",
		}).
		Section("EXAMPLE:").
		Example("skill repository list", "").
		Build()
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
	yellow := color.New(color.FgYellow).SprintFunc()

	return NewHelpBuilder().
		Description("Change the active repository.").
		Section("PARAMETERS:").
		Item(yellow("<name>"), "Name of the repository to set as active (required)").
		Section("BEHAVIOR:").
		BulletList([]string{
			"Sets the specified repository as the active one",
			"Commands like 'skill list' and 'skill install' use the active repo by default",
			"The repository must already be added to configuration",
		}).
		Section("EXAMPLES:").
		Example("skill repository set-current myrepo", "").
		Example("skill repository set-current company", "").
		Build()
}

// repositorySetCurrentCmd sets the active repository
var repositorySetCurrentCmd = &cobra.Command{
	Use:   "set-current <name>",
	Short: "Set the current active repository",
	Long:  getRepositorySetCurrentHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runRepositorySetCurrent,
}

// getRepositoryUpdateHelp returns colored help text
func getRepositoryUpdateHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()

	return NewHelpBuilder().
		Description("Update configuration settings for an existing repository.").
		Section("PARAMETERS:").
		Item(yellow("<name>"), "Name of the repository to update (required)").
		Section("FLAGS:").
		Item("--skills-path <path>", "Update the relative path to skills directory (use '/' or empty for root)").
		Section("BEHAVIOR:").
		BulletList([]string{
			"Updates repository configuration in ~/.config/skill/config.yaml",
			"Validates skills path exists in the local repository",
			"Use --skills-path '' or --skills-path / to reset to root",
		}).
		Section("EXAMPLES:").
		Example("skill repository update myrepo --skills-path skills", "# Set skills path to 'skills' directory").
		Example("skill repository update myrepo --skills-path examples/claude-skills", "# Set nested path").
		Example("skill repository update myrepo --skills-path /", "# Reset to root directory").
		Build()
}

// repositoryUpdateCmd updates a repository configuration
var repositoryUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update repository configuration",
	Long:  getRepositoryUpdateHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runRepositoryUpdate,
}

func init() {
	// Add subcommands to repository
	repositoryCmd.AddCommand(repositoryAddCmd)
	repositoryCmd.AddCommand(repositoryRemoveCmd)
	repositoryCmd.AddCommand(repositoryListCmd)
	repositoryCmd.AddCommand(repositorySetCurrentCmd)
	repositoryCmd.AddCommand(repositoryUpdateCmd)

	// Flags for add
	repositoryAddCmd.Flags().BoolVarP(&forceRepo, "force", "f", false, "overwrite existing repository with same name")
	repositoryAddCmd.Flags().BoolVar(&skipVerify, "skip-verify", false, "skip repository verification (not recommended)")
	repositoryAddCmd.Flags().BoolVar(&setCurrent, "set-current", false, "set this repository as active")

	// Flags for remove
	repositoryRemoveCmd.Flags().BoolVar(&keepLocal, "keep-local", false, "keep local repository (only remove from config)")

	// Flags for update
	repositoryUpdateCmd.Flags().StringVar(&skillsPath, "skills-path", "", "relative path to skills directory within repo (use '/' or empty for root)")
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
		fmt.Printf("%s %s\n", red("✗"), formatError("%v", err))
		return err
	}

	// 2. Load or create configuration
	var cfg *config.Config
	var isNewConfig bool

	if config.Exists() {
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Printf("%s %s\n", red("✗"), formatError("Error loading configuration: %v", err))
			return err
		}
		isNewConfig = false

		// Check if repository name already exists
		if _, err := config.GetRepo(cfg, repoName); err == nil && !forceRepo {
			fmt.Printf("%s %s\n", red("✗"), formatError("Repository '%s' already exists", repoName))
			fmt.Println("  Use --force to overwrite")
			return errors.ErrRepoAlreadyExists
		}

		// Check if repository URL already exists (duplicate validation)
		if existingRepo, found := config.FindRepoByURL(cfg, repoURL); found && !forceRepo {
			fmt.Printf("%s %s\n", red("✗"), formatError("Repository URL already exists"))
			fmt.Printf("  %s\n", formatError("URL: %s", repoURL))
			fmt.Printf("  %s\n", formatError("Existing name: %s", existingRepo.Name))
			fmt.Println()
			fmt.Println(formatError("This repository is already configured. You can:"))
			fmt.Printf("  - Use the existing repository: skill repository set-current %s\n", existingRepo.Name)
			fmt.Printf("  - Remove and re-add: skill repository remove %s\n", existingRepo.Name)
			fmt.Println("  - Use --force to add anyway (creates duplicate)")
			return fmt.Errorf("duplicate repository URL")
		}

		// Warning if forcing duplicate URL
		if existingRepo, found := config.FindRepoByURL(cfg, repoURL); found && forceRepo {
			fmt.Printf("%s %s\n", yellow("⚠"), formatWarning("Adding duplicate repository URL (--force)"))
			fmt.Printf("  %s\n", formatWarning("Existing repository '%s' uses the same URL", existingRepo.Name))
			fmt.Println()
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
		fmt.Printf("%s %s\n", red("✗"), formatError("%v", err))
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
			fmt.Printf("%s %s\n", red("✗"), formatError("%v", err))

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
		fmt.Printf("%s %s\n", yellow("⚠"), formatWarning("Verification skipped (--skip-verify)"))
	}

	// 6. Get local path for repository
	repoPath := config.GetRepoPathForRepo(repoName)

	// 7. If local repo exists and using --force, remove it
	if git.RepoExists(repoPath) && forceRepo {
		fmt.Printf("Removing previous local repository...\n")
		if err := os.RemoveAll(repoPath); err != nil {
			fmt.Printf("%s %s\n", red("✗"), formatError("Error removing previous repository: %v", err))
			return err
		}
	}

	// 8. Clone the repository
	if !git.RepoExists(repoPath) {
		fmt.Printf("Cloning repository to %s...\n", repoPath)
		if err := git.Clone(repoURL, repoPath); err != nil {
			fmt.Printf("%s %s\n", red("✗"), formatError("Error cloning repository: %v", err))
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
		fmt.Printf("%s %s\n", red("✗"), formatError("Error adding repository: %v", err))
		return err
	}

	// If --set-current, change the active one
	if setCurrent {
		if err := config.SetActiveRepo(cfg, repoName); err != nil {
			fmt.Printf("%s %s\n", red("✗"), formatError("Error setting active repository: %v", err))
			return err
		}
	}

	// 10. Validate configuration
	if err := config.Validate(cfg); err != nil {
		fmt.Printf("%s %s\n", red("✗"), formatError("Invalid configuration: %v", err))
		return err
	}

	// 11. Save configuration
	fmt.Printf("Saving configuration...\n")
	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s %s\n", red("✗"), formatError("Error saving configuration: %v", err))
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

func runRepositoryUpdate(cmd *cobra.Command, args []string) error {
	repoName := args[0]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Check if at least one flag was provided
	if !cmd.Flags().Changed("skills-path") {
		fmt.Printf("%s No changes specified\n", red("✗"))
		fmt.Println("  Use --skills-path to update the skills directory path")
		fmt.Println("\nExample:")
		fmt.Printf("  skill repository update %s --skills-path skills\n", repoName)
		return fmt.Errorf("no update flags provided")
	}

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

	// Get repository
	repo, err := config.GetRepo(cfg, repoName)
	if err != nil {
		fmt.Printf("%s Repository '%s' not found\n", red("✗"), repoName)
		fmt.Println("\nAvailable repositories:")
		for name := range cfg.Repositories {
			fmt.Printf("  - %s\n", name)
		}
		return errors.ErrRepoNotFound
	}

	// Check if repository is cloned locally
	if _, err := os.Stat(repo.LocalPath); os.IsNotExist(err) {
		fmt.Printf("%s Repository '%s' not cloned locally\n", red("✗"), repoName)
		fmt.Printf("  Run 'skill repository add %s %s' to clone it\n", repoName, repo.URL)
		return fmt.Errorf("repository not cloned")
	}

	// Process --skills-path flag
	if cmd.Flags().Changed("skills-path") {
		fmt.Printf("Updating skills path...\n")

		// Normalize path
		newSkillsPath := skillsPath

		// Handle "/" or empty as root
		if newSkillsPath == "/" {
			newSkillsPath = ""
		}

		// If not empty, validate
		if newSkillsPath != "" {
			// Normalize path (convert backslashes, clean up)
			newSkillsPath = filepath.Clean(newSkillsPath)

			// Ensure it's a relative path
			if filepath.IsAbs(newSkillsPath) {
				fmt.Printf("%s %s\n", red("✗"), formatError("Skills path must be relative, not absolute: %s", newSkillsPath))
				return fmt.Errorf("skills path must be relative")
			}

			// Check if the path exists within the cloned repo
			fullSkillsPath := filepath.Join(repo.LocalPath, newSkillsPath)
			if _, err := os.Stat(fullSkillsPath); os.IsNotExist(err) {
				fmt.Printf("%s %s\n", red("✗"), formatError("Skills path not found in repository: %s", newSkillsPath))
				fmt.Printf("  Full path: %s\n", fullSkillsPath)
				fmt.Println("\nThe specified skills path does not exist in the repository.")
				fmt.Println("Available directories:")

				// Show available directories as hint
				entries, _ := os.ReadDir(repo.LocalPath)
				for _, entry := range entries {
					if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
						fmt.Printf("  - %s\n", entry.Name())
					}
				}

				return fmt.Errorf("skills path not found")
			}

			// Check if path is a directory
			info, err := os.Stat(fullSkillsPath)
			if err != nil {
				fmt.Printf("%s %s\n", red("✗"), formatError("Error accessing skills path: %v", err))
				return err
			}
			if !info.IsDir() {
				fmt.Printf("%s %s\n", red("✗"), formatError("Skills path is not a directory: %s", newSkillsPath))
				return fmt.Errorf("skills path must be a directory")
			}
		}

		// Update the repository
		repo.SkillsPath = newSkillsPath
		cfg.Repositories[repoName] = *repo

		displayPath := newSkillsPath
		if displayPath == "" {
			displayPath = "/"
		}
		fmt.Printf("%s Skills path updated: %s\n", green("✓"), displayPath)
	}

	// Validate and save configuration
	if err := config.Validate(cfg); err != nil {
		fmt.Printf("%s Invalid configuration: %v\n", red("✗"), err)
		return err
	}

	if err := config.Save(cfg); err != nil {
		fmt.Printf("%s Error saving configuration: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Configuration saved\n", green("✓"))
	fmt.Println()
	fmt.Printf("Repository '%s' updated successfully\n", repoName)

	return nil
}

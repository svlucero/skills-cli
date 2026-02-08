package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	sharePath string
	shareRepo string
)

// getShareHelp returns colored help text for share command
func getShareHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()

	return NewHelpBuilder().
		Description("Share a skill by forking a repository and creating a PR.").
		Section("PARAMETERS:").
		Item(yellow("<path>"), "Path to the skill folder to share (required)").
		Section("FLAGS:").
		Item("--path <path>", "Path to the skill folder (required)").
		Item("--repo <url>", "Target repository URL (optional, defaults to current)").
		Section("WORKFLOW:").
		BulletList([]string{
			"1. Validates the skill structure (must have SKILL.md)",
			"2. Forks the target repository using gh CLI",
			"3. Creates a new branch",
			"4. Copies the skill folder to the forked repository",
			"5. Commits and pushes the changes",
			"6. Creates a pull request against main branch",
			"7. Returns the PR URL",
		}).
		Section("SKILL STRUCTURE:").
		Text("  my-skill/").
		Text("  ├── SKILL.md          # Required: instructions + metadata").
		Text("  ├── scripts/          # Optional: executable code").
		Text("  ├── references/       # Optional: documentation").
		Text("  └── assets/           # Optional: templates, resources").
		Section("REQUIREMENTS:").
		BulletList([]string{
			"gh CLI must be installed and authenticated",
			"You must have permissions to fork the target repository",
		}).
		Section("EXAMPLES:").
		Example("skill share --path ./my-skill", "# Share to current repo").
		Example("skill share --path ./my-skill --repo https://github.com/org/skills", "# Share to specific repo").
		Example("skill share --path ~/skills/deploy-app --repo git@github.com:org/skills.git", "").
		Build()
}

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share --path <path> [--repo <url>]",
	Short: "Share a skill by forking a repo and creating a PR",
	Long:  getShareHelp(),
	RunE:  runShare,
}

func init() {
	shareCmd.Flags().StringVar(&sharePath, "path", "", "path to the skill folder (required)")
	shareCmd.Flags().StringVar(&shareRepo, "repo", "", "target repository URL (optional, defaults to current)")
	shareCmd.MarkFlagRequired("path")
}

func runShare(cmd *cobra.Command, args []string) error {
	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Convert path to absolute path
	absPath, err := filepath.Abs(sharePath)
	if err != nil {
		fmt.Printf("%s Error resolving path: %v\n", red("✗"), err)
		return err
	}

	// Step 1: Validate skill structure
	fmt.Printf("%s Validating skill structure...\n", cyan("→"))
	skillName := filepath.Base(absPath)

	valid, warnings, err := validateSkillStructure(absPath)
	if err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return err
	}

	if !valid {
		fmt.Printf("%s Skill structure is invalid\n", red("✗"))
		return fmt.Errorf("invalid skill structure")
	}

	// Show warnings if any
	if len(warnings) > 0 {
		fmt.Printf("%s Warnings:\n", yellow("⚠"))
		for _, warning := range warnings {
			fmt.Printf("  - %s\n", warning)
		}
		fmt.Println()

		// Ask user if they want to continue
		fmt.Print("Do you want to continue anyway? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Printf("%s Share cancelled by user\n", yellow("⚠"))
			return nil
		}
	}

	fmt.Printf("%s Skill structure validated: %s\n", green("✓"), skillName)
	fmt.Println()

	// Step 2: Determine target repository
	var targetRepo string
	if shareRepo != "" {
		targetRepo = shareRepo
	} else {
		// Use current/active repository
		if !config.Exists() {
			fmt.Printf("%s Configuration not found\n", red("✗"))
			fmt.Println("  Run 'skill repository add <name> <repo-url>' to initialize")
			return fmt.Errorf("configuration not found")
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("%s Error loading configuration: %v\n", red("✗"), err)
			return err
		}

		if cfg.ActiveRepo == "" {
			fmt.Printf("%s No active repository configured\n", red("✗"))
			return fmt.Errorf("no active repository")
		}

		repo, err := config.GetRepo(cfg, cfg.ActiveRepo)
		if err != nil {
			fmt.Printf("%s Error getting active repository: %v\n", red("✗"), err)
			return err
		}

		targetRepo = repo.URL
	}

	fmt.Printf("%s Target repository: %s\n", cyan("→"), targetRepo)
	fmt.Println()

	// Step 3: Check gh CLI is installed
	fmt.Printf("%s Checking gh CLI...\n", cyan("→"))
	if err := checkGHInstalled(); err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		fmt.Println("  Install gh CLI: https://cli.github.com/")
		return err
	}
	fmt.Printf("%s gh CLI is available\n", green("✓"))
	fmt.Println()

	// Step 4: Check for existing fork or create new one
	fmt.Printf("%s Checking for existing fork or creating new one...\n", cyan("→"))
	forkURL, err := forkRepository(targetRepo)
	if err != nil {
		fmt.Printf("%s Error: %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Fork ready: %s\n", green("✓"), forkURL)
	fmt.Println()

	// Step 5: Clone the fork to a temporary directory
	fmt.Printf("%s Cloning fork...\n", cyan("→"))
	tempDir, err := os.MkdirTemp("", "skill-share-*")
	if err != nil {
		fmt.Printf("%s Error creating temporary directory: %v\n", red("✗"), err)
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := cloneRepository(forkURL, tempDir); err != nil {
		fmt.Printf("%s Error cloning fork: %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Fork cloned to temporary directory\n", green("✓"))
	fmt.Println()

	// Step 6: Create a new branch
	branchName := fmt.Sprintf("add-skill-%s", skillName)
	fmt.Printf("%s Creating branch '%s'...\n", cyan("→"), branchName)
	if err := createBranch(tempDir, branchName); err != nil {
		fmt.Printf("%s Error creating branch: %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Branch created\n", green("✓"))
	fmt.Println()

	// Step 7: Copy skill folder to the repository
	fmt.Printf("%s Copying skill folder...\n", cyan("→"))
	destPath := filepath.Join(tempDir, skillName)
	if err := copyDir(absPath, destPath); err != nil {
		fmt.Printf("%s Error copying skill: %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Skill copied\n", green("✓"))
	fmt.Println()

	// Step 8: Commit and push changes
	fmt.Printf("%s Committing and pushing changes...\n", cyan("→"))
	commitMsg := fmt.Sprintf("Add skill: %s", skillName)
	if err := commitAndPush(tempDir, branchName, commitMsg); err != nil {
		fmt.Printf("%s Error committing/pushing: %v\n", red("✗"), err)
		return err
	}
	fmt.Printf("%s Changes pushed\n", green("✓"))
	fmt.Println()

	// Step 9: Create pull request
	fmt.Printf("%s Creating pull request...\n", cyan("→"))
	prURL, err := createPullRequest(tempDir, targetRepo, branchName, skillName)
	if err != nil {
		fmt.Printf("%s Error creating PR: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Pull request created successfully!\n\n", green("✓"))
	fmt.Printf("  PR URL: %s\n", green(prURL))

	return nil
}

// validateSkillStructure validates that the folder is a valid skill
func validateSkillStructure(path string) (bool, []string, error) {
	var warnings []string

	// Check if path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, fmt.Errorf("path does not exist: %s", path)
		}
		return false, nil, err
	}

	if !info.IsDir() {
		return false, nil, fmt.Errorf("path is not a directory: %s", path)
	}

	// Check for required SKILL.md file
	skillMDPath := filepath.Join(path, "SKILL.md")
	if _, err := os.Stat(skillMDPath); os.IsNotExist(err) {
		return false, nil, fmt.Errorf("SKILL.md is required but not found")
	}

	// Check for optional directories
	optionalDirs := []string{"scripts", "references", "assets"}
	hasOptionalDir := false

	for _, dir := range optionalDirs {
		dirPath := filepath.Join(path, dir)
		if _, err := os.Stat(dirPath); err == nil {
			hasOptionalDir = true
			break
		}
	}

	if !hasOptionalDir {
		warnings = append(warnings, "No optional directories found (scripts/, references/, assets/)")
	}

	return true, warnings, nil
}

// checkGHInstalled checks if gh CLI is installed
func checkGHInstalled() error {
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh CLI is not installed or not in PATH")
	}
	return nil
}

// checkForkExists checks if the user already has a fork of the repository
func checkForkExists(owner, repo string) (bool, string, error) {
	// Get current user
	userCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	userOutput, err := userCmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("error getting current user: %w", err)
	}
	currentUser := strings.TrimSpace(string(userOutput))

	// Check if fork exists using gh API
	// gh api repos/{owner}/{repo}/forks --jq '.[] | select(.owner.login == "username") | .html_url'
	checkCmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/%s/forks", owner, repo),
		"--jq", fmt.Sprintf(".[] | select(.owner.login == \"%s\") | .clone_url", currentUser))

	output, err := checkCmd.Output()
	if err != nil {
		// If error, assume fork doesn't exist
		return false, "", nil
	}

	forkURL := strings.TrimSpace(string(output))
	if forkURL != "" {
		return true, forkURL, nil
	}

	return false, "", nil
}

// checkIsOwner checks if the current user is the owner of the repository
func checkIsOwner(owner, repo string) (bool, error) {
	// Get current user
	userCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	userOutput, err := userCmd.Output()
	if err != nil {
		return false, fmt.Errorf("error getting current user: %w", err)
	}
	currentUser := strings.TrimSpace(string(userOutput))

	// Compare owner with current user
	return strings.EqualFold(owner, currentUser), nil
}

// forkRepository forks a repository and returns the fork URL
func forkRepository(repoURL string) (string, error) {
	// Extract owner/repo from URL
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	// Check if user is already the owner
	isOwner, err := checkIsOwner(owner, repo)
	if err != nil {
		return "", fmt.Errorf("error checking repository ownership: %w", err)
	}

	if isOwner {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("%s You are the owner of this repository, using it directly\n", cyan("ℹ"))
		fmt.Printf("  Repository URL: %s\n", repoURL)
		return repoURL, nil
	}

	// Check if fork already exists
	exists, existingForkURL, err := checkForkExists(owner, repo)
	if err != nil {
		return "", fmt.Errorf("error checking for existing fork: %w", err)
	}

	if exists {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s Fork already exists, using existing fork\n", yellow("ℹ"))
		fmt.Printf("  Fork URL: %s\n", existingForkURL)
		return existingForkURL, nil
	}

	// Fork using gh CLI
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s Creating new fork...\n", green("→"))

	cmd := exec.Command("gh", "repo", "fork", fmt.Sprintf("%s/%s", owner, repo), "--clone=false")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check if already forked (backup check)
		if strings.Contains(string(output), "already exists") {
			userCmd := exec.Command("gh", "api", "user", "--jq", ".login")
			userOutput, userErr := userCmd.Output()
			if userErr != nil {
				return "", fmt.Errorf("error getting current user: %w", userErr)
			}
			currentUser := strings.TrimSpace(string(userOutput))
			return fmt.Sprintf("https://github.com/%s/%s.git", currentUser, repo), nil
		}
		return "", fmt.Errorf("error forking repository: %s", string(output))
	}

	// Extract fork URL from output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "github.com") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Contains(part, "github.com") {
					return part, nil
				}
			}
		}
	}

	// Fallback: construct fork URL
	userCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	userOutput, err := userCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}
	currentUser := strings.TrimSpace(string(userOutput))

	return fmt.Sprintf("https://github.com/%s/%s.git", currentUser, repo), nil
}

// parseGitHubURL extracts owner and repo from a GitHub URL
func parseGitHubURL(url string) (string, string, error) {
	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	// Handle different URL formats
	var parts []string

	if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
		parts = strings.Split(url, "/")
	} else if strings.HasPrefix(url, "git@github.com:") {
		url = strings.TrimPrefix(url, "git@github.com:")
		parts = strings.Split(url, "/")
	} else {
		return "", "", fmt.Errorf("unsupported URL format: %s", url)
	}

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL: %s", url)
	}

	return parts[0], parts[1], nil
}

// cloneRepository clones a repository to the specified path
func cloneRepository(url, destPath string) error {
	cmd := exec.Command("git", "clone", url, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

// createBranch creates a new git branch
func createBranch(repoPath, branchName string) error {
	cmd := exec.Command("git", "-C", repoPath, "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git checkout failed: %w", err)
	}
	return nil
}

// commitAndPush commits changes and pushes to remote
func commitAndPush(repoPath, branchName, message string) error {
	// Add all files
	addCmd := exec.Command("git", "-C", repoPath, "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Commit
	commitCmd := exec.Command("git", "-C", repoPath, "commit", "-m", message)
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Push
	pushCmd := exec.Command("git", "-C", repoPath, "push", "-u", "origin", branchName)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}

// createPullRequest creates a pull request and returns the PR URL
func createPullRequest(repoPath, targetRepo, branchName, skillName string) (string, error) {
	// Extract owner/repo from target URL
	owner, repo, err := parseGitHubURL(targetRepo)
	if err != nil {
		return "", err
	}

	prTitle := fmt.Sprintf("Add skill: %s", skillName)
	prBody := fmt.Sprintf("This PR adds the skill '%s' to the repository.\n\nSkill structure:\n- SKILL.md: Skill documentation and metadata\n- Additional files and directories as needed", skillName)

	// Create PR using gh CLI
	cmd := exec.Command("gh", "pr", "create",
		"-R", fmt.Sprintf("%s/%s", owner, repo),
		"--base", "main",
		"--head", branchName,
		"--title", prTitle,
		"--body", prBody,
	)
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error creating PR: %s", string(output))
	}

	// Extract PR URL from output
	prURL := strings.TrimSpace(string(output))
	lines := strings.Split(prURL, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "https://github.com/") {
			return strings.TrimSpace(line), nil
		}
	}

	return prURL, nil
}

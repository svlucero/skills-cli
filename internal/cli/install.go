package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
	"github.com/silvinalucero/skill_cli/internal/skill"
	"github.com/spf13/cobra"
)

var (
	installRepo     string
	installProvider string
)

// getInstallHelp returns colored help text for install command
func getInstallHelp() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return fmt.Sprintf(`%s

%s
  %s            %s

%s
  %s           %s
  %s       %s
                          %s

%s
  - Copies the skill directory to the provider's skills location
  - For Claude: ~/.claude/skills/<skill-name>/
  - For Cursor: ~/.cursor/skills/<skill-name>/
  - Overwrites existing installation if skill already exists

%s
  %s                            %s
  %s              %s
  %s          %s
  %s`,
		"Install a skill from a repository to a provider.",

		cyan("PARAMETERS:"),
		yellow("<skill-name>"), dim("Name of the skill to install (required)"),

		cyan("FLAGS:"),
		green("--repo <name>"), dim("Repository to install from (default: active repo)"),
		green("--provider <name>"), dim("Provider to install to (default: claude)"),
		dim("Supported: claude, cursor"),

		cyan("BEHAVIOR:"),

		cyan("EXAMPLES:"),
		green("skill install explain-code"), dim("# Install to Claude from active repo"),
		green("skill install explain-code --repo myrepo"), dim("# Install from specific repo"),
		green("skill install explain-code --provider cursor"), dim("# Install to Cursor"),
		green("skill install deploy-app --repo company --provider claude"),
	)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <skill-name> [flags]",
	Short: "Install a skill to a provider",
	Long:  getInstallHelp(),
	Args:  cobra.ExactArgs(1),
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installRepo, "repo", "", "repository to install from (default: current)")
	installCmd.Flags().StringVar(&installProvider, "provider", "claude", "provider to install to (claude, cursor, etc.)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	skillName := args[0]

	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

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

	// Determine which repository to use
	var repoName string
	if installRepo != "" {
		repoName = installRepo
	} else {
		if cfg.ActiveRepo == "" {
			fmt.Printf("%s No active repository configured\n", red("✗"))
			return errors.ErrNoActiveRepo
		}
		repoName = cfg.ActiveRepo
	}

	// Verify repository exists
	if _, err := config.GetRepo(cfg, repoName); err != nil {
		fmt.Printf("%s Repository '%s' not found\n", red("✗"), repoName)
		return errors.ErrRepoNotFound
	}

	// Get all skills from the repository
	fmt.Printf("%s Looking for skill '%s' in repository '%s'...\n", cyan("→"), skillName, repoName)
	skills, err := skill.GetSkillsFromRepo(cfg, repoName)
	if err != nil {
		fmt.Printf("%s Error reading skills: %v\n", red("✗"), err)
		return err
	}

	// Find the skill
	var targetSkill *skill.Skill
	for i := range skills {
		if skills[i].Name == skillName {
			targetSkill = &skills[i]
			break
		}
	}

	if targetSkill == nil {
		fmt.Printf("%s Skill '%s' not found in repository '%s'\n", red("✗"), skillName, repoName)
		fmt.Println("\nAvailable skills:")
		for _, s := range skills {
			fmt.Printf("  - %s\n", s.Name)
		}
		return fmt.Errorf("skill not found")
	}

	fmt.Printf("%s Found skill: %s\n", green("✓"), targetSkill.Name)
	fmt.Printf("  Description: %s\n", targetSkill.Description)
	fmt.Printf("  Version: %s\n", targetSkill.Version)
	fmt.Println()

	// Install to provider
	fmt.Printf("%s Installing to provider '%s'...\n", cyan("→"), installProvider)

	switch installProvider {
	case "claude":
		err = installToClaude(targetSkill)
	case "cursor":
		err = installToCursor(targetSkill)
	default:
		fmt.Printf("%s Unknown provider: %s\n", red("✗"), installProvider)
		fmt.Println("\nSupported providers:")
		fmt.Println("  - claude")
		fmt.Println("  - cursor")
		return fmt.Errorf("unsupported provider")
	}

	if err != nil {
		fmt.Printf("%s Installation failed: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Skill '%s' installed successfully to %s!\n", green("✓"), skillName, installProvider)

	// Show next steps based on provider
	fmt.Println()
	fmt.Println("Next steps:")
	switch installProvider {
	case "claude":
		fmt.Println("  - Restart Claude Desktop to load the skill")
		fmt.Println("  - The skill will appear in your skills menu")
	case "cursor":
		fmt.Println("  - Restart Cursor to load the skill")
		fmt.Println("  - The skill will appear in your skills menu")
	}

	return nil
}

// installToClaude installs a skill to Claude (~/.claude/skills/)
func installToClaude(s *skill.Skill) error {
	// Get Claude skills directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	claudeSkillsDir := filepath.Join(homeDir, ".claude", "skills")
	destPath := filepath.Join(claudeSkillsDir, s.Name)

	// Create Claude skills directory if it doesn't exist
	if err := os.MkdirAll(claudeSkillsDir, 0755); err != nil {
		return fmt.Errorf("error creating Claude skills directory: %w", err)
	}

	// Check if skill already exists
	if _, err := os.Stat(destPath); err == nil {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s Skill already exists at: %s\n", yellow("⚠"), destPath)
		fmt.Println("  Removing existing installation...")
		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("error removing existing skill: %w", err)
		}
	}

	// Copy skill directory
	fmt.Printf("  Copying skill from: %s\n", s.Path)
	fmt.Printf("  To: %s\n", destPath)

	if err := copyDir(s.Path, destPath); err != nil {
		return fmt.Errorf("error copying skill: %w", err)
	}

	return nil
}

// installToCursor installs a skill to Cursor (~/.cursor/skills/)
func installToCursor(s *skill.Skill) error {
	// Get Cursor skills directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	cursorSkillsDir := filepath.Join(homeDir, ".cursor", "skills")
	destPath := filepath.Join(cursorSkillsDir, s.Name)

	// Create Cursor skills directory if it doesn't exist
	if err := os.MkdirAll(cursorSkillsDir, 0755); err != nil {
		return fmt.Errorf("error creating Cursor skills directory: %w", err)
	}

	// Check if skill already exists
	if _, err := os.Stat(destPath); err == nil {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s Skill already exists at: %s\n", yellow("⚠"), destPath)
		fmt.Println("  Removing existing installation...")
		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("error removing existing skill: %w", err)
		}
	}

	// Copy skill directory
	fmt.Printf("  Copying skill from: %s\n", s.Path)
	fmt.Printf("  To: %s\n", destPath)

	if err := copyDir(s.Path, destPath); err != nil {
		return fmt.Errorf("error copying skill: %w", err)
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Set permissions
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	return nil
}

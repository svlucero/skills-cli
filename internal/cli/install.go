package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
	"github.com/silvinalucero/skill_cli/internal/skill"
	"github.com/spf13/cobra"
)

var (
	installRepo        string
	installProvider    string
	installForce       bool
	installSkipExisting bool
)

// getInstallHelp returns colored help text for install command
func getInstallHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return NewHelpBuilder().
		Description("Install a skill from a repository to a provider.").
		Section("PARAMETERS:").
		Item(yellow("<skill-name>"), "Name of the skill to install (optional - interactive if omitted)").
		Section("FLAGS:").
		Item("--repo <name>", "Repository to install from (default: active repo)").
		Item("--provider <name>", "Provider to install to (default: claude)").
		Text("                          "+dim("Supported: claude, cursor")).
		Item("-f, --force", "Force reinstall even if already installed").
		Item("--skip-existing", "Skip installation if already installed (useful for batch installs)").
		Section("BEHAVIOR:").
		BulletList([]string{
			"If no skill name provided, shows interactive list to select from",
			"Use arrow keys to navigate and Enter to select a skill",
			"Checks if skill is already installed before installing",
			"Shows version information if available",
			"Prompts before overwriting existing installations",
			"Copies the skill directory to the provider's skills location",
			"For Claude: ~/.claude/skills/<skill-name>/",
			"For Cursor: ~/.cursor/skills/<skill-name>/",
		}).
		Section("EXAMPLES:").
		Example("skills install", "# Interactive - select from active repo").
		Example("skills install --repo myrepo", "# Interactive - select from specific repo").
		Example("skills install explain-code", "# Install to Claude from active repo").
		Example("skills install explain-code --repo myrepo", "# Install from specific repo").
		Example("skills install explain-code --provider cursor", "# Install to Cursor").
		Example("skills install deploy-app --repo company --provider claude", "").
		Build()
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [skill-name] [flags]",
	Short: "Install a skill to a provider",
	Long:  getInstallHelp(),
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installRepo, "repo", "", "repository to install from (default: current)")
	installCmd.Flags().StringVar(&installProvider, "provider", "claude", "provider to install to (claude, cursor, etc.)")
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "force reinstall even if already installed")
	installCmd.Flags().BoolVar(&installSkipExisting, "skip-existing", false, "skip installation if already installed")
}

// checkSkillInstalled checks if a skill is already installed for a provider
// Returns true if installed, false otherwise
func checkSkillInstalled(skillName, provider string) (bool, string, error) {
	skillsDir, err := getProviderSkillsDir(provider)
	if err != nil {
		return false, "", err
	}

	skillPath := filepath.Join(skillsDir, skillName)

	// Check if directory exists
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return false, "", nil
	} else if err != nil {
		return false, "", fmt.Errorf("error checking skill path: %w", err)
	}

	return true, skillPath, nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check configuration exists
	if !config.Exists() {
		fmt.Printf("%s Configuration not found\n", red("✗"))
		fmt.Println("  Run 'skills repository add <name> <repo-url>' to initialize")
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
	var skillName string
	var targetSkill *skill.Skill

	if len(args) == 0 {
		// No skill name provided - interactive mode
		fmt.Printf("%s Loading skills from repository '%s'...\n", cyan("→"), repoName)
		skills, err := skill.GetSkillsFromRepo(cfg, repoName)
		if err != nil {
			fmt.Printf("%s Error reading skills: %v\n", red("✗"), err)
			return err
		}

		if len(skills) == 0 {
			fmt.Printf("%s No skills found in repository '%s'\n", yellow("⚠"), repoName)
			return fmt.Errorf("no skills available")
		}

		// Interactive selection
		selected, err := selectSkillInteractive(skills, repoName)
		if err != nil {
			return err
		}

		targetSkill = selected
		skillName = selected.Name
	} else {
		// Skill name provided - traditional mode
		skillName = args[0]

		fmt.Printf("%s Looking for skill '%s' in repository '%s'...\n", cyan("→"), skillName, repoName)
		skills, err := skill.GetSkillsFromRepo(cfg, repoName)
		if err != nil {
			fmt.Printf("%s Error reading skills: %v\n", red("✗"), err)
			return err
		}

		// Find the skill
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
	}

	fmt.Printf("%s Found skill: %s\n", green("✓"), targetSkill.Name)
	fmt.Printf("  Description: %s\n", targetSkill.Description)
	fmt.Printf("  Version: %s\n", targetSkill.Version)
	fmt.Println()

	// Check if skill is already installed
	fmt.Printf("%s Checking installation status...\n", cyan("→"))
	isInstalled, installedPath, err := checkSkillInstalled(skillName, installProvider)
	if err != nil {
		fmt.Printf("%s %s\n", red("✗"), formatError("Error checking installation: %v", err))
		return err
	}

	if isInstalled {
		// Skill is already installed
		cyan := color.New(color.FgCyan).SprintFunc()

		if installSkipExisting {
			// Skip installation as requested
			fmt.Printf("%s Skill '%s' is already installed, skipping\n", cyan("ℹ"), skillName)
			fmt.Printf("  Location: %s\n", installedPath)
			fmt.Printf("  Version: %s\n", targetSkill.Version)
			return nil
		}

		if !installForce {
			// Not forcing, inform user and ask what to do
			fmt.Printf("%s Skill '%s' is already installed\n", cyan("ℹ"), skillName)
			fmt.Printf("  Location: %s\n", installedPath)
			fmt.Printf("  Version: %s\n", targetSkill.Version)
			fmt.Println()
			fmt.Printf("%s Nothing to do! Skill is already installed.\n", green("✓"))
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  • Use --force to reinstall anyway")
			fmt.Println("  • Use --skip-existing to skip without this message")
			return nil
		}

		// Force reinstall requested
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s %s\n", yellow("⚠"), formatWarning("Skill is already installed, reinstalling (--force)"))
		fmt.Printf("  Location: %s\n", installedPath)
		fmt.Println()
	} else {
		fmt.Printf("%s Skill not currently installed\n", green("✓"))
	}

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

// installToProvider installs a skill to a provider's skills directory
func installToProvider(s *skill.Skill, provider string) error {
	// Get provider skills directory
	skillsDir, err := getProviderSkillsDir(provider)
	if err != nil {
		return err
	}

	destPath := filepath.Join(skillsDir, s.Name)

	// Create skills directory if it doesn't exist
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("error creating %s skills directory: %w", provider, err)
	}

	// If skill already exists, remove it first (we already checked and got user consent)
	if _, err := os.Stat(destPath); err == nil {
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

// installToClaude installs a skill to Claude (~/.claude/skills/)
func installToClaude(s *skill.Skill) error {
	return installToProvider(s, "claude")
}

// installToCursor installs a skill to Cursor (~/.cursor/skills/)
func installToCursor(s *skill.Skill) error {
	return installToProvider(s, "cursor")
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

// selectSkillInteractive displays an interactive list to select a skill
func selectSkillInteractive(skills []skill.Skill, repoName string) (*skill.Skill, error) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }} {{ if eq .Status \"installed\" }}{{ \"✓\" | green }}{{ end }}",
		Inactive: "  {{ .Name | white }} {{ if eq .Status \"installed\" }}{{ \"✓\" | green }}{{ end }}",
		Selected: "{{ \"Selected:\" | green }} {{ .Name | cyan }}",
		Details: `
--------- Skill Details ---------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}
{{ "Version:" | faint }}	{{ .Version }}
{{ "Repository:" | faint }}	{{ .RepoName }}
{{ "Status:" | faint }}	{{ .Status }}`,
	}

	searcher := func(input string, index int) bool {
		skill := skills[index]
		name := skill.Name
		description := skill.Description

		input = promptui.Styler(promptui.FGBold)(input)

		if name == input || description == input {
			return true
		}

		return false
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("%s from %s (Use ↑/↓ arrows, / to search, Enter to select, Ctrl+C to exit)", cyan("Select skill to install"), repoName),
		Items:     skills,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	fmt.Println()
	idx, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Printf("\n%s Installation cancelled\n", yellow("⚠"))
			return nil, fmt.Errorf("installation cancelled")
		}
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	selected := &skills[idx]

	fmt.Println()
	fmt.Printf("%s Skill selected: %s\n", green("✓"), cyan(selected.Name))
	fmt.Printf("  Description: %s\n", selected.Description)
	if selected.Version != "" && selected.Version != "unknown" {
		fmt.Printf("  Version: %s\n", dim(selected.Version))
	}
	fmt.Println()

	return selected, nil
}

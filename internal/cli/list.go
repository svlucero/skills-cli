package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
	"github.com/silvinalucero/skill_cli/internal/git"
	"github.com/silvinalucero/skill_cli/internal/skill"
	"github.com/spf13/cobra"
)

var (
	listRepo      string
	listAll       bool
	listCompact   bool
	listNoUpdate  bool
	listInstalled bool
	listProvider  string
)

// getListHelp returns colored help text for list command
func getListHelp() string {
	return NewHelpBuilder().
		Description("List available skills from repositories or installed skills from a provider.").
		Section("FLAGS:").
		Item("--repo <name>", "List skills from a specific repository").
		Item("--all", "List skills from all repositories").
		Item("--installed", "List installed skills (requires --provider)").
		Item("--provider <name>", "Provider to list from (claude, cursor)").
		Item("-c, --compact", "Compact non-interactive format").
		Item("--no-update", "Skip repository update before listing").
		Section("BEHAVIOR:").
		BulletList([]string{
			"By default, lists skills from the active repository in interactive mode",
			"Use arrow keys to navigate and Enter to select a skill",
			"Repository is automatically updated (git pull) before listing",
			"Use --no-update to skip the repository update",
			"Use --compact for non-interactive output",
			"Use --installed with --provider to see installed skills",
		}).
		Section("EXAMPLES:").
		Example("skills list", "# List from active repo").
		Example("skills list --repo myrepo", "# List from specific repo").
		Example("skills list --all", "# List from all repos").
		Example("skills list --installed --provider claude", "# List installed in Claude").
		Example("skills list --compact", "# Compact format").
		Example("skills list --no-update", "# Skip repo update").
		Build()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List available or installed skills",
	Long:  getListHelp(),
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&listRepo, "repo", "", "list skills from a specific repository")
	listCmd.Flags().BoolVar(&listAll, "all", false, "list skills from all repositories")
	listCmd.Flags().BoolVarP(&listCompact, "compact", "c", false, "non-interactive compact format")
	listCmd.Flags().BoolVar(&listNoUpdate, "no-update", false, "skip repository update before listing")
	listCmd.Flags().BoolVar(&listInstalled, "installed", false, "list installed skills (requires --provider)")
	listCmd.Flags().StringVar(&listProvider, "provider", "", "provider to list installed skills from (claude, cursor)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Colors for output
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	// If --installed is used, --provider is required
	if listInstalled {
		if listProvider == "" {
			fmt.Printf("%s --installed requires --provider to be specified\n", red("✗"))
			fmt.Println("\nSupported providers:")
			fmt.Println("  - claude")
			fmt.Println("  - cursor")
			fmt.Println("\nExample: skills list --installed --provider claude")
			return fmt.Errorf("missing required flag: --provider")
		}

		// List installed skills
		return runListInstalled(listProvider, listCompact)
	}

	// Check configuration exists (only for repository listing)
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

	// Validate mutually exclusive flags
	if listRepo != "" && listAll {
		fmt.Printf("%s Cannot use --repo and --all simultaneously\n", red("✗"))
		return fmt.Errorf("conflicting flags")
	}

	// Determine which repositories to update and list from
	var reposToUpdate []string
	var sourceDesc string

	if listAll {
		// Update and list from all repositories
		for name := range cfg.Repositories {
			reposToUpdate = append(reposToUpdate, name)
		}
		sourceDesc = "all repositories"
	} else if listRepo != "" {
		// Update and list from specific repository
		if _, err := config.GetRepo(cfg, listRepo); err != nil {
			fmt.Printf("%s Repository '%s' not found\n", red("✗"), listRepo)
			return errors.ErrRepoNotFound
		}
		reposToUpdate = []string{listRepo}
		sourceDesc = listRepo
	} else {
		// Update and list from active repository
		if cfg.ActiveRepo == "" {
			fmt.Printf("%s No active repository configured\n", red("✗"))
			return errors.ErrNoActiveRepo
		}
		reposToUpdate = []string{cfg.ActiveRepo}
		sourceDesc = fmt.Sprintf("%s (current)", cfg.ActiveRepo)
	}

	// Update repositories before listing (unless --no-update)
	if !listNoUpdate {
		if len(reposToUpdate) == 1 {
			fmt.Printf("%s Updating repository...\n", dim("→"))
		} else {
			fmt.Printf("%s Updating %d repositories...\n", dim("→"), len(reposToUpdate))
		}

		for _, repoName := range reposToUpdate {
			repo, err := config.GetRepo(cfg, repoName)
			if err != nil {
				continue
			}

			// Check if repository is cloned
			if !git.RepoExists(repo.LocalPath) {
				fmt.Printf("  %s %s: not cloned locally, skipping update\n", yellow("⚠"), repoName)
				continue
			}

			// Pull latest changes
			if err := git.Pull(repo.LocalPath); err != nil {
				fmt.Printf("  %s %s: failed to update (%v)\n", yellow("⚠"), repoName, err)
				// Continue anyway - we'll list what we have
			} else {
				fmt.Printf("  %s %s updated\n", green("✓"), repoName)
			}
		}
		fmt.Println()
	}

	// Get skills to list
	var skills []skill.Skill

	if listAll {
		// List from all repositories
		skills, err = skill.GetAllSkills(cfg)
		if err != nil {
			fmt.Printf("%s Error getting skills: %v\n", red("✗"), err)
			return err
		}
	} else if listRepo != "" {
		// List from specific repository
		skills, err = skill.GetSkillsFromRepo(cfg, listRepo)
		if err != nil {
			fmt.Printf("%s Error getting skills from '%s': %v\n", red("✗"), listRepo, err)
			return err
		}
	} else {
		// List from active repository
		skills, err = skill.GetSkillsFromRepo(cfg, cfg.ActiveRepo)
		if err != nil {
			fmt.Printf("%s Error getting skills from active repository: %v\n", red("✗"), err)
			return err
		}
	}

	// Check if there are skills
	if len(skills) == 0 {
		fmt.Printf("%s No skills found in %s\n", yellow("⚠"), sourceDesc)
		return nil
	}

	// Display skills according to format
	// Interactive mode is default, --compact disables it
	if listCompact {
		if listAll {
			printSkillsByRepo(skills)
		} else {
			printSkillsCompact(skills, sourceDesc)
		}

		// Summary for non-interactive mode
		fmt.Println()
		fmt.Printf("%s Found %d skill(s)\n", dim("→"), len(skills))
	} else {
		// Interactive mode (default)
		return runListInteractive(skills, sourceDesc)
	}

	return nil
}

// runListInteractive displays skills in an interactive list with cursor navigation
func runListInteractive(skills []skill.Skill, source string) error {
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
		Label:     fmt.Sprintf("%s from %s (Use ↑/↓ arrows, / to search, Enter to select, Ctrl+C to exit)", cyan("Skills"), source),
		Items:     skills,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	fmt.Println()
	idx, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Printf("\n%s Selection cancelled\n", yellow("⚠"))
			return nil
		}
		return fmt.Errorf("prompt failed: %w", err)
	}

	selected := skills[idx]

	fmt.Println()
	fmt.Printf("%s Skill selected: %s\n", green("✓"), cyan(selected.Name))
	fmt.Printf("  Description: %s\n", selected.Description)
	if selected.Version != "" && selected.Version != "unknown" {
		fmt.Printf("  Version: %s\n", dim(selected.Version))
	}
	fmt.Printf("  Repository: %s\n", selected.RepoName)
	if selected.Status == skill.StatusInstalled {
		fmt.Printf("  Status: %s\n", green(string(selected.Status)))
	} else {
		fmt.Printf("  Status: %s\n", selected.Status)
	}

	return nil
}

// printSkillsCompact prints skills in compact format (one line each)
func printSkillsCompact(skills []skill.Skill, source string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("%s from %s:\n\n", cyan("Skills"), source)

	for _, s := range skills {
		statusStr := ""
		if s.Status == skill.StatusInstalled {
			statusStr = green(" [installed]")
		}

		// Format: name (repo) - description [status]
		if source == "all repositories" {
			fmt.Printf("  %s %s - %s%s\n", s.Name, dim(fmt.Sprintf("(%s)", s.RepoName)), s.Description, statusStr)
		} else {
			fmt.Printf("  %s - %s%s\n", s.Name, s.Description, statusStr)
		}
	}
}

// printSkillsByRepo prints skills grouped by repository
func printSkillsByRepo(skills []skill.Skill) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	// Group skills by repository
	skillsByRepo := make(map[string][]skill.Skill)
	for _, s := range skills {
		skillsByRepo[s.RepoName] = append(skillsByRepo[s.RepoName], s)
	}

	// Print each repository
	first := true
	for repoName, repoSkills := range skillsByRepo {
		if !first {
			fmt.Println()
		}
		first = false

		fmt.Printf("%s from %s:\n", cyan("Skills"), repoName)

		for _, s := range repoSkills {
			fmt.Println()
			fmt.Printf("  %s\n", cyan(s.Name))
			fmt.Printf("    %s\n", s.Description)

			// Version if available
			if s.Version != "" && s.Version != "unknown" {
				fmt.Printf("    Version: %s\n", dim(s.Version))
			}

			// Status
			statusStr := string(s.Status)
			if s.Status == skill.StatusInstalled {
				statusStr = green(statusStr)
			}
			fmt.Printf("    Status: %s\n", statusStr)
		}
	}
}

// runListInstalled lists installed skills from a provider
func runListInstalled(provider string, compact bool) error {
	// Colors for output
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	// Get provider directory
	providerDir, err := getProviderSkillsDir(provider)
	if err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return err
	}

	// Check if directory exists
	if _, err := os.Stat(providerDir); os.IsNotExist(err) {
		fmt.Printf("%s No skills installed in %s\n", yellow("⚠"), provider)
		fmt.Printf("Skills directory does not exist: %s\n", providerDir)
		return nil
	}

	// Discover installed skills
	fmt.Printf("%s Scanning installed skills in %s...\n", cyan("→"), provider)
	installedSkills, err := discoverInstalledSkills(providerDir, provider)
	if err != nil {
		fmt.Printf("%s Error scanning installed skills: %v\n", red("✗"), err)
		return err
	}

	if len(installedSkills) == 0 {
		fmt.Printf("%s No skills installed in %s\n", yellow("⚠"), provider)
		return nil
	}

	fmt.Println()

	// Display skills
	if compact {
		printInstalledSkillsCompact(installedSkills, provider)
	} else {
		printInstalledSkillsDetailed(installedSkills, provider)
	}

	// Summary
	fmt.Println()
	fmt.Printf("%s Found %d installed skill(s) in %s\n", dim("→"), len(installedSkills), provider)

	return nil
}

// getProviderSkillsDir returns the skills directory for a provider
func getProviderSkillsDir(provider string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	switch provider {
	case "claude":
		return filepath.Join(homeDir, ".claude", "skills"), nil
	case "cursor":
		return filepath.Join(homeDir, ".cursor", "skills"), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s (supported: claude, cursor)", provider)
	}
}

// discoverInstalledSkills discovers skills installed in a provider directory
func discoverInstalledSkills(providerDir, provider string) ([]skill.Skill, error) {
	entries, err := os.ReadDir(providerDir)
	if err != nil {
		return nil, fmt.Errorf("error reading provider directory: %w", err)
	}

	var skills []skill.Skill
	for _, entry := range entries {
		// Only process directories
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()

		// Skip hidden directories
		if skillName[0] == '.' {
			continue
		}

		skillPath := filepath.Join(providerDir, skillName)

		// Try to parse SKILL.md frontmatter
		markerPath := filepath.Join(skillPath, "SKILL.md")
		skillCfg, err := skill.ParseSkillMarkdown(markerPath)

		// If SKILL.md doesn't exist or fails, use defaults
		if err != nil {
			skillCfg = &skill.SkillConfig{
				Name:        skillName,
				Description: "No description available",
				Version:     "unknown",
			}
		}

		// Get file info for last modification time
		info, err := entry.Info()
		if err != nil {
			continue
		}

		installedSkill := skill.Skill{
			Name:        skillCfg.Name,
			Description: skillCfg.Description,
			Version:     skillCfg.Version,
			RepoName:    provider,
			Path:        skillPath,
			Status:      skill.StatusInstalled,
			UpdatedAt:   info.ModTime(),
		}

		skills = append(skills, installedSkill)
	}

	// Sort by name
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}

// printInstalledSkillsCompact prints installed skills in compact format
func printInstalledSkillsCompact(skills []skill.Skill, provider string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("%s installed in %s:\n\n", cyan("Skills"), provider)

	for _, s := range skills {
		fmt.Printf("  %s %s - %s\n", green("✓"), s.Name, s.Description)
	}
}

// printInstalledSkillsDetailed prints installed skills in detailed format
func printInstalledSkillsDetailed(skills []skill.Skill, provider string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	fmt.Printf("%s installed in %s:\n", cyan("Skills"), provider)

	for _, s := range skills {
		fmt.Println()
		fmt.Printf("  %s %s\n", green("✓"), cyan(s.Name))
		fmt.Printf("    %s\n", s.Description)

		// Version if available
		if s.Version != "" && s.Version != "unknown" {
			fmt.Printf("    Version: %s\n", dim(s.Version))
		}

		fmt.Printf("    Location: %s\n", dim(s.Path))
	}
}

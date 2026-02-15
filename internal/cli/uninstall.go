package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/silvinalucero/skill_cli/internal/skill"
	"github.com/spf13/cobra"
)

var (
	uninstallProvider string
	uninstallForce    bool
)

// getUninstallHelp returns colored help text for uninstall command
func getUninstallHelp() string {
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	return NewHelpBuilder().
		Description("Uninstall a skill from a provider.").
		Section("PARAMETERS:").
		Item(yellow("<skill-name>"), "Name of the skill to uninstall (optional - interactive if omitted)").
		Section("FLAGS:").
		Item("--provider <name>", "Provider to uninstall from (default: claude)").
		Text("                          "+dim("Supported: claude, cursor")).
		Item("-f, --force", "Force uninstall without confirmation").
		Section("BEHAVIOR:").
		BulletList([]string{
			"If no skill name provided, shows interactive list of installed skills",
			"Use arrow keys to navigate and Enter to select a skill",
			"Checks if skill is installed before uninstalling",
			"Prompts for confirmation before removing (unless --force)",
			"Removes the skill directory from the provider's skills location",
			"For Claude: ~/.claude/skills/<skill-name>/",
			"For Cursor: ~/.cursor/skills/<skill-name>/",
		}).
		Section("EXAMPLES:").
		Example("skills uninstall", "# Interactive - select from Claude").
		Example("skills uninstall --provider cursor", "# Interactive - select from Cursor").
		Example("skills uninstall explain-code", "# Uninstall from Claude").
		Example("skills uninstall explain-code --provider cursor", "# Uninstall from Cursor").
		Example("skills uninstall deploy-app --force", "# Uninstall without confirmation").
		Build()
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [skill-name] [flags]",
	Short: "Uninstall a skill from a provider",
	Long:  getUninstallHelp(),
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().StringVar(&uninstallProvider, "provider", "claude", "provider to uninstall from (claude, cursor)")
	uninstallCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "force uninstall without confirmation")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	providerDir, err := getProviderSkillsDir(uninstallProvider)
	if err != nil {
		fmt.Printf("%s %v\n", red("✗"), err)
		return err
	}

	if _, err := os.Stat(providerDir); os.IsNotExist(err) {
		fmt.Printf("%s No skills directory found for %s\n", yellow("⚠"), uninstallProvider)
		fmt.Printf("Directory: %s\n", providerDir)
		return fmt.Errorf("skills directory does not exist")
	}

	var skillName string

	if len(args) == 0 {
		// No skill name provided - interactive mode
		fmt.Printf("%s Loading installed skills from %s...\n", cyan("→"), uninstallProvider)

		installedSkills, err := discoverInstalledSkills(providerDir, uninstallProvider)
		if err != nil {
			fmt.Printf("%s Error scanning installed skills: %v\n", red("✗"), err)
			return err
		}

		if len(installedSkills) == 0 {
			fmt.Printf("%s No skills installed in %s\n", yellow("⚠"), uninstallProvider)
			return fmt.Errorf("no skills to uninstall")
		}

		selected, err := selectSkillToUninstall(installedSkills, uninstallProvider)
		if err != nil {
			return err
		}

		skillName = selected.Name
	} else {
		// Skill name provided - traditional mode
		skillName = args[0]
	}

	skillPath := filepath.Join(providerDir, skillName)

	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		fmt.Printf("%s Skill '%s' is not installed in %s\n", red("✗"), skillName, uninstallProvider)

		installedSkills, err := discoverInstalledSkills(providerDir, uninstallProvider)
		if err == nil && len(installedSkills) > 0 {
			fmt.Println("\nInstalled skills:")
			for _, s := range installedSkills {
				fmt.Printf("  - %s\n", s.Name)
			}
		}
		return fmt.Errorf("skill not installed")
	}

	fmt.Println()
	fmt.Printf("%s Skill found: %s\n", green("✓"), cyan(skillName))
	fmt.Printf("  Location: %s\n", skillPath)
	fmt.Printf("  Provider: %s\n", uninstallProvider)

	if !uninstallForce {
		fmt.Println()
		fmt.Printf("%s Are you sure you want to uninstall '%s'? [y/N]: ", yellow("?"), skillName)

		var response string
		_, _ = fmt.Scanln(&response)

		if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
			fmt.Printf("\n%s Uninstall cancelled\n", yellow("⚠"))
			return nil
		}
	}

	fmt.Println()
	fmt.Printf("%s Uninstalling skill '%s'...\n", cyan("→"), skillName)

	if err := os.RemoveAll(skillPath); err != nil {
		fmt.Printf("%s Failed to uninstall: %v\n", red("✗"), err)
		return err
	}

	fmt.Printf("%s Skill '%s' uninstalled successfully from %s!\n", green("✓"), skillName, uninstallProvider)

	fmt.Println()
	fmt.Println("Next steps:")
	switch uninstallProvider {
	case "claude":
		fmt.Println("  - Restart Claude Desktop to remove the skill from menu")
	case "cursor":
		fmt.Println("  - Restart Cursor to remove the skill from menu")
	}

	return nil
}

// selectSkillToUninstall displays an interactive list to select a skill for uninstallation
func selectSkillToUninstall(skills []skill.Skill, provider string) (*skill.Skill, error) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }}",
		Inactive: "  {{ .Name | white }}",
		Selected: "{{ \"Selected:\" | green }} {{ .Name | cyan }}",
		Details: `
--------- Skill Details ---------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}
{{ "Version:" | faint }}	{{ .Version }}
{{ "Location:" | faint }}	{{ .Path }}`,
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
		Label:     fmt.Sprintf("%s from %s (Use ↑/↓ arrows, / to search, Enter to select, Ctrl+C to exit)", cyan("Select skill to uninstall"), provider),
		Items:     skills,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	fmt.Println()
	idx, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Printf("\n%s Uninstall cancelled\n", yellow("⚠"))
			return nil, fmt.Errorf("uninstall cancelled")
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

package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/silvinalucero/skill_cli/internal/config"
	"gopkg.in/yaml.v3"
)

const (
	skillConfigName  = "config.yaml"
	skillMarkerName  = "SKILL.md"
	installedDirName = "installed"
)

// DiscoverSkills scans all directories in a local repository as potential skills
func DiscoverSkills(repoPath, repoName, skillsPath string) ([]Skill, error) {
	// Check if repository path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("repository path not found: %s", repoPath)
	}

	// Combine repoPath with skillsPath
	scanPath := repoPath
	if skillsPath != "" {
		scanPath = filepath.Join(repoPath, skillsPath)

		// Validate that skills path exists
		if _, err := os.Stat(scanPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("skills path not found: %s (full path: %s)", skillsPath, scanPath)
		}
	}

	// Read all directories in the skills directory
	entries, err := os.ReadDir(scanPath)
	if err != nil {
		return nil, fmt.Errorf("error reading skills directory: %w", err)
	}

	var skills []Skill
	for _, entry := range entries {
		// Only process directories
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()

		// Skip .git directory and other hidden directories
		if skillName == ".git" || skillName[0] == '.' {
			continue
		}

		skillPath := filepath.Join(scanPath, skillName)

		// Check if SKILL.md exists - this marks a directory as a skill
		markerPath := filepath.Join(skillPath, skillMarkerName)
		if _, err := os.Stat(markerPath); os.IsNotExist(err) {
			// No SKILL.md file, skip silently (not a skill)
			continue
		}

		// Try to parse SKILL.md frontmatter
		skillCfg, err := ParseSkillMarkdown(markerPath)

		// If SKILL.md frontmatter parsing fails, try config.yaml as fallback
		if err != nil {
			configPath := filepath.Join(skillPath, skillConfigName)
			skillCfg, err = ParseSkillConfig(configPath)

			// If both fail, use defaults based on folder name
			if err != nil {
				skillCfg = &SkillConfig{
					Name:        skillName, // Use folder name as skill name
					Description: "No description available",
					Version:     "unknown",
				}
			}
		}

		// Get file info for last modification time
		info, err := entry.Info()
		if err != nil {
			continue // Skip if we can't get info
		}

		// Determine skill status
		status := GetSkillStatus(skillCfg.Name)

		skill := Skill{
			Name:        skillCfg.Name,
			Description: skillCfg.Description,
			Version:     skillCfg.Version,
			RepoName:    repoName,
			Path:        skillPath,
			Status:      status,
			UpdatedAt:   info.ModTime(),
		}

		skills = append(skills, skill)
	}

	// Sort by name
	sort.Sort(ByName(skills))

	return skills, nil
}

// ParseSkillMarkdown parses the YAML frontmatter from a SKILL.md file
func ParseSkillMarkdown(markdownPath string) (*SkillConfig, error) {
	// Read file
	data, err := os.ReadFile(markdownPath)
	if err != nil {
		return nil, fmt.Errorf("error reading SKILL.md: %w", err)
	}

	content := string(data)

	// Check if file starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return nil, fmt.Errorf("SKILL.md missing frontmatter (should start with ---)")
	}

	// Find the closing frontmatter delimiter
	lines := strings.Split(content, "\n")
	var frontmatterLines []string
	inFrontmatter := false
	frontmatterEnd := -1

	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.TrimSpace(line) == "---" {
			frontmatterEnd = i
			break
		}
		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
		}
	}

	if frontmatterEnd == -1 {
		return nil, fmt.Errorf("SKILL.md frontmatter not properly closed (missing closing ---)")
	}

	// Parse YAML frontmatter
	frontmatter := strings.Join(frontmatterLines, "\n")
	var cfg SkillConfig
	if err := yaml.Unmarshal([]byte(frontmatter), &cfg); err != nil {
		return nil, fmt.Errorf("error parsing SKILL.md frontmatter: %w", err)
	}

	// Validate required fields
	if cfg.Name == "" {
		return nil, fmt.Errorf("name is required in SKILL.md frontmatter")
	}

	// Description is optional but good to have
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}

	// Version is optional
	if cfg.Version == "" {
		cfg.Version = "unknown"
	}

	return &cfg, nil
}

// ParseSkillConfig parses a skill's config.yaml file (legacy/optional)
func ParseSkillConfig(configPath string) (*SkillConfig, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.yaml not found")
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config.yaml: %w", err)
	}

	// Parse YAML
	var cfg SkillConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config.yaml: %w", err)
	}

	// Validate required fields
	if cfg.Name == "" {
		return nil, fmt.Errorf("skill name is required in config.yaml")
	}

	// Description is optional but good to have
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}

	// Version is optional
	if cfg.Version == "" {
		cfg.Version = "unknown"
	}

	return &cfg, nil
}

// GetSkillStatus determines if a skill is installed
func GetSkillStatus(skillName string) Status {
	// For now, all are "available"
	// In the future, this will check in ~/.local/share/skill/installed/
	// if the skill is installed
	return StatusAvailable
}

// GetSkillsFromRepo gets all skills from a specific repository
func GetSkillsFromRepo(cfg *config.Config, repoName string) ([]Skill, error) {
	repo, err := config.GetRepo(cfg, repoName)
	if err != nil {
		return nil, err
	}

	// Check if repository is cloned locally
	if _, err := os.Stat(repo.LocalPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("repository '%s' not cloned locally\nRun 'skill repository add %s %s' to clone it", repoName, repoName, repo.URL)
	}

	return DiscoverSkills(repo.LocalPath, repoName, repo.SkillsPath)
}

// GetAllSkills gets all skills from all repositories
func GetAllSkills(cfg *config.Config) ([]Skill, error) {
	var allSkills []Skill

	for name := range cfg.Repositories {
		skills, err := GetSkillsFromRepo(cfg, name)
		if err != nil {
			// Not a fatal error, just warn
			fmt.Fprintf(os.Stderr, "Warning: could not read skills from '%s': %v\n", name, err)
			continue
		}

		allSkills = append(allSkills, skills...)
	}

	// Sort by repository then by name
	sort.Sort(ByRepo(allSkills))

	return allSkills, nil
}

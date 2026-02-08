package skill

import "time"

// Status represents the state of a skill
type Status string

const (
	StatusAvailable Status = "available"
	StatusInstalled Status = "installed"
)

// Skill represents a skill discovered in a repository
type Skill struct {
	Name        string    // Name of the skill (directory)
	Description string    // Description of the skill
	Version     string    // Version of the skill
	RepoName    string    // Name of the source repository
	Path        string    // Full path to the skill directory
	Status      Status    // Status of the skill (available/installed)
	UpdatedAt   time.Time // Last update of the skill
}

// SkillConfig represents the structure of a skill's config.yaml
type SkillConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// ByName implements sort.Interface to sort skills by name
type ByName []Skill

func (s ByName) Len() int           { return len(s) }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ByRepo implements sort.Interface to sort skills by repo then by name
type ByRepo []Skill

func (s ByRepo) Len() int { return len(s) }
func (s ByRepo) Less(i, j int) bool {
	if s[i].RepoName != s[j].RepoName {
		return s[i].RepoName < s[j].RepoName
	}
	return s[i].Name < s[j].Name
}
func (s ByRepo) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
}

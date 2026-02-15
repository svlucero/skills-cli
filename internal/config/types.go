package config

import "time"

// Config represents the complete CLI configuration
type Config struct {
	Version      string                `yaml:"version"`
	ActiveRepo   string                `yaml:"active_repo,omitempty"`
	Repositories map[string]Repository `yaml:"repositories,omitempty"`
	// Obsolete field for v1 compatibility
	Repository Repository `yaml:"repository,omitempty"`
}

// Repository contains information about the skills repository
type Repository struct {
	Name         string    `yaml:"name"`
	URL          string    `yaml:"url"`
	LocalPath    string    `yaml:"local_path"`
	SkillsPath   string    `yaml:"skills_path,omitempty"`
	LastVerified time.Time `yaml:"last_verified"`
	AuthType     string    `yaml:"auth_type"` // "https" or "ssh"
}

// AuthType represents the authentication type
type AuthType string

const (
	AuthHTTPS AuthType = "https"
	AuthSSH   AuthType = "ssh"
)

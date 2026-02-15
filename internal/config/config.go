package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	configDirName  = "skill"
	configFileName = "config.yaml"
	currentVersion = "2"
	reposDirName   = "repos"
)

// GetConfigPath returns the full path of the configuration file
func GetConfigPath() (string, error) {
	configPath := filepath.Join(xdg.ConfigHome, configDirName, configFileName)
	return configPath, nil
}

// GetRepoPath returns the path where the local repository should be cloned (v1 legacy)
func GetRepoPath() (string, error) {
	repoPath := filepath.Join(xdg.DataHome, configDirName, "repo")
	return repoPath, nil
}

// GetReposPath returns the base path for multiple repositories
func GetReposPath() string {
	return filepath.Join(xdg.DataHome, configDirName, reposDirName)
}

// GetRepoPathForRepo returns the local path for a specific repository
func GetRepoPathForRepo(name string) string {
	return filepath.Join(GetReposPath(), name)
}

// Exists checks if the configuration file exists
func Exists() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(configPath)
	return err == nil
}

// Load loads the configuration from the file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %w", err)
	}

	// Check that the file exists
	if !Exists() {
		return nil, fmt.Errorf("configuration file not found at %s", configPath)
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Automatically migrate if version 1
	if cfg.Version == "1" {
		fmt.Println("Migrating configuration to version 2...")
		if err := Migrate(&cfg); err != nil {
			return nil, fmt.Errorf("error migrating config: %w", err)
		}
		// Save migrated configuration
		if err := Save(&cfg); err != nil {
			return nil, fmt.Errorf("error saving migrated config: %w", err)
		}
		fmt.Println("✓ Configuration migrated successfully")
	}

	return &cfg, nil
}

// Save saves the configuration to the file
func Save(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("error getting config path: %w", err)
	}

	// Create the configuration directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Ensure version is set
	if cfg.Version == "" {
		cfg.Version = currentVersion
	}

	// Serialize to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write the file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// Validate verifies that the configuration is valid
func Validate(cfg *Config) error {
	// Validate v2
	if cfg.Version == "2" {
		if len(cfg.Repositories) == 0 {
			return fmt.Errorf("no repositories configured")
		}

		if cfg.ActiveRepo == "" {
			return fmt.Errorf("no active repository set")
		}

		if _, exists := cfg.Repositories[cfg.ActiveRepo]; !exists {
			return fmt.Errorf("active repository '%s' not found in repositories", cfg.ActiveRepo)
		}

		// Validate each repository
		for name, repo := range cfg.Repositories {
			if repo.Name != name {
				return fmt.Errorf("repository name mismatch: key='%s', name='%s'", name, repo.Name)
			}

			if repo.URL == "" {
				return fmt.Errorf("repository '%s' has empty URL", name)
			}

			if repo.LocalPath == "" {
				return fmt.Errorf("repository '%s' has empty local path", name)
			}

			if repo.AuthType != string(AuthHTTPS) && repo.AuthType != string(AuthSSH) {
				return fmt.Errorf("repository '%s' has invalid auth type: %s", name, repo.AuthType)
			}
		}

		return nil
	}

	// Validate v1 (legacy)
	if cfg.Repository.URL == "" {
		return fmt.Errorf("repository URL is required")
	}

	if cfg.Repository.LocalPath == "" {
		return fmt.Errorf("repository local path is required")
	}

	if cfg.Repository.AuthType != string(AuthHTTPS) && cfg.Repository.AuthType != string(AuthSSH) {
		return fmt.Errorf("invalid auth type: %s (must be 'https' or 'ssh')", cfg.Repository.AuthType)
	}

	return nil
}

// EnsureDataDir creates the data directory if it doesn't exist
func EnsureDataDir() error {
	repoPath, err := GetRepoPath()
	if err != nil {
		return err
	}

	dataDir := filepath.Dir(repoPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("error creating data directory: %w", err)
	}

	return nil
}

// LoadWithViper loads the configuration using Viper (alternative for future)
func LoadWithViper() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// Migrate migrates the configuration from v1 to v2
func Migrate(cfg *Config) error {
	if cfg.Version != "1" {
		return nil // Already in v2 or higher
	}

	// Convert the single repository to a map with name "default"
	if cfg.Repository.URL != "" {
		cfg.Repositories = make(map[string]Repository)
		repo := cfg.Repository
		repo.Name = "default"
		cfg.Repositories["default"] = repo
		cfg.ActiveRepo = "default"
	}

	// Clear obsolete field
	cfg.Repository = Repository{}
	cfg.Version = "2"

	return nil
}

// GetActiveRepo gets the active repository
func GetActiveRepo(cfg *Config) (*Repository, error) {
	if cfg.ActiveRepo == "" {
		return nil, fmt.Errorf("no active repository configured")
	}

	repo, exists := cfg.Repositories[cfg.ActiveRepo]
	if !exists {
		return nil, fmt.Errorf("active repository '%s' not found", cfg.ActiveRepo)
	}

	return &repo, nil
}

// GetRepo gets a repository by name
func GetRepo(cfg *Config, name string) (*Repository, error) {
	repo, exists := cfg.Repositories[name]
	if !exists {
		return nil, fmt.Errorf("repository '%s' not found", name)
	}

	return &repo, nil
}

// AddRepository adds a new repository
func AddRepository(cfg *Config, repo Repository) error {
	if err := ValidateRepoName(repo.Name); err != nil {
		return err
	}

	if cfg.Repositories == nil {
		cfg.Repositories = make(map[string]Repository)
	}

	if _, exists := cfg.Repositories[repo.Name]; exists {
		return fmt.Errorf("repository '%s' already exists", repo.Name)
	}

	cfg.Repositories[repo.Name] = repo

	// If it's the first repository, make it active
	if cfg.ActiveRepo == "" {
		cfg.ActiveRepo = repo.Name
	}

	return nil
}

// RemoveRepository removes a repository
func RemoveRepository(cfg *Config, name string) error {
	if _, exists := cfg.Repositories[name]; !exists {
		return fmt.Errorf("repository '%s' not found", name)
	}

	if cfg.ActiveRepo == name {
		return fmt.Errorf("cannot remove active repository, switch to another first")
	}

	delete(cfg.Repositories, name)
	return nil
}

// SetActiveRepo changes the active repository
func SetActiveRepo(cfg *Config, name string) error {
	if _, exists := cfg.Repositories[name]; !exists {
		return fmt.Errorf("repository '%s' not found", name)
	}

	cfg.ActiveRepo = name
	return nil
}

// ValidateRepoName validates a repository name
func ValidateRepoName(name string) error {
	if name == "" {
		return fmt.Errorf("repository name cannot be empty")
	}

	// Only allow alphanumeric characters, hyphens and underscores
	for _, ch := range name {
		isLower := ch >= 'a' && ch <= 'z'
		isUpper := ch >= 'A' && ch <= 'Z'
		isDigit := ch >= '0' && ch <= '9'
		isHyphen := ch == '-'
		isUnderscore := ch == '_'

		if !isLower && !isUpper && !isDigit && !isHyphen && !isUnderscore {
			return fmt.Errorf("repository name can only contain letters, numbers, hyphens and underscores")
		}
	}

	return nil
}

// NormalizeGitURL normalizes a Git URL for comparison purposes
// This handles:
// - Removing .git suffix
// - Converting SSH to HTTPS format
// - Lowercasing for case-insensitive comparison
func NormalizeGitURL(url string) string {
	normalized := strings.TrimSpace(url)

	// Remove .git suffix if present
	normalized = strings.TrimSuffix(normalized, ".git")

	// Convert SSH format to HTTPS format for comparison
	// git@github.com:org/repo -> https://github.com/org/repo
	if strings.HasPrefix(normalized, "git@github.com:") {
		normalized = strings.Replace(normalized, "git@github.com:", "https://github.com/", 1)
	}

	// Handle ssh://git@github.com/org/repo format
	if strings.HasPrefix(normalized, "ssh://git@github.com/") {
		normalized = strings.Replace(normalized, "ssh://git@github.com/", "https://github.com/", 1)
	}

	// Generic SSH URL pattern: git@host:org/repo -> https://host/org/repo
	if strings.HasPrefix(normalized, "git@") && strings.Contains(normalized, ":") {
		parts := strings.SplitN(normalized, ":", 2)
		if len(parts) == 2 {
			host := strings.TrimPrefix(parts[0], "git@")
			normalized = fmt.Sprintf("https://%s/%s", host, parts[1])
		}
	}

	// Lowercase for case-insensitive comparison
	return strings.ToLower(normalized)
}

// FindRepoByURL searches for a repository by its URL (normalized)
// Returns the repository and true if found, nil and false otherwise
func FindRepoByURL(cfg *Config, url string) (*Repository, bool) {
	if cfg == nil || cfg.Repositories == nil {
		return nil, false
	}

	normalizedURL := NormalizeGitURL(url)

	for _, repo := range cfg.Repositories {
		if NormalizeGitURL(repo.URL) == normalizedURL {
			repoCopy := repo
			return &repoCopy, true
		}
	}

	return nil, false
}

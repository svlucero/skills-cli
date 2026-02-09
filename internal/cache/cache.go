package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/silvinalucero/skill_cli/internal/skill"
)

const (
	cacheVersion  = "1"
	cacheDirName  = "skill"
	cacheFileName = "skills.json"
)

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetCachePath returns the path to the cache file
func GetCachePath() string {
	return filepath.Join(xdg.CacheHome, cacheDirName, cacheFileName)
}

// Load loads the skills cache from disk
// Returns a new empty cache if file doesn't exist
func Load() (*SkillsCache, error) {
	cachePath := GetCachePath()

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		// Return new empty cache
		return &SkillsCache{
			Version:      cacheVersion,
			LastUpdated:  time.Now(),
			Repositories: make(map[string]*RepoCache),
		}, nil
	}

	// Read cache file
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %w", err)
	}

	// Parse JSON
	var cache SkillsCache
	if err := json.Unmarshal(data, &cache); err != nil {
		// Cache corrupted, return new cache
		return &SkillsCache{
			Version:      cacheVersion,
			LastUpdated:  time.Now(),
			Repositories: make(map[string]*RepoCache),
		}, nil
	}

	// Validate cache version
	if cache.Version != cacheVersion {
		// Version mismatch, return new cache
		return &SkillsCache{
			Version:      cacheVersion,
			LastUpdated:  time.Now(),
			Repositories: make(map[string]*RepoCache),
		}, nil
	}

	return &cache, nil
}

// Save saves the skills cache to disk
func (c *SkillsCache) Save() error {
	cachePath := GetCachePath()

	// Create cache directory if it doesn't exist
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("error creating cache directory: %w", err)
	}

	// Update last updated timestamp
	c.LastUpdated = time.Now()

	// Serialize to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling cache: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("error writing cache file: %w", err)
	}

	return nil
}

// IndexRepository scans a repository and caches all skills metadata
func (c *SkillsCache) IndexRepository(repoName, repoURL, repoPath string) error {
	// Discover skills from repository
	skills, err := skill.DiscoverSkills(repoPath, repoName)
	if err != nil {
		return fmt.Errorf("error discovering skills: %w", err)
	}

	// Convert to cache entries
	entries := make([]SkillEntry, len(skills))
	for i, s := range skills {
		entries[i] = SkillEntry{
			Name:          s.Name,
			Description:   s.Description,
			Version:       s.Version,
			Path:          s.Path,
			HasScripts:    dirExists(filepath.Join(s.Path, "scripts")),
			HasReferences: dirExists(filepath.Join(s.Path, "references")),
			HasAssets:     dirExists(filepath.Join(s.Path, "assets")),
			UpdatedAt:     s.UpdatedAt,
		}
	}

	// Store in cache
	if c.Repositories == nil {
		c.Repositories = make(map[string]*RepoCache)
	}

	c.Repositories[repoName] = &RepoCache{
		URL:         repoURL,
		LastIndexed: time.Now(),
		Skills:      entries,
	}

	return nil
}

// GetSkills retrieves cached skills for a specific repository
func (c *SkillsCache) GetSkills(repoName string) ([]SkillEntry, error) {
	repo, exists := c.Repositories[repoName]
	if !exists {
		return nil, fmt.Errorf("repository '%s' not found in cache", repoName)
	}

	return repo.Skills, nil
}

// GetAllSkills retrieves all cached skills from all repositories
func (c *SkillsCache) GetAllSkills() []SkillEntry {
	var allSkills []SkillEntry

	for _, repo := range c.Repositories {
		allSkills = append(allSkills, repo.Skills...)
	}

	return allSkills
}

// RemoveRepository removes a repository's skills from the cache
func (c *SkillsCache) RemoveRepository(repoName string) {
	delete(c.Repositories, repoName)
}

// HasRepository checks if a repository is cached
func (c *SkillsCache) HasRepository(repoName string) bool {
	_, exists := c.Repositories[repoName]
	return exists
}

// Clear removes all cached data
func (c *SkillsCache) Clear() {
	c.Repositories = make(map[string]*RepoCache)
	c.LastUpdated = time.Now()
}

// GetSkillByName finds a skill by name across all repositories
func (c *SkillsCache) GetSkillByName(skillName string) (*SkillEntry, string, error) {
	for repoName, repo := range c.Repositories {
		for i := range repo.Skills {
			if repo.Skills[i].Name == skillName {
				return &repo.Skills[i], repoName, nil
			}
		}
	}

	return nil, "", fmt.Errorf("skill '%s' not found in cache", skillName)
}

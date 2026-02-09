package cache

import "time"

// SkillsCache represents the complete skills cache structure
type SkillsCache struct {
	Version      string                 `json:"version"`
	LastUpdated  time.Time              `json:"last_updated"`
	Repositories map[string]*RepoCache  `json:"repositories"`
}

// RepoCache represents cached data for a single repository
type RepoCache struct {
	URL         string        `json:"url"`
	LastIndexed time.Time     `json:"last_indexed"`
	Skills      []SkillEntry  `json:"skills"`
}

// SkillEntry represents a cached skill with metadata
type SkillEntry struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Version       string    `json:"version"`
	Path          string    `json:"path"`
	HasScripts    bool      `json:"has_scripts"`
	HasReferences bool      `json:"has_references"`
	HasAssets     bool      `json:"has_assets"`
	UpdatedAt     time.Time `json:"updated_at"`
}

package config

import (
	"strings"
	"testing"
)

func TestValidateRepoName(t *testing.T) {
	tests := []struct {
		name      string
		repoName  string
		wantError bool
	}{
		{"valid lowercase", "myrepo", false},
		{"valid uppercase", "MyRepo", false},
		{"valid with hyphen", "my-repo", false},
		{"valid with underscore", "my_repo", false},
		{"valid alphanumeric", "repo123", false},
		{"empty name", "", true},
		{"with space", "my repo", true},
		{"with special char", "my@repo", true},
		{"with dot", "my.repo", true},
		{"with slash", "my/repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRepoName(tt.repoName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateRepoName(%q) error = %v, wantError %v", tt.repoName, err, tt.wantError)
			}
		})
	}
}

func TestNormalizeGitURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			"https with .git",
			"https://github.com/user/repo.git",
			"https://github.com/user/repo",
		},
		{
			"https without .git",
			"https://github.com/user/repo",
			"https://github.com/user/repo",
		},
		{
			"ssh format converted to https",
			"git@github.com:user/repo.git",
			"https://github.com/user/repo",
		},
		{
			"ssh without .git converted to https",
			"git@github.com:user/repo",
			"https://github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeGitURL(tt.url)
			if result != tt.expected {
				t.Errorf("NormalizeGitURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestGetRepoPathForRepo(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
	}{
		{"simple name", "myrepo"},
		{"with hyphen", "my-repo"},
		{"with underscore", "my_repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := GetRepoPathForRepo(tt.repoName)
			if path == "" {
				t.Error("GetRepoPathForRepo returned empty string")
			}
			if !strings.Contains(path, tt.repoName) {
				t.Errorf("Path %q does not contain repo name %q", path, tt.repoName)
			}
		})
	}
}

func TestGetRepo(t *testing.T) {
	cfg := &Config{
		ActiveRepo: "myrepo",
		Repositories: map[string]Repository{
			"myrepo": {
				Name:     "myrepo",
				URL:      "https://github.com/user/repo.git",
				AuthType: "https",
			},
		},
	}

	tests := []struct {
		name      string
		repoName  string
		wantError bool
	}{
		{"existing repo", "myrepo", false},
		{"non-existing repo", "notfound", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetRepo(cfg, tt.repoName)
			if (err != nil) != tt.wantError {
				t.Errorf("GetRepo(%q) error = %v, wantError %v", tt.repoName, err, tt.wantError)
			}
		})
	}
}

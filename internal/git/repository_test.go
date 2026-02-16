package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/silvinalucero/skill_cli/internal/config"
)

func TestRepoExists(t *testing.T) {
	tmpDir := t.TempDir()

	gitDir := filepath.Join(tmpDir, "repo-with-git")
	if err := os.MkdirAll(filepath.Join(gitDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	emptyDir := filepath.Join(tmpDir, "empty-dir")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"git repo exists", gitDir, true},
		{"empty dir", emptyDir, false},
		{"non-existent path", filepath.Join(tmpDir, "does-not-exist"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RepoExists(tt.path)
			if result != tt.expected {
				t.Errorf("RepoExists(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestDetectAuthType(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected config.AuthType
	}{
		{"https url", "https://github.com/user/repo.git", config.AuthHTTPS},
		{"http url", "http://github.com/user/repo.git", config.AuthHTTPS},
		{"ssh url", "git@github.com:user/repo.git", config.AuthSSH},
		{"ssh protocol", "ssh://git@github.com/user/repo.git", config.AuthSSH},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAuthType(tt.url)
			if result != tt.expected {
				t.Errorf("DetectAuthType(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		{"valid https", "https://github.com/user/repo.git", false},
		{"valid ssh", "git@github.com:user/repo.git", false},
		{"empty url", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateURL(%q) error = %v, wantError %v", tt.url, err, tt.wantError)
			}
		})
	}
}

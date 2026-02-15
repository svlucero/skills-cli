package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Clone clones a repository to the specified destination
func Clone(url, destPath string) error {
	// Check that git is installed
	if err := checkGitInstalled(); err != nil {
		return err
	}

	// Check that parent directory exists
	parentDir := filepath.Dir(destPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("error creating parent directory: %w", err)
	}

	// If destination already exists, return error
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("destination path already exists: %s", destPath)
	}

	// Execute git clone
	cmd := exec.Command("git", "clone", url, destPath)

	// Capture both stdout and stderr to show progress
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Clean up directory if clone failed
		_ = os.RemoveAll(destPath)
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

// Pull updates the local repository with remote changes
func Pull(repoPath string) error {
	// Check that git is installed
	if err := checkGitInstalled(); err != nil {
		return err
	}

	// Check that directory exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	// Check that it's a git repository
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Execute git pull
	cmd := exec.Command("git", "-C", repoPath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull failed: %w", err)
	}

	return nil
}

// GetStatus gets the local repository status
func GetStatus(repoPath string) (string, error) {
	// Check that git is installed
	if err := checkGitInstalled(); err != nil {
		return "", err
	}

	// Check that directory exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	// Execute git status
	cmd := exec.Command("git", "-C", repoPath, "status", "--short")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("git status failed: %w", err)
	}

	return string(output), nil
}

// IsClean checks if the local repository has no pending changes
func IsClean(repoPath string) (bool, error) {
	status, err := GetStatus(repoPath)
	if err != nil {
		return false, err
	}

	// If status is empty, the repo is clean
	return status == "", nil
}

// GetCurrentBranch gets the current branch of the repository
func GetCurrentBranch(repoPath string) (string, error) {
	// Check that git is installed
	if err := checkGitInstalled(); err != nil {
		return "", err
	}

	// Execute git branch --show-current
	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("git branch failed: %w", err)
	}

	return string(output), nil
}

// RepoExists checks if a repository exists at the specified path
func RepoExists(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

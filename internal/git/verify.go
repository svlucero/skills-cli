package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/silvinalucero/skill_cli/internal/errors"
)

// VerificationLevel defines the repository verification level
type VerificationLevel int

const (
	LevelBasic VerificationLevel = iota // ls-remote only
	LevelClone                          // Full test clone
)

// VerifyBasic verifies that the repository is accessible using git ls-remote
func VerifyBasic(url string) error {
	// First check that git is installed
	if err := checkGitInstalled(); err != nil {
		return err
	}

	// Execute git ls-remote
	cmd := exec.Command("git", "ls-remote", "--heads", url)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return parseGitError(string(output), err)
	}

	// If we got here, the repository is accessible
	return nil
}

// VerifyRepository verifies the repository with the specified level
func VerifyRepository(url string, level VerificationLevel) error {
	switch level {
	case LevelBasic:
		return VerifyBasic(url)
	case LevelClone:
		// For now we only implement VerifyBasic
		// VerifyClone can be implemented in the future if needed
		return VerifyBasic(url)
	default:
		return fmt.Errorf("unknown verification level: %d", level)
	}
}

// DetectAuthType detects the authentication type based on the URL
func DetectAuthType(url string) config.AuthType {
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://") {
		return config.AuthSSH
	}
	return config.AuthHTTPS
}

// ValidateURL validates the basic format of the repository URL
func ValidateURL(url string) error {
	if url == "" {
		return errors.ErrInvalidURL
	}

	// Check that it has a valid format (very basic)
	hasHTTPS := strings.HasPrefix(url, "https://")
	hasHTTP := strings.HasPrefix(url, "http://")
	hasSSH := strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://")

	if !hasHTTPS && !hasHTTP && !hasSSH {
		return fmt.Errorf("%w: expected https://, http://, git@, or ssh://", errors.ErrInvalidURL)
	}

	// Check that it contains at least one dot (domain)
	if !strings.Contains(url, ".") && !strings.Contains(url, "localhost") {
		return fmt.Errorf("%w: invalid domain", errors.ErrInvalidURL)
	}

	return nil
}

// checkGitInstalled checks that git is installed and available
func checkGitInstalled() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return errors.ErrGitNotInstalled
	}
	return nil
}

// parseGitError parses git errors and returns a more specific error
func parseGitError(output string, originalErr error) error {
	outputLower := strings.ToLower(output)

	// Authentication error
	if strings.Contains(outputLower, "authentication failed") ||
		strings.Contains(outputLower, "permission denied") ||
		strings.Contains(outputLower, "could not read from remote repository") {
		return errors.ErrAuthenticationFailed
	}

	// Network error
	if strings.Contains(outputLower, "could not resolve host") ||
		strings.Contains(outputLower, "network is unreachable") ||
		strings.Contains(outputLower, "failed to connect") {
		return errors.ErrNetworkUnreachable
	}

	// Repository not found
	if strings.Contains(outputLower, "repository not found") ||
		strings.Contains(outputLower, "does not appear to be a git repository") ||
		strings.Contains(outputLower, "not found") {
		return errors.ErrRepositoryNotFound
	}

	// Generic error with output
	if output != "" {
		return fmt.Errorf("git error: %s", output)
	}

	// Unknown error
	return fmt.Errorf("git command failed: %w", originalErr)
}

package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors for common cases
var (
	ErrConfigNotFound       = errors.New("configuration not found, run 'skill init' first")
	ErrRepositoryNotFound   = errors.New("repository not found or not accessible")
	ErrAuthenticationFailed = errors.New("authentication failed, check credentials")
	ErrInvalidURL           = errors.New("invalid repository URL format")
	ErrNetworkUnreachable   = errors.New("network unreachable")
	ErrGitNotInstalled      = errors.New("git is not installed or not in PATH")
	ErrConfigExists         = errors.New("configuration already exists")
	ErrInvalidConfig        = errors.New("invalid configuration")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrRepoNotFound         = errors.New("repository not found in configuration")
	ErrRepoAlreadyExists    = errors.New("repository with this name already exists")
	ErrInvalidRepoName      = errors.New("invalid repository name")
	ErrNoActiveRepo         = errors.New("no active repository configured")
	ErrCannotRemoveActive   = errors.New("cannot remove active repository, switch to another first")
)

// Wrap wraps an error with additional context message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsConfigNotFound checks if the error is configuration not found
func IsConfigNotFound(err error) bool {
	return errors.Is(err, ErrConfigNotFound)
}

// IsRepositoryNotFound checks if the error is repository not found
func IsRepositoryNotFound(err error) bool {
	return errors.Is(err, ErrRepositoryNotFound)
}

// IsAuthenticationFailed checks if the error is authentication failed
func IsAuthenticationFailed(err error) bool {
	return errors.Is(err, ErrAuthenticationFailed)
}

// IsInvalidURL checks if the error is invalid URL
func IsInvalidURL(err error) bool {
	return errors.Is(err, ErrInvalidURL)
}

// IsNetworkError checks if the error is network-related
func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetworkUnreachable)
}

// IsGitNotInstalled checks if the error is git not installed
func IsGitNotInstalled(err error) bool {
	return errors.Is(err, ErrGitNotInstalled)
}

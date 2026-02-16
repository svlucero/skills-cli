package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/silvinalucero/skill_cli/internal/config"
	"github.com/spf13/cobra"
)

func setupTestConfig(t *testing.T) (string, func()) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `active_repository: ""
repositories: {}
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	oldConfigPath := os.Getenv("SKILL_CONFIG")
	os.Setenv("SKILL_CONFIG", configPath)

	cleanup := func() {
		os.Setenv("SKILL_CONFIG", oldConfigPath)
	}

	return configPath, cleanup
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestValidateRepoName(t *testing.T) {
	tests := []struct {
		name      string
		repoName  string
		wantError bool
	}{
		{"valid name", "myrepo", false},
		{"with hyphen", "my-repo", false},
		{"with underscore", "my_repo", false},
		{"empty name", "", true},
		{"with space", "my repo", true},
		{"with special char", "my@repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateRepoName(tt.repoName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateRepoName(%q) error = %v, wantError %v", tt.repoName, err, tt.wantError)
			}
		})
	}
}

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd is nil")
	}

	if !strings.HasPrefix(rootCmd.Use, "skills") {
		t.Errorf("rootCmd.Use = %q, want to start with %q", rootCmd.Use, "skills")
	}
}

func TestSubcommandsExist(t *testing.T) {
	expectedCommands := []string{
		"repository",
		"list",
		"install",
		"uninstall",
		"config",
		"share",
	}

	commands := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		commands[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commands[expected] {
			t.Errorf("Expected command %q not found", expected)
		}
	}
}

func TestRepositorySubcommands(t *testing.T) {
	var repoCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "repository" {
			repoCmd = cmd
			break
		}
	}

	if repoCmd == nil {
		t.Fatal("repository command not found")
	}

	expectedSubcommands := []string{
		"add",
		"list",
		"remove",
		"set-current",
		"update",
	}

	subcommands := make(map[string]bool)
	for _, cmd := range repoCmd.Commands() {
		subcommands[cmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("Expected subcommand %q not found in repository", expected)
		}
	}
}

func TestConfigSubcommands(t *testing.T) {
	var configCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "config" {
			configCmd = cmd
			break
		}
	}

	if configCmd == nil {
		t.Fatal("config command not found")
	}

	expectedSubcommands := []string{
		"show",
		"verify",
	}

	subcommands := make(map[string]bool)
	for _, cmd := range configCmd.Commands() {
		subcommands[cmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("Expected subcommand %q not found in config", expected)
		}
	}
}

func TestHelpMessages(t *testing.T) {
	commands := []*cobra.Command{
		rootCmd,
		listCmd,
		installCmd,
		uninstallCmd,
		shareCmd,
	}

	for _, cmd := range commands {
		t.Run(cmd.Name(), func(t *testing.T) {
			if cmd.Short == "" {
				t.Errorf("Command %q has no short description", cmd.Name())
			}

			if cmd.Long == "" && cmd.Name() != "skills" {
				t.Logf("Warning: Command %q has no long description", cmd.Name())
			}
		})
	}
}

func TestFlagDefinitions(t *testing.T) {
	tests := []struct {
		cmd      *cobra.Command
		flagName string
	}{
		{listCmd, "all"},
		{listCmd, "repo"},
		{listCmd, "compact"},
		{listCmd, "provider"},
		{installCmd, "repo"},
		{installCmd, "provider"},
		{installCmd, "force"},
		{uninstallCmd, "provider"},
		{uninstallCmd, "force"},
	}

	for _, tt := range tests {
		t.Run(tt.cmd.Name()+"/"+tt.flagName, func(t *testing.T) {
			flag := tt.cmd.Flags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %q not found in %q command", tt.flagName, tt.cmd.Name())
			}
		})
	}
}

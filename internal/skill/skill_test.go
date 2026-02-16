package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSkillMarkdown(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		expected    *SkillConfig
	}{
		{
			name: "valid frontmatter",
			content: `---
name: test-skill
description: A test skill
version: 1.0.0
---

# Test Skill

This is a test skill.
`,
			expectError: false,
			expected: &SkillConfig{
				Name:        "test-skill",
				Description: "A test skill",
				Version:     "1.0.0",
			},
		},
		{
			name: "missing frontmatter",
			content: `# Test Skill

This is a test skill without frontmatter.
`,
			expectError: true,
			expected:    nil,
		},
		{
			name: "invalid yaml",
			content: `---
name: test-skill
description: A test skill
version: 1.0.0
invalid yaml here
---
`,
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			skillFile := filepath.Join(tmpDir, "SKILL.md")

			if err := os.WriteFile(skillFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			cfg, err := ParseSkillMarkdown(skillFile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cfg.Name != tt.expected.Name {
				t.Errorf("Name = %q, want %q", cfg.Name, tt.expected.Name)
			}

			if cfg.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", cfg.Description, tt.expected.Description)
			}

			if cfg.Version != tt.expected.Version {
				t.Errorf("Version = %q, want %q", cfg.Version, tt.expected.Version)
			}
		})
	}
}

func TestDiscoverSkills(t *testing.T) {
	tmpDir := t.TempDir()

	skill1Dir := filepath.Join(tmpDir, "skill1")
	skill2Dir := filepath.Join(tmpDir, "skill2")
	invalidDir := filepath.Join(tmpDir, "not-a-skill")

	if err := os.MkdirAll(skill1Dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(skill2Dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(invalidDir, 0755); err != nil {
		t.Fatal(err)
	}

	skill1Content := `---
name: skill1
description: First skill
version: 1.0.0
---

# Skill 1
`

	skill2Content := `---
name: skill2
description: Second skill
version: 2.0.0
---

# Skill 2
`

	if err := os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2Content), 0644); err != nil {
		t.Fatal(err)
	}

	skills, err := DiscoverSkills(tmpDir, "test-repo", "")
	if err != nil {
		t.Fatalf("DiscoverSkills failed: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(skills))
	}

	found := make(map[string]bool)
	for _, s := range skills {
		found[s.Name] = true
	}

	if !found["skill1"] || !found["skill2"] {
		t.Error("Expected to find skill1 and skill2")
	}
}

func TestSkillStatus(t *testing.T) {
	skill := Skill{
		Name:   "test-skill",
		Status: StatusAvailable,
	}

	if skill.Status != StatusAvailable {
		t.Errorf("Status = %q, want %q", skill.Status, StatusAvailable)
	}

	skill.Status = StatusInstalled
	if skill.Status != StatusInstalled {
		t.Errorf("Status = %q, want %q", skill.Status, StatusInstalled)
	}
}

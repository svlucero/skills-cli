package cli

import (
	"strings"

	"github.com/fatih/color"
)

// HelpBuilder provides a fluent interface for building command help text
type HelpBuilder struct {
	parts []string
	cyan  func(a ...interface{}) string
	green func(a ...interface{}) string
	dim   func(a ...interface{}) string
}

// NewHelpBuilder creates a new help text builder
func NewHelpBuilder() *HelpBuilder {
	return &HelpBuilder{
		cyan:  color.New(color.FgCyan, color.Bold).SprintFunc(),
		green: color.New(color.FgGreen).SprintFunc(),
		dim:   color.New(color.Faint).SprintFunc(),
	}
}

// Description adds the main description of the command
func (h *HelpBuilder) Description(text string) *HelpBuilder {
	h.parts = append(h.parts, text)
	return h
}

// Section adds a titled section (e.g., "USAGE:", "FLAGS:")
func (h *HelpBuilder) Section(title string) *HelpBuilder {
	h.parts = append(h.parts, "", h.cyan(title))
	return h
}

// Text adds plain text content
func (h *HelpBuilder) Text(text string) *HelpBuilder {
	h.parts = append(h.parts, text)
	return h
}

// Item adds a list item with label and description
// Example: Item("--repo <name>", "Repository to install from")
func (h *HelpBuilder) Item(label, description string) *HelpBuilder {
	h.parts = append(h.parts, "  "+h.green(label)+"  "+h.dim(description))
	return h
}

// SubItem adds an indented sub-item (for nested lists)
func (h *HelpBuilder) SubItem(label, description string) *HelpBuilder {
	h.parts = append(h.parts, "    "+h.green(label)+"  "+h.dim(description))
	return h
}

// BulletList adds a bulleted list of items
func (h *HelpBuilder) BulletList(items []string) *HelpBuilder {
	for _, item := range items {
		h.parts = append(h.parts, "  - "+item)
	}
	return h
}

// Example adds a command example
func (h *HelpBuilder) Example(command, description string) *HelpBuilder {
	if description != "" {
		h.parts = append(h.parts, "  "+h.green(command)+"  "+h.dim(description))
	} else {
		h.parts = append(h.parts, "  "+h.green(command))
	}
	return h
}

// Examples adds multiple examples at once
func (h *HelpBuilder) Examples(examples map[string]string) *HelpBuilder {
	h.Section("EXAMPLES:")
	for cmd, desc := range examples {
		h.Example(cmd, desc)
	}
	return h
}

// EmptyLine adds a blank line for spacing
func (h *HelpBuilder) EmptyLine() *HelpBuilder {
	h.parts = append(h.parts, "")
	return h
}

// Build returns the final formatted help text
func (h *HelpBuilder) Build() string {
	return strings.Join(h.parts, "\n")
}

// QuickHelp is a convenience method for simple help text
// Usage: QuickHelp("Description", map[string]string{"cmd": "desc"}, []string{"example"})
func QuickHelp(description string, sections map[string][]string, examples []string) string {
	builder := NewHelpBuilder().Description(description)

	// Add sections
	for title, items := range sections {
		builder.Section(title)
		builder.BulletList(items)
	}

	// Add examples
	if len(examples) > 0 {
		builder.Section("EXAMPLES:")
		for _, ex := range examples {
			builder.Example(ex, "")
		}
	}

	return builder.Build()
}

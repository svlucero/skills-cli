# HelpBuilder

A fluent interface for building command help text without manual string formatting and placeholder counting.

## Problem

Previously, help text was built using `fmt.Sprintf` with many `%s` placeholders:

```go
func getOldHelp() string {
    cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()

    return fmt.Sprintf(`%s

%s
  %s   %s
  %s   %s

%s
  %s
  %s`,
        "Description here",
        cyan("FLAGS:"),
        green("--flag1"), dim("Description 1"),
        green("--flag2"), dim("Description 2"),
        cyan("EXAMPLES:"),
        green("skill command example1"),
        green("skill command example2"),
    )
}
```

**Issues:**
- ❌ Easy to miscount placeholders
- ❌ Hard to read and maintain
- ❌ Compiler doesn't catch placeholder count errors
- ❌ Difficult to add/remove sections

## Solution: HelpBuilder

```go
func getNewHelp() string {
    return NewHelpBuilder().
        Description("Description here").
        Section("FLAGS:").
        Item("--flag1", "Description 1").
        Item("--flag2", "Description 2").
        Section("EXAMPLES:").
        Example("skill command example1", "").
        Example("skill command example2", "").
        Build()
}
```

**Benefits:**
- ✅ No placeholder counting
- ✅ Readable and declarative
- ✅ Type-safe
- ✅ Easy to modify

## API Reference

### Basic Methods

#### `NewHelpBuilder() *HelpBuilder`
Creates a new help builder instance.

```go
builder := NewHelpBuilder()
```

#### `Description(text string) *HelpBuilder`
Sets the main description of the command.

```go
builder.Description("Add a new repository to the configuration.")
```

#### `Section(title string) *HelpBuilder`
Adds a titled section (e.g., "FLAGS:", "EXAMPLES:").

```go
builder.Section("FLAGS:")
```

#### `Item(label, description string) *HelpBuilder`
Adds a flag or parameter with its description.

```go
builder.Item("--force", "Overwrite existing configuration")
builder.Item(yellow("<name>"), "Repository name (required)")
```

#### `SubItem(label, description string) *HelpBuilder`
Adds an indented sub-item for nested lists.

```go
builder.
    Item("repository <subcommand>", "Manage repositories").
    SubItem("add", "Add a new repository").
    SubItem("remove", "Remove a repository")
```

#### `BulletList(items []string) *HelpBuilder`
Adds a bulleted list.

```go
builder.BulletList([]string{
    "Validates repository URL",
    "Clones to local directory",
    "Saves configuration",
})
```

#### `Example(command, description string) *HelpBuilder`
Adds a command example with optional description.

```go
builder.Example("skill repository add myrepo https://...", "")
builder.Example("skill install skill-name", "# Install a skill")
```

#### `EmptyLine() *HelpBuilder`
Adds a blank line for spacing.

```go
builder.EmptyLine()
```

#### `Build() string`
Returns the final formatted help text.

```go
helpText := builder.Build()
```

## Usage Examples

### Simple Command

```go
func getSimpleHelp() string {
    return NewHelpBuilder().
        Description("List all available skills.").
        Section("EXAMPLES:").
        Example("skill list", "").
        Example("skill list --installed", "").
        Build()
}
```

### Command with Flags

```go
func getCommandWithFlags() string {
    return NewHelpBuilder().
        Description("Install a skill to a provider.").
        Section("FLAGS:").
        Item("--provider <name>", "Target provider (claude, cursor)").
        Item("--force", "Overwrite existing installation").
        Section("EXAMPLES:").
        Example("skill install my-skill", "").
        Example("skill install my-skill --provider cursor", "").
        Build()
}
```

### Command with Parameters and Behavior

```go
func getComplexHelp() string {
    yellow := color.New(color.FgYellow).SprintFunc()

    return NewHelpBuilder().
        Description("Add a new Git repository.").
        Section("PARAMETERS:").
        Item(yellow("<name>"), "Repository identifier (required)").
        Item(yellow("<url>"), "Git repository URL (required)").
        Section("FLAGS:").
        Item("-f, --force", "Overwrite existing repository").
        Item("--skip-verify", "Skip verification").
        Section("BEHAVIOR:").
        BulletList([]string{
            "Validates URL format",
            "Verifies repository access",
            "Clones to local directory",
            "Saves configuration",
        }).
        Section("EXAMPLES:").
        Example("skill repository add myrepo https://github.com/org/repo.git", "").
        Example("skill repository add myrepo https://... --force", "").
        Build()
}
```

### Using QuickHelp Helper

For very simple help text:

```go
func getQuickHelp() string {
    return QuickHelp(
        "Shows configuration details",
        map[string][]string{
            "OUTPUT INCLUDES:": {
                "Repository URL",
                "Local path",
                "Active status",
            },
        },
        []string{
            "skill config show",
        },
    )
}
```

## Migration Guide

### Before (Old Style)

```go
func getOldHelp() string {
    cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()
    dim := color.New(color.Faint).SprintFunc()

    return fmt.Sprintf(`%s

%s
  %s              %s
  %s    %s

%s
  %s         %s
  %s       %s

%s
  %s
  %s`,
        "Add a new repository.",
        cyan("PARAMETERS:"),
        yellow("<name>"), dim("Name (required)"),
        yellow("<url>"), dim("URL (required)"),
        cyan("FLAGS:"),
        green("-f, --force"), dim("Force overwrite"),
        green("--skip-verify"), dim("Skip verification"),
        cyan("EXAMPLES:"),
        green("skill repository add myrepo https://..."),
        green("skill repository add myrepo https://... --force"),
    )
}
```

### After (HelpBuilder)

```go
func getNewHelp() string {
    yellow := color.New(color.FgYellow).SprintFunc()

    return NewHelpBuilder().
        Description("Add a new repository.").
        Section("PARAMETERS:").
        Item(yellow("<name>"), "Name (required)").
        Item(yellow("<url>"), "URL (required)").
        Section("FLAGS:").
        Item("-f, --force", "Force overwrite").
        Item("--skip-verify", "Skip verification").
        Section("EXAMPLES:").
        Example("skill repository add myrepo https://...", "").
        Example("skill repository add myrepo https://... --force", "").
        Build()
}
```

**Changes:**
1. Replace `fmt.Sprintf` with `NewHelpBuilder()`
2. Use method chaining for each section
3. Remove manual `%s` placeholders
4. Only create color functions that are actually used
5. End with `.Build()`

## Best Practices

1. **Only create color functions you need**
   ```go
   // If you only use yellow for parameters:
   yellow := color.New(color.FgYellow).SprintFunc()
   ```

2. **Chain methods for readability**
   ```go
   // Good - one method per line
   builder.
       Description("...").
       Section("FLAGS:").
       Item("--flag", "description")

   // Avoid - hard to read
   builder.Description("...").Section("FLAGS:").Item("--flag", "description")
   ```

3. **Use consistent section names**
   - PARAMETERS:
   - FLAGS:
   - BEHAVIOR:
   - EXAMPLES:
   - OUTPUT INCLUDES:

4. **Empty descriptions in examples**
   ```go
   // When examples are self-explanatory:
   Example("skill list", "")

   // When examples need clarification:
   Example("skill list --installed", "# Show only installed skills")
   ```

## Color Functions

The builder automatically uses these colors:
- **Section titles**: Cyan bold
- **Items/Examples**: Green
- **Descriptions**: Dim (gray)

For parameters, you can create a yellow function:
```go
yellow := color.New(color.FgYellow).SprintFunc()
builder.Item(yellow("<name>"), "Description")
```

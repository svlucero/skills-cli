# skill CLI

Skills manager from Git repositories.

## Description

`skill` is a command-line tool written in Go that allows you to manage "skills" stored in Git repositories. Skills are reusable code snippets, scripts, and documentation organized in a standardized directory structure with a `SKILL.md` file containing metadata and instructions.

Each skill is a self-contained directory with:
- **SKILL.md**: Documentation and metadata (required)
- **scripts/**: Executable code (optional)
- **references/**: Additional documentation (optional)
- **assets/**: Templates and resources (optional)

## Features

- 📦 **Multiple repositories**: Manage skills from different Git repositories
- 🔍 **Automatic discovery**: Scans repositories for skills (any directory with `SKILL.md`)
- 🚀 **Easy installation**: Install skills to Claude Desktop or Cursor with one command
- 🤝 **Sharing**: Fork repositories and create PRs to contribute skills
- 🔐 **Flexible authentication**: Supports both HTTPS and SSH for Git operations
- ✅ **Validation**: Automatically validates skill structure before sharing
- 🎨 **Colored output**: Beautiful terminal UI for better user experience
- ⚙️ **XDG compliant**: Configuration follows XDG Base Directory Specification

## Installation

### From source

```bash
# Clone the repository
git clone https://github.com/silvinalucero/skill_cli.git
cd skill_cli

# Build the binary
make build

# Install to $GOPATH/bin
make install
```

### Requirements

- Go 1.23 or higher
- Git installed and available in PATH

## Usage

### Initialize with a repository

```bash
# With HTTPS repository
skill init https://github.com/org/skills-repo.git

# With SSH repository
skill init git@github.com:org/skills-repo.git

# Force overwrite existing configuration
skill init https://github.com/org/skills-repo.git --force

# Skip repository verification (not recommended)
skill init https://github.com/org/skills-repo.git --skip-verify
```

### View current configuration

```bash
skill config show
```

Shows:
- Repository URL
- Local path of cloned repository
- Authentication type (https/ssh)
- Last verification
- Local repository status

### Change repository

```bash
skill config set-repo https://github.com/org/new-repo.git

# Without verifying the new repository
skill config set-repo https://github.com/org/new-repo.git --no-verify
```

### Verify repository access

```bash
skill config verify
```

Verifies that the configured repository is accessible with your current credentials.

## File Structure

### Configuration

Configuration is stored in:
- **Config:** `~/.config/skill/config.yaml` (XDG_CONFIG_HOME)
- **Data:** `~/.local/share/skill/repo` (XDG_DATA_HOME)

### Configuration format (config.yaml)

```yaml
version: "1"
repository:
  url: "https://github.com/org/skills-repo.git"
  local_path: "/Users/username/.local/share/skill/repo"
  last_verified: "2026-02-07T10:30:00Z"
  auth_type: "https"
```

### Skills repository structure

Skills are stored as individual directories in the root of the repository. Each skill must follow this structure:

#### Required Structure

```
my-skill/
├── SKILL.md          # Required: Skill documentation with frontmatter metadata
├── scripts/          # Optional: Executable code (bash, python, etc.)
├── references/       # Optional: Additional documentation, guides
└── assets/           # Optional: Templates, resources, configuration files
```

#### SKILL.md Format

The `SKILL.md` file must start with YAML frontmatter containing metadata:

```markdown
---
name: my-skill
description: Brief description of what this skill does
version: 1.0.0
---

# My Skill

Detailed instructions and documentation for the skill.

## Usage

How to use this skill...

## Examples

Examples of using the skill...
```

#### Complete Repository Example

```
skills-repo/
├── deploy-app/
│   ├── SKILL.md              # Deployment skill documentation
│   ├── scripts/
│   │   ├── deploy.sh         # Deployment script
│   │   └── rollback.sh       # Rollback script
│   ├── references/
│   │   └── deployment-guide.md
│   └── assets/
│       └── config.template.yml
├── explain-code/
│   ├── SKILL.md              # Code explanation skill
│   └── references/
│       └── examples.md
└── code-review/
    ├── SKILL.md              # Code review skill
    ├── scripts/
    │   └── check-style.py
    └── assets/
        └── review-checklist.md
```

**Note:** The CLI automatically discovers any directory containing a `SKILL.md` file as a skill. Directories without `SKILL.md` are ignored.

## Available Commands

### Main Commands

| Command | Description |
|---------|-------------|
| `skill init <repo-url>` | Initialize configuration with a repository |
| `skill config show` | Show current configuration |
| `skill config set-repo <url>` | Change configured repository |
| `skill config verify` | Verify repository access |
| `skill list` | List available or installed skills |
| `skill install <name>` | Install a skill to a provider (Claude, Cursor) |
| `skill share --path <path>` | Share a skill by forking and creating a PR |
| `skill --help` | Show general help |
| `skill --version` | Show CLI version |

### Share a skill

The `share` command allows you to contribute skills to a repository by automatically forking, branching, and creating a pull request.

```bash
# Share a skill to the current active repository
skill share --path ./my-skill

# Share a skill to a specific repository
skill share --path ./my-skill --repo https://github.com/org/skills-repo.git

# Share a skill with SSH URL
skill share --path ~/skills/deploy-app --repo git@github.com:org/skills.git
```

#### How it works

The share command:
1. **Validates** the skill structure (requires `SKILL.md`)
2. **Warns** if optional directories (`scripts/`, `references/`, `assets/`) are missing
3. **Forks** the target repository using `gh` CLI
4. **Creates** a new branch (`add-skill-<name>`)
5. **Copies** the skill folder to the repository
6. **Commits** and pushes changes
7. **Creates** a pull request against main
8. **Returns** the PR URL

#### Skill structure validation

Your skill directory must contain:
- ✅ **Required**: `SKILL.md` with frontmatter metadata (see [Skills repository structure](#skills-repository-structure))
- ⚠️ **Optional**: `scripts/`, `references/`, `assets/` directories (warnings shown if missing)

Example of a valid skill:
```
my-skill/
├── SKILL.md          # Required
├── scripts/          # Optional but recommended
│   └── setup.sh
├── references/       # Optional
│   └── guide.md
└── assets/           # Optional
    └── template.yml
```

**Requirements:**
- `gh` CLI must be installed and authenticated (`gh auth login`)
- You must have permissions to fork the target repository
- Your skill must have a valid `SKILL.md` file with proper frontmatter

### Coming Soon

- `skill list` - List available skills in the repository
- `skill install <name>` - Install a skill locally
- `skill uninstall <name>` - Uninstall a skill
- `skill update` - Update local repository (git pull)
- `skill sync` - Sync local changes with remote
- `skill create <name>` - Create a new skill

## Development

### Build

```bash
make build
```

The binary is generated in `bin/skill`.

### Run tests

```bash
make test
```

### Run without building

```bash
make run ARGS="init https://github.com/org/repo.git"
```

### Clean generated files

```bash
make clean
```

## Architecture

```
skill_cli/
├── cmd/skill/          # Entry point
├── internal/
│   ├── cli/           # Cobra commands
│   ├── config/        # Configuration management
│   ├── git/           # Git operations
│   └── errors/        # Custom errors
├── go.mod
├── Makefile
└── README.md
```

### Technologies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [XDG](https://github.com/adrg/xdg) - XDG compliant paths
- [Color](https://github.com/fatih/color) - Colored output

## Troubleshooting

### Error: "git is not installed or not in PATH"

Make sure you have Git installed:

```bash
git --version
```

### Error: "authentication failed"

**For SSH repositories:**
```bash
# Verify your SSH connection
ssh -T git@github.com
```

**For private HTTPS repositories:**
- Configure a personal access token
- Ensure you have read permissions on the repository

### Error: "configuration already exists"

If you want to overwrite the existing configuration:

```bash
skill init <repo-url> --force
```

## License

MIT

## Author

Silvina Lucero
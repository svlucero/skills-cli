# Skills CLI

> A powerful command-line tool for managing skills stored in Git repositories

[![GitHub release](https://img.shields.io/github/v/release/svlucero/skills-cli?style=for-the-badge&logo=github)](https://github.com/svlucero/skills-cli/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/svlucero/skills-cli/total?style=for-the-badge&logo=github)](https://github.com/svlucero/skills-cli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

[![Test](https://github.com/svlucero/skills-cli/actions/workflows/test.yml/badge.svg)](https://github.com/svlucero/skills-cli/actions/workflows/test.yml)
[![Release](https://github.com/svlucero/skills-cli/actions/workflows/release.yml/badge.svg)](https://github.com/svlucero/skills-cli/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)

## 📖 Table of Contents

- [Overview](#overview)
- [Demo](#demo)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Repository Management](#repository-management)
  - [Skill Management](#skill-management)
  - [Configuration](#configuration)
- [Skill Structure](#skill-structure)
- [Interactive Mode](#interactive-mode)
- [Contributing](#contributing)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Overview

**Skills CLI** is a modern command-line tool written in Go that helps you manage reusable code snippets, scripts, and documentation (called "skills") stored in Git repositories. It features an interactive UI, supports multiple repositories, and integrates seamlessly with Claude Desktop and Cursor.

Skills are self-contained directories with a standardized structure, making them easy to share, install, and manage across projects and teams.

> 💡 **Ready to use!** Pre-built binaries are available for Linux, macOS, and Windows. No compilation needed - just [download and run](https://github.com/svlucero/skills-cli/releases/latest)!

## Demo

<!-- Replace this comment with your demo GIF -->
<!-- ![Demo](assets/demo.gif) -->

**Demo coming soon!** - A GIF demonstration of the interactive features will be added here.

## Features

✨ **Interactive UI** - Navigate through skills and repositories with arrow keys and real-time search
📦 **Multiple Repositories** - Manage skills from different Git repositories simultaneously
🔍 **Auto-discovery** - Automatically finds all skills (directories with `SKILL.md`)
🚀 **Easy Installation** - Install skills to Claude Desktop or Cursor with one command
🗑️ **Uninstall Support** - Remove installed skills interactively or directly
🎯 **Repository Switching** - Quickly switch between different skill repositories
🔐 **Flexible Auth** - Supports both HTTPS and SSH for Git operations
✅ **Validation** - Validates repository access before adding
🎨 **Beautiful UI** - Colored output and intuitive prompts
⚙️ **XDG Compliant** - Follows XDG Base Directory Specification

## Installation

### 📦 Download Pre-built Binary (Recommended)

The easiest way to install Skills CLI is to download a pre-built binary from the [releases page](https://github.com/svlucero/skills-cli/releases/latest).

#### macOS

```bash
# Intel Mac (x86_64)
curl -L https://github.com/svlucero/skills-cli/releases/latest/download/skills_Darwin_x86_64 -o skills
sudo mv skills /usr/local/bin/

# Apple Silicon (arm64)
curl -L https://github.com/svlucero/skills-cli/releases/latest/download/skills_Darwin_arm64 -o skills
sudo mv skills /usr/local/bin/
```

#### Linux

```bash
# x86_64
curl -L https://github.com/svlucero/skills-cli/releases/latest/download/skills_Linux_x86_64.tar.gz | tar xz
sudo mv skills /usr/local/bin/

# arm64
curl -L https://github.com/svlucero/skills-cli/releases/latest/download/skills_Linux_arm64.tar.gz | tar xz
sudo mv skills /usr/local/bin/
```

#### Windows

Download the `.zip` file for your architecture from the [releases page](https://github.com/svlucero/skills-cli/releases/latest) and extract it to a directory in your PATH.

#### Verify Installation

```bash
skills --version
```

---

### 🛠️ Alternative Installation Methods

<details>
<summary>Install using Go (requires Go 1.24+)</summary>

```bash
go install github.com/svlucero/skills-cli/cmd/skill@latest
```

Note: The binary will be installed as `skill` (not `skills`) in your `$GOPATH/bin`.
</details>

<details>
<summary>Build from source (for development)</summary>

```bash
# Clone the repository
git clone https://github.com/svlucero/skills-cli.git
cd skills-cli

# Build the binary
make build

# Install to $GOPATH/bin
make install
```

**Prerequisites for building from source:**
- Go 1.24 or higher
- Git
- Make

</details>

---

### 📋 Requirements

To use Skills CLI, you need:

- **Git** - For cloning and managing skill repositories
- **gh CLI** (optional) - Required only for the `share` command to create PRs

## Quick Start

```bash
# 1. Download and install (see Installation section above)
curl -L https://github.com/svlucero/skills-cli/releases/latest/download/skills_*_Darwin_arm64.tar.gz | tar xz
sudo mv skills /usr/local/bin/

# 2. Verify installation
skills --version

# 3. Add your first repository
skills repository add myrepo https://github.com/org/skills-repo.git

# 4. List available skills (interactive mode)
skills list

# 5. Install a skill (interactive mode)
skills install

# 6. View installed skills
skills list --installed --provider claude
```

## Usage

### Repository Management

#### Add a Repository

```bash
# Add with HTTPS
skills repository add myrepo https://github.com/org/skills-repo.git

# Add with SSH
skills repository add myrepo git@github.com:org/skills-repo.git

# Add and set as active
skills repository add myrepo https://github.com/org/repo.git --set-current

# Add with custom skills path (if skills are in a subdirectory)
skills repository add myrepo https://github.com/org/repo.git --skills-path skills

# Force overwrite existing
skills repository add myrepo https://github.com/org/repo.git --force
```

#### List Repositories (Interactive)

```bash
# Interactive selection
skills repository list
```

Navigate with ↑/↓ arrows, search with `/`, press Enter to select.

#### Set Current Repository (Interactive)

```bash
# Interactive selection
skills repository set-current

# Direct selection
skills repository set-current myrepo
```

#### Remove Repository (Interactive)

```bash
# Interactive selection
skills repository remove

# Direct removal
skills repository remove oldrepo

# Remove from config but keep local files
skills repository remove oldrepo --keep-local
```

#### Update Repository

```bash
# Update skills path for a repository
skills repository update myrepo --skills-path skills
```

### Skill Management

#### List Skills (Interactive)

```bash
# List from active repository (interactive)
skills list

# List from specific repository
skills list --repo myrepo

# List from all repositories
skills list --all

# Compact format (non-interactive)
skills list --compact

# Skip git pull before listing
skills list --no-update

# List installed skills
skills list --installed --provider claude
skills list --installed --provider cursor
```

#### Install Skills (Interactive)

```bash
# Interactive selection from active repository
skills install

# Interactive selection from specific repository
skills install --repo myrepo

# Direct install
skills install explain-code

# Install to specific provider
skills install explain-code --provider cursor

# Install from specific repository
skills install explain-code --repo myrepo

# Force reinstall
skills install explain-code --force
```

#### Uninstall Skills (Interactive)

```bash
# Interactive selection
skills uninstall

# Interactive from specific provider
skills uninstall --provider cursor

# Direct uninstall
skills uninstall explain-code

# Uninstall without confirmation
skills uninstall explain-code --force
```

### Configuration

#### Show Configuration

```bash
skills config show
```

Displays:
- Active repository
- All configured repositories
- Repository details (URL, auth type, skills path, status)
- Configuration file location

#### Verify Repository

```bash
skills config verify
```

Verifies that the active repository is accessible with current credentials.

## Skill Structure

### Repository Structure

Skills are stored as individual directories in a Git repository. Each skill must contain a `SKILL.md` file.

```
skills-repo/
├── deploy-app/
│   ├── SKILL.md              # Required: Skill documentation with metadata
│   ├── scripts/              # Optional: Executable scripts
│   │   ├── deploy.sh
│   │   └── rollback.sh
│   ├── references/           # Optional: Additional documentation
│   │   └── deployment-guide.md
│   └── assets/               # Optional: Templates and resources
│       └── config.template.yml
├── explain-code/
│   ├── SKILL.md
│   └── references/
│       └── examples.md
└── code-review/
    ├── SKILL.md
    ├── scripts/
    │   └── check-style.py
    └── assets/
        └── review-checklist.md
```

### SKILL.md Format

The `SKILL.md` file must start with YAML frontmatter:

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

### Custom Skills Path

If your skills are in a subdirectory (e.g., `skills/` or `examples/claude-skills/`), specify it when adding the repository:

```bash
skills repository add myrepo https://github.com/org/repo.git --skills-path skills
```

## Interactive Mode

Most commands support interactive mode with:

- **Arrow keys (↑/↓)** - Navigate through items
- **Search (/)** - Filter items by name or description
- **Enter** - Select current item
- **Ctrl+C** - Cancel operation

Interactive mode features:
- ✓ Real-time preview of item details
- ✓ Visual indicators for status (installed, active, etc.)
- ✓ Sorted and searchable lists
- ✓ Colored output for better readability

## Contributing

Contributions are welcome! Here's how you can help:

### Reporting Issues

- Use the [issue tracker](https://github.com/svlucero/skills-cli/issues)
- Check if the issue already exists
- Provide detailed information (OS, Go version, steps to reproduce)

### Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

## Development

### Building

```bash
make build
```

Binary will be generated in `bin/skills`.

### Running Tests

```bash
make test
```

### Running Without Building

```bash
make run ARGS="list"
```

### Cleaning

```bash
make clean
```

### Project Structure

```
skills-cli/
├── cmd/skill/              # Entry point
├── internal/
│   ├── cli/               # Cobra commands
│   ├── config/            # Configuration management
│   ├── git/               # Git operations
│   ├── skill/             # Skill operations
│   └── errors/            # Custom errors
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Technologies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [XDG](https://github.com/adrg/xdg) - XDG Base Directory support
- [Color](https://github.com/fatih/color) - Terminal colors

## Troubleshooting

### Git not found

```bash
# Verify Git is installed
git --version
```

If not installed, [install Git](https://git-scm.com/downloads).

### Authentication Failed (SSH)

```bash
# Test SSH connection
ssh -T git@github.com

# If it fails, set up SSH keys
ssh-keygen -t ed25519 -C "your_email@example.com"
ssh-add ~/.ssh/id_ed25519
```

Add the public key to GitHub: Settings → SSH and GPG keys

### Authentication Failed (HTTPS)

For private repositories, you may need to configure credentials:

- **GitHub**: Use a [Personal Access Token](https://github.com/settings/tokens)
- Store credentials: `git config --global credential.helper store`

### Repository Not Found

Verify:
1. The repository URL is correct
2. You have read access to the repository
3. The repository exists and is not deleted

### Configuration Already Exists

To overwrite existing configuration:

```bash
skills repository add myrepo https://... --force
```

## File Locations

Configuration follows the XDG Base Directory Specification:

- **Config**: `~/.config/skill/config.yaml` (or `$XDG_CONFIG_HOME/skill/config.yaml`)
- **Data**: `~/.local/share/skill/repos/` (or `$XDG_DATA_HOME/skill/repos/`)
- **Claude Skills**: `~/.claude/skills/`
- **Cursor Skills**: `~/.cursor/skills/`

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Created by **Silvina Lucero**

- GitHub: [@svlucero](https://github.com/svlucero)

---

**⭐ If you find this project useful, please consider giving it a star!**

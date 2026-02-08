# skill CLI

Skills manager from Git repositories.

## Description

`skill` is a command-line tool written in Go that allows you to manage "skills" stored in Git repositories. Skills are structured as individual directories with multiple files (config.yaml, scripts, README.md, etc.) and are stored in a shared Git repository.

## Features

- Initialization with remote Git repository (HTTPS or SSH)
- Automatic repository access verification
- XDG Base Directory Specification compliant configuration
- Local cloning of skills repository
- Configuration management with simple commands
- Colored output for better user experience

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

```
skills-repo/
└── skills/
    ├── my-skill/
    │   ├── config.yaml
    │   ├── script.sh
    │   └── README.md
    └── another-skill/
        ├── config.yaml
        ├── script.py
        └── README.md
```

## Available Commands

### Main Commands

| Command | Description |
|---------|-------------|
| `skill init <repo-url>` | Initialize configuration with a repository |
| `skill config show` | Show current configuration |
| `skill config set-repo <url>` | Change configured repository |
| `skill config verify` | Verify repository access |
| `skill --help` | Show general help |
| `skill --version` | Show CLI version |

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
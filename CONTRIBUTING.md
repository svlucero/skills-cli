# Contributing to Skills CLI

First off, thank you for considering contributing to Skills CLI! It's people like you that make Skills CLI such a great tool.

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find that you don't need to create one. When you are creating a bug report, please include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps which reproduce the problem**
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed after following the steps**
- **Explain which behavior you expected to see instead and why**
- **Include screenshots if relevant**
- **Include your environment details** (OS, Go version, Git version)

#### Template for Bug Reports

```markdown
**Description**
A clear and concise description of the bug.

**Steps to Reproduce**
1. Run `skills ...`
2. Do ...
3. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g., macOS 13.0, Ubuntu 22.04]
- Go Version: [e.g., 1.24.3]
- Skills CLI Version: [e.g., 0.1.0]

**Additional Context**
Add any other context about the problem here.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Use a clear and descriptive title**
- **Provide a step-by-step description of the suggested enhancement**
- **Provide specific examples to demonstrate the steps**
- **Describe the current behavior and explain which behavior you expected to see instead**
- **Explain why this enhancement would be useful**

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. Ensure the test suite passes
4. Make sure your code follows the existing style
5. Write a clear commit message following our commit convention
6. Submit the pull request!

## Development Process

### Setting Up Your Development Environment

```bash
# Clone your fork
git clone https://github.com/your-username/skills-cli.git
cd skills-cli

# Add upstream remote
git remote add upstream https://github.com/svlucero/skills-cli.git

# Create a branch
git checkout -b feature/my-feature

# Build
make build

# Run tests
make test
```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

#### Examples

```
feat(install): add interactive skill selection

Add interactive mode to install command when no skill name is provided.
Users can now browse and select skills using arrow keys.

Closes #123
```

```
fix(repository): prevent removal of active repository

Add validation to prevent users from accidentally removing the active
repository without switching to another one first.
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Keep functions focused and small
- Add comments for complex logic
- Use meaningful variable names

### Testing

- Write tests for new features
- Ensure existing tests pass
- Test on multiple platforms if possible
- Include edge cases in your tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestSpecificFunction ./internal/cli
```

### Documentation

- Update README.md if you change functionality
- Add inline comments for complex code
- Update help text for CLI commands
- Add examples for new features

## Project Structure

```
skills-cli/
├── cmd/skill/              # Entry point
├── internal/
│   ├── cli/               # CLI commands (Cobra)
│   ├── config/            # Configuration management
│   ├── git/               # Git operations
│   ├── skill/             # Skill operations
│   └── errors/            # Custom errors
├── Makefile               # Build commands
└── README.md              # Project documentation
```

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## Recognition

Contributors will be recognized in the project README and release notes.

Thank you for contributing! 🎉

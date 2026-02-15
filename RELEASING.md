# Releasing

This document describes the process for creating new releases of Skills CLI.

## Prerequisites

- Push access to the main repository
- GoReleaser is configured (`.goreleaser.yml`)
- GitHub Actions is enabled
- All tests are passing

## Release Process

### 1. Update Version Information

Ensure all version references are updated if needed (though GoReleaser handles this automatically).

### 2. Update Changelog

Review recent commits and update if you maintain a manual CHANGELOG.md (optional, as GoReleaser auto-generates release notes).

### 3. Create and Push a Tag

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# Create a tag (use semantic versioning)
git tag -a v0.2.0 -m "Release v0.2.0"

# Push the tag
git push origin v0.2.0
```

### 4. GitHub Actions Automation

Once the tag is pushed, GitHub Actions will automatically:

1. Run tests across platforms (Linux, macOS, Windows)
2. Build binaries for all platforms
3. Generate checksums
4. Create release notes from commits
5. Publish the release on GitHub

### 5. Verify the Release

1. Go to https://github.com/svlucero/skills-cli/releases
2. Verify the latest release is published
3. Check that all binaries are attached
4. Review the auto-generated release notes
5. Test downloading and running a binary

### 6. Announce (Optional)

- Update README badges if needed
- Announce on social media
- Update any documentation sites

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0 → v2.0.0): Incompatible API changes
- **MINOR** version (v1.0.0 → v1.1.0): New features, backwards compatible
- **PATCH** version (v1.0.0 → v1.0.1): Bug fixes, backwards compatible

### Pre-releases

For pre-release versions, use suffixes with proper semver format:

```bash
# Release candidates (CORRECT format with dash and dot)
git tag -a v0.2.0-rc.1 -m "Release candidate 1"
git tag -a v0.2.0-rc.2 -m "Release candidate 2"

# Beta releases
git tag -a v0.2.0-beta.1 -m "Beta release 1"
git tag -a v0.2.0-beta.2 -m "Beta release 2"

# Alpha releases
git tag -a v0.2.0-alpha.1 -m "Alpha release 1"
```

**⚠️ Important**: Pre-release tags MUST follow this format:
- ✅ CORRECT: `v0.2.0-rc.1`, `v1.0.0-beta.1`, `v2.0.0-alpha.1`
- ❌ WRONG: `v0.2.0rc1`, `v1.0.0beta`, `v2.0.0-rc` (missing version number after dash)

## What GoReleaser Does

When a tag is pushed, GoReleaser:

1. **Builds** binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)

2. **Creates** archives:
   - `.tar.gz` for Linux and macOS
   - `.zip` for Windows

3. **Generates**:
   - SHA256 checksums
   - Release notes from commits
   - GitHub release

4. **Includes** in each archive:
   - Binary (`skills`)
   - LICENSE
   - README.md
   - CONTRIBUTING.md

## Troubleshooting

### Release Failed

Check GitHub Actions logs:
1. Go to Actions tab
2. Click on the failed workflow
3. Review error logs

Common issues:
- Tests failing: Fix tests and re-tag
- Build errors: Check Go version compatibility
- Token issues: Verify GITHUB_TOKEN permissions

### Deleting a Failed Release

```bash
# Delete the tag locally
git tag -d v0.2.0

# Delete the tag remotely
git push origin :refs/tags/v0.2.0

# Delete the release on GitHub (via web UI or gh CLI)
gh release delete v0.2.0
```

Then fix the issues and create the tag again.

## Testing Releases Locally

Before creating a real release, test locally:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Test the release process (doesn't publish)
goreleaser release --snapshot --clean

# Check generated artifacts in dist/
ls -la dist/
```

## Release Checklist

- [ ] All tests passing
- [ ] Main branch is up to date
- [ ] Version number follows semver
- [ ] Tag created with proper message
- [ ] Tag pushed to GitHub
- [ ] GitHub Actions workflow succeeded
- [ ] Release published on GitHub
- [ ] Binaries tested on at least one platform
- [ ] Release notes reviewed
- [ ] Documentation updated if needed

## Manual Release (Fallback)

If automated release fails, you can create a manual release:

```bash
# Build locally
make build

# Create release on GitHub manually
gh release create v0.2.0 \
  --title "v0.2.0" \
  --notes "Release notes here" \
  ./bin/skills
```

## Future Enhancements

Potential improvements to the release process:

- [ ] Homebrew formula auto-update
- [ ] Docker images
- [ ] Shell completion files in releases
- [ ] Snap/APT packages
- [ ] Chocolatey package for Windows

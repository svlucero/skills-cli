# Tag Format Guide

GoReleaser requires tags to follow [Semantic Versioning 2.0.0](https://semver.org/).

## Valid Tag Formats

### Stable Releases

```
v1.0.0       # Major release
v1.1.0       # Minor release (new features)
v1.1.1       # Patch release (bug fixes)
v2.0.0       # Breaking changes
```

### Pre-releases

```
v1.0.0-alpha.1    # Alpha release 1
v1.0.0-alpha.2    # Alpha release 2
v1.0.0-beta.1     # Beta release 1
v1.0.0-beta.2     # Beta release 2
v1.0.0-rc.1       # Release candidate 1
v1.0.0-rc.2       # Release candidate 2
```

## Invalid Tag Formats

❌ These will cause GoReleaser to fail:

```
v1.0.0rc1        # Missing dash before 'rc'
v1.0.0-rc        # Missing number after 'rc'
v1.0.0beta       # Missing dash and number
v1.0.0-beta      # Missing number after 'beta'
0.1.0            # Missing 'v' prefix
v1.0             # Incomplete version (needs patch number)
v1               # Incomplete version
```

## Quick Reference

```bash
# Stable releases
git tag -a v0.1.0 -m "Initial release"
git tag -a v0.2.0 -m "New features"
git tag -a v0.2.1 -m "Bug fixes"
git tag -a v1.0.0 -m "First stable release"

# Pre-releases
git tag -a v0.2.0-alpha.1 -m "Alpha 1"
git tag -a v0.2.0-beta.1 -m "Beta 1"
git tag -a v0.2.0-rc.1 -m "Release candidate 1"

# Push tag
git push origin v0.1.0
```

## Testing Locally

Before pushing a tag, you can test the release locally:

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Create a test tag
git tag v0.1.0-test

# Run goreleaser in snapshot mode (doesn't publish)
goreleaser release --snapshot --clean

# Check the output
ls -la dist/

# Delete test tag
git tag -d v0.1.0-test
```

## Semantic Versioning Rules

- **MAJOR** (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** (v1.0.0 → v1.1.0): New features (backwards compatible)
- **PATCH** (v1.0.0 → v1.0.1): Bug fixes (backwards compatible)

Pre-release identifiers:
- **alpha**: Early testing, unstable
- **beta**: Feature complete, testing
- **rc** (release candidate): Final testing before release

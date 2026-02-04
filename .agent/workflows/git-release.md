---
description: Create a Git release with proper versioning and tagging
---

# Git Release Management

This workflow guides creating a new release for Watchman.

## Prerequisites
- All tests passing
- Clean working directory
- Semantic versioning understanding

## Steps

### 1. Verify Clean State
// turbo
```bash
git status
make lint
make test
```

### 2. Determine Version
Follow Semantic Versioning (SemVer):
- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features, backward compatible
- **PATCH** (0.0.X): Bug fixes, backward compatible

Check current version:
// turbo
```bash
git describe --tags --abbrev=0
```

### 3. Update Version References (if any)
- [ ] Check `README.md` for version badges
- [ ] Update CHANGELOG.md (if exists)

### 4. Commit and Tag in One Command
```bash
git add -A && git commit -m "chore: release vX.Y.Z" && git tag -a vX.Y.Z -m "Release vX.Y.Z - Brief description"
```

### 5. Push with Tags
```bash
git push origin main --follow-tags
```

### 6. Verify GitHub Actions
The CI workflow will:
1. Run tests
2. Build binaries
3. Create GitHub release (if tag matches `v*`)

### 7. Verify Release
- Check GitHub Releases page
- Verify binary artifacts uploaded
- Test auto-update from previous version

## Tag Naming Convention
- Format: `vX.Y.Z`
- Examples: `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`

## Commit Message Convention
```
type: description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `chore`: Maintenance tasks
- `docs`: Documentation
- `refactor`: Code refactoring
- `test`: Adding tests
- `ci`: CI/CD changes

## ldflags Injection
The Makefile automatically injects version info:
```makefile
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)
```

## Rollback (if needed)
```bash
# Delete local tag
git tag -d vX.Y.Z

# Delete remote tag
git push origin :refs/tags/vX.Y.Z

# Reset commit (if not pushed)
git reset --soft HEAD~1
```

## Checklist
- [ ] All tests pass
- [ ] Lint clean
- [ ] Version follows SemVer
- [ ] CHANGELOG updated (if exists)
- [ ] Commit message follows convention
- [ ] Tag annotated with description
- [ ] Pushed with `--follow-tags`
- [ ] GitHub Release created
- [ ] Binaries uploaded

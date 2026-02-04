---
description: Fix golangci-lint errors in the Watchman project
---

# Fix Lint Errors

This workflow guides fixing common golangci-lint issues.

## Prerequisites
- golangci-lint v2.8.0+ installed
- Understand the project's `.golangci.yml` configuration

## Steps

### 1. Run Linter to Identify Issues
// turbo
```bash
make lint
```

### 2. Categorize Errors
Common error types and fixes:

#### goconst: String Literal Repeated
```go
// Before (goconst error)
if status == "failed" { ... }
if status == "failed" { ... }

// After: Define constant
const StatusFailed = "failed"
if status == StatusFailed { ... }
```

#### errcheck: Error Not Checked
```go
// Before (errcheck error)
file.Close()

// After
if err := file.Close(); err != nil {
    return fmt.Errorf("close file: %w", err)
}
// Or for deferred calls
defer func() { _ = file.Close() }()
```

#### gosimple: Simplify Code
```go
// Before (gosimple)
if x == true { ... }

// After
if x { ... }
```

#### govet: Suspicious Constructs
```go
// Before (govet: printf)
fmt.Printf("value: %s", intValue)

// After
fmt.Printf("value: %d", intValue)
```

#### ineffassign: Ineffective Assignment
```go
// Before (ineffassign)
err := doSomething()
err = doAnotherThing()  // Previous err never used

// After
if err := doSomething(); err != nil {
    return err
}
err := doAnotherThing()
```

### 3. Apply Auto-Fix (if available)
// turbo
```bash
make lint-fix
```

### 4. Run Format
// turbo
```bash
make fmt
```

### 5. Verify All Issues Resolved
// turbo
```bash
make lint
```

### 6. Run Tests to Ensure No Regression
// turbo
```bash
make test
```

## Common Linter Rules Reference

| Linter | Purpose |
|--------|---------|
| `goconst` | Find repeated strings that should be constants |
| `errcheck` | Check for unchecked errors |
| `gosimple` | Simplify code constructs |
| `govet` | Report suspicious constructs |
| `ineffassign` | Detect ineffective assignments |
| `staticcheck` | Advanced static analysis |
| `unused` | Find unused code |

## Configuration
The project uses `.golangci.yml` with v2 schema. Key settings:
- Timeout: 5m
- Build tags: `windows`
- Exclusions defined for generated code

## Checklist
- [ ] All lint errors identified
- [ ] Fixes applied without changing behavior
- [ ] Auto-fix used where applicable
- [ ] Code formatted
- [ ] Tests pass
- [ ] No new lint warnings introduced

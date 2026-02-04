---
description: Add a new internal package to the Watchman project
---

# Add New Internal Package

This workflow guides creating a new package in the `internal/` directory.

## Prerequisites
- Understand Go package patterns and internal/ directory semantics
- Review existing packages for consistency

## Steps

### 1. Create Package Directory
```powershell
mkdir internal/<packagename>
```

### 2. Create Main Package File
Create `internal/<packagename>/<packagename>.go`:

```go
// Package <packagename> provides <brief description>.
package <packagename>

// Define exported types and interfaces first
type Config struct {
    // Configuration fields
}

// Define sentinel errors
var (
    Err<Something> = errors.New("<packagename>: descriptive error")
)
```

### 3. Create Interface (if needed)
Define interfaces for testability and loose coupling:

```go
// <PackageName>er defines the contract for <packagename> operations.
type <PackageName>er interface {
    Method(ctx context.Context) error
}
```

### 4. Create Package AGENTS.md (Optional but Recommended)
Create `internal/<packagename>/AGENTS.md`:

```markdown
# <PackageName> Package Guidelines

## Purpose
Brief description of what this package does.

## Key Patterns
- Pattern 1
- Pattern 2

## Testing
- Required test coverage: 80%
- Use table-driven tests
```

### 5. Create Test File
Create `internal/<packagename>/<packagename>_test.go`:

```go
package <packagename>

import (
    "testing"
)

func Test<Function>(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "test",
            want:  "expected",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 6. Update go.mod (if needed)
// turbo
```bash
go mod tidy
```

### 7. Verify Build and Lint
// turbo
```bash
make lint
make test
```

## Package Design Principles

### Do
- ✅ Use guard clauses for early returns
- ✅ Wrap errors with context: `fmt.Errorf("packagename: action: %w", err)`
- ✅ Define interfaces at the consumer side
- ✅ Keep helper functions private (lowercase)
- ✅ Use `context.Context` for cancellation

### Don't
- ❌ Create circular dependencies
- ❌ Log and return the same error
- ❌ Export types unless necessary
- ❌ Use `panic()` for recoverable errors

## Checklist
- [ ] Package follows Go naming conventions
- [ ] Package documentation (godoc comment)
- [ ] Error wrapping with context
- [ ] Interface defined if needed
- [ ] Test file created
- [ ] AGENTS.md created (recommended)
- [ ] No circular dependencies
- [ ] Lint passes

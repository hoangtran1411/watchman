---
description: Perform AI-assisted code review for Watchman
---

# Code Review Workflow

This workflow guides performing a thorough code review.

## Prerequisites
- Changes ready for review
- Tests passing
- Lint clean

## Steps

### 1. Pre-Review Checks
// turbo
```bash
make lint
make test
```

### 2. Review Dimensions

Score each dimension 1-10. Flag any score below 8.

#### Code Quality (Weight: 30%)
- [ ] Follows Go idioms and Effective Go
- [ ] Proper error handling (wrap with context)
- [ ] Guard clauses for early returns
- [ ] No magic numbers/strings (use constants)
- [ ] Clear naming (functions, variables)
- [ ] Appropriate comments (why, not what)

#### Performance (Weight: 20%)
- [ ] No unnecessary allocations
- [ ] Efficient database queries
- [ ] Proper use of context for cancellation
- [ ] Connection pooling utilized
- [ ] No blocking operations in hot paths

#### Security (Weight: 25%)
- [ ] Parameterized SQL queries (no injection)
- [ ] No sensitive data in logs
- [ ] Proper error messages (no internal details)
- [ ] Input validation
- [ ] Credentials not hardcoded

#### Maintainability (Weight: 15%)
- [ ] Single responsibility principle
- [ ] Appropriate abstraction level
- [ ] Testable design (interfaces)
- [ ] Documentation for public APIs
- [ ] Consistent with existing patterns

#### Scalability (Weight: 10%)
- [ ] Handles multiple servers gracefully
- [ ] Resilient to server failures
- [ ] Configurable timeouts
- [ ] Graceful degradation

### 3. Review Checklist by File Type

#### For `cmd/` files:
- [ ] Command has Short and Long descriptions
- [ ] Examples provided
- [ ] JSON output supported
- [ ] Exit codes documented
- [ ] Flags bound to Viper

#### For `internal/` files:
- [ ] Package documentation
- [ ] Interfaces defined for testability
- [ ] No circular dependencies
- [ ] Errors wrapped with context
- [ ] Tests cover edge cases

#### For config changes:
- [ ] snake_case for YAML keys
- [ ] Default values set
- [ ] Validation added
- [ ] Example config updated

### 4. Common Issues to Check

| Category | Issue | Standard |
|----------|-------|----------|
| Error | Log and return | ❌ Do one, not both |
| Error | Generic message | ❌ Wrap with context |
| SQL | String concatenation | ❌ Use parameters |
| Naming | Abbreviations | ❌ Use full words |
| Testing | No edge cases | ❌ Test boundaries |
| Docs | No examples | ❌ Add usage examples |

### 5. Generate Review Report

```markdown
## Code Review Summary

**Files Reviewed:** X
**Review Date:** YYYY-MM-DD

### Scores
| Dimension | Score | Notes |
|-----------|-------|-------|
| Code Quality | X/10 | |
| Performance | X/10 | |
| Security | X/10 | |
| Maintainability | X/10 | |
| Scalability | X/10 | |

### Issues Found
1. [Critical/Major/Minor] Description
2. ...

### Recommendations
1. Suggestion for improvement
2. ...

### Approval
- [ ] Approved
- [ ] Approved with suggestions
- [ ] Changes requested
```

### 6. Verify Fixes (if changes requested)
// turbo
```bash
make lint
make test
```

## Review Principles

1. **Be Constructive**: Focus on code, not person
2. **Explain Why**: Provide rationale for feedback
3. **Prioritize**: Focus on important issues first
4. **Acknowledge Good Code**: Positive feedback matters
5. **Be Specific**: Reference exact lines/functions

## Automation Support

AI agents should check:
- `make lint` output
- `make test` output
- `coverage.out` for coverage gaps
- `.golangci.yml` for enabled linters

## Checklist
- [ ] All dimensions scored
- [ ] Critical issues identified
- [ ] Security checked
- [ ] Tests adequate
- [ ] Documentation complete
- [ ] Review report generated

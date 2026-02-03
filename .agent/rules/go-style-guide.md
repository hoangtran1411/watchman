---
trigger: always_on
---

# Go Style Guide - Watchmen

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is a **Windows Service for SQL Server Agent Job Monitoring** built with:
- **Cobra CLI** for command-line interface
- **Viper** for YAML configuration management
- **Windows Service** (`golang.org/x/sys/windows/svc`)
- **Internal packages** for modular business logic

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint` (v2.8.0+) `run ./...` before commit.
- **Linting Configuration**: MUST use `golangci-lint` v2 configuration schema (v2.8.x+). 
  - Top-level `version: "2"` is mandatory.
  - Use kebab-case for all linter settings.
  - Exclusions move to `linters: exclusions: rules`.
- Adhere to [Effective Go](https://go.dev/doc/effective_go).
- Core logic in `internal/` packages. CLI commands in `cmd/watchmen/`.

## Project Structure

```
watchmen/
  cmd/
    watchmen/
      main.go            - CLI entry point
      root.go            - Root command
      install.go         - Install service command
      check.go           - Manual check command
      update.go          - Auto-update command
      reload.go          - Reload config command
  internal/
    config/              - Configuration management (Viper/YAML)
    database/            - SQL Server connection & queries
    jobs/                - Job monitoring logic
    notification/        - Windows Toast Notification
    scheduler/           - Cron-like scheduler
    service/             - Windows Service wrapper
    updater/             - Auto-update logic
  pkg/
    logger/              - Structured logging (zerolog)
  configs/               - Example configuration files
  scripts/               - PowerShell installation scripts
```

## Error Handling

- Wrap errors: `fmt.Errorf("context: %w", err)`
- Guard clauses for fail-fast
- Do not log and return the same error
- One responsibility per layer
- **Skip unavailable servers silently** - SysAdmin responsibility

## Windows Service Integration

- Service methods on `*Service` struct
- Use `svc.Handler` interface for lifecycle management
- Graceful shutdown with context cancellation
- Log to Windows Event Log for critical errors

## SQL Server Patterns

- Use `database/sql` with `go-mssqldb` driver
- Always use parameterized queries (prevent SQL injection)
- Include `@@SERVERNAME` in queries for server identification
- Respect connection timeouts and query timeouts
- Close connections properly with `defer`

## Testing & Linting

- Table-driven tests with `t.Run`
- Target 80% coverage for `internal/`
- Use `make test` and `make lint`
- Mock interfaces for Windows Service and SQL connections

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first
- **Windows-first**: This is a Windows-only application

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Cobra CLI](https://github.com/spf13/cobra)
- [Viper Config](https://github.com/spf13/viper)
- [go-mssqldb](https://github.com/microsoft/go-mssqldb)
- [go-toast](https://github.com/go-toast/toast)
- [Windows Service](https://pkg.go.dev/golang.org/x/sys/windows/svc)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [selfupdate](https://github.com/rhysd/go-github-selfupdate)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.

---

## Project-Specific Conventions

### Configuration File Locations

| Environment | Path |
|-------------|------|
| **Production** | `%ProgramData%\Watchmen\config.yaml` |
| **Development** | `./config.yaml` (current directory) |
| **Logs** | `%ProgramData%\Watchmen\logs\` |

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Commands | lowercase verb | `install`, `check`, `update` |
| Config keys | snake_case | `check_times`, `max_retries` |
| Go structs | PascalCase | `ServerConfig`, `JobResult` |
| Errors | `Err` prefix | `ErrServerUnavailable` |

### Version & Release

- Semantic versioning: `vX.Y.Z`
- Version injected via ldflags: `-X main.Version=$VERSION`
- Auto-update checks GitHub releases on service start

---

## Version History

| Version | Date       | Changes                              |
|---------|------------|--------------------------------------|
| 1.0.0   | 2026-02-03 | Initial rules for Watchmen project   |

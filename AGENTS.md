# AGENTS.md - Watchmen Project Context for AI Agents

Welcome, AI Agent! This document provides essential context and instructions tailored for working on the **Watchmen** project. Follow these guidelines to ensure your contributions align with our standards and architecture.

## 🚀 Project Overview
**Watchmen** is a Windows-first application designed to monitor SQL Server Agent Jobs across multiple servers and send notifications (Windows Toast) when failures or specific conditions occur.

- **Role**: You are a Senior Go Developer specializing in Windows Service programming and SQL Server integrations.
- **Tech Stack**:
  - **Language**: Go 1.25.6+
  - **CLI Framework**: Cobra
  - **Configuration**: Viper (YAML)
  - **Database**: SQL Server (`github.com/microsoft/go-mssqldb`)
  - **Windows Service**: `golang.org/x/sys/windows/svc`
  - **Logging**: `zerolog` (housed in `pkg/logger`)

## 🛠️ Common Commands
Always run these commands to verify your work:

- **Build**: `make build` (creates `watchmen.exe`)
- **Format**: `make fmt` (runs `gofmt` and `goimports`)
- **Lint**: `make lint` (uses `golangci-lint`)
- **Test**: `make test` (runs all tests with race detection)
- **Coverage**: `make coverage` (checks statement coverage)

## 📁 Project Structure
- `cmd/watchman/`: Entry point and Cobra commands.
- `internal/`:
  - `config/`: Configuration parsing and validation.
  - `database/`: SQL Server connection handling and queries.
  - `jobs/`: Logic for analyzing Agent Job history.
  - `notification/`: Windows Toast notification implementation.
  - `scheduler/`: Cron-based check scheduling.
  - `service/`: Windows Service lifecycle management.
- `pkg/logger/`: Structured logging wrapper.

## ⚖️ Core Rules & Conventions
1. **Windows First**: This application is strictly for Windows. Use `golang.org/x/sys/windows` for OS-specific logic.
2. **Error Handling**: Use `fmt.Errorf("context: %w", err)`. **Never** log and return the same error.
3. **Fail-Fast**: Use guard clauses.
4. **Resilience**: If a SQL Server is unavailable, skip it silently (log a warning, but don't crash).
5. **SQL Safety**: 
   - Always use parameterized queries.
   - Include `@@SERVERNAME` in results to identify the source.
6. **Linting**: We use `golangci-lint` with a v2 schema. Ensure `.golangci.yml` is respected.
7. **Style**: Follow `go-style-guide.md` located in the root or `.agent/rules/`.

## 🔒 Boundaries
- **Do not** modify the `.git` directory.
- **Do not** change the project structure without explicit instruction.
- **Do not** introduce heavy external dependencies unless absolutely necessary and approved.

## 💡 Pro-Tips for Agents
- When making changes to `internal/service`, remember it handles the Windows Service signals.
- If you add new configuration keys, update the `config.yaml` example in `configs/` and the struct in `internal/config`.
- Check `coverage.out` regularly. We aim for high statement coverage. (Current target > 60%).

---
*This file is updated periodically. If you notice inconsistencies, please point them out.*

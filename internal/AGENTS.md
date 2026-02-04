# Internal Logic Guidelines

This directory contains the core business logic for Watchmen. All code here must follow these internal-specific rules.

## 🏗️ Architecture
- **No Circular Dependencies**: Ensure packages don't import each other in a loop.
- **Interface Segregation**: Use interfaces to define behavior between internal packages.
- **Privacy**: Keep helper functions and internal types private (lowercase) unless they MUST be accessed by other packages.

## 📦 Package Responsibilities
- `config`: Handles YAML parsing and environment variable expansion.
- `database`: SQL Server connections. **MUST** use parameterized queries.
- `jobs`: Logic for job history analysis.
- `notification`: Windows Toast UI logic.
- `scheduler`: Management of check intervals.
- `service`: Windows Service Handler (`svc.Handler`) implementation.
- `updater`: GitHub self-update integration.

## 🧪 Testing
- Every internal package should have a `_test.go` file.
- Use table-driven tests.
- Target: 80% coverage for logic in this directory.

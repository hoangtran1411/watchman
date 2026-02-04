# CLI Command Guidelines

This directory contains the Cobra-based CLI implementation.

## 🛠️ Command Structure
- `main.go`: Application entry point. Keep it minimal.
- `root.go`: Definition of the root `watchmen` command.
- Command files should be named after the verb (e.g., `check.go`, `install.go`).

## 📋 Requirements for New Commands
1. **Help Text**: Always include a `Short` and `Long` description.
2. **Examples**: Provide usage examples in the `Example` field.
3. **Flags**:
   - Use `snake_case` for flag names in configuration, but `kebab-case` for CLI flags.
   - Bind flags to Viper for unified configuration.
4. **Output**:
   - Default output should be human-readable.
   - Use `-o json` or `--output json` for machine-readable results (essential for AI agents).
5. **Exit Codes**:
   - `0`: Success
   - `1`: Operational failure (e.g., failed jobs found)
   - `2+`: System/Usage error

## 🤖 AI Agent Friendly
Ensure flags and outputs are predictable. Check `root.go` for global persistent flags like `--config` or `--debug`.

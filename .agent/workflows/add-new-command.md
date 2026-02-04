---
description: Add a new Cobra CLI command to the Watchman application
---

# Add New CLI Command

This workflow guides adding a new Cobra command to Watchman.

## Prerequisites
- Understand Go and Cobra CLI framework
- Review existing commands in `cmd/watchman/` for patterns

## Steps

### 1. Create Command File
Create a new file `cmd/watchman/<verb>.go` following naming conventions (lowercase verb).

```go
// cmd/watchman/<verb>.go
package main

import (
    "github.com/spf13/cobra"
)

var <verb>Cmd = &cobra.Command{
    Use:   "<verb>",
    Short: "Brief description (shown in help list)",
    Long:  `Detailed description with usage context.`,
    Example: `  watchmen <verb> --flag value`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(<verb>Cmd)
    
    // Add flags
    <verb>Cmd.Flags().StringP("output", "o", "text", "Output format (text|json)")
}
```

### 2. Register Command in root.go
// turbo
Verify the command is auto-registered via `init()` function. No manual registration needed.

### 3. Add Viper Binding (if needed)
```go
func init() {
    <verb>Cmd.Flags().StringP("config-key", "k", "", "Description")
    viper.BindPFlag("config_key", <verb>Cmd.Flags().Lookup("config-key"))
}
```

### 4. Implement JSON Output Support
```go
if outputFormat == "json" {
    return json.NewEncoder(os.Stdout).Encode(result)
}
```

### 5. Update Documentation
- [ ] Update `AGENTS.md` with new command
- [ ] Add example to `README.md`
- [ ] Update `Makefile` if adding a convenience target

### 6. Verify
// turbo
```bash
make build
./watchmen.exe <verb> --help
```

### 7. Run Tests
// turbo
```bash
make lint
make test
```

## Exit Code Guidelines
- `0`: Success
- `1`: Operational failure (e.g., failed jobs found)
- `2+`: System/Usage error (invalid flags, missing config)

## Checklist
- [ ] Command follows naming convention (lowercase verb)
- [ ] Short and Long descriptions provided
- [ ] Example usage included
- [ ] JSON output supported with `-o json`
- [ ] Flags use kebab-case
- [ ] Config keys use snake_case
- [ ] Tests written with table-driven pattern

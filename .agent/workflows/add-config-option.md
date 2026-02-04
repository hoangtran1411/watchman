---
description: Add a new configuration option to Watchman
---

# Add New Configuration Option

This workflow guides adding a new YAML configuration option.

## Prerequisites
- Understand Viper configuration management
- Review `internal/config/` for existing patterns

## Steps

### 1. Define the Config Field
Update `internal/config/config.go`:

```go
type Config struct {
    // ... existing fields
    NewOption string `mapstructure:"new_option"` // Use snake_case
}
```

### 2. Add Validation (if needed)
```go
func (c *Config) Validate() error {
    if c.NewOption == "" {
        return errors.New("config: new_option is required")
    }
    // Add range/format validation as needed
    return nil
}
```

### 3. Set Default Value
In the config loading function:

```go
func init() {
    viper.SetDefault("new_option", "default_value")
}
```

### 4. Update Example Config
Edit `configs/config.yaml`:

```yaml
# Description of what this option does
# Valid values: option1, option2 (or range, format)
new_option: "default_value"
```

### 5. Add CLI Flag Override (if needed)
In the relevant command file:

```go
func init() {
    cmd.Flags().String("new-option", "", "CLI description")
    viper.BindPFlag("new_option", cmd.Flags().Lookup("new-option"))
}
```

### 6. Add Environment Variable Support
Viper automatically maps:
- Config key: `new_option`
- Env var: `WATCHMEN_NEW_OPTION`

```go
viper.SetEnvPrefix("WATCHMEN")
viper.AutomaticEnv()
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
```

### 7. Update Documentation
- [ ] Update `README.md` with new option
- [ ] Update `AGENTS.md` if agent-relevant
- [ ] Document in `configs/config.yaml` with comments

### 8. Add Tests
Create test in `internal/config/config_test.go`:

```go
func TestConfig_NewOption(t *testing.T) {
    tests := []struct {
        name      string
        yamlInput string
        want      string
        wantErr   bool
    }{
        {
            name:      "valid value",
            yamlInput: "new_option: valid",
            want:      "valid",
        },
        {
            name:      "empty uses default",
            yamlInput: "",
            want:      "default_value",
        },
    }
    // ... test implementation
}
```

### 9. Verify
// turbo
```bash
make lint
make test
```

## Config Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| YAML keys | `snake_case` | `check_interval` |
| CLI flags | `kebab-case` | `--check-interval` |
| Env vars | `SCREAMING_SNAKE` | `WATCHMEN_CHECK_INTERVAL` |
| Go fields | `PascalCase` | `CheckInterval` |

## Priority Order (Viper)
1. CLI flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

## Checklist
- [ ] Struct field added with mapstructure tag
- [ ] Validation added if needed
- [ ] Default value set
- [ ] Example config updated
- [ ] CLI flag added (if applicable)
- [ ] Environment variable documented
- [ ] README.md updated
- [ ] Tests written
- [ ] Lint passes

# Go Idioms Reference - Watchmen

> **Full Reference Document** - Contains detailed idioms, code examples, and best practices.  
> For compact core rules, see `go-style-guide.md`

---

## Naming Conventions (Idiomatic Go)

- Use short but meaningful names, scoped by context:
  - `r`, `w`, `ctx`, `db`, `tx`, `cfg`, `srv` are acceptable in small scopes.
  - Avoid `data`, `info`, `obj`, `temp`, `value` unless unavoidable.

- Prefer noun-based names for structs, verb-based names for functions:
  - `Monitor.Check()`, `Notifier.Send()`, `Scheduler.Start()`

- Boolean names should read naturally:
  - `isValid`, `hasError`, `enabled`, `isAvailable`

- Avoid stuttering:
  - ❌ `config.ConfigServer`
  - ✅ `config.Server`

---

## Function Design

- Prefer small functions (≤ 40 lines).
- One function = one responsibility.
- Avoid flags that change behavior dramatically:

```go
// bad
func CheckJobs(server string, silent bool)

// good
func CheckJobs(server string) ([]FailedJob, error)
func CheckJobsSilent(server string) error
```

- Return early (guard clauses):

```go
if err != nil {
    return nil, err
}

if !server.Enabled {
    return nil, nil // Skip disabled servers
}
```

---

## Error Handling Idioms

- Never ignore errors explicitly:

```go
_ = f.Close() // ❌ unless justified in comment
```

- Wrap errors only at package boundaries:

```go
return nil, fmt.Errorf("failed to connect to %s: %w", server.Name, err)
```

- Do not wrap errors multiple times in the same layer.

- Prefer `errors.Is` / `errors.As` for comparisons:

```go
if errors.Is(err, ErrServerUnavailable) {
    // Skip silently - SysAdmin responsibility
    continue
}
```

- Define typed errors for specific cases:

```go
var (
    ErrServerUnavailable = errors.New("server unavailable")
    ErrNoFailedJobs      = errors.New("no failed jobs found")
    ErrConfigInvalid     = errors.New("configuration invalid")
)
```

---

## Package Design & Boundaries

```
internal/
├── config/       # YAML parsing, validation, ENV substitution
├── database/     # SQL Server connections, queries
├── jobs/         # Business logic: job monitoring, filtering
├── notification/ # Windows Toast, future: Email/Teams
├── scheduler/    # Cron scheduling, timer management
├── service/      # Windows Service lifecycle
└── updater/      # GitHub release check, self-update
```

- `internal` packages must be:
  - UI-agnostic
  - Framework-agnostic (no Cobra imports in business logic)

- Each package should expose minimal API surface:

```go
// good - clean public API
package jobs

func CheckFailedJobs(ctx context.Context, servers []config.Server) ([]FailedJob, error)

// avoid exposing helpers like buildQuery(), parseResult()
```

- Avoid circular dependencies at all cost.
- If two packages depend on each other → redesign.

---

## Context Usage Idioms

- `context.Context` must:
  - Be the first parameter
  - Never be stored in struct fields

- Do not pass nil context:

```go
ctx := context.Background() // For main entry
ctx := context.TODO()       // For temporary code
```

- Respect cancellation in loops:

```go
for _, server := range servers {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    result, err := checkServer(ctx, server)
    // ...
}
```

- Use timeouts for database operations:

```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

---

## Concurrency Patterns

- Prefer worker pool over unbounded goroutines:

```go
// Use errgroup for parallel server checks
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(cfg.Monitoring.Parallel.MaxConcurrent)

for _, server := range servers {
    server := server // capture
    g.Go(func() error {
        return checkServer(ctx, server)
    })
}

return g.Wait()
```

- Always define ownership of goroutines:
  - Who starts?
  - Who stops?

- Use `errgroup.Group` for concurrent tasks with error propagation.

- Channels should have clear direction:

```go
func worker(in <-chan ServerConfig, out chan<- JobResult)
```

- Avoid closing channels you did not create.

---

## Struct & Interface Idioms

- Accept interfaces, return concrete types:

```go
func NewMonitor(db Database, notifier Notifier) *Monitor
```

- Interfaces should be small (1–3 methods):

```go
type Database interface {
    QueryFailedJobs(ctx context.Context, server ServerConfig) ([]FailedJob, error)
    Ping(ctx context.Context) error
}

type Notifier interface {
    Notify(ctx context.Context, failures []FailedJob) error
}
```

- Do not define interfaces prematurely.
- Define interfaces where they are used, not where they are implemented.

---

## Zero Value Philosophy

- Design structs so zero value is usable:

```go
var m Monitor // should work with defaults
m.Check(ctx)  // uses default config
```

- Avoid constructors unless needed for invariants.

- Prefer empty slices over nil slices for JSON output:

```go
jobs := make([]FailedJob, 0) // Returns [] not null in JSON
```

---

## Slice & Map Best Practices

- Pre-allocate when size is known:

```go
results := make([]FailedJob, 0, len(servers)*10)
```

- Check map existence properly:

```go
server, ok := serverMap[name]
if !ok {
    return ErrServerNotFound
}
```

- Do not modify slices while ranging over them.

---

## Testing Idioms

- Test behavior, not implementation.
- Table-driven tests with descriptive names:

```go
tests := []struct {
    name        string
    servers     []ServerConfig
    wantJobs    int
    wantErr     bool
}{
    {
        name:     "single server with failed jobs",
        servers:  []ServerConfig{{Name: "PROD-01", Enabled: true}},
        wantJobs: 3,
        wantErr:  false,
    },
    {
        name:     "disabled server is skipped",
        servers:  []ServerConfig{{Name: "DEV-01", Enabled: false}},
        wantJobs: 0,
        wantErr:  false,
    },
}
```

- Avoid `t.Fatal` inside loops.
- Use `cmp.Diff` or `reflect.DeepEqual` consistently.
- Tests must not depend on execution order.

### Mocking for Windows Service

```go
// Mock interface for testing
type MockNotifier struct {
    notifications []Notification
}

func (m *MockNotifier) Notify(ctx context.Context, failures []FailedJob) error {
    m.notifications = append(m.notifications, Notification{Jobs: failures})
    return nil
}
```

---

## Logging (zerolog)

- Do not log inside core business logic.
- Log at boundaries (CLI / Service layer).
- Logs must be structured and actionable:

```go
log.Info().
    Str("server", server.Name).
    Int("failed_jobs", len(jobs)).
    Msg("job check completed")

log.Error().
    Err(err).
    Str("server", server.Name).
    Msg("failed to connect to server")
```

- Use appropriate log levels:
  - `Debug`: Detailed debugging info
  - `Info`: Normal operations
  - `Warn`: Recoverable issues
  - `Error`: Failed operations

---

## Comments & Documentation

- Comments explain **why**, not **what**.
- Avoid redundant comments:

```go
i++ // increment i ❌
```

- Exported comments must start with identifier name:

```go
// CheckFailedJobs queries all enabled servers for failed SQL Agent jobs.
// It returns aggregated results from all servers.
func CheckFailedJobs(ctx context.Context, servers []ServerConfig) ([]FailedJob, error)
```

- Use TODO with owner & reason:

```go
// TODO(username): add retry logic for transient failures
```

---

## Defensive Programming

- Validate inputs at package boundary:

```go
func NewMonitor(cfg Config) (*Monitor, error) {
    if len(cfg.Servers) == 0 {
        return nil, ErrNoServersConfigured
    }
    // ...
}
```

- Never trust external data types.
- Fail fast on schema mismatch.
- Prefer explicit errors over silent correction.

---

## SQL Server Specific Idioms

### Connection String Building

```go
query := url.Values{}
query.Add("database", server.Database)
query.Add("encrypt", strconv.FormatBool(server.Options.Encrypt))
query.Add("TrustServerCertificate", strconv.FormatBool(server.Options.TrustServerCertificate))
query.Add("connection timeout", strconv.Itoa(server.Options.ConnectionTimeout))

u := &url.URL{
    Scheme:   "sqlserver",
    User:     url.UserPassword(server.Auth.Username, server.Auth.Password),
    Host:     fmt.Sprintf("%s:%d", server.Host, server.Port),
    RawQuery: query.Encode(),
}
```

### Query with @@SERVERNAME

```go
const queryFailedJobs = `
SELECT 
    @@SERVERNAME AS ServerName,
    j.name AS JobName,
    h.message AS ErrorMessage,
    -- ...
FROM msdb.dbo.sysjobs j
INNER JOIN msdb.dbo.sysjobhistory h ON j.job_id = h.job_id
WHERE h.run_status = 0
`
```

### Resource Cleanup

```go
rows, err := db.QueryContext(ctx, queryFailedJobs)
if err != nil {
    return nil, fmt.Errorf("query failed: %w", err)
}
defer rows.Close()

for rows.Next() {
    // ...
}

if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("row iteration failed: %w", err)
}
```

---

## Build & Tooling Practices

- `go.mod` must be tidy:

```bash
go mod tidy
```

- CI must fail on:
  - lint
  - test
  - formatting

- **Linting**: `.golangci.yml` MUST use version 2 schema (`version: "2"`).
  - Use kebab-case for linter settings (e.g., `ignore-sigs`, `ignore-package-globs`).
  - Exclusions must be configured under `linters: exclusions: rules` instead of `issues: exclude-rules`.
  - Prefer global exclusions in config over redundant `//nolint` comments in test files.

- Avoid build tags unless justified.

---

## Reference Links

### Official Go Documentation
- Reference: https://go.dev/doc
- Primary source for Go syntax, tooling, modules.

### Effective Go
- Reference: https://go.dev/doc/effective_go
- Idiomatic Go practices. All `internal/` packages must comply.

### Go Modules
- Reference: https://go.dev/ref/mod
- Dependency management. Avoid unnecessary `replace` directives.

### Go Testing
- Reference: https://go.dev/doc/testing
- Standard patterns for unit tests, benchmarks, coverage.

### Go Context
- Reference: https://pkg.go.dev/context
- Mandatory for cancellation and timeouts in I/O operations.

### Go Error Handling
- Reference: https://go.dev/blog/error-handling-and-go
- Errors are values. Wrap errors; avoid panic in business logic.

### Go Concurrency
- Reference: https://go.dev/doc/effective_go#concurrency
- Use goroutines and channels deliberately.

### Go Standard Library
- Reference: https://pkg.go.dev/std
- Prefer stdlib before third-party dependencies.

### Cobra CLI
- Reference: https://github.com/spf13/cobra
- Follow Cobra patterns for command structure.

### Viper Configuration
- Reference: https://github.com/spf13/viper
- YAML configuration with ENV variable support.

### SQL Server Driver
- Reference: https://github.com/microsoft/go-mssqldb
- Official Microsoft driver for Go.

### Windows Service
- Reference: https://pkg.go.dev/golang.org/x/sys/windows/svc
- Windows Service lifecycle management.

### Linting
- Reference: https://github.com/golangci/golangci-lint
- Run before committing. Fix all issues to pass CI.

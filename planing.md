# 📋 Watchmen - SQL Server Agent Job Monitor

## 🎯 Mục tiêu dự án

Xây dựng một **Windows Service** bằng Go sử dụng Cobra CLI để:
- Chạy nền (background) trên Windows
- **Monitoring nhiều SQL Server instances** (multi-server support)
- Đọc lỗi của SQL Server Agent Jobs từ tất cả servers đã cấu hình
- Mỗi ngày lúc **8:00 AM** kiểm tra và thông báo qua **Windows Toast Notification** nếu có job lỗi
- Sử dụng **YAML file** để quản lý cấu hình (config.yaml)

---

## 📁 Cấu trúc dự án

```
watchman/
├── .github/
│   └── workflows/
│       ├── ci.yml              # CI workflow (lint, test, build)
│       └── release.yml         # Release workflow (build & publish)
├── cmd/
│   └── watchman/
│       ├── main.go             # CLI entry point
│       ├── root.go             # Root command
│       ├── install.go          # Install service command
│       ├── uninstall.go        # Uninstall service command
│       ├── check.go            # Manual check command
│       ├── update.go           # Auto-update command
│       └── reload.go           # Reload config command
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management (YAML/ENV)
│   ├── database/
│   │   └── sqlserver.go        # SQL Server connection & queries
│   ├── jobs/
│   │   ├── monitor.go          # Job monitoring logic
│   │   └── types.go            # Job types & structs
│   ├── notification/
│   │   └── windows.go          # Windows Toast Notification
│   ├── scheduler/
│   │   └── scheduler.go        # Cron-like scheduler (8:00 AM daily)
│   ├── service/
│   │   └── windows_service.go  # Windows Service implementation
│   └── updater/
│       └── updater.go          # GitHub release auto-update
├── pkg/
│   └── logger/
│       └── logger.go           # Structured logging (zerolog)
├── scripts/
│   ├── install.ps1             # PowerShell installer (main)
│   ├── install.bat             # Batch wrapper for install.ps1
│   ├── uninstall.ps1           # PowerShell uninstaller
│   └── uninstall.bat           # Batch wrapper for uninstall.ps1
├── configs/
│   └── config.example.yaml     # Example configuration file
├── .golangci.yml               # GolangCI-Lint configuration
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                    # Build automation
├── README.md
└── planing.md                  # This file
```

---

## 🛠️ Tech Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **CLI Framework** | [Cobra](https://github.com/spf13/cobra) | Industry standard, subcommand support |
| **Config Management** | [Viper](https://github.com/spf13/viper) | YAML, ENV, flags support |
| **SQL Server Driver** | [go-mssqldb](https://github.com/microsoft/go-mssqldb) | Official Microsoft driver |
| **Scheduler** | [gocron](https://github.com/go-co-op/gocron) | Simple, reliable cron scheduler |
| **Windows Service** | [golang.org/x/sys/windows/svc](https://pkg.go.dev/golang.org/x/sys/windows/svc) | Official Go Windows service |
| **Toast Notification** | [go-toast](https://github.com/go-toast/toast) | Native Windows notifications |
| **Logging** | [zerolog](https://github.com/rs/zerolog) | High-performance structured logging |
| **Testing** | [testify](https://github.com/stretchr/testify) | Assertions & mocking |

---

## 📋 Features & Commands

### Cobra CLI Commands

```bash
# Install as Windows Service
watchman install

# Uninstall Windows Service
watchman uninstall

# Start service
watchman start

# Stop service
watchman stop

# Run once (manual check)
watchman check

# Show version
watchman version

# Show configuration
watchman config show

# Validate configuration
watchman config validate

# Reload configuration (without restart service)
watchman reload

# Check for updates and apply
watchman update        # Check for new version
watchman update -y     # Auto-apply update without confirmation
```

### Auto-Update Feature (Required)

| Trigger | Behavior |
|---------|----------|
| **Service Start** | Check GitHub releases for new version |
| **Manual** | `watchman update -y` to force update |
| **Notification** | Toast notification when update available |

---

## 🤖 Agent Experience (AX) - AI Agent Friendly CLI

> Thiết kế CLI thân thiện với AI Agents, lấy cảm hứng từ `golangci-lint`

### Tại sao AX quan trọng?

| User Type | Needs |
|-----------|-------|
| **Human** | Readable output, colors, emojis |
| **AI Agent** | Structured output (JSON), predictable format, clear exit codes |
| **CI/CD** | Machine-parseable, no interactive prompts |

### Global Flags

```bash
# Output format (applies to most commands)
--output, -o    Output format: text, json (default "text")

# Quiet mode
--quiet, -q     Suppress all output except errors

# Verbose mode  
--verbose, -v   Enable verbose/debug logging

# Config path
--config, -c    Config file path (default "%ProgramData%\Watchmen\config.yaml")

# Help
--help, -h      Show help for command
```

### Predictable Exit Codes

```go
const (
    ExitSuccess           = 0  // Everything OK, no failed jobs
    ExitFailedJobs        = 1  // Found failed jobs (expected behavior)
    ExitConfigError       = 2  // Configuration issue
    ExitConnectionError   = 3  // Cannot connect to any server
    ExitInternalError     = 4  // Unexpected internal error
)
```

### Commands với AI Agent Support

| Command | AI-Friendly Flags | Output |
|---------|-------------------|--------|
| `check` | `--output json`, `--server`, `--notify` | Failed jobs list |
| `config show` | `--output json` | Full config (sanitized) |
| `config validate` | `--output json` | Validation result |
| `version` | `--output json` | Version info |
| `update` | `--yes`, `--check-only` | Update status |
| `install` | `--silent`, `--config` | Install status |
| `uninstall` | `--keep-config`, `--yes` | Uninstall status |
| `reload` | `--output json` | Reload status |

### JSON Output Examples

**`watchman check --output json`**

```json
{
  "status": "success",
  "timestamp": "2026-02-03T08:00:00+07:00",
  "servers_checked": 3,
  "servers_available": 2,
  "servers_unavailable": ["DEV-SQL01"],
  "failed_jobs": [
    {
      "server": "PROD-SQL01",
      "job_name": "Backup_Database",
      "failed_at": "2026-02-03T07:30:00+07:00",
      "error_message": "Timeout expired",
      "duration_seconds": 3600
    }
  ],
  "summary": "1 failed job on 1 server"
}
```

**`watchman version --output json`**

```json
{
  "version": "1.2.0",
  "commit": "abc123def",
  "build_date": "2026-02-03T10:00:00Z",
  "go_version": "go1.25.6",
  "os": "windows",
  "arch": "amd64"
}
```

**`watchman config validate --output json`**

```json
{
  "valid": true,
  "servers": [
    {"name": "PROD-SQL01", "enabled": true, "reachable": true},
    {"name": "STAGING-SQL01", "enabled": true, "reachable": false, "error": "connection timeout"}
  ],
  "warnings": [
    "Server STAGING-SQL01 is configured but not reachable"
  ],
  "errors": []
}
```

### Comprehensive Help Format

**Root command: `watchman --help`**

```
Watchman - SQL Server Agent Job Monitor

A Windows service that monitors SQL Server Agent jobs and sends 
Windows Toast notifications when jobs fail.

Usage:
  watchman [command]

Available Commands:
  check       Check for failed jobs (manual run)
  config      Manage configuration
  install     Install as Windows Service
  reload      Reload configuration without restart
  start       Start the service
  stop        Stop the service
  uninstall   Remove Windows Service
  update      Check for and apply updates
  version     Show version information

Flags:
  -c, --config string   Config file path (default "%ProgramData%\Watchman\config.yaml")
  -h, --help            Show help for command
  -o, --output string   Output format: text, json (default "text")
  -q, --quiet           Suppress all output except errors
  -v, --verbose         Enable verbose logging

Examples:
  # Check for failed jobs
  watchman check

  # Check with JSON output (for AI Agents/scripting)
  watchman check --output json

  # Install service with custom config
  watchman install --config D:\configs\watchman.yaml

  # Force update without confirmation
  watchman update --yes

Exit Codes:
  0  Success / No failed jobs
  1  Failed jobs found (check completed successfully)
  2  Configuration error
  3  Connection error (all servers unreachable)
  4  Internal error

Use "watchman [command] --help" for more information about a command.
```

**Subcommand: `watchman check --help`**

```
Check for failed SQL Server Agent jobs

Queries all configured and enabled SQL Server instances for failed 
jobs within the lookback period. By default, shows results in 
human-readable format. Use --output json for machine-readable output.

Usage:
  watchman check [flags]

Flags:
  -h, --help              Show help
  -o, --output string     Output format: text, json (default "text")
  -s, --server string     Check specific server only (by name)
      --lookback int      Hours to look back for failures (default: from config)
      --notify            Send notification if failures found
      --no-color          Disable colored output

Global Flags:
  -c, --config string     Config file path
  -q, --quiet             Suppress output (exit code only)
  -v, --verbose           Verbose output

Examples:
  # Check all servers
  watchman check

  # Check specific server
  watchman check --server PROD-SQL01

  # Check and send notification
  watchman check --notify

  # JSON output for scripting/AI Agents
  watchman check --output json

  # JSON output piped to jq
  watchman check --output json | jq '.failed_jobs[] | .job_name'

  # Check with custom lookback period
  watchman check --lookback 48

  # Quiet mode for scripts (check exit code only)
  watchman check --quiet && echo "No failures" || echo "Has failures"
```

### Design Principles

| Principle | Implementation |
|-----------|----------------|
| **Predictable** | Same input → Same output format |
| **Parseable** | JSON output với consistent schema |
| **Non-interactive** | `--yes`, `--silent` flags for automation |
| **Self-documenting** | Comprehensive `--help` với examples |
| **Exit codes** | Clear, documented exit codes |
| **Error messages** | Structured errors in JSON mode |

---

## ⚙️ Configuration (YAML)

### `config.yaml` - Multi-Server Support

```yaml
# =============================================================================
# WATCHMEN CONFIGURATION FILE
# =============================================================================
# Sử dụng YAML format với hỗ trợ:
# - Environment variables: ${ENV_VAR} hoặc ${ENV_VAR:default}
# - Multiple SQL Server instances
# - Job filtering per server
# =============================================================================

# -----------------------------------------------------------------------------
# SQL Server Instances (MULTI-SERVER SUPPORT)
# -----------------------------------------------------------------------------
servers:
  # Production Server
  - name: "PROD-SQL01"
    enabled: true
    host: "sql-prod-01.company.local"
    port: 1433
    database: "msdb"
    auth:
      type: "sql"  # sql | windows
      username: "watchman_svc"
      password: "${PROD_SQL_PASSWORD}"  # Environment variable
    options:
      encrypt: true
      trust_server_certificate: false
      connection_timeout: 30
      query_timeout: 60
    jobs:
      include: []  # Empty = all jobs
      exclude:
        - "test_*"
        - "dev_*"

  # Staging Server
  - name: "STAGING-SQL01"
    enabled: true
    host: "sql-staging-01.company.local"
    port: 1433
    database: "msdb"
    auth:
      type: "windows"  # Windows Authentication
      username: ""
      password: ""
    options:
      encrypt: false
      trust_server_certificate: true
      connection_timeout: 30
      query_timeout: 60
    jobs:
      include:
        - "ETL_*"
        - "Backup_*"
      exclude: []

  # Development Server (disabled)
  - name: "DEV-SQL01"
    enabled: false  # Không monitor server này
    host: "localhost"
    port: 1433
    database: "msdb"
    auth:
      type: "sql"
      username: "sa"
      password: "${DEV_SQL_PASSWORD:P@ssw0rd}"
    options:
      encrypt: false
      trust_server_certificate: true
      connection_timeout: 15
      query_timeout: 30
    jobs:
      include: []
      exclude: []

# -----------------------------------------------------------------------------
# Scheduler Configuration
# -----------------------------------------------------------------------------
scheduler:
  # Check time (24-hour format)
  check_times:
    - "08:00"  # Morning check
    # - "14:00"  # Afternoon check (optional)
    # - "20:00"  # Evening check (optional)
  timezone: "Asia/Ho_Chi_Minh"
  
  # Retry configuration if check fails
  retry:
    enabled: true
    max_attempts: 3
    delay_seconds: 60

# -----------------------------------------------------------------------------
# Notification Configuration
# -----------------------------------------------------------------------------
notification:
  app_id: "Watchmen"
  icon_path: ""  # Optional: absolute path to .ico file
  
  # Grouping: combine multiple failures into single notification
  grouping:
    enabled: true
    max_jobs_per_notification: 5  # Show max 5 jobs, then "and X more..."
  
  # Sound
  sound:
    enabled: true
    type: "default"  # default, mail, reminder, sms, alarm

# -----------------------------------------------------------------------------
# Logging Configuration
# -----------------------------------------------------------------------------
logging:
  level: "info"  # trace, debug, info, warn, error, fatal
  format: "json"  # json, text
  
  # File logging
  file:
    enabled: true
    path: "logs/watchman.log"
    max_size_mb: 10
    max_backups: 5
    max_age_days: 30
    compress: true
  
  # Windows Event Log
  event_log:
    enabled: true
    source: "Watchmen"

# -----------------------------------------------------------------------------
# Job Monitoring Settings
# -----------------------------------------------------------------------------
monitoring:
  # Look back period for failed jobs
  lookback_hours: 24
  
  # Job statuses to report
  report_statuses:
    - failed      # run_status = 0
    - cancelled   # run_status = 3
    # - retried   # run_status = 2 (uncomment to include)
  
  # Parallel checking (check multiple servers concurrently)
  parallel:
    enabled: true
    max_concurrent: 5
```

---

## 🗄️ SQL Server Queries

### Query failed jobs (Last 24 hours)

```sql
-- Sử dụng @@SERVERNAME để identify server trong kết quả
SELECT 
    @@SERVERNAME AS ServerName,          -- Server identifier
    j.name AS JobName,
    h.run_date AS RunDate,
    h.run_time AS RunTime,
    h.run_status AS Status,
    h.message AS ErrorMessage,
    h.run_duration AS Duration,
    CONVERT(datetime, 
        CONVERT(varchar(8), h.run_date) + ' ' + 
        STUFF(STUFF(RIGHT('000000' + CONVERT(varchar(6), h.run_time), 6), 5, 0, ':'), 3, 0, ':')
    ) AS FailedAt
FROM msdb.dbo.sysjobs j
INNER JOIN msdb.dbo.sysjobhistory h 
    ON j.job_id = h.job_id
WHERE h.step_id = 0  -- Job outcome
    AND h.run_status = 0  -- Failed
    AND CONVERT(datetime, 
        CONVERT(varchar(8), h.run_date) + ' ' + 
        STUFF(STUFF(RIGHT('000000' + CONVERT(varchar(6), h.run_time), 6), 5, 0, ':'), 3, 0, ':')
    ) >= DATEADD(hour, -@LookbackHours, GETDATE())
ORDER BY h.run_date DESC, h.run_time DESC
```

### Server Availability Check

```sql
-- Quick ping to check if server is available
-- If fails → Skip server silently (SysAdmin responsibility)
SELECT @@SERVERNAME AS ServerName, GETDATE() AS ServerTime
```

---

## 🔔 Windows Notification

### Toast Notification Sample (with Server Name)

```go
// Single server failure
notification := toast.Notification{
    AppID:   "Watchmen",
    Title:   "⚠️ [PROD-SQL01] Job Failed",  // Server name từ @@SERVERNAME
    Message: "Job 'Backup_Database' failed at 07:30 AM\nError: Timeout expired",
    Icon:    "", // Optional
    Actions: []toast.Action{
        {Type: "protocol", Label: "Open SSMS", Arguments: "ssms://"},
        {Type: "protocol", Label: "Dismiss", Arguments: "dismiss"},
    },
}

// Multiple failures grouped
notification := toast.Notification{
    AppID:   "Watchmen",
    Title:   "⚠️ 3 Jobs Failed on 2 Servers",
    Message: "[PROD-SQL01] Backup_Database, ETL_Daily\n[STAGING-SQL01] Report_Gen",
    Icon:    "",
}

// Update available notification
notification := toast.Notification{
    AppID:   "Watchmen",
    Title:   "🔄 Update Available",
    Message: "Watchmen v1.2.0 is available (current: v1.1.0)\nRun 'watchman update -y' to apply",
    Actions: []toast.Action{
        {Type: "protocol", Label: "Update Now", Arguments: "watchman://update"},
    },
}
```

---

## 📦 Installation Scripts

### Installation Configuration

| Setting | Value | Reason |
|---------|-------|--------|
| **Install Directory** | `%ProgramData%\Watchmen` | Standard app data, writable by service |
| **Service Account** | `LocalSystem` | Full network access, simple setup |
| **Startup Type** | `Automatic (Delayed Start)` | Không block boot, đợi network ready |
| **Service Name** | `Watchmen` | Short, memorable |
| **Display Name** | `Watchmen - SQL Agent Monitor` | Descriptive in Services.msc |

### Scripts Overview

| File | Description |
|------|-------------|
| `install.ps1` | Main PowerShell installer |
| `install.bat` | Batch wrapper (double-click friendly) |
| `uninstall.ps1` | Main PowerShell uninstaller |
| `uninstall.bat` | Batch wrapper (double-click friendly) |

### `install.ps1` - Features

```powershell
# Usage:
.\install.ps1                           # Interactive mode
.\install.ps1 -Silent                    # Silent mode (no prompts)
.\install.ps1 -ConfigPath "D:\config.yaml"  # Custom config path
```

**Installation Flow:**

```
┌─────────────────────────────────────────────────────────────┐
│  1. Check Administrator privileges                          │
│  2. Check if service already exists → Upgrade mode          │
│  3. Create installation directory                           │
│     %ProgramData%\Watchmen\                                  │
│  4. Copy watchman.exe to installation directory             │
│  5. Copy config.example.yaml → config.yaml (if not exists)  │
│  6. Create logs directory                                   │
│  7. Register Windows Service (sc.exe create)                │
│  8. Set service to Automatic (Delayed Start)                │
│  9. Start service                                           │
│  10. Verify installation & show status                      │
└─────────────────────────────────────────────────────────────┘
```

**Upgrade Behavior:**
- Detect existing installation
- Stop service
- Backup config.yaml → config.yaml.backup
- Replace watchman.exe
- Start service
- Verify upgrade

### `uninstall.ps1` - Features

```powershell
# Usage:
.\uninstall.ps1                    # Interactive mode (asks to keep config)
.\uninstall.ps1 -KeepConfig        # Keep config and logs
.\uninstall.ps1 -RemoveAll         # Remove everything including config
```

**Uninstallation Flow:**

```
┌─────────────────────────────────────────────────────────────┐
│  1. Check Administrator privileges                          │
│  2. Stop service (if running)                               │
│  3. Delete Windows Service (sc.exe delete)                  │
│  4. Ask: Keep config & logs? (Interactive mode)             │
│  5. Remove installation directory (based on choice)         │
│  6. Confirm uninstallation complete                         │
└─────────────────────────────────────────────────────────────┘
```

### Batch Wrappers

**`install.bat`**
```batch
@echo off
cd /d "%~dp0"
PowerShell -ExecutionPolicy Bypass -File ".\install.ps1" %*
pause
```

**`uninstall.bat`**
```batch
@echo off
cd /d "%~dp0"
PowerShell -ExecutionPolicy Bypass -File ".\uninstall.ps1" %*
pause
```

### Sample Output

```
═══════════════════════════════════════════════════════════════
   🔧 WATCHMEN INSTALLER v1.0.0
   SQL Server Agent Job Monitor
═══════════════════════════════════════════════════════════════

[✓] Administrator privileges confirmed
[✓] Creating installation directory...
[✓] Copying watchman.exe...
[✓] Creating default configuration...
[✓] Registering Windows Service...
[✓] Setting service to Auto-Start (Delayed)...
[✓] Starting service...

═══════════════════════════════════════════════════════════════
   ✅ INSTALLATION COMPLETE!
   
   Service Name:  Watchmen
   Status:        Running
   Startup:       Automatic (Delayed Start)
   Config:        C:\ProgramData\Watchmen\config.yaml
   Logs:          C:\ProgramData\Watchmen\logs\
   
   Next steps:
   1. Edit config.yaml to add your SQL Server(s)
   2. Run: watchman reload (to apply changes)
═══════════════════════════════════════════════════════════════
```

---

## 🧪 Testing Strategy

### Unit Tests
- `internal/config/config_test.go` - Config parsing
- `internal/database/sqlserver_test.go` - DB connection (mock)
- `internal/jobs/monitor_test.go` - Job monitoring logic
- `internal/notification/windows_test.go` - Notification (mock)
- `internal/scheduler/scheduler_test.go` - Scheduler logic

### Integration Tests
- Database connectivity
- End-to-end job monitoring

### Test Coverage Target: **80%+**

---

## 📊 GolangCI-Lint Configuration

### `.golangci.yml`

```yaml
version: "2"

run:
  timeout: 5m
  go: "1.25.6"
  modules-download-mode: readonly

linters:
  enable:
    # Bugs
    - bodyclose
    - durationcheck
    - errcheck
    - exportloopref
    - gosec
    - nilerr
    - noctx
    - rowserrcheck
    - sqlclosecheck
    - staticcheck
    - typecheck
    # Performance
    - prealloc
    # Style
    - gofmt
    - goimports
    - govet
    - ineffassign
    - misspell
    - unconvert
    - unused
    # Complexity
    - cyclop
    - funlen
    - gocognit
    - goconst
    - gocyclo
    # Error Handling
    - errorlint
    - wrapcheck
    # Code Quality
    - dupl
    - gocritic
    - revive
    - stylecheck

linters-settings:
  cyclop:
    max-complexity: 15
  funlen:
    lines: 100
    statements: 50
  gocognit:
    min-complexity: 20
  goconst:
    min-len: 3
    min-occurrences: 3
  gocyclo:
    min-complexity: 15
  govet:
    enable-all: true
  revive:
    rules:
      - name: exported
        severity: warning
  stylecheck:
    checks: ["all", "-ST1000"]
  errorlint:
    errorf: true
  gosec:
    excludes:
      - G104  # Unhandled errors (handled by errcheck)

issues:
  exclude-use-default: false
  max-issues-per-linter: 50
  max-same-issues: 10
  exclude-dirs:
    - vendor
    - .git
```

---

## 🚀 GitHub Actions

### CI Workflow (`.github/workflows/ci.yml`)

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  GO_VERSION: "1.24"

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v7
        with:
          version: latest
          args: --timeout=5m

  test:
    name: Test
    runs-on: windows-latest
    needs: lint
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: coverage.out

  build:
    name: Build
    runs-on: windows-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - name: Build
        run: go build -ldflags="-s -w" -o watchman.exe ./cmd/watchman
      - uses: actions/upload-artifact@v4
        with:
          name: watchman-windows
          path: watchman.exe
```

### Release Workflow (`.github/workflows/release.yml`)

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

env:
  GO_VERSION: "1.24"

jobs:
  release:
    name: Release
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build
        run: |
          $VERSION = "${{ github.ref_name }}"
          go build -ldflags="-s -w -X main.Version=$VERSION" -o watchman.exe ./cmd/watchman

      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          files: watchman.exe
          generate_release_notes: true
```

---

## 📅 Implementation Timeline

### Phase 1: Foundation (Day 1-2)
- [x] Tạo `planing.md`
- [x] Khởi tạo Go module
- [x] Setup `.golangci.yml`
- [x] Setup GitHub Actions (CI/Release)
- [x] Tạo cấu trúc thư mục

### Phase 2: Core Logic (Day 3-5)
- [x] Config management (Viper)
- [x] SQL Server connection
- [x] Job monitoring queries
- [x] Logger setup

### Phase 3: Features (Day 6-8)
- [x] Cobra CLI commands
- [x] Windows Toast Notification
- [x] Scheduler (8:00 AM daily)
- [x] Windows Service wrapper

### Phase 4: Testing & Polish (Day 9-10)
- [x] Unit tests (80%+ coverage)
- [x] Integration tests (Manual)
- [x] Documentation
- [x] First release

---

## ✅ Core Features (MVP)

| Feature | Status | Description |
|---------|--------|-------------|
| Multi-server monitoring | ✅ | Monitor nhiều SQL Server instances |
| YAML configuration | ✅ | File config dễ đọc, dễ chỉnh sửa |
| Scheduled check (8:00 AM) | ✅ | Kiểm tra job failures hàng ngày |
| Windows Toast Notification | ✅ | Thông báo với server name (@@SERVERNAME) |
| Auto-update on startup | ✅ | Check GitHub releases khi khởi động |
| Manual update (`update -y`) | ✅ | Cập nhật thủ công khi cần |
| Config reload (`reload`) | ✅ | Tải lại config không cần restart |
| Windows Service | ✅ | Chạy nền như Windows Service |
| Graceful shutdown | ✅ | Tắt đúng cách khi stop service |

---

## 🔮 Future Enhancements (Backlog)

> Các tính năng có thể phát triển sau khi MVP hoàn thành

### Priority: High
| Feature | Description | Rationale |
|---------|-------------|----------|
| Email notifications | Gửi email khi có job fail | Backup cho Toast notification |
| Teams/Slack webhooks | Notify qua chat apps | Team collaboration |
| Custom notification templates | User tự định nghĩa message format | Flexibility |

### Priority: Medium
| Feature | Description | Rationale |
|---------|-------------|----------|
| Server health check notification | Thông báo khi server không available | Hiện tại skip silently (SysAdmin responsibility) |
| Web dashboard | UI để xem history của failed jobs | Better visibility |
| MSI installer | Professional installer cho Windows | Easier deployment |
| Chocolatey package | Publish lên Chocolatey | Auto-updates via choco |

### Priority: Low
| Feature | Description | Rationale |
|---------|-------------|----------|
| Prometheus metrics | Export metrics cho monitoring | Enterprise environments |
| Windows Credential Manager | Store passwords securely | Better security |
| Real-time monitoring | WebSocket-based live updates | Overkill cho use case hiện tại |
| Job step details | Hiển thị chi tiết từng step fail | More granular info |

---

## ⚠️ Known Limitations & Design Decisions

| Decision | Reasoning |
|----------|----------|
| **Server unavailable → Skip silently** | SysAdmin responsibility; app không cần notify vì họ có monitoring riêng |
| **No hot-reload** | Config ít thay đổi; `reload` command đủ dùng |
| **Check only on schedule** | Không real-time vì job failures không critical enough |
| **Update check on startup only** | Tránh spam GitHub API; user có thể manual update |

---

## 🧪 Testing Strategy

### Challenges & Solutions

| Challenge | Solution |
|-----------|----------|
| Mock Windows Service | Interface abstraction với `ServiceManager` |
| Mock Toast Notifications | Interface `Notifier` với mock implementation |
| SQL Server test data | Docker container hoặc in-memory mock |
| Auto-update testing | Mock GitHub API responses |

---

## ✅ Next Steps

1. **Review this plan** - Xác nhận các features và timeline
2. **Initialize project** - `go mod init github.com/username/watchman`
3. **Setup CI/CD** - Push to GitHub, verify workflows
4. **Start implementation** - Begin with Phase 1 (Foundation)

---

## 📝 Technical Notes

| Item | Value |
|------|-------|
| **Go Version** | 1.24+ |
| **Target OS** | Windows 10/11, Windows Server 2016+ |
| **SQL Server** | 2012+ |
| **Build Command** | `go build -ldflags="-s -w" -o watchman.exe ./cmd/watchman` |
| **Config Location** | `%ProgramData%\Watchmen\config.yaml` |
| **Log Location** | `%ProgramData%\Watchmen\logs\` |

---

## 📚 References

- [Cobra CLI](https://github.com/spf13/cobra)
- [Viper Config](https://github.com/spf13/viper)
- [go-mssqldb](https://github.com/microsoft/go-mssqldb)
- [go-toast](https://github.com/go-toast/toast)
- [Windows Service in Go](https://pkg.go.dev/golang.org/x/sys/windows/svc)
- [selfupdate](https://github.com/rhysd/go-github-selfupdate)

---

*Last Updated: 2026-02-03*

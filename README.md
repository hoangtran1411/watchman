# 🔍 Watchman

[![CI](https://github.com/hoangtran1411/watchman/actions/workflows/ci.yml/badge.svg)](https://github.com/hoangtran1411/watchman/actions/workflows/ci.yml)
[![Release](https://github.com/hoangtran1411/watchman/actions/workflows/release.yml/badge.svg)](https://github.com/hoangtran1411/watchman/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Watchman** is a Windows Service that monitors SQL Server Agent jobs and sends Windows Toast notifications when jobs fail.

## ✨ Features

- 🖥️ **Windows Service** - Runs in background as a Windows Service
- 🗄️ **Multi-Server Support** - Monitor multiple SQL Server instances
- ⏰ **Scheduled Checks** - Check for failed jobs at specified times (default: 8:00 AM)
- 🔔 **Toast Notifications** - Native Windows 10/11 notifications with server name
- 🔄 **Auto-Update** - Automatic updates from GitHub releases
- 🤖 **AI Agent Friendly** - JSON output, predictable exit codes, comprehensive `--help`

## 📥 Installation

### Quick Install

1. Download the latest release from [Releases](https://github.com/hoangtran1411/watchman/releases)
2. Extract the ZIP file
3. Run `install.bat` as Administrator
4. Edit `%ProgramData%\Watchman\config.yaml` with your SQL Server details
5. Run `watchman reload`

### Manual Install

```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/hoangtran1411/watchman/releases/latest/download/watchman.exe" -OutFile "watchman.exe"

# Install as service
.\watchman.exe install

# Or use PowerShell script
.\install.ps1 -Silent
```

## ⚙️ Configuration

Configuration file location: `%ProgramData%\Watchman\config.yaml`

```yaml
# SQL Server instances to monitor
servers:
  - name: "PROD-SQL01"
    enabled: true
    host: "sql-prod-01.company.local"
    port: 1433
    database: "msdb"
    auth:
      type: "sql"  # sql | windows
      username: "watchman_svc"
      password: "${PROD_SQL_PASSWORD}"  # Environment variable
    jobs:
      include: []  # Empty = all jobs
      exclude:
        - "test_*"

# Schedule
scheduler:
  check_times:
    - "08:00"
  timezone: "Asia/Ho_Chi_Minh"

# Notification
notification:
  app_id: "Watchmen"
  grouping:
    enabled: true
    max_jobs_per_notification: 5
```

See [config.example.yaml](configs/config.example.yaml) for full configuration options.

## 🚀 Usage

### CLI Commands

```bash
# Check for failed jobs manually
watchman check

# Check with JSON output (for AI Agents/scripts)
watchman check --output json

# Check specific server
watchman check --server PROD-SQL01

# Show version
watchman version

# Show/validate configuration
watchman config show
watchman config validate

# Reload configuration without restart
watchman reload

# Update to latest version
watchman update
watchman update --yes  # Auto-apply without confirmation

# Service management
watchman install    # Install as Windows Service
watchman uninstall  # Remove Windows Service
watchman start      # Start service
watchman stop       # Stop service
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success / No failed jobs |
| 1 | Failed jobs found |
| 2 | Configuration error |
| 3 | Connection error |
| 4 | Internal error |

### JSON Output (AI Agent Friendly)

```bash
watchman check --output json
```

```json
{
  "status": "success",
  "timestamp": "2026-02-03T08:00:00+07:00",
  "servers_checked": 2,
  "servers_available": 2,
  "failed_jobs": [
    {
      "server": "PROD-SQL01",
      "job_name": "Backup_Database",
      "failed_at": "2026-02-03T07:30:00+07:00",
      "error_message": "Timeout expired"
    }
  ],
  "summary": "1 failed job on 1 server"
}
```

## 🛠️ Development

### Prerequisites

- Go 1.24+
- golangci-lint v2.0+
- Windows 10/11 or Windows Server 2016+

### Build

```bash
# Clone repository
git clone https://github.com/hoangtran1411/watchman.git
cd watchman

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run linter
make lint

# See all commands
make help
```

### Project Structure

```
watchman/
├── cmd/watchman/          # CLI commands (Cobra)
├── internal/
│   ├── config/            # Configuration (Viper/YAML)
│   ├── database/          # SQL Server connection
│   ├── jobs/              # Job monitoring logic
│   ├── notification/      # Windows Toast
│   ├── scheduler/         # Cron scheduler
│   ├── service/           # Windows Service
│   └── updater/           # Auto-update
├── pkg/logger/            # Structured logging
├── scripts/               # Install/Uninstall scripts
└── configs/               # Example configuration
```

## 📋 Roadmap

- [x] Multi-server monitoring
- [x] YAML configuration
- [x] Windows Toast notifications
- [x] Auto-update from GitHub
- [x] AI Agent friendly CLI
- [ ] Email notifications
- [ ] Microsoft Teams webhook
- [ ] Web dashboard

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration
- [go-mssqldb](https://github.com/microsoft/go-mssqldb) - SQL Server driver
- [go-toast](https://github.com/go-toast/toast) - Windows notifications
- [selfupdate](https://github.com/rhysd/go-github-selfupdate) - Auto-update

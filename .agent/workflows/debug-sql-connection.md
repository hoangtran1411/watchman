---
description: Debug SQL Server connection issues in Watchman
---

# Debug SQL Server Connection

This workflow helps diagnose SQL Server connectivity problems.

## Prerequisites
- Access to SQL Server instance
- Network connectivity to server
- Valid credentials

## Steps

### 1. Verify Network Connectivity
// turbo
```powershell
Test-NetConnection -ComputerName <server> -Port 1433
```

### 2. Test Connection with sqlcmd (if available)
```powershell
sqlcmd -S <server> -U <user> -P <password> -Q "SELECT @@SERVERNAME, @@VERSION"
```

### 3. Check Config File
Verify `config.yaml` settings:

```yaml
servers:
  - name: "Server Display Name"
    host: "server.domain.com"      # Or IP address
    port: 1433                      # Default SQL Server port
    database: "msdb"                # For Agent job queries
    username: "watchmen_user"
    password: "secure_password"
    connection_timeout: 30          # Seconds
    query_timeout: 60               # Seconds
```

### 4. Run Manual Check
// turbo
```bash
./watchmen.exe check --config ./configs/config.yaml --output json
```

### 5. Enable Debug Logging
```yaml
logging:
  level: debug
  output: stdout
```

Or via environment:
```powershell
$env:WATCHMEN_LOG_LEVEL = "debug"
./watchmen.exe check
```

### 6. Common Connection String Issues

| Issue | Symptom | Solution |
|-------|---------|----------|
| Wrong port | Connection timeout | Use `1433` or named instance port |
| Firewall | Connection refused | Open port in Windows Firewall |
| SQL Auth disabled | Login failed | Enable SQL + Windows auth |
| Wrong database | Cannot open DB | Use `msdb` for Agent queries |
| SSL/TLS | Certificate error | Add `TrustServerCertificate=true` |

### 7. Test Connection in Code
Add temporary debug code:

```go
import (
    "database/sql"
    "fmt"
    _ "github.com/microsoft/go-mssqldb"
)

func testConnection(connStr string) error {
    db, err := sql.Open("sqlserver", connStr)
    if err != nil {
        return fmt.Errorf("open: %w", err)
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        return fmt.Errorf("ping: %w", err)
    }
    
    var serverName string
    err = db.QueryRow("SELECT @@SERVERNAME").Scan(&serverName)
    if err != nil {
        return fmt.Errorf("query: %w", err)
    }
    
    fmt.Printf("Connected to: %s\n", serverName)
    return nil
}
```

### 8. Connection String Format
```
sqlserver://user:password@server:port?database=msdb&connection+timeout=30
```

Or with named instance:
```
sqlserver://user:password@server\instance?database=msdb
```

## Resilience Pattern (Project Standard)
```go
// Skip unavailable servers silently - SysAdmin responsibility
if err := db.Ping(); err != nil {
    log.Warn().
        Str("server", config.Name).
        Err(err).
        Msg("server unavailable, skipping")
    return nil // Don't fail the entire check
}
```

## Integration with Windows Authentication
```go
// Windows Auth (Integrated Security)
connStr := "sqlserver://server:1433?database=msdb&trusted_connection=yes"
```

## Checklist
- [ ] Network connectivity verified
- [ ] Port is correct (1433 or custom)
- [ ] Credentials are valid
- [ ] Database exists (msdb for Agent)
- [ ] Firewall allows connection
- [ ] Connection timeout appropriate
- [ ] Query timeout appropriate
- [ ] Error handling follows resilience pattern

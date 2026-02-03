// Package database provides SQL Server connectivity for Watchman.
// It uses go-mssqldb driver to connect and query SQL Server Agent jobs.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	_ "github.com/microsoft/go-mssqldb" // SQL Server driver

	"github.com/hoangtran1411/watchman/internal/config"
)

// DB represents a SQL Server database connection.
type DB struct {
	conn   *sql.DB
	server config.ServerConfig
}

// FailedJob represents a failed SQL Server Agent job.
type FailedJob struct {
	ServerName   string    `json:"server"`
	JobName      string    `json:"job_name"`
	RunDate      int       `json:"run_date"`
	RunTime      int       `json:"run_time"`
	FailedAt     time.Time `json:"failed_at"`
	Status       int       `json:"status"`
	ErrorMessage string    `json:"error_message"`
	Duration     int       `json:"duration_seconds"`
}

// New creates a new database connection.
func New(server config.ServerConfig) (*DB, error) {
	connStr := buildConnectionString(server)

	conn, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(time.Duration(server.Options.ConnectionTimeout) * time.Second * 2)

	return &DB{
		conn:   conn,
		server: server,
	}, nil
}

// Ping tests the database connection.
// Ping tests the database connection.
func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.server.Options.ConnectionTimeout)*time.Second)
	defer cancel()

	if err := db.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// Close closes the database connection.
// Close closes the database connection.
func (db *DB) Close() error {
	if db.conn != nil {
		if err := db.conn.Close(); err != nil {
			return fmt.Errorf("close failed: %w", err)
		}
	}
	return nil
}

// GetServerName returns the SQL Server name using @@SERVERNAME.
func (db *DB) GetServerName(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.server.Options.QueryTimeout)*time.Second)
	defer cancel()

	var serverName string
	err := db.conn.QueryRowContext(ctx, "SELECT @@SERVERNAME").Scan(&serverName)
	if err != nil {
		return "", fmt.Errorf("failed to get server name: %w", err)
	}

	return serverName, nil
}

// QueryFailedJobs queries for failed SQL Server Agent jobs.
func (db *DB) QueryFailedJobs(ctx context.Context, lookbackHours int) ([]FailedJob, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(db.server.Options.QueryTimeout)*time.Second)
	defer cancel()

	query := `
SELECT 
    @@SERVERNAME AS ServerName,
    j.name AS JobName,
    h.run_date AS RunDate,
    h.run_time AS RunTime,
    h.run_status AS Status,
    ISNULL(h.message, '') AS ErrorMessage,
    h.run_duration AS Duration
FROM msdb.dbo.sysjobs j
INNER JOIN msdb.dbo.sysjobhistory h 
    ON j.job_id = h.job_id
WHERE h.step_id = 0
    AND h.run_status = 0
    AND CONVERT(datetime, 
        CONVERT(varchar(8), h.run_date) + ' ' + 
        STUFF(STUFF(RIGHT('000000' + CONVERT(varchar(6), h.run_time), 6), 5, 0, ':'), 3, 0, ':')
    ) >= DATEADD(hour, -@LookbackHours, GETDATE())
ORDER BY h.run_date DESC, h.run_time DESC
`

	rows, err := db.conn.QueryContext(ctx, query, sql.Named("LookbackHours", lookbackHours))
	if err != nil {
		return nil, fmt.Errorf("failed to query failed jobs: %w", err)
	}
	defer func() {
		_ = rows.Close() // Ignore validation error on close
	}()

	var jobs []FailedJob
	for rows.Next() {
		var job FailedJob
		err := rows.Scan(
			&job.ServerName,
			&job.JobName,
			&job.RunDate,
			&job.RunTime,
			&job.Status,
			&job.ErrorMessage,
			&job.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse FailedAt from RunDate and RunTime
		job.FailedAt = parseDateTime(job.RunDate, job.RunTime)

		// Apply job filters
		if !db.matchesFilter(job.JobName) {
			continue
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return jobs, nil
}

// matchesFilter checks if a job name matches the include/exclude filters.
func (db *DB) matchesFilter(jobName string) bool {
	filter := db.server.Jobs

	// If include list is specified, job must match at least one pattern
	if len(filter.Include) > 0 {
		matched := false
		for _, pattern := range filter.Include {
			if matchPattern(jobName, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// If exclude list is specified, job must not match any pattern
	for _, pattern := range filter.Exclude {
		if matchPattern(jobName, pattern) {
			return false
		}
	}

	return true
}

// matchPattern matches a job name against a pattern (supports * wildcard).
func matchPattern(name, pattern string) bool {
	// Simple wildcard matching
	if pattern == "*" {
		return true
	}

	// Prefix match (e.g., "test_*")
	if len(pattern) > 1 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(name) >= len(prefix) && name[:len(prefix)] == prefix
	}

	// Suffix match (e.g., "*_backup")
	if len(pattern) > 1 && pattern[0] == '*' {
		suffix := pattern[1:]
		return len(name) >= len(suffix) && name[len(name)-len(suffix):] == suffix
	}

	// Exact match
	return name == pattern
}

// parseDateTime converts SQL Server run_date and run_time to time.Time.
func parseDateTime(runDate, runTime int) time.Time {
	// run_date format: YYYYMMDD
	// run_time format: HHMMSS

	year := runDate / 10000
	month := (runDate % 10000) / 100
	day := runDate % 100

	hour := runTime / 10000
	minute := (runTime % 10000) / 100
	second := runTime % 100

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// buildConnectionString builds a SQL Server connection string.
func buildConnectionString(server config.ServerConfig) string {
	query := url.Values{}
	query.Add("database", server.Database)
	query.Add("encrypt", strconv.FormatBool(server.Options.Encrypt))
	query.Add("TrustServerCertificate", strconv.FormatBool(server.Options.TrustServerCertificate))
	query.Add("connection timeout", strconv.Itoa(server.Options.ConnectionTimeout))

	u := &url.URL{
		Scheme:   "sqlserver",
		Host:     fmt.Sprintf("%s:%d", server.Host, server.Port),
		RawQuery: query.Encode(),
	}

	// Set authentication
	if server.Auth.Type == "sql" {
		u.User = url.UserPassword(server.Auth.Username, server.Auth.Password)
	}
	// For Windows auth, no user info in URL (driver uses Windows credentials)

	return u.String()
}

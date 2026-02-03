package database

import (
	"testing"
	"time"

	"github.com/hoangtran1411/watchman/internal/config"
)

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name      string
		runDate   int
		runTime   int
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantHour  int
		wantMin   int
		wantSec   int
	}{
		{
			name:      "normal datetime",
			runDate:   20260203,
			runTime:   83015,
			wantYear:  2026,
			wantMonth: time.February,
			wantDay:   3,
			wantHour:  8,
			wantMin:   30,
			wantSec:   15,
		},
		{
			name:      "midnight",
			runDate:   20260101,
			runTime:   0,
			wantYear:  2026,
			wantMonth: time.January,
			wantDay:   1,
			wantHour:  0,
			wantMin:   0,
			wantSec:   0,
		},
		{
			name:      "end of day",
			runDate:   20261231,
			runTime:   235959,
			wantYear:  2026,
			wantMonth: time.December,
			wantDay:   31,
			wantHour:  23,
			wantMin:   59,
			wantSec:   59,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDateTime(tt.runDate, tt.runTime)

			if got.Year() != tt.wantYear {
				t.Errorf("Year = %d, want %d", got.Year(), tt.wantYear)
			}
			if got.Month() != tt.wantMonth {
				t.Errorf("Month = %v, want %v", got.Month(), tt.wantMonth)
			}
			if got.Day() != tt.wantDay {
				t.Errorf("Day = %d, want %d", got.Day(), tt.wantDay)
			}
			if got.Hour() != tt.wantHour {
				t.Errorf("Hour = %d, want %d", got.Hour(), tt.wantHour)
			}
			if got.Minute() != tt.wantMin {
				t.Errorf("Minute = %d, want %d", got.Minute(), tt.wantMin)
			}
			if got.Second() != tt.wantSec {
				t.Errorf("Second = %d, want %d", got.Second(), tt.wantSec)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		jobName string
		pattern string
		want    bool
	}{
		{
			name:    "exact match",
			jobName: "Backup_Database",
			pattern: "Backup_Database",
			want:    true,
		},
		{
			name:    "exact no match",
			jobName: "Backup_Database",
			pattern: "ETL_Job",
			want:    false,
		},
		{
			name:    "prefix wildcard match",
			jobName: "test_job_1",
			pattern: "test_*",
			want:    true,
		},
		{
			name:    "prefix wildcard no match",
			jobName: "prod_job_1",
			pattern: "test_*",
			want:    false,
		},
		{
			name:    "suffix wildcard match",
			jobName: "job_backup",
			pattern: "*_backup",
			want:    true,
		},
		{
			name:    "suffix wildcard no match",
			jobName: "job_restore",
			pattern: "*_backup",
			want:    false,
		},
		{
			name:    "match all",
			jobName: "any_job",
			pattern: "*",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPattern(tt.jobName, tt.pattern)
			if got != tt.want {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.jobName, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name    string
		server  config.ServerConfig
		jobName string
		want    bool
	}{
		{
			name: "no filters - allow all",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{},
					Exclude: []string{},
				},
			},
			jobName: "any_job",
			want:    true,
		},
		{
			name: "include filter match",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{"ETL_*"},
					Exclude: []string{},
				},
			},
			jobName: "ETL_Daily",
			want:    true,
		},
		{
			name: "include filter no match",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{"ETL_*"},
					Exclude: []string{},
				},
			},
			jobName: "Backup_Daily",
			want:    false,
		},
		{
			name: "exclude filter match",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{},
					Exclude: []string{"test_*"},
				},
			},
			jobName: "test_job",
			want:    false,
		},
		{
			name: "exclude filter no match",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{},
					Exclude: []string{"test_*"},
				},
			},
			jobName: "prod_job",
			want:    true,
		},
		{
			name: "include and exclude",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{"ETL_*"},
					Exclude: []string{"ETL_test_*"},
				},
			},
			jobName: "ETL_Daily",
			want:    true,
		},
		{
			name: "include match but exclude also match",
			server: config.ServerConfig{
				Jobs: config.JobsFilter{
					Include: []string{"ETL_*"},
					Exclude: []string{"ETL_test_*"},
				},
			},
			jobName: "ETL_test_job",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{server: tt.server}
			got := db.matchesFilter(tt.jobName)
			if got != tt.want {
				t.Errorf("matchesFilter(%q) = %v, want %v", tt.jobName, got, tt.want)
			}
		})
	}
}

func TestBuildConnectionString(t *testing.T) {
	server := config.ServerConfig{
		Host:     "localhost",
		Port:     1433,
		Database: "msdb",
		Auth: config.AuthConfig{
			Type:     "sql",
			Username: "sa",
			Password: "test123",
		},
		Options: config.DBOptions{
			Encrypt:                false,
			TrustServerCertificate: true,
			ConnectionTimeout:      30,
		},
	}

	connStr := buildConnectionString(server)

	// Should contain basic parts
	if connStr == "" {
		t.Error("connection string is empty")
	}

	// Should start with sqlserver://
	if len(connStr) < 12 || connStr[:12] != "sqlserver://" {
		t.Errorf("connection string should start with 'sqlserver://', got: %s", connStr)
	}
}

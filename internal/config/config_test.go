package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandEnvVar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envName  string
		envValue string
		want     string
	}{
		{
			name:  "no env var",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:     "env var exists",
			input:    "${TEST_VAR}",
			envName:  "TEST_VAR",
			envValue: "secret123",
			want:     "secret123",
		},
		{
			name:  "env var not exists with default",
			input: "${MISSING_VAR:default_value}",
			want:  "default_value",
		},
		{
			name:  "env var not exists no default",
			input: "${MISSING_VAR_NO_DEFAULT}",
			want:  "${MISSING_VAR_NO_DEFAULT}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envName != "" {
				t.Setenv(tt.envName, tt.envValue)
			}

			got := expandEnvVar(tt.input)
			if got != tt.want {
				t.Errorf("expandEnvVar(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConfigValidate_Valid(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "valid config",
			config: Config{
				Servers: []ServerConfig{
					{
						Name:     "TEST-SQL",
						Enabled:  true,
						Host:     "localhost",
						Port:     1433,
						Database: "msdb",
						Auth:     AuthConfig{Type: "sql", Username: "sa", Password: "test"},
					},
				},
				Scheduler: SchedulerConfig{
					CheckTimes: []string{"08:00"},
				},
				Monitoring: MonitoringConfig{
					LookbackHours: 24,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestConfigValidate_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		errMsg string
	}{
		{
			name: "no servers",
			config: Config{
				Servers: []ServerConfig{},
			},
			errMsg: "no servers configured",
		},
		{
			name: "missing server name",
			config: Config{
				Servers: []ServerConfig{
					{Host: "localhost", Port: 1433, Auth: AuthConfig{Type: "sql"}},
				},
			},
			errMsg: "name is required",
		},
		{
			name: "invalid port",
			config: Config{
				Servers: []ServerConfig{
					{Name: "TEST", Host: "localhost", Port: 0, Auth: AuthConfig{Type: "sql"}},
				},
			},
			errMsg: "invalid port",
		},
		{
			name: "invalid auth type",
			config: Config{
				Servers: []ServerConfig{
					{Name: "TEST", Host: "localhost", Port: 1433, Auth: AuthConfig{Type: "invalid"}},
				},
			},
			errMsg: "auth type must be",
		},
		{
			name: "invalid check time",
			config: Config{
				Servers: []ServerConfig{
					{Name: "TEST", Host: "localhost", Port: 1433, Auth: AuthConfig{Type: "sql"}},
				},
				Scheduler: SchedulerConfig{
					CheckTimes: []string{"invalid"},
				},
			},
			errMsg: "invalid check time format",
		},
		{
			name: "no check times",
			config: Config{
				Servers: []ServerConfig{
					{Name: "TEST", Host: "localhost", Port: 1433, Auth: AuthConfig{Type: "sql"}},
				},
				Scheduler: SchedulerConfig{
					CheckTimes: []string{},
				},
			},
			errMsg: "no check times configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Errorf("Validate() expected error containing %q, got nil", tt.errMsg)
				return
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, want substring %q", err, tt.errMsg)
			}
		})
	}
}

func TestGetEnabledServers(t *testing.T) {
	cfg := &Config{
		Servers: []ServerConfig{
			{Name: "ENABLED-1", Enabled: true},
			{Name: "DISABLED-1", Enabled: false},
			{Name: "ENABLED-2", Enabled: true},
		},
	}

	enabled := cfg.GetEnabledServers()
	if len(enabled) != 2 {
		t.Errorf("GetEnabledServers() returned %d servers, want 2", len(enabled))
	}

	for _, srv := range enabled {
		if !srv.Enabled {
			t.Errorf("GetEnabledServers() returned disabled server: %s", srv.Name)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Scheduler.CheckTimes[0] != "08:00" {
		t.Errorf("default check time = %q, want %q", cfg.Scheduler.CheckTimes[0], "08:00")
	}

	if cfg.Monitoring.LookbackHours != 24 {
		t.Errorf("default lookback_hours = %d, want 24", cfg.Monitoring.LookbackHours)
	}

	if cfg.Notification.AppID != "Watchman" {
		t.Errorf("default app_id = %q, want %q", cfg.Notification.AppID, "Watchman")
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
servers:
  - name: "TEST-SQL"
    enabled: true
    host: "localhost"
    port: 1433
    database: "msdb"
    auth:
      type: "sql"
      username: "sa"
      password: "test123"
    options:
      encrypt: false
      trust_server_certificate: true
      connection_timeout: 30
      query_timeout: 60

scheduler:
  check_times:
    - "08:00"
  timezone: "Local"

monitoring:
  lookback_hours: 24
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(cfg.Servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(cfg.Servers))
	}

	if cfg.Servers[0].Name != "TEST-SQL" {
		t.Errorf("server name = %q, want %q", cfg.Servers[0].Name, "TEST-SQL")
	}
}

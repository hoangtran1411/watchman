// Package config provides configuration management for Watchman.
// It uses Viper to load YAML configuration with environment variable support.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration.
type Config struct {
	Servers      []ServerConfig     `mapstructure:"servers"`
	Scheduler    SchedulerConfig    `mapstructure:"scheduler"`
	Notification NotificationConfig `mapstructure:"notification"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Monitoring   MonitoringConfig   `mapstructure:"monitoring"`
	Update       UpdateConfig       `mapstructure:"update"`
}

// ServerConfig represents a SQL Server instance configuration.
type ServerConfig struct {
	Name     string     `mapstructure:"name"`
	Enabled  bool       `mapstructure:"enabled"`
	Host     string     `mapstructure:"host"`
	Port     int        `mapstructure:"port"`
	Database string     `mapstructure:"database"`
	Auth     AuthConfig `mapstructure:"auth"`
	Options  DBOptions  `mapstructure:"options"`
	Jobs     JobsFilter `mapstructure:"jobs"`
}

// AuthConfig represents authentication configuration.
type AuthConfig struct {
	Type     string `mapstructure:"type"` // "sql" or "windows"
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// DBOptions represents database connection options.
type DBOptions struct {
	Encrypt                bool `mapstructure:"encrypt"`
	TrustServerCertificate bool `mapstructure:"trust_server_certificate"`
	ConnectionTimeout      int  `mapstructure:"connection_timeout"`
	QueryTimeout           int  `mapstructure:"query_timeout"`
}

// JobsFilter represents job filtering configuration.
type JobsFilter struct {
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

// SchedulerConfig represents scheduler configuration.
type SchedulerConfig struct {
	CheckTimes []string    `mapstructure:"check_times"`
	Timezone   string      `mapstructure:"timezone"`
	Retry      RetryConfig `mapstructure:"retry"`
}

// RetryConfig represents retry configuration.
type RetryConfig struct {
	Enabled      bool `mapstructure:"enabled"`
	MaxAttempts  int  `mapstructure:"max_attempts"`
	DelaySeconds int  `mapstructure:"delay_seconds"`
}

// NotificationConfig represents notification configuration.
type NotificationConfig struct {
	AppID    string         `mapstructure:"app_id"`
	IconPath string         `mapstructure:"icon_path"`
	Grouping GroupingConfig `mapstructure:"grouping"`
	Sound    SoundConfig    `mapstructure:"sound"`
}

// GroupingConfig represents notification grouping configuration.
type GroupingConfig struct {
	Enabled                bool `mapstructure:"enabled"`
	MaxJobsPerNotification int  `mapstructure:"max_jobs_per_notification"`
}

// SoundConfig represents notification sound configuration.
type SoundConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Type    string `mapstructure:"type"`
}

// LoggingConfig represents logging configuration.
type LoggingConfig struct {
	Level    string         `mapstructure:"level"`
	Format   string         `mapstructure:"format"`
	File     FileLogConfig  `mapstructure:"file"`
	EventLog EventLogConfig `mapstructure:"event_log"`
}

// FileLogConfig represents file logging configuration.
type FileLogConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
}

// EventLogConfig represents Windows Event Log configuration.
type EventLogConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Source  string `mapstructure:"source"`
}

// MonitoringConfig represents monitoring configuration.
type MonitoringConfig struct {
	LookbackHours  int            `mapstructure:"lookback_hours"`
	ReportStatuses []string       `mapstructure:"report_statuses"`
	Parallel       ParallelConfig `mapstructure:"parallel"`
}

// ParallelConfig represents parallel checking configuration.
type ParallelConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	MaxConcurrent int  `mapstructure:"max_concurrent"`
}

// UpdateConfig represents auto-update configuration.
type UpdateConfig struct {
	CheckOnStartup    bool   `mapstructure:"check_on_startup"`
	GithubRepo        string `mapstructure:"github_repo"`
	IncludePrerelease bool   `mapstructure:"include_prerelease"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Servers: []ServerConfig{},
		Scheduler: SchedulerConfig{
			CheckTimes: []string{"08:00"},
			Timezone:   "Local",
			Retry: RetryConfig{
				Enabled:      true,
				MaxAttempts:  3,
				DelaySeconds: 60,
			},
		},
		Notification: NotificationConfig{
			AppID: "Watchman",
			Grouping: GroupingConfig{
				Enabled:                true,
				MaxJobsPerNotification: 5,
			},
			Sound: SoundConfig{
				Enabled: true,
				Type:    "default",
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			File: FileLogConfig{
				Enabled:    true,
				Path:       "logs/watchman.log",
				MaxSizeMB:  10,
				MaxBackups: 5,
				MaxAgeDays: 30,
				Compress:   true,
			},
			EventLog: EventLogConfig{
				Enabled: true,
				Source:  "Watchman",
			},
		},
		Monitoring: MonitoringConfig{
			LookbackHours:  24,
			ReportStatuses: []string{"failed"},
			Parallel: ParallelConfig{
				Enabled:       true,
				MaxConcurrent: 5,
			},
		},
		Update: UpdateConfig{
			CheckOnStartup:    true,
			GithubRepo:        "hoangtran1411/watchman",
			IncludePrerelease: false,
		},
	}
}

// Load loads configuration from file.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Determine config path
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Set config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Enable environment variable substitution
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal to struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand environment variables in passwords
	for i := range cfg.Servers {
		cfg.Servers[i].Auth.Password = expandEnvVar(cfg.Servers[i].Auth.Password)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Check for at least one server
	if len(c.Servers) == 0 {
		return fmt.Errorf("no servers configured")
	}

	// Validate servers
	for i, srv := range c.Servers {
		if srv.Name == "" {
			return fmt.Errorf("server[%d]: name is required", i)
		}
		if srv.Host == "" {
			return fmt.Errorf("server[%d] (%s): host is required", i, srv.Name)
		}
		if srv.Port <= 0 || srv.Port > 65535 {
			return fmt.Errorf("server[%d] (%s): invalid port: %d", i, srv.Name, srv.Port)
		}
		if srv.Auth.Type != "sql" && srv.Auth.Type != "windows" {
			return fmt.Errorf("server[%d] (%s): auth type must be 'sql' or 'windows'", i, srv.Name)
		}
	}

	// Validate scheduler
	if len(c.Scheduler.CheckTimes) == 0 {
		return fmt.Errorf("no check times configured")
	}
	for _, t := range c.Scheduler.CheckTimes {
		if _, err := time.Parse("15:04", t); err != nil {
			return fmt.Errorf("invalid check time format: %s (expected HH:MM)", t)
		}
	}

	// Validate monitoring
	if c.Monitoring.LookbackHours <= 0 {
		return fmt.Errorf("lookback_hours must be positive")
	}

	return nil
}

// GetEnabledServers returns only enabled servers.
func (c *Config) GetEnabledServers() []ServerConfig {
	var enabled []ServerConfig
	for _, srv := range c.Servers {
		if srv.Enabled {
			enabled = append(enabled, srv)
		}
	}
	return enabled
}

// GetLocation returns the timezone location.
func (c *Config) GetLocation() (*time.Location, error) {
	if c.Scheduler.Timezone == "" || c.Scheduler.Timezone == "Local" {
		return time.Local, nil
	}
	loc, err := time.LoadLocation(c.Scheduler.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone '%s': %w", c.Scheduler.Timezone, err)
	}
	return loc, nil
}

// setDefaults sets default values in viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("scheduler.check_times", []string{"08:00"})
	v.SetDefault("scheduler.timezone", "Local")
	v.SetDefault("scheduler.retry.enabled", true)
	v.SetDefault("scheduler.retry.max_attempts", 3)
	v.SetDefault("scheduler.retry.delay_seconds", 60)

	v.SetDefault("notification.app_id", "Watchman")
	v.SetDefault("notification.grouping.enabled", true)
	v.SetDefault("notification.grouping.max_jobs_per_notification", 5)
	v.SetDefault("notification.sound.enabled", true)
	v.SetDefault("notification.sound.type", "default")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.file.enabled", true)
	v.SetDefault("logging.file.path", "logs/watchman.log")
	v.SetDefault("logging.file.max_size_mb", 10)
	v.SetDefault("logging.file.max_backups", 5)
	v.SetDefault("logging.file.max_age_days", 30)
	v.SetDefault("logging.file.compress", true)
	v.SetDefault("logging.event_log.enabled", true)
	v.SetDefault("logging.event_log.source", "Watchman")

	v.SetDefault("monitoring.lookback_hours", 24)
	v.SetDefault("monitoring.report_statuses", []string{"failed"})
	v.SetDefault("monitoring.parallel.enabled", true)
	v.SetDefault("monitoring.parallel.max_concurrent", 5)

	v.SetDefault("update.check_on_startup", true)
	v.SetDefault("update.github_repo", "hoangtran1411/watchman")
	v.SetDefault("update.include_prerelease", false)
}

// getDefaultConfigPath returns the default config file path.
func getDefaultConfigPath() string {
	// Try current directory first
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}

	// Try ProgramData
	programData := os.Getenv("ProgramData")
	if programData != "" {
		configPath := filepath.Join(programData, "Watchman", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Default to ProgramData path (even if doesn't exist)
	if programData != "" {
		return filepath.Join(programData, "Watchman", "config.yaml")
	}

	return "config.yaml"
}

// expandEnvVar expands environment variables in format ${VAR} or ${VAR:default}.
func expandEnvVar(s string) string {
	if !strings.HasPrefix(s, "${") || !strings.HasSuffix(s, "}") {
		return s
	}

	// Remove ${ and }
	inner := s[2 : len(s)-1]

	// Check for default value
	parts := strings.SplitN(inner, ":", 2)
	varName := parts[0]

	value := os.Getenv(varName)
	if value != "" {
		return value
	}

	// Return default if provided
	if len(parts) > 1 {
		return parts[1]
	}

	return s // Return original if no env var found
}

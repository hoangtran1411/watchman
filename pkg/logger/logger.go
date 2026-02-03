// Package logger provides structured logging for Watchman.
// It uses zerolog for JSON/text logging with file rotation.
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/hoangtran1411/watchman/internal/config"
)

// Logger wraps zerolog.Logger with additional functionality.
type Logger struct {
	zerolog.Logger
	writers []io.Writer
}

// New creates a new logger based on configuration.
func New(cfg config.LoggingConfig) (*Logger, error) {
	var writers []io.Writer

	// Set log level
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Console output (always enabled for development)
	if cfg.Format == "text" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
		writers = append(writers, consoleWriter)
	} else {
		writers = append(writers, os.Stdout)
	}

	// File output
	if cfg.File.Enabled {
		fileWriter, err := newFileWriter(cfg.File)
		if err != nil {
			return nil, err
		}
		writers = append(writers, fileWriter)
	}

	// Create multi-writer
	multi := io.MultiWriter(writers...)

	// Create logger
	logger := zerolog.New(multi).With().Timestamp().Logger()

	return &Logger{
		Logger:  logger,
		writers: writers,
	}, nil
}

// newFileWriter creates a file writer with rotation.
func newFileWriter(cfg config.FileLogConfig) (io.Writer, error) {
	// Ensure log directory exists
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create rotating file writer
	writer := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	return writer, nil
}

// parseLevel parses log level string to zerolog.Level.
func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// WithServer returns a logger with server context.
func (l *Logger) WithServer(serverName string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("server", serverName).Logger(),
	}
}

// WithJob returns a logger with job context.
func (l *Logger) WithJob(jobName string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("job", jobName).Logger(),
	}
}

// LogCheckResult logs the result of a job check.
func (l *Logger) LogCheckResult(serversChecked, serversAvailable, failedJobs int, duration time.Duration) {
	l.Info().
		Int("servers_checked", serversChecked).
		Int("servers_available", serversAvailable).
		Int("failed_jobs", failedJobs).
		Dur("duration", duration).
		Msg("check completed")
}

// LogServerUnavailable logs a server connection failure.
func (l *Logger) LogServerUnavailable(serverName string, err error) {
	l.Warn().
		Str("server", serverName).
		Err(err).
		Msg("server unavailable")
}

// LogFailedJob logs a failed job.
func (l *Logger) LogFailedJob(serverName, jobName string, failedAt time.Time) {
	l.Warn().
		Str("server", serverName).
		Str("job", jobName).
		Time("failed_at", failedAt).
		Msg("job failed")
}

// LogNotificationSent logs a notification being sent.
func (l *Logger) LogNotificationSent(jobCount int) {
	l.Info().
		Int("job_count", jobCount).
		Msg("notification sent")
}

// LogServiceStart logs service start.
func (l *Logger) LogServiceStart(version string) {
	l.Info().
		Str("version", version).
		Msg("service started")
}

// LogServiceStop logs service stop.
func (l *Logger) LogServiceStop() {
	l.Info().Msg("service stopped")
}

// LogConfigReload logs configuration reload.
func (l *Logger) LogConfigReload() {
	l.Info().Msg("configuration reloaded")
}

// LogUpdateAvailable logs available update.
func (l *Logger) LogUpdateAvailable(currentVersion, newVersion string) {
	l.Info().
		Str("current_version", currentVersion).
		Str("new_version", newVersion).
		Msg("update available")
}

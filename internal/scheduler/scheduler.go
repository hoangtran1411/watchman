// Package scheduler provides job scheduling for Watchman.
// It uses gocron to schedule checks at specified times.
package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/hoangtran1411/watchman/internal/config"
)

// Scheduler handles scheduled job checks.
type Scheduler struct {
	scheduler gocron.Scheduler
	cfg       *config.Config
	location  *time.Location
	handler   func(ctx context.Context) error
}

// NewScheduler creates a new scheduler.
func NewScheduler(cfg *config.Config, handler func(ctx context.Context) error) (*Scheduler, error) {
	// Get timezone location
	loc, err := cfg.GetLocation()
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	// Create gocron scheduler
	s, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scheduler{
		scheduler: s,
		cfg:       cfg,
		location:  loc,
		handler:   handler,
	}, nil
}

// Start starts the scheduler.
func (s *Scheduler) Start(ctx context.Context) error {
	// Schedule jobs for each check time
	for _, checkTime := range s.cfg.Scheduler.CheckTimes {
		hour, minute, err := parseTime(checkTime)
		if err != nil {
			return fmt.Errorf("invalid check time %s: %w", checkTime, err)
		}
		if hour < 0 || minute < 0 {
			return fmt.Errorf("time values cannot be negative")
		}

		_, err = s.scheduler.NewJob(
			gocron.DailyJob(1, gocron.NewAtTimes(
				gocron.NewAtTime(uint(hour), uint(minute), 0),
			)),
			gocron.NewTask(s.runCheck, ctx),
			gocron.WithName(fmt.Sprintf("check_%s", checkTime)),
		)
		if err != nil {
			return fmt.Errorf("failed to schedule job for %s: %w", checkTime, err)
		}
	}

	// Start the scheduler
	s.scheduler.Start()
	return nil
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() error {
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown scheduler: %w", err)
	}
	return nil
}

// runCheck runs the handler with retry logic.
func (s *Scheduler) runCheck(ctx context.Context) {
	cfg := s.cfg.Scheduler.Retry

	var lastErr error
	attempts := 1
	if cfg.Enabled {
		attempts = cfg.MaxAttempts
	}

	for i := 0; i < attempts; i++ {
		if err := s.handler(ctx); err != nil {
			lastErr = err
			if cfg.Enabled && i < attempts-1 {
				time.Sleep(time.Duration(cfg.DelaySeconds) * time.Second)
				continue
			}
		}
		return // Success
	}

	// Log error after all retries failed
	if lastErr != nil {
		// TODO: Log error using logger package
		_ = lastErr
	}
}

// NextRun returns the next scheduled run time.
func (s *Scheduler) NextRun() (time.Time, error) {
	jobs := s.scheduler.Jobs()
	if len(jobs) == 0 {
		return time.Time{}, fmt.Errorf("no scheduled jobs")
	}

	var nextRun time.Time
	for _, job := range jobs {
		next, err := job.NextRun()
		if err != nil {
			continue
		}
		if nextRun.IsZero() || next.Before(nextRun) {
			nextRun = next
		}
	}

	if nextRun.IsZero() {
		return time.Time{}, fmt.Errorf("no upcoming runs scheduled")
	}

	return nextRun, nil
}

// parseTime parses a time string in HH:MM format.
func parseTime(s string) (hour, minute int, err error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, 0, fmt.Errorf("format must be HH:MM: %w", err)
	}
	return t.Hour(), t.Minute(), nil
}

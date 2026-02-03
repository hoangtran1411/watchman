// Package service provides Windows Service implementation for Watchman.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"

	"github.com/hoangtran1411/watchman/internal/config"
)

const (
	// ServiceName is the Windows service name.
	ServiceName = "Watchman"

	// ServiceDisplayName is the display name shown in Services.
	ServiceDisplayName = "Watchman - SQL Agent Monitor"

	// ServiceDescription is the service description.
	ServiceDescription = "Monitors SQL Server Agent jobs and sends Windows Toast notifications when jobs fail."
)

// Service represents the Windows service.
type Service struct {
	cfg          *config.Config
	startHandler func(ctx context.Context) error
	stopHandler  func() error
	logger       zerolog.Logger
}

// NewService creates a new Windows service handler.
func NewService(cfg *config.Config, start func(ctx context.Context) error, stop func() error, logger zerolog.Logger) *Service {
	return &Service{
		cfg:          cfg,
		startHandler: start,
		stopHandler:  stop,
		logger:       logger,
	}
}

// Run runs the service.
func (s *Service) Run(isDebug bool) error {
	var err error

	if isDebug {
		// Run in interactive/debug mode
		err = debug.Run(ServiceName, s)
	} else {
		// Run as Windows service
		err = svc.Run(ServiceName, s)
	}

	if err != nil {
		return fmt.Errorf("service failed: %w", err)
	}

	return nil
}

// Execute implements svc.Handler interface.
func (s *Service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	// Report starting status
	changes <- svc.Status{State: svc.StartPending}

	// Create context for the service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the service logic in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.startHandler(ctx)
	}()

	// Report running status
	changes <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}

	// Main service loop
	for {
		select {
		case err := <-errChan:
			if err != nil {
				s.logger.Error().
					Err(err).
					Msg("service start handler failed")
				return false, 1
			}
			return false, 0

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus

			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}

				// Cancel context to signal stop
				cancel()

				// Call stop handler
				if s.stopHandler != nil {
					if err := s.stopHandler(); err != nil {
						s.logger.Warn().
							Err(err).
							Msg("stop handler returned error")
					}
				}

				// Give some time for cleanup
				time.Sleep(time.Second)
				changes <- svc.Status{State: svc.Stopped}
				return false, 0

			default:
				// Ignore unknown commands
			}
		}
	}
}

// IsInteractive checks if running interactively (not as service).
func IsInteractive() (bool, error) {
	isService, err := svc.IsWindowsService()
	if err != nil {
		return false, fmt.Errorf("failed to check service status: %w", err)
	}
	return isService, nil
}

// Install installs the service.
func Install(exePath, configPath string) error {
	// Use Windows sc.exe to install service
	// This is a placeholder - actual implementation would use mgr.Connect()
	return fmt.Errorf("install not implemented - use scripts/install.ps1")
}

// Uninstall removes the service.
func Uninstall() error {
	// Use Windows sc.exe to remove service
	return fmt.Errorf("uninstall not implemented - use scripts/uninstall.ps1")
}

// Start starts the service.
func Start() error {
	return fmt.Errorf("start not implemented - use 'sc.exe start Watchman'")
}

// Stop stops the service.
func Stop() error {
	return fmt.Errorf("stop not implemented - use 'sc.exe stop Watchman'")
}

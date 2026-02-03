package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows/svc"

	"github.com/hoangtran1411/watchman/internal/config"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{}
	start := func(ctx context.Context) error { return nil }
	stop := func() error { return nil }

	s := NewService(cfg, start, stop)
	assert.NotNil(t, s)
	assert.Equal(t, cfg, s.cfg)
}

func TestExecute_Lifecycle(t *testing.T) {
	// Setup channels to simulate Windows Service Manager
	reqChan := make(chan svc.ChangeRequest)
	statusChan := make(chan svc.Status, 5) // Buffer to prevent blocking

	// Handlers
	startCalled := false
	stopCalled := false
	start := func(ctx context.Context) error {
		startCalled = true
		<-ctx.Done() // Block until canceled
		return nil
	}
	stop := func() error {
		stopCalled = true
		return nil
	}

	s := NewService(&config.Config{}, start, stop)

	// Run Execute in a goroutine
	done := make(chan bool)
	go func() {
		s.Execute([]string{}, reqChan, statusChan)
		done <- true
	}()

	// Verify StartPending
	status := <-statusChan
	assert.Equal(t, svc.StartPending, status.State)

	// Verify Running
	status = <-statusChan
	assert.Equal(t, svc.Running, status.State)

	// Wait a bit to ensure start handler is called
	time.Sleep(100 * time.Millisecond)
	assert.True(t, startCalled)

	// Send Stop command
	reqChan <- svc.ChangeRequest{Cmd: svc.Stop, CurrentStatus: status}

	// Verify StopPending
	status = <-statusChan
	assert.Equal(t, svc.StopPending, status.State)

	// Verify Stopped
	status = <-statusChan
	assert.Equal(t, svc.Stopped, status.State)

	// Ensure Execute returns
	<-done
	assert.True(t, stopCalled)
}

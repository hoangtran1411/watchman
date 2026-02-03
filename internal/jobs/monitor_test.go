package jobs

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangtran1411/watchman/internal/config"
	"github.com/hoangtran1411/watchman/internal/database"
)

// MockJobQuerier is a mock implementation of JobQuerier.
type MockJobQuerier struct {
	mock.Mock
}

func (m *MockJobQuerier) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock: %w", err)
	}
	return nil
}

func (m *MockJobQuerier) Close() error {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return fmt.Errorf("mock: %w", err)
	}
	return nil
}

func (m *MockJobQuerier) QueryFailedJobs(ctx context.Context, lookbackHours int) ([]database.FailedJob, error) {
	args := m.Called(ctx, lookbackHours)
	err := args.Error(1)
	if err != nil {
		err = fmt.Errorf("mock: %w", err)
	}
	return args.Get(0).([]database.FailedJob), err
}

func TestCheckAll(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			LookbackHours: 24,
			Parallel: config.ParallelConfig{
				Enabled:       false, // Sequential for easier testing
				MaxConcurrent: 1,
			},
		},
		Servers: []config.ServerConfig{
			{Name: "Server1", Enabled: true},
			{Name: "Server2", Enabled: true},
		},
	}

	mockDB1 := new(MockJobQuerier)
	mockDB2 := new(MockJobQuerier)

	monitor := NewMonitor(cfg)
	monitor.dbFactory = func(s config.ServerConfig) (JobQuerier, error) {
		if s.Name == "Server1" {
			return mockDB1, nil
		}
		return mockDB2, nil
	}

	// Expectations
	mockDB1.On("Ping", mock.Anything).Return(nil)
	mockDB1.On("QueryFailedJobs", mock.Anything, 24).Return([]database.FailedJob{}, nil)
	mockDB1.On("Close").Return(nil)

	// Server2 has a failed job
	failedJob := database.FailedJob{
		ServerName: "Server2",
		JobName:    "TestJob",
		Status:     0,
		FailedAt:   time.Now(),
	}
	mockDB2.On("Ping", mock.Anything).Return(nil)
	mockDB2.On("QueryFailedJobs", mock.Anything, 24).Return([]database.FailedJob{failedJob}, nil)
	mockDB2.On("Close").Return(nil)

	// Execute
	result, err := monitor.CheckAll(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "failed_jobs", result.Status)
	assert.Equal(t, 2, result.ServersChecked)
	assert.Equal(t, 2, result.ServersAvailable)
	assert.Equal(t, 1, len(result.FailedJobs))
	assert.Equal(t, "Server2", result.FailedJobs[0].ServerName)

	mockDB1.AssertExpectations(t)
	mockDB2.AssertExpectations(t)
}

func TestCheckAll_ConnectionError(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			LookbackHours: 24,
			Parallel:      config.ParallelConfig{Enabled: false},
		},
		Servers: []config.ServerConfig{
			{Name: "Server1", Enabled: true},
		},
	}

	mockDB := new(MockJobQuerier)

	monitor := NewMonitor(cfg)
	monitor.dbFactory = func(s config.ServerConfig) (JobQuerier, error) {
		return mockDB, nil
	}

	// Expectations
	mockDB.On("Ping", mock.Anything).Return(errors.New("connection failed"))
	mockDB.On("Close").Return(nil)

	// Execute
	result, err := monitor.CheckAll(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "error", result.Status)
	assert.Equal(t, 1, result.ServersChecked)
	assert.Equal(t, 0, result.ServersAvailable)
	assert.Equal(t, 1, len(result.ServersUnavailable))
	assert.Equal(t, "Server1", result.ServersUnavailable[0])

	mockDB.AssertExpectations(t)
	// QueryFailedJobs should not be called
	mockDB.AssertNotCalled(t, "QueryFailedJobs", mock.Anything, mock.Anything)
}

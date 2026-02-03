package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangtran1411/watchman/internal/config"
)

func TestNewScheduler(t *testing.T) {
	cfg := &config.Config{
		Scheduler: config.SchedulerConfig{
			Timezone: "UTC",
		},
	}
	handler := func(ctx context.Context) error { return nil }

	s, err := NewScheduler(cfg, handler)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestNewScheduler_InvalidTimezone(t *testing.T) {
	cfg := &config.Config{
		Scheduler: config.SchedulerConfig{
			Timezone: "Invalid/Timezone",
		},
	}
	handler := func(ctx context.Context) error { return nil }

	s, err := NewScheduler(cfg, handler)
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestStart_InvalidTime(t *testing.T) {
	cfg := &config.Config{
		Scheduler: config.SchedulerConfig{
			CheckTimes: []string{"25:00"},
			Timezone:   "UTC",
		},
	}
	handler := func(ctx context.Context) error { return nil }

	s, err := NewScheduler(cfg, handler)
	assert.NoError(t, err)

	err = s.Start(context.Background())
	assert.Error(t, err)
}

// Mocking function execution for retry test
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRunCheck_Retry(t *testing.T) {
	cfg := &config.Config{
		Scheduler: config.SchedulerConfig{
			Retry: config.RetryConfig{
				Enabled:      true,
				MaxAttempts:  3,
				DelaySeconds: 0, // Instant retry for test
			},
			Timezone: "UTC",
		},
	}

	mockHandler := new(MockHandler)
	// Fail twice, then succeed
	mockHandler.On("Handle", mock.Anything).Return(errors.New("fail 1")).Once()
	mockHandler.On("Handle", mock.Anything).Return(errors.New("fail 2")).Once()
	mockHandler.On("Handle", mock.Anything).Return(nil).Once()

	s, _ := NewScheduler(cfg, mockHandler.Handle)

	s.runCheck(context.Background())

	mockHandler.AssertNumberOfCalls(t, "Handle", 3)
}

func TestRunCheck_NoRetry(t *testing.T) {
	cfg := &config.Config{
		Scheduler: config.SchedulerConfig{
			Retry: config.RetryConfig{
				Enabled:     false,
				MaxAttempts: 3,
			},
			Timezone: "UTC",
		},
	}

	mockHandler := new(MockHandler)
	mockHandler.On("Handle", mock.Anything).Return(errors.New("fail")).Once()

	s, _ := NewScheduler(cfg, mockHandler.Handle)

	s.runCheck(context.Background())

	mockHandler.AssertNumberOfCalls(t, "Handle", 1)
}

func TestParseTime(t *testing.T) {
	h, m, err := parseTime("08:30")
	assert.NoError(t, err)
	assert.Equal(t, 8, h)
	assert.Equal(t, 30, m)

	_, _, err = parseTime("invalid")
	assert.Error(t, err)
}

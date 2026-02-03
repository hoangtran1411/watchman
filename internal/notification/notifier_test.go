package notification

import (
	"testing"
	"time"

	"github.com/go-toast/toast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangtran1411/watchman/internal/config"
	"github.com/hoangtran1411/watchman/internal/database"
)

// MockToastPusher is a mock implementation of ToastPusher
type MockToastPusher struct {
	mock.Mock
}

func (m *MockToastPusher) Push(notification toast.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func TestNotifyFailedJobs_NoJobs(t *testing.T) {
	cfg := config.NotificationConfig{}
	pusher := new(MockToastPusher)
	notifier := NewNotifier(cfg)
	notifier.pusher = pusher

	err := notifier.NotifyFailedJobs([]database.FailedJob{})
	assert.NoError(t, err)
	pusher.AssertNotCalled(t, "Push")
}

func TestNotifyFailedJobs_Individual(t *testing.T) {
	cfg := config.NotificationConfig{
		AppID: "TestApp",
		Grouping: config.GroupingConfig{
			Enabled: false,
		},
	}
	pusher := new(MockToastPusher)
	notifier := NewNotifier(cfg)
	notifier.pusher = pusher

	jobs := []database.FailedJob{
		{ServerName: "S1", JobName: "J1", FailedAt: time.Now()},
		{ServerName: "S2", JobName: "J2", FailedAt: time.Now()},
	}

	pusher.On("Push", mock.MatchedBy(func(n toast.Notification) bool {
		return n.AppID == "TestApp" && (n.Title == "❌ Job Failed on S1" || n.Title == "❌ Job Failed on S2")
	})).Return(nil).Times(2)

	err := notifier.NotifyFailedJobs(jobs)
	assert.NoError(t, err)
	pusher.AssertExpectations(t)
}

func TestNotifyFailedJobs_Grouped(t *testing.T) {
	cfg := config.NotificationConfig{
		AppID: "TestApp",
		Grouping: config.GroupingConfig{
			Enabled: true,
		},
	}
	pusher := new(MockToastPusher)
	notifier := NewNotifier(cfg)
	notifier.pusher = pusher

	jobs := []database.FailedJob{
		{ServerName: "S1", JobName: "J1", FailedAt: time.Now()},
		{ServerName: "S1", JobName: "J2", FailedAt: time.Now()},
	}

	pusher.On("Push", mock.MatchedBy(func(n toast.Notification) bool {
		return n.AppID == "TestApp" && n.Title == "❌ 2 SQL Agent Jobs Failed"
	})).Return(nil).Once()

	err := notifier.NotifyFailedJobs(jobs)
	assert.NoError(t, err)
	pusher.AssertExpectations(t)
}

func TestNotifyUpdateAvailable(t *testing.T) {
	cfg := config.NotificationConfig{AppID: "TestApp"}
	pusher := new(MockToastPusher)
	notifier := NewNotifier(cfg)
	notifier.pusher = pusher

	pusher.On("Push", mock.MatchedBy(func(n toast.Notification) bool {
		return n.Title == "🔄 Watchman Update Available"
	})).Return(nil).Once()

	err := notifier.NotifyUpdateAvailable("v1.0.0", "v1.1.0")
	assert.NoError(t, err)
	pusher.AssertExpectations(t)
}

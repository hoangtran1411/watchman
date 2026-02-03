package updater

import (
	"context"
	"runtime"
	"testing"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hoangtran1411/watchman/internal/config"
)

// MockSelfUpdater is a mock implementation of SelfUpdater
type MockSelfUpdater struct {
	mock.Mock
}

func (m *MockSelfUpdater) DetectLatest(slug string) (*selfupdate.Release, bool, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).(*selfupdate.Release), args.Bool(1), args.Error(2)
}

func (m *MockSelfUpdater) UpdateTo(url, cmdPath string) error {
	args := m.Called(url, cmdPath)
	return args.Error(0)
}

func TestCheckForUpdate_Available(t *testing.T) {
	cfg := config.UpdateConfig{GithubRepo: "test/repo"}
	updater := NewUpdater(cfg, "v1.0.0")
	mockSelfUpdater := new(MockSelfUpdater)
	updater.selfUpdater = mockSelfUpdater

	latest := &selfupdate.Release{
		Version: semver.MustParse("1.1.0"),
		URL:     "http://example.com/release",
	}

	mockSelfUpdater.On("DetectLatest", "test/repo").Return(latest, true, nil)

	result, err := updater.CheckForUpdate(context.Background())
	assert.NoError(t, err)
	assert.True(t, result.UpdateAvailable)
	assert.Equal(t, "1.1.0", result.LatestVersion)
}

func TestCheckForUpdate_NotAvailable(t *testing.T) {
	cfg := config.UpdateConfig{GithubRepo: "test/repo"}
	updater := NewUpdater(cfg, "v1.0.0")
	mockSelfUpdater := new(MockSelfUpdater)
	updater.selfUpdater = mockSelfUpdater

	latest := &selfupdate.Release{
		Version: semver.MustParse("1.0.0"),
	}

	mockSelfUpdater.On("DetectLatest", "test/repo").Return(latest, true, nil)

	result, err := updater.CheckForUpdate(context.Background())
	assert.NoError(t, err)
	assert.False(t, result.UpdateAvailable)
}

func TestUpdate_Success(t *testing.T) {
	// Skip on non-windows for now as the logic checks GOOS
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	cfg := config.UpdateConfig{GithubRepo: "test/repo"}
	updater := NewUpdater(cfg, "v1.0.0")
	mockSelfUpdater := new(MockSelfUpdater)
	updater.selfUpdater = mockSelfUpdater

	latest := &selfupdate.Release{
		Version:  semver.MustParse("1.1.0"),
		AssetURL: "http://example.com/asset",
		URL:      "http://example.com/release",
	}

	mockSelfUpdater.On("DetectLatest", "test/repo").Return(latest, true, nil)
	mockSelfUpdater.On("UpdateTo", "http://example.com/asset", "").Return(nil)

	result, err := updater.Update(context.Background())
	assert.NoError(t, err)
	assert.True(t, result.Applied)
	assert.Equal(t, "1.1.0", result.LatestVersion)
}

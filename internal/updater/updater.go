// Package updater provides auto-update functionality for Watchman.
package updater

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/hoangtran1411/watchman/internal/config"
)

// UpdateResult represents the result of an update check.
type UpdateResult struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	UpdateAvailable bool   `json:"update_available"`
	ReleaseURL      string `json:"release_url,omitempty"`
	ReleaseNotes    string `json:"release_notes,omitempty"`
	Applied         bool   `json:"applied"`
	Error           string `json:"error,omitempty"`
}

// SelfUpdater defines the interface for self-update operations.
type SelfUpdater interface {
	DetectLatest(slug string) (*selfupdate.Release, bool, error)
	UpdateTo(url, cmdPath string) error
}

// DefaultSelfUpdater implements SelfUpdater using the selfupdate package.
type DefaultSelfUpdater struct{}

// DetectLatest finds the latest release for the given slug.
func (u *DefaultSelfUpdater) DetectLatest(slug string) (*selfupdate.Release, bool, error) {
	rel, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return nil, false, fmt.Errorf("failed to detect latest release: %w", err)
	}
	return rel, found, nil
}

// UpdateTo applies the update from the given URL.
func (u *DefaultSelfUpdater) UpdateTo(url, cmdPath string) error {
	if err := selfupdate.UpdateTo(url, cmdPath); err != nil {
		return fmt.Errorf("failed to update binary: %w", err)
	}
	return nil
}

// Updater handles auto-update functionality.
type Updater struct {
	cfg            config.UpdateConfig
	currentVersion string
	selfUpdater    SelfUpdater
}

// NewUpdater creates a new updater.
func NewUpdater(cfg config.UpdateConfig, currentVersion string) *Updater {
	return &Updater{
		cfg:            cfg,
		currentVersion: currentVersion,
		selfUpdater:    &DefaultSelfUpdater{},
	}
}

// CheckForUpdate checks if a new version is available.
func (u *Updater) CheckForUpdate(ctx context.Context) (*UpdateResult, error) {
	result := &UpdateResult{
		CurrentVersion: u.currentVersion,
	}

	// Get the latest release
	latest, found, err := u.selfUpdater.DetectLatest(u.cfg.GithubRepo)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	if !found {
		return result, nil
	}

	result.LatestVersion = latest.Version.String()
	result.ReleaseURL = latest.URL
	result.ReleaseNotes = latest.ReleaseNotes

	// Compare versions
	currentVer := cleanVersion(u.currentVersion)
	if currentVer != "" && latest.Version.String() != currentVer {
		result.UpdateAvailable = true
	}

	return result, nil
}

// Update downloads and applies the update.
func (u *Updater) Update(ctx context.Context) (*UpdateResult, error) {
	result := &UpdateResult{
		CurrentVersion: u.currentVersion,
	}

	// Get the latest release
	latest, found, err := u.selfUpdater.DetectLatest(u.cfg.GithubRepo)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	if !found {
		return result, fmt.Errorf("no release found")
	}

	result.LatestVersion = latest.Version.String()
	result.ReleaseURL = latest.URL

	// Check if update is needed
	currentVer := cleanVersion(u.currentVersion)
	if currentVer == latest.Version.String() {
		return result, nil // Already up to date
	}

	result.UpdateAvailable = true

	// Check OS/Arch compatibility
	if runtime.GOOS != "windows" || runtime.GOARCH != "amd64" {
		result.Error = fmt.Sprintf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
		return result, fmt.Errorf("%s", result.Error)
	}

	// Apply update
	if err := u.selfUpdater.UpdateTo(latest.AssetURL, ""); err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Applied = true
	return result, nil
}

// cleanVersion removes 'v' prefix from version string.
func cleanVersion(v string) string {
	if v != "" && v[0] == 'v' {
		return v[1:]
	}
	return v
}

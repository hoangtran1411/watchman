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

// Updater handles auto-update functionality.
type Updater struct {
	cfg            config.UpdateConfig
	currentVersion string
}

// NewUpdater creates a new updater.
func NewUpdater(cfg config.UpdateConfig, currentVersion string) *Updater {
	return &Updater{
		cfg:            cfg,
		currentVersion: currentVersion,
	}
}

// CheckForUpdate checks if a new version is available.
func (u *Updater) CheckForUpdate(ctx context.Context) (*UpdateResult, error) {
	result := &UpdateResult{
		CurrentVersion: u.currentVersion,
	}

	// Get the latest release
	latest, found, err := selfupdate.DetectLatest(ctx, u.cfg.GithubRepo)
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
	latest, found, err := selfupdate.DetectLatest(ctx, u.cfg.GithubRepo)
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
		return result, fmt.Errorf(result.Error)
	}

	// Apply update
	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, ""); err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Applied = true
	return result, nil
}

// cleanVersion removes 'v' prefix from version string.
func cleanVersion(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}
	return v
}

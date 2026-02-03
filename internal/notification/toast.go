// Package notification provides Windows Toast notification support.
package notification

import (
	"fmt"
	"strings"

	"github.com/go-toast/toast"

	"github.com/hoangtran1411/watchman/internal/config"
	"github.com/hoangtran1411/watchman/internal/database"
)

// ToastPusher abstracts the toast notification sending.
type ToastPusher interface {
	Push(notification toast.Notification) error
}

// DefaultToastPusher is the default implementation that sends actual toasts.
type DefaultToastPusher struct{}

// Push sends the toast notification.
func (p *DefaultToastPusher) Push(notification toast.Notification) error {
	if err := notification.Push(); err != nil {
		return fmt.Errorf("failed to push notification: %w", err)
	}
	return nil
}

// Notifier handles Windows Toast notifications.
type Notifier struct {
	cfg    config.NotificationConfig
	pusher ToastPusher
}

// NewNotifier creates a new notification handler.
func NewNotifier(cfg config.NotificationConfig) *Notifier {
	return &Notifier{
		cfg:    cfg,
		pusher: &DefaultToastPusher{},
	}
}

// NotifyFailedJobs sends a notification about failed jobs.
func (n *Notifier) NotifyFailedJobs(jobs []database.FailedJob) error {
	if len(jobs) == 0 {
		return nil
	}

	// Group jobs by server if grouping is enabled
	if n.cfg.Grouping.Enabled {
		return n.sendGroupedNotification(jobs)
	}

	// Send individual notifications
	for _, job := range jobs {
		if err := n.sendSingleNotification(job); err != nil {
			return err
		}
	}

	return nil
}

// sendGroupedNotification sends a single notification for multiple failed jobs.
func (n *Notifier) sendGroupedNotification(jobs []database.FailedJob) error {
	// Group by server
	serverJobs := make(map[string][]database.FailedJob)
	for _, job := range jobs {
		serverJobs[job.ServerName] = append(serverJobs[job.ServerName], job)
	}

	// Build notification content
	title := n.buildTitle(len(jobs), len(serverJobs))
	body := n.buildBody(jobs, serverJobs)

	notification := toast.Notification{
		AppID:   n.cfg.AppID,
		Title:   title,
		Message: body,
	}

	// Set icon if specified
	if n.cfg.IconPath != "" {
		notification.Icon = n.cfg.IconPath
	}

	// Set sound
	// Set sound
	n.setAudio(&notification)

	return n.pusher.Push(notification)
}

// sendSingleNotification sends a notification for a single failed job.
func (n *Notifier) sendSingleNotification(job database.FailedJob) error {
	title := fmt.Sprintf("❌ Job Failed on %s", job.ServerName)
	body := fmt.Sprintf("Job: %s\nFailed at: %s\n%s",
		job.JobName,
		job.FailedAt.Format("2006-01-02 15:04:05"),
		truncateMessage(job.ErrorMessage, 100),
	)

	notification := toast.Notification{
		AppID:   n.cfg.AppID,
		Title:   title,
		Message: body,
	}

	if n.cfg.IconPath != "" {
		notification.Icon = n.cfg.IconPath
	}

	n.setAudio(&notification)

	return n.pusher.Push(notification)
}

// buildTitle builds the notification title.
func (n *Notifier) buildTitle(jobCount, serverCount int) string {
	if jobCount == 1 {
		return "❌ SQL Agent Job Failed"
	}

	if serverCount == 1 {
		return fmt.Sprintf("❌ %d SQL Agent Jobs Failed", jobCount)
	}

	return fmt.Sprintf("❌ %d Jobs Failed on %d Servers", jobCount, serverCount)
}

// buildBody builds the notification body.
func (n *Notifier) buildBody(jobs []database.FailedJob, serverJobs map[string][]database.FailedJob) string {
	var lines []string
	maxJobs := n.cfg.Grouping.MaxJobsPerNotification
	if maxJobs <= 0 {
		maxJobs = 5
	}

	shown := 0
	for server, srvJobs := range serverJobs {
		lines = append(lines, fmt.Sprintf("🖥️ %s:", server))

		for _, job := range srvJobs {
			if shown >= maxJobs {
				remaining := len(jobs) - shown
				if remaining > 0 {
					lines = append(lines, fmt.Sprintf("... and %d more", remaining))
				}
				break
			}
			lines = append(lines, fmt.Sprintf("  • %s", job.JobName))
			shown++
		}

		if shown >= maxJobs {
			break
		}
	}

	return strings.Join(lines, "\n")
}

// setAudio sets the audio for the notification based on config.
func (n *Notifier) setAudio(notification *toast.Notification) {
	if !n.cfg.Sound.Enabled {
		return
	}

	switch n.cfg.Sound.Type {
	case "mail":
		notification.Audio = toast.Mail
	case "reminder":
		notification.Audio = toast.Reminder
	case "sms":
		notification.Audio = toast.SMS
	case "alarm":
		// toast.Alarm is not available in this version, using Default
		notification.Audio = toast.Default
	case "alarm2":
		// toast.Alarm2 is not available in this version, using Default
		notification.Audio = toast.Default
	default:
		notification.Audio = toast.Default
	}
}

// NotifyUpdateAvailable sends a notification about available update.
func (n *Notifier) NotifyUpdateAvailable(currentVersion, newVersion string) error {
	notification := toast.Notification{
		AppID:   n.cfg.AppID,
		Title:   "🔄 Watchman Update Available",
		Message: fmt.Sprintf("Version %s is available (current: %s)\nRun 'watchman update' to upgrade.", newVersion, currentVersion),
	}

	if n.cfg.IconPath != "" {
		notification.Icon = n.cfg.IconPath
	}

	return n.pusher.Push(notification)
}

// truncateMessage truncates a message to max length.
func truncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

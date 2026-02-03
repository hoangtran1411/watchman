// Package jobs provides job monitoring logic for Watchman.
// It coordinates checking multiple servers for failed jobs.
package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hoangtran1411/watchman/internal/config"
	"github.com/hoangtran1411/watchman/internal/database"
)

// CheckResult represents the result of checking all servers.
type CheckResult struct {
	Status             string               `json:"status"`
	Timestamp          time.Time            `json:"timestamp"`
	ServersChecked     int                  `json:"servers_checked"`
	ServersAvailable   int                  `json:"servers_available"`
	ServersUnavailable []string             `json:"servers_unavailable"`
	FailedJobs         []database.FailedJob `json:"failed_jobs"`
	Summary            string               `json:"summary"`
	Duration           time.Duration        `json:"duration_ms"`
}

// ServerResult represents the result of checking a single server.
type ServerResult struct {
	ServerName string
	Available  bool
	FailedJobs []database.FailedJob
	Error      error
}

// Monitor handles job monitoring operations.
type Monitor struct {
	cfg *config.Config
}

// NewMonitor creates a new job monitor.
func NewMonitor(cfg *config.Config) *Monitor {
	return &Monitor{cfg: cfg}
}

// CheckAll checks all enabled servers for failed jobs.
func (m *Monitor) CheckAll(ctx context.Context) (*CheckResult, error) {
	startTime := time.Now()
	servers := m.cfg.GetEnabledServers()

	if len(servers) == 0 {
		return &CheckResult{
			Status:    "error",
			Timestamp: startTime,
			Summary:   "No enabled servers configured",
		}, nil
	}

	// Check servers (parallel or sequential based on config)
	var results []ServerResult
	if m.cfg.Monitoring.Parallel.Enabled {
		results = m.checkParallel(ctx, servers)
	} else {
		results = m.checkSequential(ctx, servers)
	}

	// Aggregate results
	return m.aggregateResults(startTime, results), nil
}

// CheckServer checks a single server for failed jobs.
func (m *Monitor) CheckServer(ctx context.Context, serverName string) (*CheckResult, error) {
	startTime := time.Now()

	// Find server config
	var serverCfg *config.ServerConfig
	for _, srv := range m.cfg.Servers {
		if srv.Name == serverName {
			serverCfg = &srv
			break
		}
	}

	if serverCfg == nil {
		return nil, fmt.Errorf("server not found: %s", serverName)
	}

	result := m.checkSingleServer(ctx, *serverCfg)
	return m.aggregateResults(startTime, []ServerResult{result}), nil
}

// checkParallel checks servers in parallel with concurrency limit.
func (m *Monitor) checkParallel(ctx context.Context, servers []config.ServerConfig) []ServerResult {
	maxConcurrent := m.cfg.Monitoring.Parallel.MaxConcurrent
	if maxConcurrent <= 0 {
		maxConcurrent = 5
	}

	// Semaphore for limiting concurrency
	sem := make(chan struct{}, maxConcurrent)
	results := make([]ServerResult, len(servers))
	var wg sync.WaitGroup

	for i, srv := range servers {
		wg.Add(1)
		go func(idx int, server config.ServerConfig) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			results[idx] = m.checkSingleServer(ctx, server)
		}(i, srv)
	}

	wg.Wait()
	return results
}

// checkSequential checks servers one by one.
func (m *Monitor) checkSequential(ctx context.Context, servers []config.ServerConfig) []ServerResult {
	results := make([]ServerResult, 0, len(servers))

	for _, srv := range servers {
		result := m.checkSingleServer(ctx, srv)
		results = append(results, result)
	}

	return results
}

// checkSingleServer checks a single server for failed jobs.
func (m *Monitor) checkSingleServer(ctx context.Context, server config.ServerConfig) ServerResult {
	result := ServerResult{
		ServerName: server.Name,
	}

	// Create database connection
	db, err := database.New(server)
	if err != nil {
		result.Error = err
		return result
	}
	defer db.Close()

	// Ping to check connectivity
	if err := db.Ping(ctx); err != nil {
		result.Error = err
		return result
	}

	result.Available = true

	// Query failed jobs
	jobs, err := db.QueryFailedJobs(ctx, m.cfg.Monitoring.LookbackHours)
	if err != nil {
		result.Error = err
		return result
	}

	result.FailedJobs = jobs
	return result
}

// aggregateResults aggregates results from all servers.
func (m *Monitor) aggregateResults(startTime time.Time, results []ServerResult) *CheckResult {
	cr := &CheckResult{
		Status:             "success",
		Timestamp:          startTime,
		ServersChecked:     len(results),
		ServersUnavailable: []string{},
		FailedJobs:         []database.FailedJob{},
	}

	for _, r := range results {
		if r.Available {
			cr.ServersAvailable++
			cr.FailedJobs = append(cr.FailedJobs, r.FailedJobs...)
		} else {
			cr.ServersUnavailable = append(cr.ServersUnavailable, r.ServerName)
		}
	}

	// Generate summary
	cr.Summary = m.generateSummary(cr)
	cr.Duration = time.Since(startTime)

	// Set status based on results
	if cr.ServersAvailable == 0 && cr.ServersChecked > 0 {
		cr.Status = "error"
	} else if len(cr.FailedJobs) > 0 {
		cr.Status = "failed_jobs"
	}

	return cr
}

// generateSummary generates a human-readable summary.
func (m *Monitor) generateSummary(cr *CheckResult) string {
	if cr.ServersAvailable == 0 && cr.ServersChecked > 0 {
		return fmt.Sprintf("All %d servers unavailable", cr.ServersChecked)
	}

	if len(cr.FailedJobs) == 0 {
		return fmt.Sprintf("No failed jobs on %d servers", cr.ServersAvailable)
	}

	// Count unique servers with failures
	serverMap := make(map[string]struct{})
	for _, job := range cr.FailedJobs {
		serverMap[job.ServerName] = struct{}{}
	}

	jobWord := "job"
	if len(cr.FailedJobs) > 1 {
		jobWord = "jobs"
	}

	serverWord := "server"
	if len(serverMap) > 1 {
		serverWord = "servers"
	}

	return fmt.Sprintf("%d failed %s on %d %s",
		len(cr.FailedJobs), jobWord, len(serverMap), serverWord)
}

// HasFailedJobs returns true if there are failed jobs in the result.
func (cr *CheckResult) HasFailedJobs() bool {
	return len(cr.FailedJobs) > 0
}

// GetExitCode returns the appropriate exit code based on results.
func (cr *CheckResult) GetExitCode() int {
	switch {
	case cr.Status == "error":
		return 3 // Connection error
	case cr.HasFailedJobs():
		return 1 // Failed jobs found
	default:
		return 0 // Success
	}
}

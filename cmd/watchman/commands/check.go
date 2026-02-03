package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command.
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for failed jobs (manual run)",
	Long: `Check for failed SQL Server Agent jobs.

Queries all configured and enabled SQL Server instances for failed 
jobs within the lookback period. By default, shows results in 
human-readable format. Use --output json for machine-readable output.`,
	Example: `  # Check all servers
  watchmen check

  # Check specific server
  watchmen check --server PROD-SQL01

  # Check and send notification
  watchmen check --notify

  # JSON output for scripting/AI Agents
  watchmen check --output json

  # JSON output piped to jq
  watchmen check --output json | jq '.failed_jobs[] | .job_name'

  # Check with custom lookback period
  watchmen check --lookback 48

  # Quiet mode for scripts (check exit code only)
  watchmen check --quiet && echo "No failures" || echo "Has failures"`,
	RunE: runCheck,
}

var (
	checkServer   string
	checkLookback int
	checkNotify   bool
	checkNoColor  bool
)

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&checkServer, "server", "s", "",
		"check specific server only (by name)")
	checkCmd.Flags().IntVar(&checkLookback, "lookback", 0,
		"hours to look back for failures (default: from config)")
	checkCmd.Flags().BoolVar(&checkNotify, "notify", false,
		"send notification if failures found")
	checkCmd.Flags().BoolVar(&checkNoColor, "no-color", false,
		"disable colored output")
}

func runCheck(cmd *cobra.Command, args []string) error {
	// TODO: Implement check logic
	// This is a placeholder that will be implemented in Phase 2

	if isQuiet() {
		return nil
	}

	if getOutput() == "json" {
		result := map[string]interface{}{
			"status":              "success",
			"message":             "Check command not yet implemented",
			"servers_checked":     0,
			"servers_available":   0,
			"servers_unavailable": []string{},
			"failed_jobs":         []interface{}{},
			"summary":             "Not implemented",
		}
		printJSON(result)
		return nil
	}

	fmt.Println("Check command not yet implemented")
	fmt.Println("This will query SQL Server Agent jobs for failures")

	if checkServer != "" {
		fmt.Printf("Server filter: %s\n", checkServer)
	}
	if checkLookback > 0 {
		fmt.Printf("Lookback: %d hours\n", checkLookback)
	}
	if checkNotify {
		fmt.Println("Notification: enabled")
	}

	return nil
}

// Package commands contains all Cobra CLI commands for Watchmen.
package commands

import (
	"github.com/spf13/cobra"
)

// Build info (set by main.go).
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// Global flags.
var (
	cfgFile string
	output  string
	quiet   bool
	verbose bool
)

// SetBuildInfo sets build information from main package.
func SetBuildInfo(v, c, d string) {
	version = v
	commit = c
	buildDate = d
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "watchmen",
	Short: "SQL Server Agent Job Monitor",
	Long: `Watchmen - SQL Server Agent Job Monitor

A Windows service that monitors SQL Server Agent jobs and sends 
Windows Toast notifications when jobs fail.

Features:
  • Multi-server monitoring
  • Scheduled checks (configurable times)
  • Windows Toast notifications with server name
  • Auto-update from GitHub releases
  • AI Agent friendly (JSON output, exit codes)`,
	Example: `  # Check for failed jobs
  watchmen check

  # Check with JSON output (for AI Agents/scripting)
  watchmen check --output json

  # Install service with custom config
  watchmen install --config D:\configs\watchmen.yaml

  # Force update without confirmation
  watchmen update --yes`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"config file path (default \"%ProgramData%\\Watchmen\\config.yaml\")")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text",
		"output format: text, json")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress all output except errors")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose logging")

	// Add exit codes to help
	rootCmd.SetUsageTemplate(rootCmd.UsageTemplate() + `
Exit Codes:
  0  Success / No failed jobs
  1  Failed jobs found (check completed successfully)
  2  Configuration error
  3  Connection error (all servers unreachable)
  4  Internal error
`)
}

// getOutput returns the current output format.
func getOutput() string {
	return output
}

// isQuiet returns whether quiet mode is enabled.
func isQuiet() bool {
	return quiet
}

// getConfigFile returns the config file path.
func getConfigFile() string {
	return cfgFile
}

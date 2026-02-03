// Package main is the entry point for Watchmen CLI.
// Watchmen is a Windows Service that monitors SQL Server Agent jobs
// and sends Windows Toast notifications when jobs fail.
package main

import (
	"os"

	"github.com/hoangtran1411/watchman/cmd/watchman/commands"
)

// Build-time variables (injected via ldflags).
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Pass build info to commands package
	commands.SetBuildInfo(Version, Commit, BuildDate)

	// Execute root command
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

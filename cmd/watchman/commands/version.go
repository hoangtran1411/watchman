package commands

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// VersionInfo holds version information.
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Show version information including build details.

Use --output json for machine-readable output.`,
	Example: `  # Show version
  watchmen version

  # JSON output
  watchmen version --output json`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	info := VersionInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	if getOutput() == "json" {
		printJSON(info)
		return
	}

	// Text output
	fmt.Printf("Watchmen %s\n", info.Version)
	fmt.Printf("  Commit:     %s\n", info.Commit)
	fmt.Printf("  Built:      %s\n", info.BuildDate)
	fmt.Printf("  Go version: %s\n", info.GoVersion)
	fmt.Printf("  OS/Arch:    %s/%s\n", info.OS, info.Arch)
}

// printJSON prints data as JSON.
func printJSON(v interface{}) {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}

// cmd is a reference to access stdout (used by printJSON).
var cmd = rootCmd

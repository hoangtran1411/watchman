package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install as Windows Service",
	Long: `Install Watchmen as a Windows Service.

The service will be configured to start automatically (delayed start)
and will run under the LocalSystem account.`,
	Example: `  # Install with default settings
  watchmen install

  # Install with custom config
  watchmen install --config D:\configs\watchmen.yaml

  # Silent install (no prompts)
  watchmen install --silent`,
	RunE: runInstall,
}

var (
	installSilent bool
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&installSilent, "silent", false,
		"run without prompts (for automation)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// TODO: Implement install logic
	if !isQuiet() {
		fmt.Println("Install command not yet implemented")
		fmt.Println("Use scripts/install.ps1 for now")
	}
	return nil
}

// uninstallCmd represents the uninstall command.
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Windows Service",
	Long: `Remove Watchmen Windows Service.

This will stop the service if running and remove it from Windows.`,
	Example: `  # Interactive uninstall
  watchmen uninstall

  # Keep configuration files
  watchmen uninstall --keep-config

  # Remove everything without prompts
  watchmen uninstall --yes`,
	RunE: runUninstall,
}

var (
	uninstallKeepConfig bool
	uninstallYes        bool
)

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().BoolVar(&uninstallKeepConfig, "keep-config", false,
		"keep configuration and log files")
	uninstallCmd.Flags().BoolVar(&uninstallYes, "yes", false,
		"skip confirmation prompt")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	// TODO: Implement uninstall logic
	if !isQuiet() {
		fmt.Println("Uninstall command not yet implemented")
		fmt.Println("Use scripts/uninstall.ps1 for now")
	}
	return nil
}

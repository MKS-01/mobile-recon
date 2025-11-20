// Package cmd implements all CLI commands using the Cobra framework.
// It provides the root command and shared functionality for all subcommands.
package cmd

import (
	"fmt"
	"os"

	"github.com/mks/adb-toolkit/pkg/adb"
	"github.com/mks/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	// deviceSerial stores the user-specified device serial from the -d/--device flag.
	// If empty, commands will use the first available device.
	deviceSerial string

	// rootCmd is the base command for the CLI application.
	// All other commands (device, app, recon, etc.) are added as subcommands to this.
	rootCmd = &cobra.Command{
		Use:   "adb-toolkit",
		Short: "Advanced ADB toolkit for Android development and reverse engineering",
		Long: `A comprehensive toolkit for Android Debug Bridge (ADB) operations.
Perfect for day-to-day development tasks, debugging, and reverse engineering.`,
		// PersistentPreRun executes before any command runs, validating ADB installation.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !adb.IsADBInstalled() {
				utils.PrintError("ADB is not installed or not in PATH")
				os.Exit(1)
			}
		},
	}
)

// Execute runs the root command and handles any errors.
// This is called from main.go to start the CLI application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.PrintError("%v", err)
		os.Exit(1)
	}
}

// init sets up persistent flags that are available to all subcommands.
// The --device/-d flag allows targeting a specific device when multiple are connected.
func init() {
	rootCmd.PersistentFlags().StringVarP(&deviceSerial, "device", "d", "", "Target device serial number")
}

// getTargetDevice returns the device serial to use for command execution.
// It prioritizes the user-specified --device flag, falling back to the first
// available device if no specific device is specified.
//
// Returns:
//   - string: Device serial number to target
//   - error: Error if no devices are available or cannot be detected
//
// This helper is used by all device-targeting commands to determine which
// device to execute operations on.
func getTargetDevice() (string, error) {
	if deviceSerial != "" {
		return deviceSerial, nil
	}

	device, err := adb.GetDefaultDevice()
	if err != nil {
		return "", fmt.Errorf("failed to get default device: %v", err)
	}

	return device.Serial, nil
}

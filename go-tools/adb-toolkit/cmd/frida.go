// Package cmd/frida implements Frida dynamic instrumentation utilities.
// This file provides helper commands for working with the Frida framework
// on Android devices, including server management and process inspection.
//
// Commands:
//   - frida ps:            List running processes in Frida-friendly format
//   - frida trace:         Generate command for tracing method calls
//   - frida server check:  Verify if Frida server is running on device
//   - frida server start:  Start Frida server on device
//
// Frida is a dynamic instrumentation toolkit for reverse engineering,
// security research, and debugging. It allows runtime manipulation of
// application behavior without modifying the APK.
//
// Prerequisites:
//   - frida-server binary pushed to device at /data/local/tmp/frida-server
//   - Frida tools installed locally (pip install frida-tools)
//   - Device with root access (for most use cases)
package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// fridaCmd is the parent command for Frida-related operations.
var fridaCmd = &cobra.Command{
	Use:   "frida",
	Short: "Frida dynamic instrumentation utilities",
	Long:  "Helper commands for working with Frida framework on Android",
}

// fridaPsCmd lists running processes in a format suitable for Frida operations.
// Attempts to use 'frida-ps' if available, otherwise falls back to standard 'ps'.
var fridaPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running processes (Frida-friendly format)",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection("Running Processes")

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "frida-ps", "-U")
		if err != nil {
			utils.PrintWarning("Frida not found, falling back to ps command")
			output, err = adb.ExecuteCommandWithDevice(serial, "shell", "ps")
			if err != nil {
				utils.PrintError("Failed to list processes: %v", err)
				return
			}
		}

		fmt.Println(output)
	},
}

// fridaTraceCmd provides guidance for using frida-trace to monitor method calls.
// This command prints the appropriate frida-trace command rather than executing it,
// as Frida trace requires interactive session management.
var fridaTraceCmd = &cobra.Command{
	Use:   "trace <package_name> [method_pattern]",
	Short: "Trace method calls in an app",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
		methodPattern := "*"
		if len(args) > 1 {
			methodPattern = args[1]
		}

		utils.PrintInfo("Tracing %s methods in %s...", methodPattern, packageName)
		utils.PrintInfo("This requires Frida to be installed on your system")
		utils.PrintInfo("Run: frida-trace -U -f %s -i '%s'", packageName, methodPattern)
	},
}

// fridaServerCmd is the parent command for Frida server management.
var fridaServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Frida server management",
}

// fridaServerCheckCmd verifies if the Frida server is running on the device.
// Checks for 'frida-server' process in the process list.
var fridaServerCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if Frida server is running",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "ps | grep frida-server")
		if err != nil || output == "" {
			utils.PrintError("Frida server is not running")
			utils.PrintInfo("Push frida-server and run: adb shell '/data/local/tmp/frida-server &'")
			return
		}

		utils.PrintSuccess("Frida server is running")
		fmt.Println(output)
	},
}

// fridaServerStartCmd attempts to start the Frida server on the device.
// Assumes frida-server binary is located at /data/local/tmp/frida-server.
// The server runs in the background to accept connections from Frida clients.
var fridaServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Frida server on device",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintInfo("Starting Frida server...")

		// Try to start frida-server
		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "/data/local/tmp/frida-server &")
		if err != nil {
			utils.PrintError("Failed to start Frida server: %v", err)
			utils.PrintInfo("Make sure frida-server is pushed to /data/local/tmp/")
			return
		}

		utils.PrintSuccess("Frida server started")
	},
}

// init registers all Frida-related subcommands.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(fridaCmd)
	fridaCmd.AddCommand(fridaPsCmd)         // Process listing
	fridaCmd.AddCommand(fridaTraceCmd)      // Method tracing helper
	fridaCmd.AddCommand(fridaServerCmd)     // Server management parent
	fridaServerCmd.AddCommand(fridaServerCheckCmd)  // Check server status
	fridaServerCmd.AddCommand(fridaServerStartCmd)  // Start server
}

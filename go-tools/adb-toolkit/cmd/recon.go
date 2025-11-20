// Package cmd/recon implements reconnaissance and reverse engineering commands.
// This file provides tools for security research, penetration testing,
// and reverse engineering of Android applications.
//
// Commands:
//   - recon logcat:     Monitor and filter device logs in real-time
//   - recon dump:       Extract comprehensive package information via dumpsys
//   - recon activities: List all activities in an application
//   - recon services:   List all services in an application
//   - recon receivers:  List all broadcast receivers in an application
//   - recon files:      List app's private files and directories (requires root)
//   - recon db:         Pull and inspect app databases (requires root)
//   - recon network:    Monitor active network connections
//   - recon processes:  List running processes with optional filtering
//
// Security Note: Many commands require root access for full functionality.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// reconCmd is the parent command for reconnaissance and reverse engineering operations.
var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Reconnaissance and reverse engineering utilities",
	Long:  "Tools for security research, penetration testing, and reverse engineering",
}

var reconLogcatCmd = &cobra.Command{
	Use:   "logcat [filter]",
	Short: "Monitor device logs (filtered)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		clear, _ := cmd.Flags().GetBool("clear")
		save, _ := cmd.Flags().GetString("save")

		if clear {
			utils.PrintInfo("Clearing logcat buffer...")
			adb.ExecuteCommandWithDevice(serial, "logcat", "-c")
			utils.PrintSuccess("Logcat cleared")
			return
		}

		logcatArgs := []string{"logcat"}
		if len(args) > 0 {
			logcatArgs = append(logcatArgs, args[0])
		}

		utils.PrintInfo("Monitoring logcat (Ctrl+C to stop)...")

		if save != "" {
			logcatArgs = append(logcatArgs, "-d")
			output, err := adb.ExecuteCommandWithDevice(serial, logcatArgs...)
			if err != nil {
				utils.PrintError("Logcat failed: %v", err)
				return
			}

			err = os.WriteFile(save, []byte(output), 0644)
			if err != nil {
				utils.PrintError("Failed to save log: %v", err)
				return
			}
			utils.PrintSuccess("Log saved to %s", save)
		} else {
			output, err := adb.ExecuteCommandWithDevice(serial, logcatArgs...)
			if err != nil {
				utils.PrintError("Logcat failed: %v", err)
				return
			}
			fmt.Println(output)
		}
	},
}

var reconDumpCmd = &cobra.Command{
	Use:   "dump <package_name>",
	Short: "Dump comprehensive package information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintSection(fmt.Sprintf("Package Dump: %s", packageName))

		// Dump package info
		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			utils.PrintError("Dump failed: %v", err)
			return
		}

		save, _ := cmd.Flags().GetString("save")
		if save != "" {
			err = os.WriteFile(save, []byte(dump), 0644)
			if err != nil {
				utils.PrintError("Failed to save dump: %v", err)
				return
			}
			utils.PrintSuccess("Dump saved to %s", save)
		} else {
			fmt.Println(dump)
		}
	},
}

var reconActivitiesCmd = &cobra.Command{
	Use:   "activities <package_name>",
	Short: "List all activities in a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintSection(fmt.Sprintf("Activities: %s", packageName))

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			utils.PrintError("Failed to get activities: %v", err)
			return
		}

		lines := strings.Split(dump, "\n")
		inActivitySection := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			if strings.Contains(trimmed, "Activity Resolver Table:") {
				inActivitySection = true
				continue
			}

			if inActivitySection {
				if trimmed == "" || !strings.HasPrefix(line, " ") {
					break
				}

				if strings.Contains(trimmed, packageName) && strings.Contains(trimmed, "/") {
					parts := strings.Fields(trimmed)
					for _, part := range parts {
						if strings.Contains(part, packageName+"/") {
							activity := strings.Split(part, " ")[0]
							fmt.Println(activity)
						}
					}
				}
			}
		}
	},
}

var reconServicesCmd = &cobra.Command{
	Use:   "services <package_name>",
	Short: "List all services in a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintSection(fmt.Sprintf("Services: %s", packageName))

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			utils.PrintError("Failed to get services: %v", err)
			return
		}

		lines := strings.Split(dump, "\n")
		inServiceSection := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			if strings.Contains(trimmed, "Service Resolver Table:") {
				inServiceSection = true
				continue
			}

			if inServiceSection {
				if trimmed == "" || !strings.HasPrefix(line, " ") {
					break
				}

				if strings.Contains(trimmed, packageName) && strings.Contains(trimmed, "/") {
					parts := strings.Fields(trimmed)
					for _, part := range parts {
						if strings.Contains(part, packageName+"/") {
							service := strings.Split(part, " ")[0]
							fmt.Println(service)
						}
					}
				}
			}
		}
	},
}

var reconBroadcastCmd = &cobra.Command{
	Use:   "receivers <package_name>",
	Short: "List all broadcast receivers in a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintSection(fmt.Sprintf("Broadcast Receivers: %s", packageName))

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			utils.PrintError("Failed to get receivers: %v", err)
			return
		}

		lines := strings.Split(dump, "\n")
		inReceiverSection := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			if strings.Contains(trimmed, "Receiver Resolver Table:") {
				inReceiverSection = true
				continue
			}

			if inReceiverSection {
				if trimmed == "" || !strings.HasPrefix(line, " ") {
					break
				}

				if strings.Contains(trimmed, packageName) && strings.Contains(trimmed, "/") {
					parts := strings.Fields(trimmed)
					for _, part := range parts {
						if strings.Contains(part, packageName+"/") {
							receiver := strings.Split(part, " ")[0]
							fmt.Println(receiver)
						}
					}
				}
			}
		}
	},
}

var reconFilesCmd = &cobra.Command{
	Use:   "files <package_name>",
	Short: "List app files and directories",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		appDir := "/data/data/" + packageName

		utils.PrintSection(fmt.Sprintf("App Files: %s", packageName))

		// Check if we have permission
		testCmd := fmt.Sprintf("ls -la %s 2>&1", appDir)
		output, err := adb.ExecuteCommandWithDevice(serial, "shell", testCmd)
		if err != nil || strings.Contains(output, "Permission denied") {
			utils.PrintWarning("Root access required. Attempting root...")
			_, err = adb.ExecuteCommandWithDevice(serial, "root")
			if err != nil {
				utils.PrintError("Failed to get root access. Cannot list app files.")
				utils.PrintInfo("Try running: adb root")
				return
			}
		}

		// List files
		listCmd := fmt.Sprintf("ls -laR %s", appDir)
		files, err := adb.ExecuteCommandWithDevice(serial, "shell", listCmd)
		if err != nil {
			utils.PrintError("Failed to list files: %v", err)
			return
		}

		fmt.Println(files)
	},
}

var reconDbCmd = &cobra.Command{
	Use:   "db <package_name> [database_name]",
	Short: "Pull and inspect app databases",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		dbDir := fmt.Sprintf("/data/data/%s/databases", packageName)

		// List databases if no specific DB provided
		if len(args) == 1 {
			utils.PrintSection(fmt.Sprintf("Databases: %s", packageName))

			// May need root
			adb.ExecuteCommandWithDevice(serial, "root")

			dbs, err := adb.ExecuteCommandWithDevice(serial, "shell", "ls", dbDir)
			if err != nil {
				utils.PrintError("Failed to list databases: %v", err)
				utils.PrintInfo("Make sure device is rooted")
				return
			}

			fmt.Println(dbs)
			return
		}

		// Pull specific database
		dbName := args[1]
		remotePath := fmt.Sprintf("%s/%s", dbDir, dbName)
		localPath := fmt.Sprintf("%s_%s", packageName, dbName)

		utils.PrintInfo("Pulling database %s...", dbName)

		adb.ExecuteCommandWithDevice(serial, "root")
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, localPath)
		if err != nil {
			utils.PrintError("Pull failed: %v", err)
			return
		}

		utils.PrintSuccess("Database saved to %s", localPath)
		utils.PrintInfo("You can inspect it with: sqlite3 %s", localPath)
	},
}

var reconNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Monitor network connections",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection("Network Connections")

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "netstat")
		if err != nil {
			utils.PrintError("Failed to get network info: %v", err)
			return
		}

		fmt.Println(output)
	},
}

var reconProcsCmd = &cobra.Command{
	Use:   "processes [filter]",
	Short: "List running processes",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection("Running Processes")

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "ps")
		if err != nil {
			utils.PrintError("Failed to get processes: %v", err)
			return
		}

		if len(args) > 0 {
			filter := args[0]
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if strings.Contains(line, filter) {
					fmt.Println(line)
				}
			}
		} else {
			fmt.Println(output)
		}
	},
}

// init registers all reconnaissance subcommands and their flags.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(reconCmd)

	// Register logcat command with clear and save options
	reconCmd.AddCommand(reconLogcatCmd)
	reconLogcatCmd.Flags().BoolP("clear", "c", false, "Clear logcat buffer")
	reconLogcatCmd.Flags().StringP("save", "s", "", "Save log to file")

	// Register dump command with save option
	reconCmd.AddCommand(reconDumpCmd)
	reconDumpCmd.Flags().StringP("save", "s", "", "Save dump to file")

	// Register app component analysis commands
	reconCmd.AddCommand(reconActivitiesCmd)  // List activities
	reconCmd.AddCommand(reconServicesCmd)    // List services
	reconCmd.AddCommand(reconBroadcastCmd)   // List broadcast receivers

	// Register file system and data extraction commands
	reconCmd.AddCommand(reconFilesCmd)       // List app files (root required)
	reconCmd.AddCommand(reconDbCmd)          // Pull databases (root required)

	// Register system monitoring commands
	reconCmd.AddCommand(reconNetworkCmd)     // Network connections
	reconCmd.AddCommand(reconProcsCmd)       // Running processes
}

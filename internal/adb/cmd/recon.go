package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MKS-01/mobile-recon/internal/adb"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Reconnaissance and reverse engineering utilities",
	Long:  "Tools for security research, penetration testing, and reverse engineering",
}

// extractComponents pulls component names (pkg/...) from a `dumpsys package`
// resolver table that begins at the given marker line.
func extractComponents(dump, packageName, marker string) []string {
	var comps []string
	inSection := false
	for _, line := range strings.Split(dump, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, marker) {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		if trimmed == "" || !strings.HasPrefix(line, " ") {
			break
		}
		if strings.Contains(trimmed, packageName) && strings.Contains(trimmed, "/") {
			for _, part := range strings.Fields(trimmed) {
				if strings.Contains(part, packageName+"/") {
					comps = append(comps, strings.Split(part, " ")[0])
				}
			}
		}
	}
	return comps
}

// renderComponents prints a component list as JSON or as a text section.
func renderComponents(section string, comps []string) {
	if output.IsJSON() {
		if err := output.JSON(comps); err != nil {
			output.Error("Failed to generate JSON: %v", err)
		}
		return
	}
	output.Section(section)
	for _, c := range comps {
		fmt.Println(c)
	}
}

var reconLogcatCmd = &cobra.Command{
	Use:   "logcat [filter]",
	Short: "Monitor device logs (filtered)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		clear, _ := cmd.Flags().GetBool("clear")
		save, _ := cmd.Flags().GetString("save")

		if clear {
			output.Info("Clearing logcat buffer...")
			adb.ExecuteCommandWithDevice(serial, "logcat", "-c")
			output.Success("Logcat cleared")
			return
		}

		logcatArgs := []string{"logcat"}
		if len(args) > 0 {
			logcatArgs = append(logcatArgs, args[0])
		}

		output.Info("Monitoring logcat (Ctrl+C to stop)...")

		if save != "" {
			logcatArgs = append(logcatArgs, "-d")
			cmdOutput, err := adb.ExecuteCommandWithDevice(serial, logcatArgs...)
			if err != nil {
				output.Error("Logcat failed: %v", err)
				return
			}

			err = os.WriteFile(save, []byte(cmdOutput), 0644)
			if err != nil {
				output.Error("Failed to save log: %v", err)
				return
			}
			output.Success("Log saved to %s", save)
		} else {
			cmdOutput, err := adb.ExecuteCommandWithDevice(serial, logcatArgs...)
			if err != nil {
				output.Error("Logcat failed: %v", err)
				return
			}
			fmt.Println(cmdOutput)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		output.Section(fmt.Sprintf("Package Dump: %s", packageName))

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			output.Error("Dump failed: %v", err)
			return
		}

		save, _ := cmd.Flags().GetString("save")
		if save != "" {
			err = os.WriteFile(save, []byte(dump), 0644)
			if err != nil {
				output.Error("Failed to save dump: %v", err)
				return
			}
			output.Success("Dump saved to %s", save)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			output.Error("Failed to get activities: %v", err)
			return
		}

		renderComponents(fmt.Sprintf("Activities: %s", packageName),
			extractComponents(dump, packageName, "Activity Resolver Table:"))
	},
}

var reconServicesCmd = &cobra.Command{
	Use:   "services <package_name>",
	Short: "List all services in a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		packageName := args[0]

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			output.Error("Failed to get services: %v", err)
			return
		}

		renderComponents(fmt.Sprintf("Services: %s", packageName),
			extractComponents(dump, packageName, "Service Resolver Table:"))
	},
}

var reconBroadcastCmd = &cobra.Command{
	Use:   "receivers <package_name>",
	Short: "List all broadcast receivers in a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		packageName := args[0]

		dump, err := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		if err != nil {
			output.Error("Failed to get receivers: %v", err)
			return
		}

		renderComponents(fmt.Sprintf("Broadcast Receivers: %s", packageName),
			extractComponents(dump, packageName, "Receiver Resolver Table:"))
	},
}

var reconFilesCmd = &cobra.Command{
	Use:   "files <package_name>",
	Short: "List app files and directories",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		appDir := "/data/data/" + packageName

		output.Section(fmt.Sprintf("App Files: %s", packageName))

		testCmd := fmt.Sprintf("ls -la %s 2>&1", appDir)
		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, "shell", testCmd)
		if err != nil || strings.Contains(cmdOutput, "Permission denied") {
			output.Warning("Root access required. Attempting root...")
			_, err = adb.ExecuteCommandWithDevice(serial, "root")
			if err != nil {
				output.Error("Failed to get root access. Cannot list app files.")
				output.Info("Try running: adb root")
				return
			}
		}

		listCmd := fmt.Sprintf("ls -laR %s", appDir)
		files, err := adb.ExecuteCommandWithDevice(serial, "shell", listCmd)
		if err != nil {
			output.Error("Failed to list files: %v", err)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		dbDir := fmt.Sprintf("/data/data/%s/databases", packageName)

		if len(args) == 1 {
			output.Section(fmt.Sprintf("Databases: %s", packageName))

			adb.ExecuteCommandWithDevice(serial, "root")

			dbs, err := adb.ExecuteCommandWithDevice(serial, "shell", "ls", dbDir)
			if err != nil {
				output.Error("Failed to list databases: %v", err)
				output.Info("Make sure device is rooted")
				return
			}

			fmt.Println(dbs)
			return
		}

		dbName := args[1]
		remotePath := fmt.Sprintf("%s/%s", dbDir, dbName)
		localPath := fmt.Sprintf("%s_%s", packageName, dbName)

		output.Info("Pulling database %s...", dbName)

		adb.ExecuteCommandWithDevice(serial, "root")
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, localPath)
		if err != nil {
			output.Error("Pull failed: %v", err)
			return
		}

		output.Success("Database saved to %s", localPath)
		output.Info("You can inspect it with: sqlite3 %s", localPath)
	},
}

var reconNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Monitor network connections",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		output.Section("Network Connections")

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, "shell", "netstat")
		if err != nil {
			output.Error("Failed to get network info: %v", err)
			return
		}

		fmt.Println(cmdOutput)
	},
}

var reconProcsCmd = &cobra.Command{
	Use:   "processes [filter]",
	Short: "List running processes",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, "shell", "ps")
		if err != nil {
			output.Error("Failed to get processes: %v", err)
			return
		}

		lines := strings.Split(strings.TrimRight(cmdOutput, "\n"), "\n")
		if len(args) > 0 {
			filter := args[0]
			var filtered []string
			for _, line := range lines {
				if strings.Contains(line, filter) {
					filtered = append(filtered, line)
				}
			}
			lines = filtered
		}

		if output.IsJSON() {
			if err := output.JSON(lines); err != nil {
				output.Error("Failed to generate JSON: %v", err)
			}
			return
		}

		output.Section("Running Processes")
		for _, line := range lines {
			fmt.Println(line)
		}
	},
}

func init() {
	RootCmd.AddCommand(reconCmd)

	reconCmd.AddCommand(reconLogcatCmd)
	reconLogcatCmd.Flags().BoolP("clear", "c", false, "Clear logcat buffer")
	reconLogcatCmd.Flags().StringP("save", "s", "", "Save log to file")

	reconCmd.AddCommand(reconDumpCmd)
	reconDumpCmd.Flags().StringP("save", "s", "", "Save dump to file")

	reconCmd.AddCommand(reconActivitiesCmd)
	reconCmd.AddCommand(reconServicesCmd)
	reconCmd.AddCommand(reconBroadcastCmd)
	reconCmd.AddCommand(reconFilesCmd)
	reconCmd.AddCommand(reconDbCmd)
	reconCmd.AddCommand(reconNetworkCmd)
	reconCmd.AddCommand(reconProcsCmd)
}

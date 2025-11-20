// Package cmd/app implements application/package management commands.
// This file provides functionality for installing, uninstalling, listing,
// and managing Android applications on connected devices.
//
// Commands:
//   - app list:      List installed packages (with filters for system/third-party apps)
//   - app install:   Install an APK file to the device
//   - app uninstall: Remove an application from the device
//   - app clear:     Clear app data and cache
//   - app info:      Get detailed package information (version, path, permissions)
//   - app start:     Launch an application
//   - app stop:      Force stop a running application
//   - app pull:      Extract APK from device to local machine
package cmd

import (
	"fmt"
	"strings"

	"github.com/mks/adb-toolkit/pkg/adb"
	"github.com/mks/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	includeSystem  bool // Flag to filter for system apps only
	thirdPartyOnly bool // Flag to filter for third-party apps only
)

// appCmd is the parent command for all app/package management operations.
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "App/Package management commands",
	Long:  "Install, uninstall, list, and manage Android applications",
}

// appListCmd lists installed packages on the device.
// Supports filtering by package name and flags for system/third-party apps.
var appListCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "List installed packages",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		shellArgs := []string{"shell", "pm", "list", "packages"}

		if thirdPartyOnly {
			shellArgs = append(shellArgs, "-3")
		} else if includeSystem {
			shellArgs = append(shellArgs, "-s")
		}

		if len(args) > 0 {
			shellArgs = append(shellArgs, args[0])
		}

		output, err := adb.ExecuteCommandWithDevice(serial, shellArgs...)
		if err != nil {
			utils.PrintError("Failed to list packages: %v", err)
			return
		}

		packages := strings.Split(output, "\n")
		utils.PrintSection("Installed Packages")

		for _, pkg := range packages {
			pkg = strings.TrimPrefix(pkg, "package:")
			if pkg != "" {
				fmt.Println(pkg)
			}
		}

		utils.PrintInfo("Total packages: %d", len(packages))
	},
}

var appInstallCmd = &cobra.Command{
	Use:   "install <apk_path>",
	Short: "Install an APK",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		apkPath := args[0]
		reinstall, _ := cmd.Flags().GetBool("reinstall")

		utils.PrintInfo("Installing %s...", apkPath)

		installArgs := []string{"install"}
		if reinstall {
			installArgs = append(installArgs, "-r")
		}
		installArgs = append(installArgs, apkPath)

		output, err := adb.ExecuteCommandWithDevice(serial, installArgs...)
		if err != nil {
			utils.PrintError("Installation failed: %v", err)
			return
		}

		if strings.Contains(output, "Success") {
			utils.PrintSuccess("APK installed successfully")
		} else {
			utils.PrintError("Installation failed: %s", output)
		}
	},
}

var appUninstallCmd = &cobra.Command{
	Use:   "uninstall <package_name>",
	Short: "Uninstall an app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		keepData, _ := cmd.Flags().GetBool("keep-data")

		utils.PrintInfo("Uninstalling %s...", packageName)

		uninstallArgs := []string{"uninstall"}
		if keepData {
			uninstallArgs = append(uninstallArgs, "-k")
		}
		uninstallArgs = append(uninstallArgs, packageName)

		output, err := adb.ExecuteCommandWithDevice(serial, uninstallArgs...)
		if err != nil {
			utils.PrintError("Uninstall failed: %v", err)
			return
		}

		if strings.Contains(output, "Success") {
			utils.PrintSuccess("App uninstalled successfully")
		} else {
			utils.PrintError("Uninstall failed: %s", output)
		}
	},
}

var appClearCmd = &cobra.Command{
	Use:   "clear <package_name>",
	Short: "Clear app data and cache",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintInfo("Clearing data for %s...", packageName)

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "clear", packageName)
		if err != nil {
			utils.PrintError("Clear failed: %v", err)
			return
		}

		if strings.Contains(output, "Success") {
			utils.PrintSuccess("App data cleared successfully")
		} else {
			utils.PrintError("Clear failed: %s", output)
		}
	},
}

var appInfoCmd = &cobra.Command{
	Use:   "info <package_name>",
	Short: "Get detailed package information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintSection(fmt.Sprintf("Package Info: %s", packageName))

		// Get package path
		path, _ := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "path", packageName)
		utils.Info.Printf("APK Path: %s\n", strings.TrimPrefix(path, "package:"))

		// Get package dump
		dump, _ := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)

		// Extract version info
		lines := strings.Split(dump, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "versionCode=") || strings.HasPrefix(line, "versionName=") {
				utils.Info.Println(line)
			}
		}

		// Get permissions
		utils.PrintInfo("\nPermissions:")
		perms, _ := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)
		permLines := strings.Split(perms, "\n")
		inPermSection := false
		for _, line := range permLines {
			if strings.Contains(line, "requested permissions:") {
				inPermSection = true
				continue
			}
			if inPermSection && strings.TrimSpace(line) != "" {
				if strings.HasPrefix(strings.TrimSpace(line), "android.permission") {
					fmt.Println("  " + strings.TrimSpace(line))
				} else if !strings.HasPrefix(strings.TrimSpace(line), " ") {
					break
				}
			}
		}
	},
}

var appStartCmd = &cobra.Command{
	Use:   "start <package_name>",
	Short: "Start/Launch an application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintInfo("Starting %s...", packageName)

		output, err := adb.ExecuteCommandWithDevice(serial, "shell", "monkey", "-p", packageName, "-c", "android.intent.category.LAUNCHER", "1")
		if err != nil {
			utils.PrintError("Failed to start app: %v", err)
			return
		}

		utils.PrintSuccess("App started")
		if cmd.Flag("verbose").Changed {
			fmt.Println(output)
		}
	},
}

var appStopCmd = &cobra.Command{
	Use:   "stop <package_name>",
	Short: "Force stop an application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]
		utils.PrintInfo("Stopping %s...", packageName)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "am", "force-stop", packageName)
		if err != nil {
			utils.PrintError("Failed to stop app: %v", err)
			return
		}

		utils.PrintSuccess("App stopped")
	},
}

var appPullCmd = &cobra.Command{
	Use:   "pull <package_name> [output_path]",
	Short: "Pull APK from device",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		packageName := args[0]

		// Get APK path
		path, err := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "path", packageName)
		if err != nil {
			utils.PrintError("Package not found: %v", err)
			return
		}

		apkPath := strings.TrimPrefix(strings.TrimSpace(path), "package:")

		outputPath := packageName + ".apk"
		if len(args) > 1 {
			outputPath = args[1]
		}

		utils.PrintInfo("Pulling APK from %s...", apkPath)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", apkPath, outputPath)
		if err != nil {
			utils.PrintError("Pull failed: %v", err)
			return
		}

		utils.PrintSuccess("APK saved to %s", outputPath)
	},
}

// init registers all app management subcommands and their flags.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(appCmd)

	// Register app list command with filtering flags
	appCmd.AddCommand(appListCmd)
	appListCmd.Flags().BoolVarP(&thirdPartyOnly, "third-party", "3", false, "Show only third-party apps")
	appListCmd.Flags().BoolVarP(&includeSystem, "system", "s", false, "Show only system apps")

	// Register app install command with reinstall option
	appCmd.AddCommand(appInstallCmd)
	appInstallCmd.Flags().BoolP("reinstall", "r", false, "Reinstall app keeping data")

	// Register app uninstall command with data retention option
	appCmd.AddCommand(appUninstallCmd)
	appUninstallCmd.Flags().BoolP("keep-data", "k", false, "Keep app data after uninstall")

	// Register other app commands
	appCmd.AddCommand(appClearCmd)
	appCmd.AddCommand(appInfoCmd)
	appCmd.AddCommand(appStartCmd)
	appStartCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
	appCmd.AddCommand(appStopCmd)
	appCmd.AddCommand(appPullCmd)
}

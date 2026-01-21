package cmd

import (
	"fmt"
	"strings"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/spf13/cobra"
)

var (
	includeSystem  bool
	thirdPartyOnly bool
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "App/Package management commands",
	Long:  "Install, uninstall, list, and manage Android applications",
}

var appListCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "List installed packages",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
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

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, shellArgs...)
		if err != nil {
			output.Error("Failed to list packages: %v", err)
			return
		}

		packages := strings.Split(cmdOutput, "\n")
		output.Section("Installed Packages")

		for _, pkg := range packages {
			pkg = strings.TrimPrefix(pkg, "package:")
			if pkg != "" {
				fmt.Println(pkg)
			}
		}

		output.Info("Total packages: %d", len(packages))
	},
}

var appInstallCmd = &cobra.Command{
	Use:   "install <apk_path>",
	Short: "Install an APK",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		apkPath := args[0]
		reinstall, _ := cmd.Flags().GetBool("reinstall")

		output.Info("Installing %s...", apkPath)

		installArgs := []string{"install"}
		if reinstall {
			installArgs = append(installArgs, "-r")
		}
		installArgs = append(installArgs, apkPath)

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, installArgs...)
		if err != nil {
			output.Error("Installation failed: %v", err)
			return
		}

		if strings.Contains(cmdOutput, "Success") {
			output.Success("APK installed successfully")
		} else {
			output.Error("Installation failed: %s", cmdOutput)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		keepData, _ := cmd.Flags().GetBool("keep-data")

		output.Info("Uninstalling %s...", packageName)

		uninstallArgs := []string{"uninstall"}
		if keepData {
			uninstallArgs = append(uninstallArgs, "-k")
		}
		uninstallArgs = append(uninstallArgs, packageName)

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, uninstallArgs...)
		if err != nil {
			output.Error("Uninstall failed: %v", err)
			return
		}

		if strings.Contains(cmdOutput, "Success") {
			output.Success("App uninstalled successfully")
		} else {
			output.Error("Uninstall failed: %s", cmdOutput)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		output.Info("Clearing data for %s...", packageName)

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "clear", packageName)
		if err != nil {
			output.Error("Clear failed: %v", err)
			return
		}

		if strings.Contains(cmdOutput, "Success") {
			output.Success("App data cleared successfully")
		} else {
			output.Error("Clear failed: %s", cmdOutput)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		output.Section(fmt.Sprintf("Package Info: %s", packageName))

		path, _ := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "path", packageName)
		output.InfoColor().Printf("APK Path: %s\n", strings.TrimPrefix(path, "package:"))

		dump, _ := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "package", packageName)

		lines := strings.Split(dump, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "versionCode=") || strings.HasPrefix(line, "versionName=") {
				output.InfoColor().Println(line)
			}
		}

		output.Info("\nPermissions:")
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		output.Info("Starting %s...", packageName)

		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, "shell", "monkey", "-p", packageName, "-c", "android.intent.category.LAUNCHER", "1")
		if err != nil {
			output.Error("Failed to start app: %v", err)
			return
		}

		output.Success("App started")
		if cmd.Flag("verbose").Changed {
			fmt.Println(cmdOutput)
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
			output.Error("%v", err)
			return
		}

		packageName := args[0]
		output.Info("Stopping %s...", packageName)

		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "am", "force-stop", packageName)
		if err != nil {
			output.Error("Failed to stop app: %v", err)
			return
		}

		output.Success("App stopped")
	},
}

var appPullCmd = &cobra.Command{
	Use:   "pull <package_name> [output_path]",
	Short: "Pull APK from device",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		packageName := args[0]

		path, err := adb.ExecuteCommandWithDevice(serial, "shell", "pm", "path", packageName)
		if err != nil {
			output.Error("Package not found: %v", err)
			return
		}

		apkPath := strings.TrimPrefix(strings.TrimSpace(path), "package:")

		outputPath := packageName + ".apk"
		if len(args) > 1 {
			outputPath = args[1]
		}

		output.Info("Pulling APK from %s...", apkPath)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", apkPath, outputPath)
		if err != nil {
			output.Error("Pull failed: %v", err)
			return
		}

		output.Success("APK saved to %s", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.AddCommand(appListCmd)
	appListCmd.Flags().BoolVarP(&thirdPartyOnly, "third-party", "3", false, "Show only third-party apps")
	appListCmd.Flags().BoolVarP(&includeSystem, "system", "s", false, "Show only system apps")

	appCmd.AddCommand(appInstallCmd)
	appInstallCmd.Flags().BoolP("reinstall", "r", false, "Reinstall app keeping data")

	appCmd.AddCommand(appUninstallCmd)
	appUninstallCmd.Flags().BoolP("keep-data", "k", false, "Keep app data after uninstall")

	appCmd.AddCommand(appClearCmd)
	appCmd.AddCommand(appInfoCmd)
	appCmd.AddCommand(appStartCmd)
	appStartCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
	appCmd.AddCommand(appStopCmd)
	appCmd.AddCommand(appPullCmd)
}

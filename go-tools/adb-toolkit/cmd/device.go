// Package cmd/device implements device management commands.
// This file contains commands for device listing, information retrieval,
// rebooting, shell access, screenshots, and screen recording.
//
// Commands:
//   - device list:        List all connected Android devices
//   - device info:        Get detailed device information (Android version, model, battery, etc.)
//   - device reboot:      Reboot device (normal, recovery, or bootloader mode)
//   - device shell:       Execute shell commands on the device
//   - device screenshot:  Capture and save a screenshot from the device
//   - device screenrecord: Record the device screen to a video file
package cmd

import (
	"fmt"

	"github.com/mks/adb-toolkit/pkg/adb"
	"github.com/mks/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// deviceCmd is the parent command for all device-related operations.
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device management commands",
	Long:  "Manage connected Android devices, get device info, reboot, and more",
}

// deviceListCmd lists all connected Android devices with their details.
// Output includes serial number, connection state, model, product name, and device codename.
var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all connected devices",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := adb.GetDevices()
		if err != nil {
			utils.PrintError("Failed to get devices: %v", err)
			return
		}

		if len(devices) == 0 {
			utils.PrintWarning("No devices connected")
			return
		}

		utils.PrintSection("Connected Devices")
		w := utils.NewTable()
		fmt.Fprintln(w, "SERIAL\tSTATE\tMODEL\tPRODUCT\tDEVICE")
		fmt.Fprintln(w, "------\t-----\t-----\t-------\t------")

		for _, device := range devices {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				device.Serial, device.State, device.Model, device.Product, device.Device)
		}

		w.Flush()
		fmt.Println()
		utils.PrintInfo("Total devices: %d", len(devices))
	},
}

// deviceInfoCmd retrieves comprehensive device information.
// Displays: Android version, SDK version, brand, model, screen resolution,
// density, and battery status using device properties (getprop) and dumpsys.
var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get detailed device information",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection(fmt.Sprintf("Device Info: %s", serial))

		// Android version
		androidVer, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.build.version.release")
		utils.Info.Printf("Android Version: %s\n", androidVer)

		// SDK version
		sdkVer, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.build.version.sdk")
		utils.Info.Printf("SDK Version: %s\n", sdkVer)

		// Brand and Model
		brand, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.brand")
		model, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.model")
		utils.Info.Printf("Brand: %s\n", brand)
		utils.Info.Printf("Model: %s\n", model)

		// Screen resolution
		resolution, _ := adb.ExecuteCommandWithDevice(serial, "shell", "wm", "size")
		utils.Info.Printf("Resolution: %s\n", resolution)

		// Density
		density, _ := adb.ExecuteCommandWithDevice(serial, "shell", "wm", "density")
		utils.Info.Printf("Density: %s\n", density)

		// Battery info
		battery, _ := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "battery")
		utils.Info.Printf("\nBattery Info:\n%s\n", battery)
	},
}

// deviceRebootCmd reboots the target device.
// Supports optional boot modes: recovery, bootloader.
// Without a mode argument, performs a normal reboot.
var deviceRebootCmd = &cobra.Command{
	Use:   "reboot [mode]",
	Short: "Reboot device (modes: recovery, bootloader)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		mode := ""
		if len(args) > 0 {
			mode = args[0]
		}

		utils.PrintInfo("Rebooting device %s...", serial)

		var output string
		if mode != "" {
			output, err = adb.ExecuteCommandWithDevice(serial, "reboot", mode)
		} else {
			output, err = adb.ExecuteCommandWithDevice(serial, "reboot")
		}

		if err != nil {
			utils.PrintError("Reboot failed: %v", err)
			return
		}

		utils.PrintSuccess("Device rebooted successfully")
		if output != "" {
			fmt.Println(output)
		}
	},
}

// deviceShellCmd executes arbitrary shell commands on the target device.
// Useful for running custom commands not covered by other subcommands.
// Example: adb-toolkit device shell "pm list packages | grep chrome"
var deviceShellCmd = &cobra.Command{
	Use:   "shell [command...]",
	Short: "Execute shell command on device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		shellArgs := append([]string{"shell"}, args...)
		output, err := adb.ExecuteCommandWithDevice(serial, shellArgs...)

		if err != nil {
			utils.PrintError("Shell command failed: %v", err)
			return
		}

		fmt.Println(output)
	},
}

// deviceScreenshotCmd captures a screenshot from the device and saves it locally.
// The screenshot is taken using screencap, stored temporarily on the device,
// then pulled to the local machine and the temporary file is deleted.
// Default filename: screenshot.png
var deviceScreenshotCmd = &cobra.Command{
	Use:   "screenshot [filename]",
	Short: "Take a screenshot and pull it to local machine",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		filename := "screenshot.png"
		if len(args) > 0 {
			filename = args[0]
		}

		remotePath := "/sdcard/screenshot_temp.png"

		utils.PrintInfo("Taking screenshot...")
		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "screencap", "-p", remotePath)
		if err != nil {
			utils.PrintError("Screenshot failed: %v", err)
			return
		}

		utils.PrintInfo("Pulling screenshot to %s...", filename)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, filename)
		if err != nil {
			utils.PrintError("Pull failed: %v", err)
			return
		}

		adb.ExecuteCommandWithDevice(serial, "shell", "rm", remotePath)
		utils.PrintSuccess("Screenshot saved to %s", filename)
	},
}

// deviceScreenRecordCmd records the device screen to a video file.
// Recording continues until manually stopped with Ctrl+C.
// The recording is saved on the device, then pulled locally and cleaned up.
// Default filename: screenrecord.mp4
var deviceScreenRecordCmd = &cobra.Command{
	Use:   "screenrecord [filename]",
	Short: "Record device screen (Ctrl+C to stop)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		filename := "screenrecord.mp4"
		if len(args) > 0 {
			filename = args[0]
		}

		remotePath := "/sdcard/screenrecord_temp.mp4"

		utils.PrintInfo("Recording screen... Press Ctrl+C to stop")
		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "screenrecord", remotePath)
		if err != nil {
			utils.PrintError("Screen recording failed: %v", err)
			return
		}

		utils.PrintInfo("Pulling recording to %s...", filename)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, filename)
		if err != nil {
			utils.PrintError("Pull failed: %v", err)
			return
		}

		adb.ExecuteCommandWithDevice(serial, "shell", "rm", remotePath)
		utils.PrintSuccess("Recording saved to %s", filename)
	},
}

// init registers all device subcommands to the root command.
// This function is automatically called when the package is imported.
func init() {
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceRebootCmd)
	deviceCmd.AddCommand(deviceShellCmd)
	deviceCmd.AddCommand(deviceScreenshotCmd)
	deviceCmd.AddCommand(deviceScreenRecordCmd)
}

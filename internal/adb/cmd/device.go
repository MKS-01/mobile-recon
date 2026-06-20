package cmd

import (
	"fmt"
	"strings"

	"github.com/MKS-01/mobile-recon/internal/adb"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Device management commands",
	Long:  "Manage connected Android devices, get device info, reboot, and more",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all connected devices",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := adb.GetDevices()
		if err != nil {
			output.Error("Failed to get devices: %v", err)
			return
		}

		if output.IsJSON() {
			if err := output.JSON(devices); err != nil {
				output.Error("Failed to generate JSON: %v", err)
			}
			return
		}

		if len(devices) == 0 {
			output.Warning("No devices connected")
			return
		}

		output.Section("Connected Devices")
		w := output.NewTable()
		fmt.Fprintln(w, "SERIAL\tSTATE\tMODEL\tPRODUCT\tDEVICE")
		fmt.Fprintln(w, "------\t-----\t-----\t-------\t------")

		for _, device := range devices {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				device.Serial, device.State, device.Model, device.Product, device.Device)
		}

		w.Flush()
		fmt.Println()
		output.Info("Total devices: %d", len(devices))
	},
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get detailed device information",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		prop := func(args ...string) string {
			v, _ := adb.ExecuteCommandWithDevice(serial, append([]string{"shell"}, args...)...)
			return strings.TrimSpace(v)
		}

		androidVer := prop("getprop", "ro.build.version.release")
		sdkVer := prop("getprop", "ro.build.version.sdk")
		brand := prop("getprop", "ro.product.brand")
		model := prop("getprop", "ro.product.model")
		resolution := prop("wm", "size")
		density := prop("wm", "density")
		battery := prop("dumpsys", "battery")

		if output.IsJSON() {
			payload := map[string]interface{}{
				"serial":          serial,
				"android_version": androidVer,
				"sdk_version":     sdkVer,
				"brand":           brand,
				"model":           model,
				"resolution":      resolution,
				"density":         density,
				"battery":         battery,
			}
			if err := output.JSON(payload); err != nil {
				output.Error("Failed to generate JSON: %v", err)
			}
			return
		}

		output.Section(fmt.Sprintf("Device Info: %s", serial))
		output.InfoColor().Printf("Android Version: %s\n", androidVer)
		output.InfoColor().Printf("SDK Version: %s\n", sdkVer)
		output.InfoColor().Printf("Brand: %s\n", brand)
		output.InfoColor().Printf("Model: %s\n", model)
		output.InfoColor().Printf("Resolution: %s\n", resolution)
		output.InfoColor().Printf("Density: %s\n", density)
		output.InfoColor().Printf("\nBattery Info:\n%s\n", battery)
	},
}

var deviceRebootCmd = &cobra.Command{
	Use:   "reboot [mode]",
	Short: "Reboot device (modes: recovery, bootloader)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		mode := ""
		if len(args) > 0 {
			mode = args[0]
		}

		output.Info("Rebooting device %s...", serial)

		var cmdOutput string
		if mode != "" {
			cmdOutput, err = adb.ExecuteCommandWithDevice(serial, "reboot", mode)
		} else {
			cmdOutput, err = adb.ExecuteCommandWithDevice(serial, "reboot")
		}

		if err != nil {
			output.Error("Reboot failed: %v", err)
			return
		}

		output.Success("Device rebooted successfully")
		if cmdOutput != "" {
			fmt.Println(cmdOutput)
		}
	},
}

var deviceShellCmd = &cobra.Command{
	Use:   "shell [command...]",
	Short: "Execute shell command on device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		shellArgs := append([]string{"shell"}, args...)
		cmdOutput, err := adb.ExecuteCommandWithDevice(serial, shellArgs...)

		if err != nil {
			output.Error("Shell command failed: %v", err)
			return
		}

		fmt.Println(cmdOutput)
	},
}

var deviceScreenshotCmd = &cobra.Command{
	Use:   "screenshot [filename]",
	Short: "Take a screenshot and pull it to local machine",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		filename := "screenshot.png"
		if len(args) > 0 {
			filename = args[0]
		}

		remotePath := "/sdcard/screenshot_temp.png"

		output.Info("Taking screenshot...")
		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "screencap", "-p", remotePath)
		if err != nil {
			output.Error("Screenshot failed: %v", err)
			return
		}

		output.Info("Pulling screenshot to %s...", filename)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, filename)
		if err != nil {
			output.Error("Pull failed: %v", err)
			return
		}

		adb.ExecuteCommandWithDevice(serial, "shell", "rm", remotePath)
		output.Success("Screenshot saved to %s", filename)
	},
}

var deviceScreenRecordCmd = &cobra.Command{
	Use:   "screenrecord [filename]",
	Short: "Record device screen (Ctrl+C to stop)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		filename := "screenrecord.mp4"
		if len(args) > 0 {
			filename = args[0]
		}

		remotePath := "/sdcard/screenrecord_temp.mp4"

		output.Info("Recording screen... Press Ctrl+C to stop")
		_, err = adb.ExecuteCommandWithDevice(serial, "shell", "screenrecord", remotePath)
		if err != nil {
			output.Error("Screen recording failed: %v", err)
			return
		}

		output.Info("Pulling recording to %s...", filename)
		_, err = adb.ExecuteCommandWithDevice(serial, "pull", remotePath, filename)
		if err != nil {
			output.Error("Pull failed: %v", err)
			return
		}

		adb.ExecuteCommandWithDevice(serial, "shell", "rm", remotePath)
		output.Success("Recording saved to %s", filename)
	},
}

func init() {
	RootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceRebootCmd)
	deviceCmd.AddCommand(deviceShellCmd)
	deviceCmd.AddCommand(deviceScreenshotCmd)
	deviceCmd.AddCommand(deviceScreenRecordCmd)
}

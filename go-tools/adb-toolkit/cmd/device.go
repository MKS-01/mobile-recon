package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/common/output"
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

		output.Section(fmt.Sprintf("Device Info: %s", serial))

		androidVer, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.build.version.release")
		output.InfoColor().Printf("Android Version: %s\n", androidVer)

		sdkVer, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.build.version.sdk")
		output.InfoColor().Printf("SDK Version: %s\n", sdkVer)

		brand, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.brand")
		model, _ := adb.ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.model")
		output.InfoColor().Printf("Brand: %s\n", brand)
		output.InfoColor().Printf("Model: %s\n", model)

		resolution, _ := adb.ExecuteCommandWithDevice(serial, "shell", "wm", "size")
		output.InfoColor().Printf("Resolution: %s\n", resolution)

		density, _ := adb.ExecuteCommandWithDevice(serial, "shell", "wm", "density")
		output.InfoColor().Printf("Density: %s\n", density)

		battery, _ := adb.ExecuteCommandWithDevice(serial, "shell", "dumpsys", "battery")
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
	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
	deviceCmd.AddCommand(deviceRebootCmd)
	deviceCmd.AddCommand(deviceShellCmd)
	deviceCmd.AddCommand(deviceScreenshotCmd)
	deviceCmd.AddCommand(deviceScreenRecordCmd)
}

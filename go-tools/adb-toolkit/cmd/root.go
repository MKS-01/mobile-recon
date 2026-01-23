package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/spf13/cobra"
)

var (
	deviceSerial string

	rootCmd = &cobra.Command{
		Use:   "adb-toolkit",
		Short: "Advanced ADB toolkit for Android development and reverse engineering",
		Long: `A comprehensive toolkit for Android Debug Bridge (ADB) operations.
Perfect for day-to-day development tasks, debugging, and reverse engineering.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !adb.IsADBInstalled() {
				output.Error("ADB is not installed or not in PATH")
				os.Exit(1)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&deviceSerial, "device", "d", "", "Target device serial number")
}

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

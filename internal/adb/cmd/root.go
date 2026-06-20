package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/internal/adb"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var (
	deviceSerial string

	// RootCmd is the root command for ADB toolkit (exported for embedding)
	RootCmd = &cobra.Command{
		Use:   "adb",
		Short: "Android Debug Bridge toolkit",
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
	if err := RootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&deviceSerial, "device", "d", "", "Target device serial number")
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

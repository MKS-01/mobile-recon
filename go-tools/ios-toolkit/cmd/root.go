package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/ios-toolkit/pkg/simctl"
	"github.com/spf13/cobra"
)

var (
	simulatorUDID string

	rootCmd = &cobra.Command{
		Use:   "ios-toolkit",
		Short: "iOS Simulator toolkit for development and security testing",
		Long: `A comprehensive toolkit for iOS Simulator operations.
Perfect for day-to-day development tasks, debugging, and security research.

Key Features:
  - Frida integration for dynamic instrumentation (no jailbreak needed!)
  - Simulator management via xcrun simctl
  - App management and inspection`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !simctl.IsXcodeInstalled() {
				output.Error("Xcode command line tools not installed")
				fmt.Println("\n  Install with: xcode-select --install")
				os.Exit(1)
			}
		},
	}
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&simulatorUDID, "udid", "u", "", "Target simulator UDID")
}

// getTargetSimulator returns the target simulator (from flag or default booted).
func getTargetSimulator() (*simctl.Simulator, error) {
	if simulatorUDID != "" {
		return simctl.GetSimulatorByUDID(simulatorUDID)
	}

	sim, err := simctl.GetDefaultSimulator()
	if err != nil {
		return nil, fmt.Errorf("no target simulator: %v", err)
	}

	return sim, nil
}

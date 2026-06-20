package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/MKS-01/mobile-recon/internal/ios"
	"github.com/spf13/cobra"
)

var (
	simulatorUDID string

	// RootCmd is the root command for iOS toolkit (exported for embedding)
	RootCmd = &cobra.Command{
		Use:   "ios",
		Short: "iOS Simulator toolkit",
		Long: `A comprehensive toolkit for iOS Simulator operations.
Perfect for day-to-day development tasks, debugging, and security research.

Key Features:
  - Frida integration for dynamic instrumentation (no jailbreak needed!)
  - Simulator management via xcrun simctl
  - App management and inspection`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !ios.IsXcodeInstalled() {
				output.Error("Xcode command line tools not installed")
				fmt.Println("\n  Install with: xcode-select --install")
				os.Exit(1)
			}
		},
	}
)

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&simulatorUDID, "udid", "u", "", "Target simulator UDID")
}

// getTargetSimulator returns the target simulator (from flag or default booted).
func getTargetSimulator() (*ios.Simulator, error) {
	if simulatorUDID != "" {
		return ios.GetSimulatorByUDID(simulatorUDID)
	}

	sim, err := ios.GetDefaultSimulator()
	if err != nil {
		return nil, fmt.Errorf("no target simulator: %v", err)
	}

	return sim, nil
}

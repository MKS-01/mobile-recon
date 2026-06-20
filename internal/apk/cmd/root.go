// Package cmd implements all CLI commands using the Cobra framework.
// It provides the root command and shared functionality for all subcommands.
package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var (
	// outputFormat stores the output format flag (text, json)
	outputFormat string

	// verbose enables verbose output
	verbose bool

	// RootCmd is the root command for APK analyzer (exported for embedding)
	RootCmd = &cobra.Command{
		Use:   "apk",
		Short: "Android APK static analysis toolkit",
		Long: `A comprehensive toolkit for static analysis of Android APK files.

Features:
  • Extract and analyze APK metadata
  • Parse AndroidManifest.xml
  • List and extract permissions
  • Extract strings from DEX and native libraries
  • Identify security issues and misconfigurations
  • Search for sensitive patterns (API keys, URLs, etc.)

Perfect for mobile security testing, app reverse engineering, and malware analysis.`,
		Version: "1.0.0",
	}
)

// Execute runs the root command and handles any errors.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

// init sets up persistent flags that are available to all subcommands.
func init() {
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

// validateAPKPath checks if the provided path is a valid APK file.
func validateAPKPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("cannot access file: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	return nil
}

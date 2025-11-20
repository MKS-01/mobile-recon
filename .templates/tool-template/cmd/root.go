// Package cmd implements all CLI commands using the Cobra framework.
// It provides the root command and shared functionality for all subcommands.
package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Color functions for consistent output formatting
	colorError   = color.New(color.FgRed, color.Bold).SprintFunc()
	colorSuccess = color.New(color.FgGreen, color.Bold).SprintFunc()
	colorWarning = color.New(color.FgYellow, color.Bold).SprintFunc()
	colorInfo    = color.New(color.FgCyan).SprintFunc()

	// rootCmd is the base command for the CLI application.
	// All other commands are added as subcommands to this.
	rootCmd = &cobra.Command{
		Use:   "{{TOOL_NAME}}",
		Short: "{{TOOL_SHORT_DESCRIPTION}}",
		Long: `{{TOOL_LONG_DESCRIPTION}}

Features:
  • Feature 1
  • Feature 2
  • Feature 3

Perfect for {{USE_CASE}}.`,
		Version: "{{VERSION}}",
	}
)

// Execute runs the root command and handles any errors.
// This is called from main.go to start the CLI application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError("%v", err)
		os.Exit(1)
	}
}

// init sets up persistent flags that are available to all subcommands.
func init() {
	// Add global flags here if needed
	// rootCmd.PersistentFlags().StringVarP(&variable, "flag", "f", "", "Description")
}

// Helper functions for consistent output formatting

func printError(format string, args ...interface{}) {
	color.New(color.FgRed, color.Bold).Printf("✗ "+format+"\n", args...)
}

func printSuccess(format string, args ...interface{}) {
	color.New(color.FgGreen, color.Bold).Printf("✓ "+format+"\n", args...)
}

func printWarning(format string, args ...interface{}) {
	color.New(color.FgYellow, color.Bold).Printf("⚠ "+format+"\n", args...)
}

func printInfo(format string, args ...interface{}) {
	color.New(color.FgCyan).Printf("ℹ "+format+"\n", args...)
}

// Package utils provides common utility functions for the nmap-toolkit CLI.
package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	// Color functions for consistent output styling
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)
	headerColor  = color.New(color.FgMagenta, color.Bold)
)

// PrintSuccess prints a success message in green
func PrintSuccess(format string, args ...interface{}) {
	successColor.Printf("✓ "+format+"\n", args...)
}

// PrintError prints an error message in red
func PrintError(format string, args ...interface{}) {
	errorColor.Printf("✗ "+format+"\n", args...)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(format string, args ...interface{}) {
	warningColor.Printf("⚠ "+format+"\n", args...)
}

// PrintInfo prints an info message in cyan
func PrintInfo(format string, args ...interface{}) {
	infoColor.Printf("ℹ "+format+"\n", args...)
}

// PrintHeader prints a header message
func PrintHeader(format string, args ...interface{}) {
	headerColor.Printf("\n"+format+"\n", args...)
	fmt.Println(color.MagentaString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}

// PrintData prints regular data output
func PrintData(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// PrintDivider prints a visual divider
func PrintDivider() {
	fmt.Println(color.HiBlackString("────────────────────────────────────────────────────────────────"))
}

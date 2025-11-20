// Package utils provides utility functions for formatted console output.
// It includes colored output for different message types and table formatting.
package utils

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Pre-configured color schemes for different message types.
// These are used throughout the application for consistent visual feedback.
var (
	Success = color.New(color.FgGreen, color.Bold)  // Green bold for success messages
	Error   = color.New(color.FgRed, color.Bold)    // Red bold for error messages
	Warning = color.New(color.FgYellow, color.Bold) // Yellow bold for warnings
	Info    = color.New(color.FgCyan)               // Cyan for informational messages
	Bold    = color.New(color.Bold)                 // Bold for emphasis
)

// PrintSuccess prints a success message in green with a checkmark symbol.
// Typically used to indicate successful completion of operations.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Optional arguments for format string
//
// Example:
//   PrintSuccess("APK installed successfully")
func PrintSuccess(format string, args ...interface{}) {
	Success.Printf("✓ "+format+"\n", args...)
}

// PrintError prints an error message in red with an X symbol.
// Used to display error conditions or failed operations to the user.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Optional arguments for format string
//
// Example:
//   PrintError("Failed to connect to device: %v", err)
func PrintError(format string, args ...interface{}) {
	Error.Printf("✗ "+format+"\n", args...)
}

// PrintWarning prints a warning message in yellow with a warning symbol.
// Used for non-critical issues that the user should be aware of.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Optional arguments for format string
//
// Example:
//   PrintWarning("No devices found, please connect a device")
func PrintWarning(format string, args ...interface{}) {
	Warning.Printf("⚠ "+format+"\n", args...)
}

// PrintInfo prints an informational message in cyan with an info symbol.
// Used for general status updates and non-critical information.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Optional arguments for format string
//
// Example:
//   PrintInfo("Connecting to device %s...", serial)
func PrintInfo(format string, args ...interface{}) {
	Info.Printf("ℹ "+format+"\n", args...)
}

// NewTable creates a new tabwriter for formatted tabular output.
// The table writer automatically aligns columns with 3-space padding.
//
// Returns:
//   - *tabwriter.Writer: Configured table writer that outputs to stdout
//
// Example:
//   w := NewTable()
//   fmt.Fprintln(w, "COLUMN1\tCOLUMN2\tCOLUMN3")
//   w.Flush()
func NewTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
}

// PrintSection prints a visually distinct section header with decorative borders.
// Used to separate different sections of output for better readability.
//
// Parameters:
//   - title: The section title to display
//
// Example:
//   PrintSection("Device Information")
func PrintSection(title string) {
	fmt.Println()
	Bold.Println("═══", title, "═══")
	fmt.Println()
}

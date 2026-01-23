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
func PrintSuccess(format string, args ...interface{}) {
	Success.Printf("✓ "+format+"\n", args...)
}

// PrintError prints an error message in red with an X symbol.
func PrintError(format string, args ...interface{}) {
	Error.Printf("✗ "+format+"\n", args...)
}

// PrintWarning prints a warning message in yellow with a warning symbol.
func PrintWarning(format string, args ...interface{}) {
	Warning.Printf("⚠ "+format+"\n", args...)
}

// PrintInfo prints an informational message in cyan with an info symbol.
func PrintInfo(format string, args ...interface{}) {
	Info.Printf("ℹ "+format+"\n", args...)
}

// NewTable creates a new tabwriter for formatted tabular output.
func NewTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
}

// PrintSection prints a visually distinct section header with decorative borders.
func PrintSection(title string) {
	fmt.Println()
	Bold.Println("═══", title, "═══")
	fmt.Println()
}

// PrintKeyValue prints a key-value pair with formatting.
func PrintKeyValue(key, value string) {
	Bold.Printf("%s: ", key)
	fmt.Println(value)
}

// PrintList prints a list of items with bullet points.
func PrintList(items []string) {
	for _, item := range items {
		fmt.Printf("  • %s\n", item)
	}
}

// Package output provides formatted console output utilities.
package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)
	headerColor  = color.New(color.FgMagenta, color.Bold)
	boldColor    = color.New(color.Bold)
)

// Success prints a green success message with checkmark.
func Success(format string, args ...interface{}) {
	successColor.Printf("✓ "+format+"\n", args...)
}

// Error prints a red error message with X symbol.
func Error(format string, args ...interface{}) {
	errorColor.Printf("✗ "+format+"\n", args...)
}

// Warning prints a yellow warning message.
func Warning(format string, args ...interface{}) {
	warningColor.Printf("⚠ "+format+"\n", args...)
}

// Info prints a cyan informational message.
func Info(format string, args ...interface{}) {
	infoColor.Printf("ℹ "+format+"\n", args...)
}

// Header prints a magenta header with underline.
func Header(format string, args ...interface{}) {
	headerColor.Printf("\n"+format+"\n", args...)
	fmt.Println(color.MagentaString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}

// Section prints a bold section title with decorative borders.
func Section(title string) {
	fmt.Println()
	boldColor.Println("═══", title, "═══")
	fmt.Println()
}

// Data prints plain formatted output.
func Data(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// KeyValue prints a key-value pair with the key emphasized.
func KeyValue(key, value string) {
	boldColor.Printf("%s: ", key)
	fmt.Println(value)
}

// List prints a bulleted list of items.
func List(items []string) {
	for _, item := range items {
		fmt.Printf("  • %s\n", item)
	}
}

// Divider prints a visual separator line.
func Divider() {
	fmt.Println(color.HiBlackString("────────────────────────────────────────────────────────────────"))
}

// NewTable creates a tabwriter for formatted tabular output.
func NewTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
}

// InfoColor returns the info color for direct use.
func InfoColor() *color.Color {
	return infoColor
}

// SuccessColor returns the success color for direct use.
func SuccessColor() *color.Color {
	return successColor
}

// ErrorColor returns the error color for direct use.
func ErrorColor() *color.Color {
	return errorColor
}

// WarningColor returns the warning color for direct use.
func WarningColor() *color.Color {
	return warningColor
}

// BoldColor returns the bold color for direct use.
func BoldColor() *color.Color {
	return boldColor
}

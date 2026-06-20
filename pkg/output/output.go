// Package output provides formatted console output utilities with a global
// text/JSON mode.
//
// In JSON mode, status and decorative messages (Info, Success, Section, …) are
// written to stderr so that stdout carries only the JSON document a command
// emits via JSON(). Errors always go to stderr. This keeps `--json` output
// pipeable (e.g. into jq) while preserving human-readable progress on stderr.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Format selects how command output is rendered.
type Format int

const (
	// FormatText is the default human-readable, colorized output.
	FormatText Format = iota
	// FormatJSON emits a JSON document on stdout; status goes to stderr.
	FormatJSON
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)
	headerColor  = color.New(color.FgMagenta, color.Bold)
	boldColor    = color.New(color.Bold)
)

var (
	format = FormatText
	quiet  bool
)

// Configure sets the global output mode. noColor disables ANSI colors; q
// (quiet) suppresses informational and decorative output (Info, Header,
// Section, Divider). It is typically called once from the root command's
// PersistentPreRun based on global flags.
func Configure(f Format, noColor, q bool) {
	format = f
	quiet = q
	if noColor {
		color.NoColor = true
	}
}

// IsJSON reports whether output is in JSON mode.
func IsJSON() bool { return format == FormatJSON }

// statusOut returns the writer for status/decorative messages: stderr in JSON
// mode (so stdout stays pure JSON), stdout otherwise.
func statusOut() io.Writer {
	if format == FormatJSON {
		return os.Stderr
	}
	return os.Stdout
}

// JSON encodes v as indented JSON to stdout. Commands use this for their result
// payload when IsJSON() is true.
func JSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Success prints a green success message with a checkmark.
func Success(format string, args ...interface{}) {
	successColor.Fprintf(statusOut(), "✓ "+format+"\n", args...)
}

// Error prints a red error message with an X symbol. Always to stderr.
func Error(format string, args ...interface{}) {
	errorColor.Fprintf(os.Stderr, "✗ "+format+"\n", args...)
}

// Warning prints a yellow warning message.
func Warning(format string, args ...interface{}) {
	warningColor.Fprintf(statusOut(), "⚠ "+format+"\n", args...)
}

// Info prints a cyan informational message. Suppressed when quiet.
func Info(format string, args ...interface{}) {
	if quiet {
		return
	}
	infoColor.Fprintf(statusOut(), "ℹ "+format+"\n", args...)
}

// Header prints a magenta header with an underline. Suppressed when quiet.
func Header(format string, args ...interface{}) {
	if quiet {
		return
	}
	w := statusOut()
	headerColor.Fprintf(w, "\n"+format+"\n", args...)
	fmt.Fprintln(w, color.MagentaString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}

// Section prints a bold section title with decorative borders. Suppressed when quiet.
func Section(title string) {
	if quiet {
		return
	}
	w := statusOut()
	fmt.Fprintln(w)
	boldColor.Fprintln(w, "═══", title, "═══")
	fmt.Fprintln(w)
}

// Data prints plain content to stdout (used for human-readable result text).
func Data(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// KeyValue prints a key-value pair with the key emphasized.
func KeyValue(key, value string) {
	w := statusOut()
	boldColor.Fprintf(w, "%s: ", key)
	fmt.Fprintln(w, value)
}

// List prints a bulleted list of items.
func List(items []string) {
	for _, item := range items {
		fmt.Fprintf(statusOut(), "  • %s\n", item)
	}
}

// Divider prints a visual separator line. Suppressed when quiet.
func Divider() {
	if quiet {
		return
	}
	fmt.Fprintln(statusOut(), color.HiBlackString("────────────────────────────────────────────────────────────────"))
}

// NewTable creates a tabwriter for formatted tabular output on stdout.
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

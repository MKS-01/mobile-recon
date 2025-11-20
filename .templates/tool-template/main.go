// Package main is the entry point for the {{TOOL_NAME_TITLE}} CLI application.
// It initializes and executes the Cobra command-line interface.
package main

import (
	"github.com/MKS-01/mobile-recon/go-tools/{{TOOL_NAME}}/cmd"
)

// main is the application entry point that executes the root command.
// It delegates to the cmd package which handles all CLI commands and subcommands.
func main() {
	cmd.Execute()
}

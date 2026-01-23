// Package main is the entry point for the APK Analyzer CLI application.
// It initializes and executes the Cobra command-line interface.
package main

import (
	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/cmd"
)

// main is the application entry point that executes the root command.
// It delegates to the cmd package which handles all CLI commands and subcommands.
func main() {
	cmd.Execute()
}

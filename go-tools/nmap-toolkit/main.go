// Package main is the entry point for the Nmap Toolkit CLI application.
// It provides powerful network reconnaissance capabilities for mobile security testing.
package main

import (
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/cmd"
)

// main is the application entry point that executes the root command.
func main() {
	cmd.Execute()
}

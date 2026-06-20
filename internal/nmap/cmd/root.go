package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/internal/nmap"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool

	// RootCmd is the root command for Nmap toolkit (exported for embedding)
	RootCmd = &cobra.Command{
		Use:   "nmap",
		Short: "Network reconnaissance toolkit",
		Long: `A lightweight toolkit for local network discovery and mobile device reconnaissance.
Focused on host discovery, service detection, and mobile security testing.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !nmap.IsNmapInstalled() {
				output.Error("Nmap is not installed or not in PATH")
				output.Info("Install nmap:")
				output.Data("  macOS: brew install nmap")
				output.Data("  Linux: sudo apt install nmap / sudo yum install nmap")
				output.Data("  Windows: https://nmap.org/download.html")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := nmap.GetVersion()
			output.Header("Nmap Toolkit - Network Reconnaissance Tool")
			output.Info("Nmap Version: %s", version)
			output.Info("Use --help to see available commands")
			fmt.Println()
			cmd.Help()
		},
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

// emitJSON writes a scan result as JSON when in JSON mode and reports whether
// it handled output (so the caller can return before printing text).
func emitJSON(result *nmap.ScanResult) bool {
	if !output.IsJSON() {
		return false
	}
	if err := output.JSON(result); err != nil {
		output.Error("Failed to generate JSON: %v", err)
	}
	return true
}

// streaming reports whether to stream live nmap output. Streaming is disabled
// in JSON mode so the raw report does not corrupt the JSON document.
func streaming() bool {
	return stream && !output.IsJSON()
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output file path")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

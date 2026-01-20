// Package cmd implements all CLI commands using the Cobra framework.
package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFormat string
	verbose      bool

	// rootCmd is the base command for the CLI application
	rootCmd = &cobra.Command{
		Use:   "nmap-toolkit",
		Short: "Advanced Nmap toolkit for network reconnaissance and mobile security testing",
		Long: `A comprehensive toolkit for network reconnaissance using Nmap.
Perfect for mobile security testing, penetration testing, and network analysis.

Features:
  • Host and network discovery
  • Port scanning (TCP, UDP, Stealth)
  • Service and version detection
  • OS fingerprinting
  • Vulnerability scanning with NSE scripts
  • SSL/TLS enumeration
  • Mobile-specific scanning (Android ADB, iOS services)
  • Custom scan templates`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !nmap.IsNmapInstalled() {
				utils.PrintError("Nmap is not installed or not in PATH")
				utils.PrintInfo("Install nmap:")
				utils.PrintData("  macOS: brew install nmap")
				utils.PrintData("  Linux: sudo apt install nmap / sudo yum install nmap")
				utils.PrintData("  Windows: https://nmap.org/download.html")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := nmap.GetVersion()
			utils.PrintHeader("🔍 Nmap Toolkit - Network Reconnaissance Tool")
			utils.PrintInfo("Nmap Version: %s", version)
			utils.PrintInfo("Use --help to see available commands")
			fmt.Println()
			cmd.Help()
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.PrintError("%v", err)
		os.Exit(1)
	}
}

// init sets up persistent flags
func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

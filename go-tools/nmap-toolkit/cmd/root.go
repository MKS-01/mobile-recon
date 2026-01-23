package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool

	rootCmd = &cobra.Command{
		Use:   "nmap-toolkit",
		Short: "Advanced Nmap toolkit for network reconnaissance and mobile security testing",
		Long: `A comprehensive toolkit for network reconnaissance using Nmap.
Perfect for mobile security testing, penetration testing, and network analysis.`,
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
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

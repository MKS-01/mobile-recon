// Package cmd implements all CLI commands
package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	toolMgr *toolmanager.ToolManager

	rootCmd = &cobra.Command{
		Use:   "mobile-recon",
		Short: "Unified CLI for mobile reconnaissance tools",
		Long: `Mobile Recon - A unified command-line interface for mobile security testing
and network reconnaissance tools.

This CLI provides easy access to all reconnaissance tools including:
  • ADB Toolkit - Android device automation and reverse engineering
  • Nmap Toolkit - Network scanning and reconnaissance

Use this tool to discover, build, and run different reconnaissance utilities
from a single interface.`,
		Run: func(cmd *cobra.Command, args []string) {
			printBanner()
			fmt.Println()
			cmd.Help()
		},
	}
)

func init() {
	// Initialize tool manager
	var err error
	toolMgr, err = toolmanager.NewToolManager()
	if err != nil {
		color.Red("Error: Failed to initialize tool manager: %v", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

func printBanner() {
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta)

	cyan.Println("╔════════════════════════════════════════════════════════╗")
	cyan.Println("║                                                        ║")
	cyan.Println("║            📱  MOBILE RECON TOOLKIT  🔍               ║")
	cyan.Println("║                                                        ║")
	cyan.Println("║     Unified CLI for Mobile Security Testing           ║")
	cyan.Println("║                                                        ║")
	cyan.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	magenta.Println("  🎯 Discover, Build, and Run Reconnaissance Tools")
	fmt.Println()
}

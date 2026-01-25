package cmd

import (
	"fmt"
	"os"

	// Import toolkit commands
	adbcmd "github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/cmd"
	apkcmd "github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/cmd"
	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	ioscmd "github.com/MKS-01/mobile-recon/go-tools/ios-toolkit/cmd"
	"github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager"
	nmapcmd "github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/cmd"
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
  - ADB Toolkit - Android device automation and reverse engineering
  - Nmap Toolkit - Network scanning and reconnaissance
  - APK Analyzer - Android APK static analysis
  - iOS Toolkit - iOS simulator management and Frida integration

Run tools directly:
  mobile-recon adb [args...]     Run ADB Toolkit
  mobile-recon nmap [args...]    Run Nmap Toolkit
  mobile-recon apk [args...]     Run APK Analyzer
  mobile-recon ios [args...]     Run iOS Toolkit

Use 'mobile-recon list' to see all available tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			printBanner()
			fmt.Println()
			cmd.Help()
		},
	}
)

func init() {
	var err error
	toolMgr, err = toolmanager.NewToolManager()
	if err != nil {
		output.Error("Failed to initialize tool manager: %v", err)
		os.Exit(1)
	}

	// Register all toolkit commands as subcommands
	rootCmd.AddCommand(adbcmd.RootCmd)
	rootCmd.AddCommand(nmapcmd.RootCmd)
	rootCmd.AddCommand(apkcmd.RootCmd)
	rootCmd.AddCommand(ioscmd.RootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
}

func printBanner() {
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	white := color.New(color.FgHiWhite)

	fmt.Println()
	fmt.Print("  ")
	cyan.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Print("  ")
	cyan.Println("║                                                              ║")
	fmt.Print("  ")
	cyan.Print("║  ")
	white.Print("📱 ")
	cyan.Print("MOBILE RECON ")
	white.Print("- ")
	magenta.Print("Security Testing Toolkit")
	cyan.Println("               ║")
	fmt.Print("  ")
	cyan.Println("║                                                              ║")
	fmt.Print("  ")
	cyan.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Print("  ")
	cyan.Print("║  ")
	yellow.Print("🔧 ")
	white.Print("Mobile Security Testing & Reconnaissance Toolkit")
	cyan.Println("         ║")
	fmt.Print("  ")
	cyan.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show available tools summary
	fmt.Print("  ")
	green.Println("📦 Available Tools:")
	fmt.Println()
	fmt.Print("     ")
	color.New(color.FgGreen).Print("🤖 ")
	fmt.Print("ADB Toolkit    ")
	color.New(color.FgBlue).Print("🌐 ")
	fmt.Print("Nmap Toolkit    ")
	color.New(color.FgYellow).Print("📦 ")
	fmt.Print("APK Analyzer    ")
	color.New(color.FgMagenta).Print("🍎 ")
	fmt.Println("iOS Toolkit")
	fmt.Println()
	fmt.Print("  ")
	color.New(color.FgHiBlack).Println("─────────────────────────────────────────────────────────────")
	fmt.Println()
}

package cli

import (
	"fmt"
	"os"

	// Import toolkit commands
	adbcmd "github.com/MKS-01/mobile-recon/internal/adb/cmd"
	apkcmd "github.com/MKS-01/mobile-recon/internal/apk/cmd"
	"github.com/MKS-01/mobile-recon/pkg/output"
	ioscmd "github.com/MKS-01/mobile-recon/internal/ios/cmd"
	nmapcmd "github.com/MKS-01/mobile-recon/internal/nmap/cmd"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// toolGroups defines the categories used to organize toolkit commands in both
// `--help` and `mobile-recon list`. Order here is the display order.
var toolGroups = []*cobra.Group{
	{ID: "mobile", Title: "Mobile Tools"},
	{ID: "network", Title: "Network Tools"},
}

var (
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
	for _, g := range toolGroups {
		rootCmd.AddGroup(g)
	}

	// Register each toolkit's root command, tagging it with a category group.
	// The toolkit binaries are compiled directly into this CLI, so a registered
	// command is always runnable — there is no separate build/install step.
	register(adbcmd.RootCmd, "mobile")
	register(apkcmd.RootCmd, "mobile")
	register(ioscmd.RootCmd, "mobile")
	register(nmapcmd.RootCmd, "network")
}

func register(cmd *cobra.Command, groupID string) {
	cmd.GroupID = groupID
	rootCmd.AddCommand(cmd)
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

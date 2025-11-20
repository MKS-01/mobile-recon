// Package cmd implements list commands
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	showAll bool

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available tools",
		Long:  `List all reconnaissance tools available in the toolkit, organized by category.`,
		Run: func(cmd *cobra.Command, args []string) {
			listTools()
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tools including those not built")
}

func listTools() {
	headerColor := color.New(color.FgCyan, color.Bold)
	successColor := color.New(color.FgGreen)
	warningColor := color.New(color.FgYellow)
	infoColor := color.New(color.FgHiBlack)

	headerColor.Println("\n📦 Available Reconnaissance Tools")
	fmt.Println()

	hasUnbuilt := false

	for _, category := range toolMgr.Categories {
		color.New(color.FgMagenta, color.Bold).Printf("▶ %s Tools\n", category.DisplayName)
		fmt.Println(color.HiBlackString("  ────────────────────────────────────────"))

		for _, tool := range category.Tools {
			if !showAll && !tool.Available {
				hasUnbuilt = true
				continue
			}

			status := ""
			if tool.Available {
				status = successColor.Sprint("✓ Available")
			} else {
				status = warningColor.Sprint("✗ Not built")
			}

			fmt.Printf("  %-20s %s\n", color.CyanString(tool.DisplayName), status)
			fmt.Printf("  %-20s %s\n", "", infoColor.Sprint(tool.Description))
			fmt.Printf("  %-20s %s %s\n", "", infoColor.Sprint("Command:"), color.HiWhiteString("mobile-recon run %s", tool.Name))
			fmt.Println()
		}
	}

	if hasUnbuilt && !showAll {
		infoColor.Println("💡 Use --all or -a flag to show tools that are not yet built")
		fmt.Println()
	}

	// Show quick start
	fmt.Println()
	headerColor.Println("🚀 Quick Start")
	fmt.Println()
	fmt.Printf("  Build a tool:     %s\n", color.HiWhiteString("mobile-recon build <tool-name>"))
	fmt.Printf("  Build all tools:  %s\n", color.HiWhiteString("mobile-recon build --all"))
	fmt.Printf("  Run a tool:       %s\n", color.HiWhiteString("mobile-recon run <tool-name> [args...]"))
	fmt.Printf("  Interactive mode: %s\n", color.HiWhiteString("mobile-recon interactive"))
	fmt.Println()
}

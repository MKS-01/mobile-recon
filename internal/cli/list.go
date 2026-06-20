package cli

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tools",
	Long:  `List all reconnaissance tools available in the toolkit, organized by category.`,
	Run: func(cmd *cobra.Command, args []string) {
		listTools()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// listTools prints the registered toolkit commands grouped by category. The
// command tree is the single source of truth: every command shown here is
// compiled into this binary and runnable, so there is no "built/not built"
// distinction to report.
func listTools() {
	infoColor := color.New(color.FgHiBlack)

	output.Header("Available Reconnaissance Tools")
	fmt.Println()

	for _, group := range toolGroups {
		color.New(color.FgMagenta, color.Bold).Printf("▶ %s\n", group.Title)
		output.Divider()

		for _, c := range rootCmd.Commands() {
			if c.GroupID != group.ID {
				continue
			}
			fmt.Printf("  %-12s %s\n", color.CyanString(c.Name()), infoColor.Sprint(c.Short))
			fmt.Printf("  %-12s %s %s\n", "", infoColor.Sprint("Command:"), color.HiWhiteString("mobile-recon %s", c.Name()))
			fmt.Println()
		}
	}

	output.Header("Quick Start")
	fmt.Println()
	fmt.Printf("  Run a tool:    %s\n", color.HiWhiteString("mobile-recon <tool> [args...]"))
	fmt.Printf("  Tool help:     %s\n", color.HiWhiteString("mobile-recon <tool> --help"))
	fmt.Println()
}

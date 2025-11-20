// Package cmd implements build commands
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	buildAll bool

	buildCmd = &cobra.Command{
		Use:   "build [tool-name]",
		Short: "Build reconnaissance tools",
		Long: `Build one or more reconnaissance tools.

Examples:
  mobile-recon build adb-toolkit
  mobile-recon build nmap-toolkit
  mobile-recon build --all`,
		Run: func(cmd *cobra.Command, args []string) {
			if buildAll {
				buildAllTools()
				return
			}

			if len(args) == 0 {
				color.Red("Error: Please specify a tool name or use --all flag")
				fmt.Println()
				cmd.Help()
				return
			}

			toolName := args[0]
			buildTool(toolName)
		},
	}

	installCmd = &cobra.Command{
		Use:   "install [tool-name]",
		Short: "Install a tool globally",
		Long: `Install a reconnaissance tool globally so it can be used from anywhere.

Examples:
  mobile-recon install adb-toolkit
  mobile-recon install nmap-toolkit`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			toolName := args[0]
			installTool(toolName)
		},
	}
)

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(installCmd)

	buildCmd.Flags().BoolVarP(&buildAll, "all", "a", false, "Build all tools")
}

func buildTool(toolName string) {
	successColor := color.New(color.FgGreen, color.Bold)
	infoColor := color.New(color.FgCyan)

	tool, err := toolMgr.GetTool(toolName)
	if err != nil {
		color.Red("Error: %v", err)
		fmt.Println()
		color.Yellow("Available tools:")
		for _, t := range toolMgr.ListTools() {
			fmt.Printf("  - %s\n", t.Name)
		}
		return
	}

	infoColor.Printf("\n🔨 Building %s...\n", tool.DisplayName)
	fmt.Println()

	if err := toolMgr.BuildTool(tool.Name); err != nil {
		color.Red("✗ Build failed: %v", err)
		return
	}

	successColor.Printf("✓ Successfully built %s\n", tool.DisplayName)
	fmt.Println()
	color.Cyan("Run with: mobile-recon run %s [args...]", tool.Name)
	fmt.Println()
}

func buildAllTools() {
	headerColor := color.New(color.FgCyan, color.Bold)
	successColor := color.New(color.FgGreen, color.Bold)

	headerColor.Println("\n🔨 Building All Tools...")
	fmt.Println()

	if err := toolMgr.BuildAllTools(); err != nil {
		color.Red("✗ Build failed: %v", err)
		return
	}

	fmt.Println()
	successColor.Println("✓ All tools built successfully!")
	fmt.Println()
}

func installTool(toolName string) {
	successColor := color.New(color.FgGreen, color.Bold)
	infoColor := color.New(color.FgCyan)

	tool, err := toolMgr.GetTool(toolName)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	infoColor.Printf("\n📦 Installing %s globally...\n", tool.DisplayName)
	fmt.Println()

	if err := toolMgr.InstallTool(tool.Name); err != nil {
		color.Red("✗ Installation failed: %v", err)
		return
	}

	successColor.Printf("✓ Successfully installed %s\n", tool.DisplayName)
	fmt.Println()
	color.Cyan("You can now run: %s [args...]", tool.Name)
	fmt.Println()
}

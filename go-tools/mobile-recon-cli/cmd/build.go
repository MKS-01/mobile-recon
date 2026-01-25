package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
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
				output.Error("Please specify a tool name or use --all flag")
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
	tool, err := toolMgr.GetTool(toolName)
	if err != nil {
		output.Error("%v", err)
		fmt.Println()
		output.Warning("Available tools:")
		for _, t := range toolMgr.ListTools() {
			fmt.Printf("  - %s\n", t.Name)
		}
		return
	}

	output.Info("Building %s...", tool.DisplayName)
	fmt.Println()

	if err := toolMgr.BuildTool(tool.Name); err != nil {
		output.Error("Build failed: %v", err)
		return
	}

	output.Success("Successfully built %s", tool.DisplayName)
	fmt.Println()
	output.Info("Run with: mobile-recon %s [args...]", getShortName(tool.Name))
	fmt.Println()
}

func buildAllTools() {
	output.Header("Building All Tools...")
	fmt.Println()

	if err := toolMgr.BuildAllTools(); err != nil {
		output.Error("Build failed: %v", err)
		return
	}

	fmt.Println()
	output.Success("All tools built successfully!")
	fmt.Println()
}

func installTool(toolName string) {
	tool, err := toolMgr.GetTool(toolName)
	if err != nil {
		output.Error("%v", err)
		return
	}

	output.Info("Installing %s globally...", tool.DisplayName)
	fmt.Println()

	if err := toolMgr.InstallTool(tool.Name); err != nil {
		output.Error("Installation failed: %v", err)
		return
	}

	output.Success("Successfully installed %s", tool.DisplayName)
	fmt.Println()
	output.Info("You can now run: %s [args...]", tool.Name)
	fmt.Println()
}

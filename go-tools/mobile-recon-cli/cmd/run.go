package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run [tool-name] [args...]",
		Short: "Run a reconnaissance tool",
		Long: `Run a reconnaissance tool with the specified arguments.

All arguments after the tool name are passed directly to the tool.

Examples:
  mobile-recon run adb-toolkit device list
  mobile-recon run nmap-toolkit scan quick 192.168.1.0/24
  mobile-recon run adb-toolkit --help`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			toolName := args[0]
			toolArgs := args[1:]
			runTool(toolName, toolArgs)
		},
	}

	adbCmd = &cobra.Command{
		Use:                "adb [args...]",
		Short:              "Run ADB Toolkit (shortcut for 'run adb-toolkit')",
		Long:               `Shortcut command to run the ADB Toolkit directly.`,
		Run:                func(cmd *cobra.Command, args []string) { runTool("adb-toolkit", args) },
		DisableFlagParsing: true,
	}

	nmapCmd = &cobra.Command{
		Use:                "nmap [args...]",
		Short:              "Run Nmap Toolkit (shortcut for 'run nmap-toolkit')",
		Long:               `Shortcut command to run the Nmap Toolkit directly.`,
		Run:                func(cmd *cobra.Command, args []string) { runTool("nmap-toolkit", args) },
		DisableFlagParsing: true,
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(adbCmd)
	rootCmd.AddCommand(nmapCmd)
}

func runTool(toolName string, args []string) {
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

	if !tool.Available {
		output.Error("%s is not built yet", tool.DisplayName)
		fmt.Println()
		output.Info("Build it with: mobile-recon build %s", tool.Name)
		return
	}

	if err := toolMgr.RunTool(tool.Name, args); err != nil {
		output.Error("Error running %s: %v", tool.DisplayName, err)
		return
	}
}

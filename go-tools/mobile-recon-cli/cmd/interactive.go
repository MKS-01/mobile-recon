package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode for tool selection",
	Long:  `Launch an interactive menu to select and run tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractive()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractive() {
	for {
		mainPrompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				"List Tools",
				"Build Tool",
				"Run Tool",
				"Browse by Category",
				"Exit",
			},
			Size: 10,
		}

		_, result, err := mainPrompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\nGoodbye!")
				os.Exit(0)
			}
			continue
		}

		switch result {
		case "List Tools":
			listToolsInteractive()
		case "Build Tool":
			buildToolInteractive()
		case "Run Tool":
			runToolInteractive()
		case "Browse by Category":
			browseCategoryInteractive()
		case "Exit":
			fmt.Println("\nGoodbye!")
			os.Exit(0)
		}
	}
}

func listToolsInteractive() {
	fmt.Println()
	listTools()
	fmt.Println()
	promptContinue()
}

func buildToolInteractive() {
	tools := toolMgr.ListTools()
	if len(tools) == 0 {
		output.Warning("No tools available")
		promptContinue()
		return
	}

	items := make([]string, 0)
	items = append(items, "Build All Tools")
	items = append(items, "← Back")

	for _, tool := range tools {
		status := "Not built"
		if tool.Available {
			status = "Built"
		}
		items = append(items, fmt.Sprintf("%s [%s]", tool.DisplayName, status))
	}

	prompt := promptui.Select{
		Label: "Select Tool to Build",
		Items: items,
		Size:  10,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return
	}

	if result == "← Back" {
		return
	}

	if strings.Contains(result, "Build All") {
		buildAllTools()
		promptContinue()
		return
	}

	for _, tool := range tools {
		if strings.Contains(result, tool.DisplayName) {
			buildTool(tool.Name)
			promptContinue()
			return
		}
	}
}

func runToolInteractive() {
	tools := toolMgr.ListAvailableTools()
	if len(tools) == 0 {
		output.Warning("No tools are built yet. Please build tools first.")
		promptContinue()
		return
	}

	items := make([]string, 0)
	items = append(items, "← Back")

	for _, tool := range tools {
		items = append(items, fmt.Sprintf("%s - %s", tool.DisplayName, tool.Description))
	}

	prompt := promptui.Select{
		Label: "Select Tool to Run",
		Items: items,
		Size:  10,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return
	}

	if result == "← Back" {
		return
	}

	for _, tool := range tools {
		if strings.Contains(result, tool.DisplayName) {
			fmt.Printf("\n%s %s\n", color.CyanString("Running:"), color.HiWhiteString(tool.DisplayName))
			fmt.Printf("%s Enter arguments (or press Enter for help):\n", color.YellowString("→"))

			argsPrompt := promptui.Prompt{
				Label: "Args",
			}

			argsInput, _ := argsPrompt.Run()
			args := []string{}
			if argsInput != "" {
				args = strings.Fields(argsInput)
			} else {
				args = []string{"--help"}
			}

			fmt.Println()
			runTool(tool.Name, args)
			fmt.Println()
			promptContinue()
			return
		}
	}
}

func browseCategoryInteractive() {
	items := make([]string, 0)
	items = append(items, "← Back")

	for _, category := range toolMgr.Categories {
		items = append(items, fmt.Sprintf("%s Tools", category.DisplayName))
	}

	prompt := promptui.Select{
		Label: "Select Category",
		Items: items,
		Size:  10,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return
	}

	if result == "← Back" {
		return
	}

	for _, category := range toolMgr.Categories {
		if strings.Contains(result, category.DisplayName) {
			showCategoryTools(category)
			return
		}
	}
}

func showCategoryTools(category toolmanager.ToolCategory) {
	fmt.Println()
	output.Header(fmt.Sprintf("%s Tools", category.DisplayName))
	fmt.Println()

	for _, tool := range category.Tools {
		status := color.GreenString("✓ Available")
		if !tool.Available {
			status = color.YellowString("✗ Not built")
		}

		fmt.Printf("  %s %s\n", color.CyanString(tool.DisplayName), status)
		fmt.Printf("  %s\n", color.HiBlackString(tool.Description))
		fmt.Println()
	}

	promptContinue()
}

func promptContinue() {
	prompt := promptui.Prompt{
		Label:     "Press Enter to continue",
		IsConfirm: false,
	}
	prompt.Run()
}

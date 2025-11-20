// Package cmd implements all CLI commands using the Cobra framework.
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// exampleCmd demonstrates a sample command structure
	exampleCmd = &cobra.Command{
		Use:   "example [args]",
		Short: "Example command description",
		Long: `Detailed description of what this command does.

This command is an example of how to structure subcommands in your tool.
Replace this with your actual command implementation.`,
		Example: `  # Example usage 1
  {{TOOL_NAME}} example arg1

  # Example usage 2
  {{TOOL_NAME}} example arg1 --flag value`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runExample(args)
		},
	}

	// Command-specific flags
	exampleFlag string
)

func init() {
	// Register this command with the root command
	rootCmd.AddCommand(exampleCmd)

	// Add command-specific flags
	exampleCmd.Flags().StringVarP(&exampleFlag, "flag", "f", "", "Example flag description")
}

func runExample(args []string) {
	printInfo("Running example command...")

	// TODO: Implement your command logic here

	printSuccess("Example command completed successfully")
}

package cli

import (
	"fmt"
	"os"

	// Import toolkit commands
	adbcmd "github.com/MKS-01/mobile-recon/internal/adb/cmd"
	apkcmd "github.com/MKS-01/mobile-recon/internal/apk/cmd"
	ioscmd "github.com/MKS-01/mobile-recon/internal/ios/cmd"
	nmapcmd "github.com/MKS-01/mobile-recon/internal/nmap/cmd"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// toolGroups defines the categories used to organize toolkit commands in both
// `--help` and `mobile-recon list`. Order here is the display order.
var toolGroups = []*cobra.Group{
	{ID: "mobile", Title: "Mobile Tools"},
	{ID: "network", Title: "Network Tools"},
}

// Global output flags, applied in PersistentPreRun.
var (
	flagJSON    bool
	flagNoColor bool
	flagQuiet   bool
)

var (
	rootCmd = &cobra.Command{
		Use:   "mobile-recon",
		Short: "Unified CLI for mobile security testing and reconnaissance",
		Long: `A unified CLI for mobile security testing and network reconnaissance.

Tools are grouped by category below. Run 'mobile-recon list' for descriptions,
or 'mobile-recon <tool> --help' for a tool's subcommands.`,
		// PersistentPreRun runs for every command. EnableTraverseRunHooks (set
		// in init) ensures this fires in addition to a toolkit's own pre-run.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f := output.FormatText
			if flagJSON {
				f = output.FormatJSON
			}
			output.Configure(f, flagNoColor, flagQuiet)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if !flagJSON {
				printBanner()
			}
			cmd.Help()
		},
	}
)

func init() {
	// Run parent PersistentPreRun hooks in addition to the matched command's,
	// so global output configuration applies even for toolkit subcommands that
	// define their own PersistentPreRun (e.g. adb's ADB-installed check).
	cobra.EnableTraverseRunHooks = true

	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output JSON instead of text (where supported)")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&flagQuiet, "quiet", false, "Suppress informational/decorative output")

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

// printBanner prints a single styled header line. Kept deliberately minimal
// (no box) so it never misaligns across terminals — the grouped command list
// from cobra does the rest.
func printBanner() {
	cyan := color.New(color.FgCyan, color.Bold)
	dim := color.New(color.FgHiBlack)

	fmt.Println()
	fmt.Print("📱 ")
	cyan.Print("mobile-recon")
	dim.Println("  ·  mobile security & reconnaissance toolkit")
	fmt.Println()
}

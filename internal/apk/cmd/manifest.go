// Package cmd/manifest implements the manifest command for AndroidManifest.xml analysis.
package cmd

import (
	"fmt"
	"os"

	"github.com/MKS-01/mobile-recon/internal/apk"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var (
	extractManifest bool
	manifestOutput  string
)

var manifestCmd = &cobra.Command{
	Use:   "manifest <apk_path>",
	Short: "Analyze or extract AndroidManifest.xml",
	Long: `Analyze the AndroidManifest.xml from an APK file.

Note: APK files contain binary XML that requires decompilation for full analysis.
This command extracts what information is available from the binary format.

For full manifest analysis, use apktool or aapt2:
  apktool d app.apk -o output/
  aapt2 dump xmltree app.apk --file AndroidManifest.xml

Examples:
  apk-analyzer manifest app.apk
  apk-analyzer manifest app.apk --extract -o manifest.xml`,
	Args: cobra.ExactArgs(1),
	Run:  runManifest,
}

func init() {
	RootCmd.AddCommand(manifestCmd)
	manifestCmd.Flags().BoolVarP(&extractManifest, "extract", "e", false, "Extract raw manifest to file")
	manifestCmd.Flags().StringVar(&manifestOutput, "file", "AndroidManifest.xml", "Output file for extraction")
}

func runManifest(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	// Check if manifest exists
	hasManifest, err := apk.HasFile(apkPath, "AndroidManifest.xml")
	if err != nil {
		output.Error("Failed to check manifest: %v", err)
		os.Exit(1)
	}

	if !hasManifest {
		output.Error("AndroidManifest.xml not found in APK")
		os.Exit(1)
	}

	if extractManifest {
		// Extract raw manifest
		err := apk.ExtractFile(apkPath, "AndroidManifest.xml", manifestOutput)
		if err != nil {
			output.Error("Failed to extract manifest: %v", err)
			os.Exit(1)
		}
		output.Success("Manifest extracted to %s", manifestOutput)
		output.Info("Note: The extracted file is in binary XML format.")
		output.Info("Use 'apktool d %s' or 'aapt2 dump xmltree' to decompile it.", apkPath)
		return
	}

	// Analyze manifest
	output.Section("AndroidManifest.xml Analysis")

	manifestData, err := apk.ReadManifestRaw(apkPath)
	if err != nil {
		output.Error("Failed to read manifest: %v", err)
		os.Exit(1)
	}

	output.KeyValue("Manifest Size", apk.FormatSize(int64(len(manifestData))))

	// Extract what we can from binary format
	info, err := apk.ParseManifestBasic(apkPath)
	if err != nil {
		output.Warning("Could not parse manifest details: %v", err)
	} else {
		if info.PackageName != "" {
			output.KeyValue("Package Name", info.PackageName)
		}

		if len(info.Permissions) > 0 {
			fmt.Println()
			output.BoldColor().Printf("Permissions Found: %d\n", len(info.Permissions))
			for _, perm := range info.Permissions {
				fmt.Printf("  • %s\n", perm)
			}
		}
	}

	fmt.Println()
	output.Info("For complete manifest analysis, use:")
	fmt.Println("  apktool d", apkPath, "-o output/")
	fmt.Println("  aapt2 dump xmltree", apkPath, "--file AndroidManifest.xml")
}

// Package cmd/info implements the info command for displaying APK metadata.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/apk"
	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <apk_path>",
	Short: "Display APK metadata and basic information",
	Long: `Extract and display basic information from an APK file.

Shows:
  • File information (name, size, path)
  • Package name and version
  • SDK requirements
  • Native library architectures
  • Number of DEX files

Examples:
  apk-analyzer info app.apk
  apk-analyzer info app.apk --output json`,
	Args: cobra.ExactArgs(1),
	Run:  runInfo,
}

func init() {
	RootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		utils.PrintError("%v", err)
		os.Exit(1)
	}

	info, err := apk.GetAPKInfo(apkPath)
	if err != nil {
		utils.PrintError("Failed to analyze APK: %v", err)
		os.Exit(1)
	}

	// Get additional info
	dexFiles, _ := apk.GetDexFiles(apkPath)
	nativeLibs, _ := apk.GetNativeLibraries(apkPath)

	if outputFormat == "json" {
		outputJSON(info, dexFiles, nativeLibs)
		return
	}

	// Text output
	utils.PrintSection("APK Information")

	utils.PrintKeyValue("File Name", filepath.Base(apkPath))
	utils.PrintKeyValue("File Path", apkPath)
	utils.PrintKeyValue("File Size", apk.FormatSize(info.FileSize))

	fmt.Println()

	// DEX files
	if len(dexFiles) > 0 {
		utils.PrintKeyValue("DEX Files", fmt.Sprintf("%d", len(dexFiles)))
		if verbose {
			for _, dex := range dexFiles {
				fmt.Printf("  • %s\n", dex)
			}
		}
	}

	// Native libraries
	if info.HasNativeLib {
		utils.PrintKeyValue("Native Code", "Yes")
		utils.PrintKeyValue("Architectures", strings.Join(info.Architectures, ", "))

		if verbose && len(nativeLibs) > 0 {
			fmt.Println()
			utils.Bold.Println("Native Libraries:")
			for arch, libs := range nativeLibs {
				fmt.Printf("  %s:\n", arch)
				for _, lib := range libs {
					fmt.Printf("    • %s\n", lib)
				}
			}
		}
	} else {
		utils.PrintKeyValue("Native Code", "No")
	}

	utils.PrintSuccess("Analysis complete")
}

func outputJSON(info *apk.APKInfo, dexFiles []string, nativeLibs map[string][]string) {
	output := map[string]interface{}{
		"file_path":     info.FilePath,
		"file_size":     info.FileSize,
		"has_native":    info.HasNativeLib,
		"architectures": info.Architectures,
		"dex_files":     dexFiles,
		"native_libs":   nativeLibs,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		utils.PrintError("Failed to generate JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

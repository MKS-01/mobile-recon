// Package cmd/permissions implements the permissions command for analyzing app permissions.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MKS-01/mobile-recon/internal/apk"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var (
	dangerousOnly bool
)

var permissionsCmd = &cobra.Command{
	Use:   "permissions <apk_path>",
	Short: "Analyze APK permissions",
	Long: `Extract and analyze permissions declared in the APK.

Identifies:
  • All requested permissions
  • Dangerous permissions requiring runtime approval
  • Permission risk categories

Examples:
  mobile-recon apk permissions app.apk
  mobile-recon apk permissions app.apk --dangerous
  mobile-recon apk permissions app.apk -o json`,
	Args: cobra.ExactArgs(1),
	Run:  runPermissions,
}

func init() {
	RootCmd.AddCommand(permissionsCmd)
	permissionsCmd.Flags().BoolVarP(&dangerousOnly, "dangerous", "d", false, "Show only dangerous permissions")
}

func runPermissions(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	info, err := apk.ParseManifestBasic(apkPath)
	if err != nil {
		output.Error("Failed to parse manifest: %v", err)
		os.Exit(1)
	}

	// Categorize permissions
	var dangerous []string
	var normal []string

	for _, perm := range info.Permissions {
		if _, isDangerous := apk.DangerousPermissions[perm]; isDangerous {
			dangerous = append(dangerous, perm)
		} else {
			normal = append(normal, perm)
		}
	}

	if outputFormat == "json" {
		outputPermissionsJSON(info.Permissions, dangerous, normal)
		return
	}

	// Text output
	output.Section("Permission Analysis")

	if len(info.Permissions) == 0 {
		output.Info("No permissions detected from binary manifest.")
		output.Info("Use 'apktool d %s' for complete permission list.", apkPath)
		return
	}

	output.KeyValue("Total Permissions", fmt.Sprintf("%d", len(info.Permissions)))

	// Dangerous permissions
	if len(dangerous) > 0 {
		fmt.Println()
		output.WarningColor().Printf("Dangerous Permissions: %d\n", len(dangerous))
		fmt.Println()

		for _, perm := range dangerous {
			permInfo := apk.DangerousPermissions[perm]
			shortName := strings.TrimPrefix(perm, "android.permission.")
			output.ErrorColor().Printf("  • %s\n", shortName)
			if verbose {
				fmt.Printf("      Description: %s\n", permInfo.Description)
				fmt.Printf("      Risk: %s\n", permInfo.Risk)
			}
		}
	}

	// Normal permissions
	if !dangerousOnly && len(normal) > 0 {
		fmt.Println()
		output.InfoColor().Printf("Other Permissions: %d\n", len(normal))
		fmt.Println()

		for _, perm := range normal {
			shortName := strings.TrimPrefix(perm, "android.permission.")
			fmt.Printf("  • %s\n", shortName)
		}
	}

	// Summary
	fmt.Println()
	if len(dangerous) > 0 {
		output.Warning("App requests %d dangerous permission(s)", len(dangerous))
	} else {
		output.Success("No dangerous permissions detected")
	}
}

func outputPermissionsJSON(all, dangerous, normal []string) {
	payload := map[string]interface{}{
		"total":     len(all),
		"dangerous": dangerous,
		"normal":    normal,
		"all":       all,
		"dangerous_details": func() []map[string]string {
			var details []map[string]string
			for _, perm := range dangerous {
				if info, ok := apk.DangerousPermissions[perm]; ok {
					details = append(details, map[string]string{
						"permission":  perm,
						"description": info.Description,
						"risk":        info.Risk,
					})
				}
			}
			return details
		}(),
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		output.Error("Failed to generate JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

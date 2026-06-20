// Package cmd/abuse_permissions implements the abuse-permissions command for detecting
// Android permissions commonly abused by malware.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MKS-01/mobile-recon/internal/apk"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	malwareOnly bool
)

var abusePermissionsCmd = &cobra.Command{
	Use:   "abuse-permissions <apk_path>",
	Short: "Detect abusive Android permissions",
	Long: `Analyze APK for permissions commonly abused by malware.

Detects two categories of abusive permissions:

Malware Permissions:
  Top permissions widely abused by known malware including spyware,
  banking trojans, ransomware, and adware. These are high-priority
  indicators of potentially malicious behavior.

Other Common Permissions:
  Permissions commonly abused by known malware but also frequently
  found in legitimate applications. Requires contextual analysis.

Each permission includes:
  • Status: normal, dangerous, or unknown
  • Info: Brief description of capability
  • Description: Detailed explanation of abuse potential

Examples:
  mobile-recon apk abuse-permissions app.apk
  mobile-recon apk abuse-permissions app.apk --malware
  mobile-recon apk abuse-permissions app.apk -o json
  mobile-recon apk abuse-permissions app.apk -v`,
	Args: cobra.ExactArgs(1),
	Run:  runAbusePermissions,
}

func init() {
	RootCmd.AddCommand(abusePermissionsCmd)
	abusePermissionsCmd.Flags().BoolVarP(&malwareOnly, "malware", "m", false, "Show only malware permissions (high-priority indicators)")
}

func runAbusePermissions(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	output.Info("Scanning for abusive permissions...")

	result, err := apk.AnalyzeAbusivePermissions(apkPath)
	if err != nil {
		output.Error("Failed to analyze permissions: %v", err)
		os.Exit(1)
	}

	if output.IsJSON() {
		outputAbusePermissionsJSON(result)
		return
	}

	// Text output
	output.Section("Abused Permissions Analysis")

	if result.TotalPermissions == 0 {
		output.Info("No permissions detected from binary manifest.")
		output.Info("Use 'apktool d %s' for complete permission list.", apkPath)
		return
	}

	output.KeyValue("Total Permissions", fmt.Sprintf("%d", result.TotalPermissions))
	fmt.Println()

	// Malware Permissions
	if len(result.MalwareMatches) > 0 {
		output.ErrorColor().Printf("Malware Permissions: %d/%d\n", len(result.MalwareMatches), len(apk.MalwarePermissions))
		output.BoldColor().Println("Top permissions that are widely abused by known malware")
		fmt.Println()

		printPermissionTable(result.MalwareMatches)
	} else {
		output.Success("No malware-associated permissions detected")
		fmt.Println()
	}

	// Other Common Abused Permissions
	if !malwareOnly && len(result.OtherMatches) > 0 {
		fmt.Println()
		output.WarningColor().Printf("Other Common Permissions: %d/%d\n", len(result.OtherMatches), len(apk.OtherCommonAbusedPermissions))
		output.BoldColor().Println("Permissions that are commonly abused by known malware")
		fmt.Println()

		printPermissionTable(result.OtherMatches)
	}

	// Summary
	fmt.Println()
	printAbuseSummary(result)
}

func printPermissionTable(permissions []apk.AbusivePermission) {
	// Print header
	output.BoldColor().Printf("  %-55s %-10s %-30s\n", "PERMISSION", "STATUS", "INFO")
	fmt.Println(strings.Repeat("-", 100))

	for _, perm := range permissions {
		shortName := formatPermissionName(perm.Permission)
		statusColor := getStatusColor(perm.Status)

		fmt.Printf("  %-55s ", shortName)
		statusColor.Printf("%-10s ", perm.Status)
		fmt.Printf("%-30s\n", truncateString(perm.Info, 30))

		if verbose {
			fmt.Printf("      %s\n", wrapText(perm.Description, 90))
			fmt.Println()
		}
	}
}

func formatPermissionName(perm string) string {
	// Shorten common prefixes for display
	if strings.HasPrefix(perm, "android.permission.") {
		return strings.TrimPrefix(perm, "android.permission.")
	}
	if strings.HasPrefix(perm, "com.google.android.") {
		return strings.Replace(perm, "com.google.android.", "google.", 1)
	}
	return perm
}

func getStatusColor(status string) *color.Color {
	switch status {
	case "dangerous":
		return output.ErrorColor()
	case "normal":
		return output.InfoColor()
	default:
		return output.WarningColor()
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wrapText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		if lineLen+len(word)+1 > maxLen && lineLen > 0 {
			result.WriteString("\n      ")
			lineLen = 0
		}
		if i > 0 && lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}
		result.WriteString(word)
		lineLen += len(word)
	}

	return result.String()
}

func printAbuseSummary(result *apk.AbusivePermissionResult) {
	malwareCount := len(result.MalwareMatches)
	otherCount := len(result.OtherMatches)

	// Calculate risk level
	var riskLevel string
	var riskColor *color.Color

	dangerousCount := 0
	for _, perm := range result.MalwareMatches {
		if perm.Status == "dangerous" {
			dangerousCount++
		}
	}

	if malwareCount >= 10 || dangerousCount >= 5 {
		riskLevel = "HIGH"
		riskColor = output.ErrorColor()
	} else if malwareCount >= 5 || dangerousCount >= 2 {
		riskLevel = "MEDIUM"
		riskColor = output.WarningColor()
	} else if malwareCount > 0 {
		riskLevel = "LOW"
		riskColor = output.InfoColor()
	} else {
		riskLevel = "MINIMAL"
		riskColor = output.SuccessColor()
	}

	output.Section("Summary")

	fmt.Printf("  Malware Permissions:        %d\n", malwareCount)
	fmt.Printf("  Other Abused Permissions:   %d\n", otherCount)
	fmt.Printf("  Dangerous Permissions:      %d\n", dangerousCount)
	fmt.Print("  Risk Assessment:            ")
	riskColor.Printf("%s\n", riskLevel)
	fmt.Println()

	// Recommendations
	if malwareCount > 0 {
		output.WarningColor().Println("Recommendations:")
		if dangerousCount > 0 {
			fmt.Println("  • Review dangerous permissions for legitimate use case")
		}
		if containsPermission(result.MalwareMatches, "SYSTEM_ALERT_WINDOW") {
			fmt.Println("  • Overlay permission detected - check for clickjacking potential")
		}
		if containsPermission(result.MalwareMatches, "RECEIVE_BOOT_COMPLETED") {
			fmt.Println("  • Boot persistence detected - verify startup behavior")
		}
		if containsPermission(result.MalwareMatches, "BIND_ACCESSIBILITY_SERVICE") ||
			containsPermissionOther(result.OtherMatches, "BIND_ACCESSIBILITY_SERVICE") {
			fmt.Println("  • Accessibility service binding - high risk for credential theft")
		}
		if containsPermission(result.MalwareMatches, "READ_SMS") ||
			containsPermission(result.MalwareMatches, "RECEIVE_SMS") {
			fmt.Println("  • SMS permissions detected - check for OTP interception capability")
		}
		fmt.Println("  • Perform dynamic analysis to observe actual behavior")
		fmt.Println("  • Cross-reference with VirusTotal or similar services")
	}
}

func containsPermission(perms []apk.AbusivePermission, shortName string) bool {
	for _, p := range perms {
		if strings.Contains(p.Permission, shortName) {
			return true
		}
	}
	return false
}

func containsPermissionOther(perms []apk.AbusivePermission, shortName string) bool {
	return containsPermission(perms, shortName)
}

func outputAbusePermissionsJSON(result *apk.AbusivePermissionResult) {
	// Calculate statistics
	dangerousCount := 0
	normalCount := 0
	unknownCount := 0

	for _, perm := range result.MalwareMatches {
		switch perm.Status {
		case "dangerous":
			dangerousCount++
		case "normal":
			normalCount++
		default:
			unknownCount++
		}
	}
	for _, perm := range result.OtherMatches {
		switch perm.Status {
		case "dangerous":
			dangerousCount++
		case "normal":
			normalCount++
		default:
			unknownCount++
		}
	}

	// Determine risk level
	var riskLevel string
	if len(result.MalwareMatches) >= 10 || dangerousCount >= 5 {
		riskLevel = "HIGH"
	} else if len(result.MalwareMatches) >= 5 || dangerousCount >= 2 {
		riskLevel = "MEDIUM"
	} else if len(result.MalwareMatches) > 0 {
		riskLevel = "LOW"
	} else {
		riskLevel = "MINIMAL"
	}

	payload := map[string]interface{}{
		"total_permissions": result.TotalPermissions,
		"malware_permissions": map[string]interface{}{
			"count":       len(result.MalwareMatches),
			"total_known": len(apk.MalwarePermissions),
			"matches":     result.MalwareMatches,
		},
		"other_common_permissions": map[string]interface{}{
			"count":       len(result.OtherMatches),
			"total_known": len(apk.OtherCommonAbusedPermissions),
			"matches":     result.OtherMatches,
		},
		"statistics": map[string]int{
			"dangerous": dangerousCount,
			"normal":    normalCount,
			"unknown":   unknownCount,
		},
		"risk_assessment": riskLevel,
		"all_permissions": result.AllPermissions,
	}

	if err := output.JSON(payload); err != nil {
		output.Error("Failed to generate JSON: %v", err)
	}
}

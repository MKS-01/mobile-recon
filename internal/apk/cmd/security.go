// Package cmd/security implements the security command for APK security analysis.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/MKS-01/mobile-recon/internal/apk"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security <apk_path>",
	Short: "Perform security analysis on APK",
	Long: `Perform static security analysis on an APK file.

Checks for:
  • Hardcoded secrets and API keys
  • Insecure configurations (debuggable, backup)
  • HTTP URLs (insecure communication)
  • Root/emulator detection bypasses
  • Third-party SDK fingerprints
  • Anti-tampering mechanisms

Examples:
  mobile-recon apk security app.apk
  mobile-recon apk security app.apk -v
  mobile-recon apk security app.apk -o json`,
	Args: cobra.ExactArgs(1),
	Run:  runSecurity,
}

func init() {
	RootCmd.AddCommand(securityCmd)
}

func runSecurity(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}

	output.Info("Analyzing APK security...")

	issues, err := apk.AnalyzeSecurity(apkPath)
	if err != nil {
		output.Error("Security analysis failed: %v", err)
		os.Exit(1)
	}

	// Also check permissions
	info, _ := apk.ParseManifestBasic(apkPath)
	var dangerousPerms []string
	if info != nil {
		for _, perm := range info.Permissions {
			if _, isDangerous := apk.DangerousPermissions[perm]; isDangerous {
				dangerousPerms = append(dangerousPerms, perm)
			}
		}
	}

	if outputFormat == "json" {
		outputSecurityJSON(issues, dangerousPerms)
		return
	}

	// Text output
	output.Section("Security Analysis Results")

	if len(issues) == 0 && len(dangerousPerms) == 0 {
		output.Success("No security issues detected")
		output.Info("Note: This is a static analysis and may not catch all vulnerabilities.")
		return
	}

	// Group issues by severity
	bySeverity := map[string][]apk.SecurityIssue{
		"HIGH":   {},
		"MEDIUM": {},
		"LOW":    {},
		"INFO":   {},
	}

	for _, issue := range issues {
		bySeverity[issue.Severity] = append(bySeverity[issue.Severity], issue)
	}

	// Print summary
	output.KeyValue("Total Issues", fmt.Sprintf("%d", len(issues)))
	if len(dangerousPerms) > 0 {
		output.KeyValue("Dangerous Permissions", fmt.Sprintf("%d", len(dangerousPerms)))
	}
	fmt.Println()

	// Print by severity
	severities := []string{"HIGH", "MEDIUM", "LOW", "INFO"}
	for _, sev := range severities {
		sevIssues := bySeverity[sev]
		if len(sevIssues) == 0 {
			continue
		}

		// Sort by category
		sort.Slice(sevIssues, func(i, j int) bool {
			return sevIssues[i].Category < sevIssues[j].Category
		})

		switch sev {
		case "HIGH":
			output.ErrorColor().Printf("═══ HIGH SEVERITY (%d) ═══\n", len(sevIssues))
		case "MEDIUM":
			output.WarningColor().Printf("═══ MEDIUM SEVERITY (%d) ═══\n", len(sevIssues))
		case "LOW":
			output.InfoColor().Printf("═══ LOW SEVERITY (%d) ═══\n", len(sevIssues))
		case "INFO":
			output.BoldColor().Printf("═══ INFORMATIONAL (%d) ═══\n", len(sevIssues))
		}
		fmt.Println()

		for _, issue := range sevIssues {
			fmt.Printf("  [%s] %s\n", issue.Category, issue.Description)
			if verbose && issue.Details != "" {
				fmt.Printf("      %s\n", issue.Details)
			}
		}
		fmt.Println()
	}

	// Print dangerous permissions
	if len(dangerousPerms) > 0 {
		output.WarningColor().Println("═══ DANGEROUS PERMISSIONS ═══")
		fmt.Println()
		for _, perm := range dangerousPerms {
			permInfo := apk.DangerousPermissions[perm]
			output.WarningColor().Printf("  • %s\n", permInfo.Name)
			if verbose {
				fmt.Printf("      %s (%s)\n", permInfo.Description, permInfo.Risk)
			}
		}
		fmt.Println()
	}

	// Summary
	highCount := len(bySeverity["HIGH"])
	medCount := len(bySeverity["MEDIUM"])

	if highCount > 0 {
		output.Error("Found %d high severity issue(s) - review recommended", highCount)
	} else if medCount > 0 {
		output.Warning("Found %d medium severity issue(s)", medCount)
	} else {
		output.Success("No critical security issues found")
	}

	output.Info("Note: Use 'mobile-recon apk strings' for deeper secrets analysis")
}

func outputSecurityJSON(issues []apk.SecurityIssue, dangerousPerms []string) {
	payload := map[string]interface{}{
		"total_issues":         len(issues),
		"issues":               issues,
		"dangerous_permissions": dangerousPerms,
		"summary": map[string]int{
			"high":   countBySeverity(issues, "HIGH"),
			"medium": countBySeverity(issues, "MEDIUM"),
			"low":    countBySeverity(issues, "LOW"),
			"info":   countBySeverity(issues, "INFO"),
		},
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		output.Error("Failed to generate JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

func countBySeverity(issues []apk.SecurityIssue, severity string) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

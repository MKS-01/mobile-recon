// Package cmd/strings implements the strings command for extracting readable strings.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/apk"
	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	minStringLength int
	searchPattern   string
	targetFile      string
	showURLs        bool
	showEmails      bool
	showIPs         bool
)

var stringsCmd = &cobra.Command{
	Use:   "strings <apk_path>",
	Short: "Extract strings from APK files",
	Long: `Extract readable strings from APK files including DEX and native libraries.

Useful for:
  • Finding hardcoded URLs and API endpoints
  • Discovering API keys and secrets
  • Identifying interesting functionality
  • Malware analysis

Examples:
  apk-analyzer strings app.apk
  apk-analyzer strings app.apk --min 10
  apk-analyzer strings app.apk --search "api|key|secret"
  apk-analyzer strings app.apk --urls
  apk-analyzer strings app.apk --file classes.dex`,
	Args: cobra.ExactArgs(1),
	Run:  runStrings,
}

func init() {
	rootCmd.AddCommand(stringsCmd)
	stringsCmd.Flags().IntVarP(&minStringLength, "min", "m", 6, "Minimum string length")
	stringsCmd.Flags().StringVarP(&searchPattern, "search", "s", "", "Regex pattern to search for")
	stringsCmd.Flags().StringVarP(&targetFile, "file", "f", "", "Specific file to extract from")
	stringsCmd.Flags().BoolVar(&showURLs, "urls", false, "Show only URLs")
	stringsCmd.Flags().BoolVar(&showEmails, "emails", false, "Show only email addresses")
	stringsCmd.Flags().BoolVar(&showIPs, "ips", false, "Show only IP addresses")
}

func runStrings(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		utils.PrintError("%v", err)
		os.Exit(1)
	}

	// Determine search pattern
	pattern := searchPattern
	if showURLs {
		pattern = `https?://[^\s"'<>]+`
	} else if showEmails {
		pattern = `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	} else if showIPs {
		pattern = `\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`
	}

	var results map[string][]string
	var err error

	if pattern != "" {
		// Search with pattern
		utils.PrintInfo("Searching for pattern: %s", pattern)
		results, err = apk.SearchStrings(apkPath, pattern)
	} else if targetFile != "" {
		// Extract from specific file
		strings, err := apk.ExtractStrings(apkPath, targetFile, minStringLength)
		if err != nil {
			utils.PrintError("Failed to extract strings: %v", err)
			os.Exit(1)
		}
		results = map[string][]string{targetFile: strings}
	} else {
		// Extract all strings
		results, err = apk.ExtractAllStrings(apkPath, minStringLength)
	}

	if err != nil {
		utils.PrintError("Failed to extract strings: %v", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		outputStringsJSON(results)
		return
	}

	// Text output
	utils.PrintSection("String Extraction Results")

	if len(results) == 0 {
		utils.PrintInfo("No strings found matching criteria")
		return
	}

	// Count total strings
	totalStrings := 0
	for _, strs := range results {
		totalStrings += len(strs)
	}

	utils.PrintKeyValue("Files Analyzed", fmt.Sprintf("%d", len(results)))
	utils.PrintKeyValue("Strings Found", fmt.Sprintf("%d", totalStrings))
	fmt.Println()

	// Sort files by name for consistent output
	files := make([]string, 0, len(results))
	for file := range results {
		files = append(files, file)
	}
	sort.Strings(files)

	for _, file := range files {
		strs := results[file]
		utils.Bold.Printf("─── %s (%d strings) ───\n", file, len(strs))

		// Deduplicate and limit output
		seen := make(map[string]bool)
		count := 0
		maxShow := 50
		if verbose {
			maxShow = 500
		}

		for _, s := range strs {
			if seen[s] {
				continue
			}
			seen[s] = true
			count++

			if count <= maxShow {
				// Truncate long strings
				if len(s) > 100 && !verbose {
					fmt.Printf("  %s...\n", s[:100])
				} else {
					fmt.Printf("  %s\n", s)
				}
			}
		}

		if count > maxShow {
			utils.PrintInfo("  ... and %d more (use -v for all)", count-maxShow)
		}
		fmt.Println()
	}

	utils.PrintSuccess("Extraction complete")
}

func outputStringsJSON(results map[string][]string) {
	// Deduplicate per file
	deduped := make(map[string][]string)
	for file, strs := range results {
		seen := make(map[string]bool)
		var unique []string
		for _, s := range strs {
			if !seen[s] {
				seen[s] = true
				unique = append(unique, s)
			}
		}
		deduped[file] = unique
	}

	data, err := json.MarshalIndent(deduped, "", "  ")
	if err != nil {
		utils.PrintError("Failed to generate JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

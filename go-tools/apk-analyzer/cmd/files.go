// Package cmd/files implements the files command for listing and extracting APK contents.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/apk"
	"github.com/MKS-01/mobile-recon/go-tools/apk-analyzer/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	filterPattern string
	extractPath   string
	listTree      bool
)

var filesCmd = &cobra.Command{
	Use:   "files <apk_path>",
	Short: "List or extract files from APK",
	Long: `List all files contained in an APK or extract specific files.

Examples:
  apk-analyzer files app.apk
  apk-analyzer files app.apk --filter "*.dex"
  apk-analyzer files app.apk --filter "lib/*"
  apk-analyzer files app.apk --tree
  apk-analyzer files app.apk --extract classes.dex -o ./output/`,
	Args: cobra.ExactArgs(1),
	Run:  runFiles,
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.Flags().StringVar(&filterPattern, "filter", "", "Filter files by pattern (glob)")
	filesCmd.Flags().StringVar(&extractPath, "extract", "", "Extract specific file from APK")
	filesCmd.Flags().BoolVar(&listTree, "tree", false, "Display as directory tree")
}

func runFiles(cmd *cobra.Command, args []string) {
	apkPath := args[0]

	if err := validateAPKPath(apkPath); err != nil {
		utils.PrintError("%v", err)
		os.Exit(1)
	}

	// Handle extraction
	if extractPath != "" {
		outputPath := filepath.Base(extractPath)
		if outputFormat != "text" && outputFormat != "json" {
			outputPath = outputFormat
		}

		err := apk.ExtractFile(apkPath, extractPath, outputPath)
		if err != nil {
			utils.PrintError("Failed to extract: %v", err)
			os.Exit(1)
		}
		utils.PrintSuccess("Extracted %s to %s", extractPath, outputPath)
		return
	}

	// List files
	files, err := apk.ListFiles(apkPath)
	if err != nil {
		utils.PrintError("Failed to list files: %v", err)
		os.Exit(1)
	}

	// Apply filter
	if filterPattern != "" {
		var filtered []string
		for _, f := range files {
			matched, _ := filepath.Match(filterPattern, f)
			if matched {
				filtered = append(filtered, f)
			}
			// Also check if pattern matches directory prefix
			if strings.HasSuffix(filterPattern, "/*") {
				prefix := strings.TrimSuffix(filterPattern, "/*")
				if strings.HasPrefix(f, prefix+"/") {
					filtered = append(filtered, f)
				}
			}
		}
		files = filtered
	}

	sort.Strings(files)

	if outputFormat == "json" {
		outputFilesJSON(files)
		return
	}

	// Text output
	utils.PrintSection("APK Contents")
	utils.PrintKeyValue("Total Files", fmt.Sprintf("%d", len(files)))
	fmt.Println()

	if listTree {
		printTree(files)
	} else {
		for _, f := range files {
			fmt.Println(f)
		}
	}

	// Show summary by type
	if !listTree {
		fmt.Println()
		showFileSummary(files)
	}
}

func printTree(files []string) {
	tree := make(map[string][]string)

	for _, f := range files {
		dir := filepath.Dir(f)
		if dir == "." {
			dir = "/"
		}
		tree[dir] = append(tree[dir], filepath.Base(f))
	}

	// Get sorted directories
	dirs := make([]string, 0, len(tree))
	for dir := range tree {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		utils.Bold.Printf("%s/\n", dir)
		for _, file := range tree[dir] {
			fmt.Printf("  %s\n", file)
		}
	}
}

func showFileSummary(files []string) {
	counts := map[string]int{
		"DEX files":      0,
		"Native libs":    0,
		"Resources":      0,
		"Assets":         0,
		"XML files":      0,
		"Other":          0,
	}

	for _, f := range files {
		switch {
		case strings.HasSuffix(f, ".dex"):
			counts["DEX files"]++
		case strings.HasPrefix(f, "lib/") && strings.HasSuffix(f, ".so"):
			counts["Native libs"]++
		case strings.HasPrefix(f, "res/"):
			counts["Resources"]++
		case strings.HasPrefix(f, "assets/"):
			counts["Assets"]++
		case strings.HasSuffix(f, ".xml"):
			counts["XML files"]++
		default:
			counts["Other"]++
		}
	}

	utils.Bold.Println("Summary:")
	for category, count := range counts {
		if count > 0 {
			fmt.Printf("  %-15s %d\n", category+":", count)
		}
	}
}

func outputFilesJSON(files []string) {
	data, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		utils.PrintError("Failed to generate JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

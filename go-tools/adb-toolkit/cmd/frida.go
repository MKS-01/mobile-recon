// Package cmd/frida implements Frida dynamic instrumentation utilities.
// This file provides commands for setting up and managing the Frida framework
// on Android devices, including automated server download, installation, and management.
//
// Commands:
//   - frida setup:          Automated Frida server setup (download, push, start)
//   - frida ps:             List running processes in Frida-friendly format
//   - frida trace:          Generate command for tracing method calls
//   - frida server check:   Verify if Frida server is running on device
//   - frida server start:   Start Frida server on device
//   - frida server stop:    Stop Frida server on device
//   - frida server status:  Detailed Frida server status
//
// Frida is a dynamic instrumentation toolkit for reverse engineering,
// security research, and debugging. It allows runtime manipulation of
// application behavior without modifying the APK.
package cmd

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/adb"
	"github.com/MKS-01/mobile-recon/go-tools/adb-toolkit/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
)

const (
	fridaServerPath     = "/data/local/tmp/frida-server"
	fridaGitHubAPI      = "https://api.github.com/repos/frida/frida/releases/latest"
	fridaDownloadURL    = "https://github.com/frida/frida/releases/download/%s/frida-server-%s-android-%s.xz"
	defaultFridaCacheDir = ".frida-server"
)

// GitHubRelease represents the GitHub API response for latest release.
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// fridaCmd is the parent command for Frida-related operations.
var fridaCmd = &cobra.Command{
	Use:   "frida",
	Short: "Frida dynamic instrumentation utilities",
	Long: `Helper commands for working with Frida framework on Android.

Frida is a dynamic instrumentation toolkit that lets you inject JavaScript
into native apps on Android. Use these commands to set up and manage Frida
server on your device.

Quick Start:
  adb-toolkit frida setup    # Automated setup (recommended)
  adb-toolkit frida ps       # List processes after setup`,
}

// fridaSetupCmd provides automated Frida server setup.
var fridaSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Automated Frida server setup (download, push, start)",
	Long: `Automatically set up Frida server on the connected device.

This command will:
  1. Verify device is connected and ready
  2. Check and enable root access
  3. Detect device architecture
  4. Download the correct Frida server version
  5. Push Frida server to the device
  6. Start Frida server

Prerequisites:
  - Android emulator with google_apis image (NOT google_apis_playstore)
  - Or a rooted physical device`,
	Run: runFridaSetup,
}

var (
	forceDownload bool
	fridaVersion  string
	forceKill     bool
)

func init() {
	rootCmd.AddCommand(fridaCmd)

	// Setup command with flags
	fridaSetupCmd.Flags().BoolVarP(&forceDownload, "force", "f", false, "Force re-download of Frida server")
	fridaSetupCmd.Flags().StringVarP(&fridaVersion, "version", "v", "", "Specific Frida version (default: latest)")
	fridaCmd.AddCommand(fridaSetupCmd)

	// Process commands
	fridaCmd.AddCommand(fridaPsCmd)
	fridaCmd.AddCommand(fridaTraceCmd)

	// Kill command (shortcut)
	fridaCmd.AddCommand(fridaKillCmd)

	// Server management
	fridaCmd.AddCommand(fridaServerCmd)
	fridaServerCmd.AddCommand(fridaServerCheckCmd)
	fridaServerCmd.AddCommand(fridaServerStartCmd)
	fridaServerStopCmd.Flags().BoolVarP(&forceKill, "force", "f", false, "Force kill with SIGKILL immediately")
	fridaServerCmd.AddCommand(fridaServerStopCmd)
	fridaServerCmd.AddCommand(fridaServerStatusCmd)
}

// runFridaSetup performs the automated Frida setup process.
func runFridaSetup(cmd *cobra.Command, args []string) {
	serial, err := getTargetDevice()
	if err != nil {
		utils.PrintError("%v", err)
		return
	}

	utils.PrintSection("Frida Server Setup")

	// Step 1: Check device is ready
	utils.PrintInfo("Checking device status...")
	if !adb.IsDeviceReady(serial) {
		utils.PrintError("Device is not fully booted. Please wait and try again.")
		return
	}
	utils.PrintSuccess("Device is ready")

	// Step 2: Verify root access
	utils.PrintInfo("Checking root access...")
	hasRoot, err := adb.RestartAsRoot(serial)
	if err != nil {
		utils.PrintError("Failed to check root access: %v", err)
		return
	}
	if !hasRoot {
		utils.PrintError("Root access not available!")
		fmt.Println()
		fmt.Println("  This device doesn't support ADB root access.")
		fmt.Println("  Make sure you're using a google_apis system image")
		fmt.Println("  (NOT google_apis_playstore).")
		fmt.Println()
		fmt.Println("  See: docs/android-root-using-frida.md")
		return
	}
	utils.PrintSuccess("Root access enabled")

	// Give ADB time to restart
	time.Sleep(2 * time.Second)

	// Step 3: Detect architecture
	utils.PrintInfo("Detecting device architecture...")
	arch, err := adb.GetDeviceArchitecture(serial)
	if err != nil {
		utils.PrintError("Failed to detect architecture: %v", err)
		return
	}

	fridaArch := mapArchToFrida(arch)
	if fridaArch == "" {
		utils.PrintError("Unsupported architecture: %s", arch)
		return
	}
	utils.PrintSuccess("Architecture: %s (Frida: %s)", arch, fridaArch)

	// Step 4: Get Frida version
	version := fridaVersion
	if version == "" {
		utils.PrintInfo("Fetching latest Frida version...")
		version, err = getLatestFridaVersion()
		if err != nil {
			utils.PrintError("Failed to get Frida version: %v", err)
			return
		}
	}
	utils.PrintSuccess("Frida version: %s", version)

	// Step 5: Download Frida server
	utils.PrintInfo("Downloading Frida server...")
	localPath, err := downloadFridaServer(version, fridaArch)
	if err != nil {
		utils.PrintError("Failed to download Frida server: %v", err)
		return
	}
	utils.PrintSuccess("Downloaded: %s", localPath)

	// Step 6: Check if already running
	pid, _ := adb.GetProcessPID(serial, "frida-server")
	if pid != "" {
		utils.PrintWarning("Frida server already running (PID: %s)", pid)
		utils.PrintInfo("Stopping existing server...")
		adb.ShellCommand(serial, fmt.Sprintf("kill %s", pid))
		time.Sleep(1 * time.Second)
	}

	// Step 7: Push to device
	utils.PrintInfo("Pushing Frida server to device...")
	if err := adb.PushFile(serial, localPath, fridaServerPath); err != nil {
		utils.PrintError("Failed to push Frida server: %v", err)
		return
	}
	utils.PrintSuccess("Pushed to %s", fridaServerPath)

	// Step 8: Set permissions
	utils.PrintInfo("Setting permissions...")
	if _, err := adb.ShellCommand(serial, fmt.Sprintf("chmod 755 %s", fridaServerPath)); err != nil {
		utils.PrintError("Failed to set permissions: %v", err)
		return
	}
	utils.PrintSuccess("Permissions set")

	// Step 9: Start Frida server
	utils.PrintInfo("Starting Frida server...")
	// Start in background - we use a goroutine since the shell command won't return
	go func() {
		adb.ShellCommand(serial, fridaServerPath)
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Step 10: Verify running
	pid, _ = adb.GetProcessPID(serial, "frida-server")
	if pid == "" {
		utils.PrintError("Frida server failed to start")
		fmt.Println()
		fmt.Println("  Try starting manually:")
		fmt.Println("  adb shell '/data/local/tmp/frida-server &'")
		return
	}
	utils.PrintSuccess("Frida server running (PID: %s)", pid)

	// Done!
	utils.PrintSection("Setup Complete!")
	fmt.Println("Frida server is running and ready for connections.")
	fmt.Println()
	fmt.Println("Quick commands:")
	fmt.Println("  frida-ps -U              # List processes")
	fmt.Println("  frida-ps -Ua             # List apps only")
	fmt.Println("  frida -U -f <package>    # Spawn and attach to app")
	fmt.Println()
}

// mapArchToFrida converts Android ABI to Frida architecture name.
func mapArchToFrida(abi string) string {
	switch abi {
	case "arm64-v8a":
		return "arm64"
	case "armeabi-v7a", "armeabi":
		return "arm"
	case "x86_64":
		return "x86_64"
	case "x86":
		return "x86"
	default:
		return ""
	}
}

// getLatestFridaVersion fetches the latest Frida release version from GitHub.
func getLatestFridaVersion() (string, error) {
	resp, err := http.Get(fridaGitHubAPI)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return release.TagName, nil
}

// getFridaCacheDir returns the directory where Frida servers are cached.
func getFridaCacheDir() (string, error) {
	// Try to use the project's .frida-server directory
	// Fall back to user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".mobile-recon", "frida-server")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return cacheDir, nil
}

// downloadFridaServer downloads the Frida server binary for the specified version and architecture.
func downloadFridaServer(version, arch string) (string, error) {
	cacheDir, err := getFridaCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %v", err)
	}

	filename := fmt.Sprintf("frida-server-%s-android-%s", version, arch)
	localPath := filepath.Join(cacheDir, filename)

	// Check if already cached
	if !forceDownload {
		if _, err := os.Stat(localPath); err == nil {
			utils.PrintInfo("Using cached version: %s", localPath)
			return localPath, nil
		}
	}

	// Download
	url := fmt.Sprintf(fridaDownloadURL, version, version, arch)
	utils.PrintInfo("Downloading from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file for download
	tmpFile, err := os.CreateTemp(cacheDir, "frida-download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Download to temp file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to save download: %v", err)
	}
	tmpFile.Close()

	// Decompress XZ
	utils.PrintInfo("Extracting...")
	if err := extractXZ(tmpFile.Name(), localPath); err != nil {
		return "", fmt.Errorf("failed to extract: %v", err)
	}

	return localPath, nil
}

// extractXZ extracts an XZ compressed file.
func extractXZ(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Try XZ decompression
	xzReader, err := xz.NewReader(srcFile)
	if err != nil {
		// Fall back to gzip if xz fails
		srcFile.Seek(0, 0)
		gzReader, gzErr := gzip.NewReader(srcFile)
		if gzErr != nil {
			return fmt.Errorf("failed to create decompressor: xz error: %v, gzip error: %v", err, gzErr)
		}
		defer gzReader.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, gzReader); err != nil {
			return err
		}
		return nil
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, xzReader); err != nil {
		return err
	}

	return nil
}

// fridaPsCmd lists running processes in a format suitable for Frida operations.
var fridaPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running processes (Frida-friendly format)",
	Long: `List running processes on the device.

If frida-ps is installed locally, it will be used for better output.
Otherwise, falls back to the standard Android ps command.`,
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection("Running Processes")

		// Try using local frida-ps first (better output)
		if isFridaToolsInstalled() {
			utils.PrintInfo("Using local frida-ps...")
			runLocalFridaPs()
			return
		}

		// Fall back to adb shell ps
		output, err := adb.ShellCommand(serial, "ps -A")
		if err != nil {
			utils.PrintError("Failed to list processes: %v", err)
			return
		}

		fmt.Println(output)
	},
}

// isFridaToolsInstalled checks if frida-tools is installed on the host.
func isFridaToolsInstalled() bool {
	_, err := os.Stat(getFridaPsPath())
	return err == nil
}

// getFridaPsPath returns the path to frida-ps binary.
func getFridaPsPath() string {
	// Check common locations
	paths := []string{
		"/usr/local/bin/frida-ps",
		"/opt/homebrew/bin/frida-ps",
	}

	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		paths = append(paths,
			filepath.Join(homeDir, ".local/bin/frida-ps"),
			filepath.Join(homeDir, "Library/Python/3.9/bin/frida-ps"),
			filepath.Join(homeDir, "Library/Python/3.10/bin/frida-ps"),
			filepath.Join(homeDir, "Library/Python/3.11/bin/frida-ps"),
			filepath.Join(homeDir, "Library/Python/3.12/bin/frida-ps"),
		)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Try PATH
	if runtime.GOOS == "windows" {
		return "frida-ps.exe"
	}
	return "frida-ps"
}

// runLocalFridaPs runs frida-ps on the host to list device processes.
func runLocalFridaPs() {
	fridaPs := getFridaPsPath()
	output, err := adb.ExecuteCommand(fridaPs, "-U")
	if err != nil {
		utils.PrintWarning("frida-ps failed: %v", err)
		utils.PrintInfo("Make sure Frida tools are installed: pip install frida-tools")
		return
	}
	fmt.Println(output)
}

// fridaTraceCmd provides guidance for using frida-trace.
var fridaTraceCmd = &cobra.Command{
	Use:   "trace <package_name> [method_pattern]",
	Short: "Trace method calls in an app",
	Long: `Generate and display the frida-trace command for tracing method calls.

This command prints the appropriate frida-trace command rather than executing it,
as Frida trace requires an interactive session.

Examples:
  adb-toolkit frida trace com.example.app
  adb-toolkit frida trace com.example.app "open*"
  adb-toolkit frida trace com.example.app "SSL*"`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
		methodPattern := "*"
		if len(args) > 1 {
			methodPattern = args[1]
		}

		utils.PrintSection("Frida Trace")
		utils.PrintInfo("Package: %s", packageName)
		utils.PrintInfo("Pattern: %s", methodPattern)
		fmt.Println()
		fmt.Println("Run this command to start tracing:")
		fmt.Printf("  frida-trace -U -f %s -i '%s'\n", packageName, methodPattern)
		fmt.Println()
		fmt.Println("Common patterns:")
		fmt.Println("  -i 'open*'        # File operations")
		fmt.Println("  -i 'SSL*'         # SSL/TLS functions")
		fmt.Println("  -i 'recv*'        # Network receive")
		fmt.Println("  -i 'send*'        # Network send")
		fmt.Println("  -j '*!*.check*'   # Java methods with 'check' in name")
	},
}

// fridaKillCmd is a shortcut to quickly kill the Frida server.
var fridaKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill Frida server (shortcut for 'frida server stop -f')",
	Long: `Immediately kill the Frida server process using SIGKILL.

This is a shortcut command equivalent to 'frida server stop --force'.
Use this when the server is unresponsive or you need to quickly terminate it.`,
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			utils.PrintInfo("Frida server is not running")
			return
		}

		utils.PrintInfo("Killing Frida server (PID: %s)...", pid)

		// Use SIGKILL directly for immediate termination
		if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid)); err != nil {
			utils.PrintError("Failed to kill Frida server: %v", err)
			return
		}

		time.Sleep(500 * time.Millisecond)

		// Verify killed
		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			utils.PrintError("Failed to kill Frida server - still running (PID: %s)", pid)
			return
		}

		utils.PrintSuccess("Frida server killed")
	},
}

// fridaServerCmd is the parent command for Frida server management.
var fridaServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Frida server management",
	Long: `Manage the Frida server running on the device.

Commands:
  check   - Quick check if server is running
  status  - Detailed server status
  start   - Start the server
  stop    - Stop the server`,
}

// fridaServerCheckCmd verifies if the Frida server is running.
var fridaServerCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if Frida server is running",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			utils.PrintError("Frida server is not running")
			utils.PrintInfo("Run: adb-toolkit frida setup")
			return
		}

		utils.PrintSuccess("Frida server is running (PID: %s)", pid)
	},
}

// fridaServerStatusCmd shows detailed Frida server status.
var fridaServerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Detailed Frida server status",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		utils.PrintSection("Frida Server Status")

		// Check if binary exists
		output, _ := adb.ShellCommand(serial, fmt.Sprintf("ls -la %s 2>/dev/null", fridaServerPath))
		if strings.TrimSpace(output) == "" {
			utils.PrintWarning("Frida server not installed at %s", fridaServerPath)
		} else {
			utils.PrintSuccess("Binary: %s", strings.TrimSpace(output))
		}

		// Check if running
		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			utils.PrintWarning("Server not running")
		} else {
			utils.PrintSuccess("Running with PID: %s", pid)

			// Get process details
			procInfo, _ := adb.ShellCommand(serial, fmt.Sprintf("cat /proc/%s/cmdline 2>/dev/null", pid))
			if procInfo != "" {
				utils.PrintInfo("Command: %s", strings.ReplaceAll(procInfo, "\x00", " "))
			}
		}

		// Check architecture
		arch, _ := adb.GetDeviceArchitecture(serial)
		utils.PrintInfo("Device architecture: %s", arch)

		// Check root status
		hasRoot, _ := adb.RestartAsRoot(serial)
		if hasRoot {
			utils.PrintSuccess("Root access: available")
		} else {
			utils.PrintWarning("Root access: not available")
		}
	},
}

// fridaServerStartCmd starts the Frida server on the device.
var fridaServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Frida server on device",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		// Check if already running
		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			utils.PrintWarning("Frida server already running (PID: %s)", pid)
			return
		}

		// Check if binary exists
		output, _ := adb.ShellCommand(serial, fmt.Sprintf("ls %s 2>/dev/null", fridaServerPath))
		if strings.TrimSpace(output) == "" {
			utils.PrintError("Frida server not found at %s", fridaServerPath)
			utils.PrintInfo("Run: adb-toolkit frida setup")
			return
		}

		utils.PrintInfo("Starting Frida server...")

		// Start in background
		go func() {
			adb.ShellCommand(serial, fridaServerPath)
		}()

		time.Sleep(2 * time.Second)

		// Verify
		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			utils.PrintError("Failed to start Frida server")
			return
		}

		utils.PrintSuccess("Frida server started (PID: %s)", pid)
	},
}

// fridaServerStopCmd stops the Frida server on the device.
var fridaServerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Frida server on device",
	Long: `Stop the Frida server running on the device.

By default, sends SIGTERM for graceful shutdown.
Use --force (-f) to immediately kill with SIGKILL.`,
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			utils.PrintError("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			utils.PrintInfo("Frida server is not running")
			return
		}

		if forceKill {
			utils.PrintInfo("Force killing Frida server (PID: %s)...", pid)
			if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid)); err != nil {
				utils.PrintError("Failed to kill Frida server: %v", err)
				return
			}
		} else {
			utils.PrintInfo("Stopping Frida server (PID: %s)...", pid)
			if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill %s", pid)); err != nil {
				utils.PrintError("Failed to stop Frida server: %v", err)
				return
			}

			time.Sleep(1 * time.Second)

			// Verify stopped, escalate to SIGKILL if needed
			pid, _ = adb.GetProcessPID(serial, "frida-server")
			if pid != "" {
				utils.PrintWarning("Server still running, trying SIGKILL...")
				adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid))
			}
		}

		time.Sleep(500 * time.Millisecond)

		// Final verification
		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			utils.PrintError("Failed to stop Frida server - still running (PID: %s)", pid)
			return
		}

		utils.PrintSuccess("Frida server stopped")
	},
}

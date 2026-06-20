package cmd

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MKS-01/mobile-recon/internal/adb"
	"github.com/MKS-01/mobile-recon/pkg/frida"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
	"github.com/ulikunitz/xz"
)

const (
	fridaServerPath  = "/data/local/tmp/frida-server"
	fridaGitHubAPI   = "https://api.github.com/repos/frida/frida/releases/latest"
	fridaDownloadURL = "https://github.com/frida/frida/releases/download/%s/frida-server-%s-android-%s.xz"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

var fridaCmd = &cobra.Command{
	Use:   "frida",
	Short: "Frida dynamic instrumentation utilities",
	Long: `Helper commands for working with Frida framework on Android.

Quick Start:
  mobile-recon adb frida setup    # Automated setup (recommended)
  mobile-recon adb frida ps       # List processes after setup`,
}

var fridaSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Automated Frida server setup (download, push, start)",
	Long: `Automatically set up Frida server on the connected device.

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
	RootCmd.AddCommand(fridaCmd)

	fridaSetupCmd.Flags().BoolVarP(&forceDownload, "force", "f", false, "Force re-download of Frida server")
	fridaSetupCmd.Flags().StringVarP(&fridaVersion, "version", "v", "", "Specific Frida version (default: latest)")
	fridaCmd.AddCommand(fridaSetupCmd)

	fridaCmd.AddCommand(fridaPsCmd)
	fridaCmd.AddCommand(fridaTraceCmd)
	fridaCmd.AddCommand(fridaKillCmd)

	fridaCmd.AddCommand(fridaServerCmd)
	fridaServerCmd.AddCommand(fridaServerCheckCmd)
	fridaServerCmd.AddCommand(fridaServerStartCmd)
	fridaServerStopCmd.Flags().BoolVarP(&forceKill, "force", "f", false, "Force kill with SIGKILL immediately")
	fridaServerCmd.AddCommand(fridaServerStopCmd)
	fridaServerCmd.AddCommand(fridaServerStatusCmd)
}

func runFridaSetup(cmd *cobra.Command, args []string) {
	serial, err := getTargetDevice()
	if err != nil {
		output.Error("%v", err)
		return
	}

	output.Section("Frida Server Setup")

	output.Info("Checking device status...")
	if !adb.IsDeviceReady(serial) {
		output.Error("Device is not fully booted. Please wait and try again.")
		return
	}
	output.Success("Device is ready")

	output.Info("Checking root access...")
	hasRoot, err := adb.RestartAsRoot(serial)
	if err != nil {
		output.Error("Failed to check root access: %v", err)
		return
	}
	if !hasRoot {
		output.Error("Root access not available!")
		fmt.Println()
		fmt.Println("  This device doesn't support ADB root access.")
		fmt.Println("  Make sure you're using a google_apis system image")
		fmt.Println("  (NOT google_apis_playstore).")
		return
	}
	output.Success("Root access enabled")

	time.Sleep(2 * time.Second)

	output.Info("Detecting device architecture...")
	arch, err := adb.GetDeviceArchitecture(serial)
	if err != nil {
		output.Error("Failed to detect architecture: %v", err)
		return
	}

	fridaArch := mapArchToFrida(arch)
	if fridaArch == "" {
		output.Error("Unsupported architecture: %s", arch)
		return
	}
	output.Success("Architecture: %s (Frida: %s)", arch, fridaArch)

	version := fridaVersion
	if version == "" {
		output.Info("Fetching latest Frida version...")
		version, err = getLatestFridaVersion()
		if err != nil {
			output.Error("Failed to get Frida version: %v", err)
			return
		}
	}
	output.Success("Frida version: %s", version)

	output.Info("Downloading Frida server...")
	localPath, err := downloadFridaServer(version, fridaArch)
	if err != nil {
		output.Error("Failed to download Frida server: %v", err)
		return
	}
	output.Success("Downloaded: %s", localPath)

	pid, _ := adb.GetProcessPID(serial, "frida-server")
	if pid != "" {
		output.Warning("Frida server already running (PID: %s)", pid)
		output.Info("Stopping existing server...")
		adb.ShellCommand(serial, fmt.Sprintf("kill %s", pid))
		time.Sleep(1 * time.Second)
	}

	output.Info("Pushing Frida server to device...")
	if err := adb.PushFile(serial, localPath, fridaServerPath); err != nil {
		output.Error("Failed to push Frida server: %v", err)
		return
	}
	output.Success("Pushed to %s", fridaServerPath)

	output.Info("Setting permissions...")
	if _, err := adb.ShellCommand(serial, fmt.Sprintf("chmod 755 %s", fridaServerPath)); err != nil {
		output.Error("Failed to set permissions: %v", err)
		return
	}
	output.Success("Permissions set")

	output.Info("Starting Frida server...")
	go func() {
		adb.ShellCommand(serial, fridaServerPath)
	}()

	time.Sleep(2 * time.Second)

	pid, _ = adb.GetProcessPID(serial, "frida-server")
	if pid == "" {
		output.Error("Frida server failed to start")
		fmt.Println()
		fmt.Println("  Try starting manually:")
		fmt.Println("  adb shell '/data/local/tmp/frida-server &'")
		return
	}
	output.Success("Frida server running (PID: %s)", pid)

	output.Section("Setup Complete!")
	fmt.Println("Frida server is running and ready for connections.")
	fmt.Println()
	fmt.Println("Quick commands:")
	fmt.Println("  frida-ps -U              # List processes")
	fmt.Println("  frida-ps -Ua             # List apps only")
	fmt.Println("  frida -U -f <package>    # Spawn and attach to app")
	fmt.Println()
}

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

func getFridaCacheDir() (string, error) {
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

func downloadFridaServer(version, arch string) (string, error) {
	cacheDir, err := getFridaCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %v", err)
	}

	filename := fmt.Sprintf("frida-server-%s-android-%s", version, arch)
	localPath := filepath.Join(cacheDir, filename)

	if !forceDownload {
		if _, err := os.Stat(localPath); err == nil {
			output.Info("Using cached version: %s", localPath)
			return localPath, nil
		}
	}

	url := fmt.Sprintf(fridaDownloadURL, version, version, arch)
	output.Info("Downloading from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp(cacheDir, "frida-download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to save download: %v", err)
	}
	tmpFile.Close()

	output.Info("Extracting...")
	if err := extractXZ(tmpFile.Name(), localPath); err != nil {
		return "", fmt.Errorf("failed to extract: %v", err)
	}

	return localPath, nil
}

func extractXZ(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	xzReader, err := xz.NewReader(srcFile)
	if err != nil {
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

var fridaPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running processes (Frida-friendly format)",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		output.Section("Running Processes")

		if frida.Installed() {
			output.Info("Using local frida-ps...")
			runLocalFridaPs()
			return
		}

		cmdOutput, err := adb.ShellCommand(serial, "ps -A")
		if err != nil {
			output.Error("Failed to list processes: %v", err)
			return
		}

		fmt.Println(cmdOutput)
	},
}

func runLocalFridaPs() {
	fridaPs := frida.Locate("frida-ps")
	if fridaPs == "" {
		// Fall back to the bare name so exec can still resolve it via $PATH.
		fridaPs = "frida-ps"
	}
	cmdOutput, err := adb.ExecuteCommand(fridaPs, "-U")
	if err != nil {
		output.Warning("frida-ps failed: %v", err)
		output.Info("Make sure Frida tools are installed: pip install frida-tools")
		return
	}
	fmt.Println(cmdOutput)
}

var fridaTraceCmd = &cobra.Command{
	Use:   "trace <package_name> [method_pattern]",
	Short: "Trace method calls in an app",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
		methodPattern := "*"
		if len(args) > 1 {
			methodPattern = args[1]
		}

		output.Section("Frida Trace")
		output.Info("Package: %s", packageName)
		output.Info("Pattern: %s", methodPattern)
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

var fridaKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill Frida server (shortcut for 'frida server stop -f')",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			output.Info("Frida server is not running")
			return
		}

		output.Info("Killing Frida server (PID: %s)...", pid)

		if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid)); err != nil {
			output.Error("Failed to kill Frida server: %v", err)
			return
		}

		time.Sleep(500 * time.Millisecond)

		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			output.Error("Failed to kill Frida server - still running (PID: %s)", pid)
			return
		}

		output.Success("Frida server killed")
	},
}

var fridaServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Frida server management",
}

var fridaServerCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if Frida server is running",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			output.Error("Frida server is not running")
			output.Info("Run: mobile-recon adb frida setup")
			return
		}

		output.Success("Frida server is running (PID: %s)", pid)
	},
}

var fridaServerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Detailed Frida server status",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		output.Section("Frida Server Status")

		cmdOutput, _ := adb.ShellCommand(serial, fmt.Sprintf("ls -la %s 2>/dev/null", fridaServerPath))
		if strings.TrimSpace(cmdOutput) == "" {
			output.Warning("Frida server not installed at %s", fridaServerPath)
		} else {
			output.Success("Binary: %s", strings.TrimSpace(cmdOutput))
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			output.Warning("Server not running")
		} else {
			output.Success("Running with PID: %s", pid)

			procInfo, _ := adb.ShellCommand(serial, fmt.Sprintf("cat /proc/%s/cmdline 2>/dev/null", pid))
			if procInfo != "" {
				output.Info("Command: %s", strings.ReplaceAll(procInfo, "\x00", " "))
			}
		}

		arch, _ := adb.GetDeviceArchitecture(serial)
		output.Info("Device architecture: %s", arch)

		hasRoot, _ := adb.RestartAsRoot(serial)
		if hasRoot {
			output.Success("Root access: available")
		} else {
			output.Warning("Root access: not available")
		}
	},
}

var fridaServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Frida server on device",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			output.Warning("Frida server already running (PID: %s)", pid)
			return
		}

		cmdOutput, _ := adb.ShellCommand(serial, fmt.Sprintf("ls %s 2>/dev/null", fridaServerPath))
		if strings.TrimSpace(cmdOutput) == "" {
			output.Error("Frida server not found at %s", fridaServerPath)
			output.Info("Run: mobile-recon adb frida setup")
			return
		}

		output.Info("Starting Frida server...")

		go func() {
			adb.ShellCommand(serial, fridaServerPath)
		}()

		time.Sleep(2 * time.Second)

		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			output.Error("Failed to start Frida server")
			return
		}

		output.Success("Frida server started (PID: %s)", pid)
	},
}

var fridaServerStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Frida server on device",
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := getTargetDevice()
		if err != nil {
			output.Error("%v", err)
			return
		}

		pid, _ := adb.GetProcessPID(serial, "frida-server")
		if pid == "" {
			output.Info("Frida server is not running")
			return
		}

		if forceKill {
			output.Info("Force killing Frida server (PID: %s)...", pid)
			if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid)); err != nil {
				output.Error("Failed to kill Frida server: %v", err)
				return
			}
		} else {
			output.Info("Stopping Frida server (PID: %s)...", pid)
			if _, err := adb.ShellCommand(serial, fmt.Sprintf("kill %s", pid)); err != nil {
				output.Error("Failed to stop Frida server: %v", err)
				return
			}

			time.Sleep(1 * time.Second)

			pid, _ = adb.GetProcessPID(serial, "frida-server")
			if pid != "" {
				output.Warning("Server still running, trying SIGKILL...")
				adb.ShellCommand(serial, fmt.Sprintf("kill -9 %s", pid))
			}
		}

		time.Sleep(500 * time.Millisecond)

		pid, _ = adb.GetProcessPID(serial, "frida-server")
		if pid != "" {
			output.Error("Failed to stop Frida server - still running (PID: %s)", pid)
			return
		}

		output.Success("Frida server stopped")
	},
}

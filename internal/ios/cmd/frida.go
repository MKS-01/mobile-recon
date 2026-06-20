package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MKS-01/mobile-recon/internal/ios"
	"github.com/MKS-01/mobile-recon/pkg/frida"
	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/spf13/cobra"
)

var fridaScript string

var fridaCmd = &cobra.Command{
	Use:   "frida",
	Short: "Frida dynamic instrumentation for iOS Simulator",
	Long: `Helper commands for working with Frida framework on iOS Simulator.

IMPORTANT: iOS Simulators do NOT need jailbreaking!
Unlike physical iOS devices, Frida can attach directly to simulator processes
because they run as regular macOS processes.

Prerequisites:
  - Frida tools installed: pip3 install frida-tools
  - A booted iOS Simulator

Quick Start:
  mobile-recon ios frida setup     # Verify Frida installation
  mobile-recon ios frida ps        # List processes
  mobile-recon ios frida attach    # Attach to an app`,
}

var fridaSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Verify Frida installation and simulator readiness",
	Long: `Verify that Frida is properly installed and can connect to the iOS Simulator.

Unlike Android, iOS Simulators don't need frida-server because:
  1. Simulator apps run as macOS processes
  2. Frida can attach directly via the macOS kernel
  3. No root/jailbreak required`,
	Run: runFridaSetup,
}

var fridaPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running apps/processes on the simulator",
	Run:   runFridaPs,
}

var fridaAttachCmd = &cobra.Command{
	Use:   "attach <bundle_id_or_process>",
	Short: "Attach Frida to a running app",
	Long: `Attach Frida to a running app on the iOS Simulator.

Examples:
  mobile-recon ios frida attach com.example.app
  mobile-recon ios frida attach Safari
  mobile-recon ios frida attach -s script.js com.example.app`,
	Args: cobra.ExactArgs(1),
	Run:  runFridaAttach,
}

var fridaSpawnCmd = &cobra.Command{
	Use:   "spawn <bundle_id>",
	Short: "Spawn and attach to an app",
	Long: `Spawn an app on the iOS Simulator and attach Frida immediately.

This is useful for hooking early initialization code.

Examples:
  mobile-recon ios frida spawn com.example.app
  mobile-recon ios frida spawn -s early-hook.js com.example.app`,
	Args: cobra.ExactArgs(1),
	Run:  runFridaSpawn,
}

var fridaTraceCmd = &cobra.Command{
	Use:   "trace <bundle_id_or_process> [pattern]",
	Short: "Trace method calls in an app",
	Long: `Trace Objective-C/Swift methods or C functions in an iOS app.

Examples:
  mobile-recon ios frida trace com.example.app
  mobile-recon ios frida trace Safari "-m '*[NSURL*]*'"
  mobile-recon ios frida trace Safari "-i 'open*'"`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runFridaTrace,
}

var fridaAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List installed apps on the simulator",
	Run:   runFridaApps,
}

var fridaKillCmd = &cobra.Command{
	Use:   "kill <bundle_id_or_process>",
	Short: "Kill an app on the simulator",
	Args:  cobra.ExactArgs(1),
	Run:   runFridaKill,
}

func init() {
	RootCmd.AddCommand(fridaCmd)

	fridaCmd.AddCommand(fridaSetupCmd)
	fridaCmd.AddCommand(fridaPsCmd)
	fridaCmd.AddCommand(fridaAppsCmd)

	fridaAttachCmd.Flags().StringVarP(&fridaScript, "script", "s", "", "Frida script to load")
	fridaCmd.AddCommand(fridaAttachCmd)

	fridaSpawnCmd.Flags().StringVarP(&fridaScript, "script", "s", "", "Frida script to load")
	fridaCmd.AddCommand(fridaSpawnCmd)

	fridaCmd.AddCommand(fridaTraceCmd)
	fridaCmd.AddCommand(fridaKillCmd)
}

func runFridaSetup(cmd *cobra.Command, args []string) {
	output.Section("Frida Setup for iOS Simulator")

	// Step 1: Check if simulator is booted
	output.Info("Checking for booted simulator...")
	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator found")
		fmt.Println("\n  Boot a simulator with:")
		fmt.Println("  xcrun simctl boot <device_udid>")
		fmt.Println("\n  List available simulators:")
		fmt.Println("  xcrun simctl list devices")
		return
	}
	output.Success("Simulator: %s (%s)", sim.Name, sim.UDID)
	output.Success("State: %s", sim.State)

	// Step 2: Check Frida installation
	output.Info("Checking Frida installation...")
	fridaPath := frida.Locate("frida")
	if fridaPath == "" {
		output.Error("Frida not found!")
		fmt.Println("\n  Install Frida with:")
		fmt.Println("  pip3 install frida-tools")
		fmt.Println("\n  Or with Homebrew:")
		fmt.Println("  brew install frida")
		return
	}
	output.Success("Frida found: %s", fridaPath)

	// Step 3: Check Frida version
	output.Info("Checking Frida version...")
	version, err := frida.Version()
	if err != nil {
		output.Warning("Could not determine Frida version: %v", err)
	} else {
		output.Success("Frida version: %s", version)
	}

	// Step 4: Test Frida connectivity
	output.Info("Testing Frida connectivity...")
	if testFridaConnectivity() {
		output.Success("Frida can enumerate processes")
	} else {
		output.Warning("Frida connectivity test failed")
		fmt.Println("\n  This might be normal if no apps are running.")
		fmt.Println("  Try launching an app first.")
	}

	// Step 5: Check architecture
	output.Info("Checking system architecture...")
	arch, _ := ios.GetArchitecture()
	output.Success("Architecture: %s", arch)
	if arch == "arm64" {
		output.Info("Running on Apple Silicon - native performance")
	} else {
		output.Info("Running on Intel - using Rosetta translation")
	}

	output.Section("Setup Complete!")
	fmt.Println("iOS Simulator is ready for Frida instrumentation.")
	fmt.Println()
	fmt.Println("Key difference from Android:")
	fmt.Println("  - NO frida-server needed on the simulator")
	fmt.Println("  - NO jailbreak/root required")
	fmt.Println("  - Frida attaches directly to macOS processes")
	fmt.Println()
	fmt.Println("Quick commands:")
	fmt.Println("  mobile-recon ios frida ps                     # List processes")
	fmt.Println("  mobile-recon ios frida apps                   # List installed apps")
	fmt.Println("  mobile-recon ios frida attach <bundle_id>     # Attach to app")
	fmt.Println("  mobile-recon ios frida spawn <bundle_id>      # Spawn and attach")
	fmt.Println()
	fmt.Println("Direct Frida commands:")
	fmt.Printf("  frida -D %s <process>     # Attach\n", sim.UDID[:8])
	fmt.Printf("  frida -D %s -f <bundle>   # Spawn\n", sim.UDID[:8])
	fmt.Println()
}

func runFridaPs(cmd *cobra.Command, args []string) {
	output.Section("Running Processes")

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	fridaPsPath := frida.Locate("frida-ps")
	if fridaPsPath != "" {
		// Use frida-ps to list processes
		output.Info("Using frida-ps for simulator: %s", sim.Name)
		fmt.Println()

		// Try with device UDID first
		cmdOutput, err := ios.ExecuteCommand(fridaPsPath, "-D", sim.UDID)
		if err != nil {
			// Fall back to listing all local processes (which includes simulator apps)
			cmdOutput, err = ios.ExecuteCommand(fridaPsPath, "-a")
			if err != nil {
				output.Error("frida-ps failed: %v", err)
				return
			}
		}
		fmt.Println(cmdOutput)
		return
	}

	// Fallback: show apps running via simctl
	output.Warning("frida-ps not found, showing basic process info")
	output.Info("Install frida-tools: pip3 install frida-tools")
	fmt.Println()

	// Show what apps might be running
	apps, err := ios.GetInstalledApps(sim.UDID)
	if err != nil {
		output.Warning("Could not list apps: %v", err)
		return
	}

	output.Info("Installed apps (may not all be running):")
	for _, app := range apps {
		fmt.Printf("  %s (%s)\n", app.Name, app.BundleID)
	}
}

func runFridaApps(cmd *cobra.Command, args []string) {
	output.Section("Installed Apps")

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	output.Info("Simulator: %s", sim.Name)
	fmt.Println()

	apps, err := ios.GetInstalledApps(sim.UDID)
	if err != nil {
		output.Error("Failed to list apps: %v", err)
		return
	}

	if len(apps) == 0 {
		output.Info("No third-party apps installed")
		fmt.Println("\n  Install an app with:")
		fmt.Println("  xcrun simctl install booted /path/to/App.app")
		return
	}

	fmt.Printf("%-40s %s\n", "BUNDLE ID", "NAME")
	fmt.Printf("%-40s %s\n", strings.Repeat("-", 40), strings.Repeat("-", 20))
	for _, app := range apps {
		name := app.Name
		if name == "" {
			name = "(unknown)"
		}
		fmt.Printf("%-40s %s\n", app.BundleID, name)
	}
}

func runFridaAttach(cmd *cobra.Command, args []string) {
	target := args[0]

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	output.Section("Frida Attach")
	output.Info("Simulator: %s", sim.Name)
	output.Info("Target: %s", target)

	fridaPath := frida.Locate("frida")
	if fridaPath == "" {
		output.Error("Frida not installed")
		fmt.Println("\n  Install with: pip3 install frida-tools")
		return
	}

	// Build Frida command
	fridaArgs := []string{"-D", sim.UDID, target}
	if fridaScript != "" {
		fridaArgs = append(fridaArgs, "-l", fridaScript)
	}

	fmt.Println()
	fmt.Printf("Running: %s %s\n", fridaPath, strings.Join(fridaArgs, " "))
	fmt.Println()

	// Execute interactively
	execCmd := exec.Command(fridaPath, fridaArgs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		output.Error("Frida exited with error: %v", err)
	}
}

func runFridaSpawn(cmd *cobra.Command, args []string) {
	bundleID := args[0]

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	output.Section("Frida Spawn")
	output.Info("Simulator: %s", sim.Name)
	output.Info("Bundle ID: %s", bundleID)

	fridaPath := frida.Locate("frida")
	if fridaPath == "" {
		output.Error("Frida not installed")
		fmt.Println("\n  Install with: pip3 install frida-tools")
		return
	}

	// Build Frida command with spawn (-f) flag
	fridaArgs := []string{"-D", sim.UDID, "-f", bundleID}
	if fridaScript != "" {
		fridaArgs = append(fridaArgs, "-l", fridaScript)
	}

	fmt.Println()
	fmt.Printf("Running: %s %s\n", fridaPath, strings.Join(fridaArgs, " "))
	fmt.Println()

	// Execute interactively
	execCmd := exec.Command(fridaPath, fridaArgs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		output.Error("Frida exited with error: %v", err)
	}
}

func runFridaTrace(cmd *cobra.Command, args []string) {
	target := args[0]
	pattern := ""
	if len(args) > 1 {
		pattern = args[1]
	}

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	output.Section("Frida Trace")
	output.Info("Simulator: %s", sim.Name)
	output.Info("Target: %s", target)

	fridaTracePath := frida.Locate("frida-trace")
	if fridaTracePath == "" {
		output.Error("frida-trace not found")
		fmt.Println("\n  Install with: pip3 install frida-tools")
		return
	}

	if pattern == "" {
		// Show example patterns
		fmt.Println()
		fmt.Println("Example trace commands:")
		fmt.Println()
		fmt.Printf("  # Trace all NSURLSession methods:\n")
		fmt.Printf("  frida-trace -D %s -f %s -m '*[NSURLSession* *]'\n", sim.UDID, target)
		fmt.Println()
		fmt.Printf("  # Trace SSL/TLS functions:\n")
		fmt.Printf("  frida-trace -D %s -f %s -i 'SSL_*'\n", sim.UDID, target)
		fmt.Println()
		fmt.Printf("  # Trace Keychain access:\n")
		fmt.Printf("  frida-trace -D %s -f %s -i 'SecItem*'\n", sim.UDID, target)
		fmt.Println()
		fmt.Printf("  # Trace file operations:\n")
		fmt.Printf("  frida-trace -D %s -f %s -i 'open*' -i 'read*' -i 'write*'\n", sim.UDID, target)
		fmt.Println()
		fmt.Printf("  # Trace UIKit view lifecycle:\n")
		fmt.Printf("  frida-trace -D %s -f %s -m '*[UIViewController viewDid*]'\n", sim.UDID, target)
		fmt.Println()
		return
	}

	// Build frida-trace command
	traceArgs := []string{"-D", sim.UDID, "-f", target}
	// Add pattern (could be -m or -i based on content)
	traceArgs = append(traceArgs, pattern)

	fmt.Println()
	fmt.Printf("Running: %s %s\n", fridaTracePath, strings.Join(traceArgs, " "))
	fmt.Println()

	// Execute interactively
	execCmd := exec.Command(fridaTracePath, traceArgs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		output.Error("frida-trace exited with error: %v", err)
	}
}

func runFridaKill(cmd *cobra.Command, args []string) {
	target := args[0]

	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No booted simulator: %v", err)
		return
	}

	output.Info("Terminating %s on %s...", target, sim.Name)

	// Try to terminate via simctl first (works for bundle IDs)
	if err := ios.TerminateApp(sim.UDID, target); err != nil {
		output.Warning("simctl terminate failed: %v", err)
		output.Info("The target might be a process name, not a bundle ID")
		return
	}

	time.Sleep(500 * time.Millisecond)
	output.Success("App terminated: %s", target)
}

func testFridaConnectivity() bool {
	fridaPsPath := frida.Locate("frida-ps")
	if fridaPsPath == "" {
		return false
	}

	_, err := ios.ExecuteCommand(fridaPsPath, "-a")
	return err == nil
}

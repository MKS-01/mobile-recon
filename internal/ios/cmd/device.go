package cmd

import (
	"fmt"
	"strings"

	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/MKS-01/mobile-recon/internal/ios"
	"github.com/spf13/cobra"
)

var (
	showAll bool
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "iOS Simulator device management",
	Long:  `Commands for managing iOS Simulator devices.`,
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available simulators",
	Run:   runDeviceList,
}

var deviceBootCmd = &cobra.Command{
	Use:   "boot <udid_or_name>",
	Short: "Boot a simulator",
	Args:  cobra.ExactArgs(1),
	Run:   runDeviceBoot,
}

var deviceShutdownCmd = &cobra.Command{
	Use:   "shutdown [udid_or_name]",
	Short: "Shutdown a simulator (default: all booted)",
	Args:  cobra.MaximumNArgs(1),
	Run:   runDeviceShutdown,
}

var deviceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed info about the current simulator",
	Run:   runDeviceInfo,
}

func init() {
	RootCmd.AddCommand(deviceCmd)

	deviceListCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all simulators (including unavailable)")
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceBootCmd)
	deviceCmd.AddCommand(deviceShutdownCmd)
	deviceCmd.AddCommand(deviceInfoCmd)
}

func runDeviceList(cmd *cobra.Command, args []string) {
	output.Section("iOS Simulators")

	simulators, err := ios.GetSimulators()
	if err != nil {
		output.Error("Failed to list simulators: %v", err)
		return
	}

	// Group by runtime
	byRuntime := make(map[string][]ios.Simulator)
	for _, sim := range simulators {
		if !showAll && !sim.IsAvailable {
			continue
		}
		byRuntime[sim.Runtime] = append(byRuntime[sim.Runtime], sim)
	}

	for runtime, sims := range byRuntime {
		// Clean up runtime name
		runtimeName := runtime
		if strings.Contains(runtime, "iOS") {
			parts := strings.Split(runtime, ".")
			if len(parts) > 0 {
				runtimeName = parts[len(parts)-1]
				runtimeName = strings.ReplaceAll(runtimeName, "-", " ")
				runtimeName = strings.ReplaceAll(runtimeName, "SimRuntime", "")
			}
		}

		fmt.Printf("\n%s:\n", runtimeName)
		for _, sim := range sims {
			stateIcon := "  "
			if sim.State == "Booted" {
				stateIcon = "* "
			}
			fmt.Printf("  %s%-30s %s (%s)\n", stateIcon, sim.Name, sim.State, sim.UDID[:8])
		}
	}

	fmt.Println()
	output.Info("* = Booted")
	output.Info("Use -a to show unavailable simulators")
}

func runDeviceBoot(cmd *cobra.Command, args []string) {
	target := args[0]

	output.Info("Booting simulator: %s", target)

	// First try to find by UDID
	sim, err := ios.GetSimulatorByUDID(target)
	if err != nil {
		// Try to find by name
		simulators, _ := ios.GetSimulators()
		for _, s := range simulators {
			if strings.EqualFold(s.Name, target) || strings.Contains(strings.ToLower(s.Name), strings.ToLower(target)) {
				sim = &s
				break
			}
		}
	}

	if sim == nil {
		output.Error("Simulator not found: %s", target)
		return
	}

	if sim.State == "Booted" {
		output.Warning("Simulator already booted: %s", sim.Name)
		return
	}

	if err := ios.BootSimulator(sim.UDID); err != nil {
		output.Error("Failed to boot simulator: %v", err)
		return
	}

	output.Success("Booted: %s (%s)", sim.Name, sim.UDID[:8])
	output.Info("Opening Simulator.app...")

	// Open Simulator.app
	ios.ExecuteCommand("open", "-a", "Simulator")
}

func runDeviceShutdown(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// Shutdown all booted simulators
		booted, err := ios.GetBootedSimulators()
		if err != nil {
			output.Error("Failed to get simulators: %v", err)
			return
		}

		if len(booted) == 0 {
			output.Info("No simulators are booted")
			return
		}

		output.Info("Shutting down %d simulator(s)...", len(booted))
		for _, sim := range booted {
			if err := ios.ShutdownSimulator(sim.UDID); err != nil {
				output.Warning("Failed to shutdown %s: %v", sim.Name, err)
			} else {
				output.Success("Shutdown: %s", sim.Name)
			}
		}
		return
	}

	target := args[0]
	sim, err := ios.GetSimulatorByUDID(target)
	if err != nil {
		simulators, _ := ios.GetSimulators()
		for _, s := range simulators {
			if strings.EqualFold(s.Name, target) {
				sim = &s
				break
			}
		}
	}

	if sim == nil {
		output.Error("Simulator not found: %s", target)
		return
	}

	if sim.State != "Booted" {
		output.Info("Simulator is not booted: %s", sim.Name)
		return
	}

	output.Info("Shutting down: %s", sim.Name)
	if err := ios.ShutdownSimulator(sim.UDID); err != nil {
		output.Error("Failed to shutdown: %v", err)
		return
	}

	output.Success("Simulator shutdown complete")
}

func runDeviceInfo(cmd *cobra.Command, args []string) {
	sim, err := getTargetSimulator()
	if err != nil {
		output.Error("No target simulator: %v", err)
		fmt.Println("\n  Boot a simulator first:")
		fmt.Println("  ios-toolkit device boot <name>")
		return
	}

	output.Section("Simulator Info")

	output.Success("Name: %s", sim.Name)
	output.Success("UDID: %s", sim.UDID)
	output.Success("State: %s", sim.State)
	output.Info("Runtime: %s", sim.Runtime)
	output.Info("Device Type: %s", sim.DeviceType)

	if sim.DataPath != "" {
		output.Info("Data Path: %s", sim.DataPath)
	}

	// Show architecture
	arch, _ := ios.GetArchitecture()
	output.Info("Host Architecture: %s", arch)

	// Show installed apps count
	apps, _ := ios.GetInstalledApps(sim.UDID)
	if len(apps) > 0 {
		output.Info("Installed Apps: %d", len(apps))
	}
}

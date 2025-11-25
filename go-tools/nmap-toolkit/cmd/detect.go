// Package cmd implements detection commands
package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// detectCmd represents the detect command group
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Service and OS detection commands",
	Long:  `Perform service version detection and operating system fingerprinting.`,
}

// serviceDetectCmd detects services and versions
var serviceDetectCmd = &cobra.Command{
	Use:   "service [target]",
	Short: "Detect services and versions running on open ports",
	Long: `Performs service and version detection (-sV) to identify what services
are running on open ports and their versions.

Examples:
  nmap-toolkit detect service 192.168.1.1
  nmap-toolkit detect service 192.168.1.1 -p 80,443,8080
  nmap-toolkit detect service 192.168.1.1 --aggressive
  nmap-toolkit detect service 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🔬 Service Version Detection")
		utils.PrintInfo("Target: %s", target)
		if ports != "" {
			utils.PrintInfo("Ports: %s", ports)
		}
		if aggressive {
			utils.PrintInfo("Mode: Aggressive (intensity 9)")
		}

		result, err := nmap.ServiceVersionScan(target, ports, aggressive, stream)
		if err != nil {
			utils.PrintError("Detection failed: %v", err)
			return
		}

		utils.PrintSuccess("Detection completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
		}
	},
}

// osDetectCmd detects operating system
var osDetectCmd = &cobra.Command{
	Use:   "os [target]",
	Short: "Detect operating system (requires root/admin)",
	Long: `Performs OS detection (-O) to identify the operating system running on the target.
This requires root/administrator privileges.

Examples:
  sudo nmap-toolkit detect os 192.168.1.1
  sudo nmap-toolkit detect os 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("💻 Operating System Detection")
		utils.PrintInfo("Target: %s", target)
		utils.PrintWarning("This requires root/administrator privileges")

		result, err := nmap.OSDetection(target, stream)
		if err != nil {
			utils.PrintError("Detection failed: %v", err)
			utils.PrintInfo("Try running with sudo/administrator privileges")
			return
		}

		utils.PrintSuccess("Detection completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
		}
	},
}

// aggressiveDetectCmd performs aggressive detection
var aggressiveDetectCmd = &cobra.Command{
	Use:   "aggressive [target]",
	Short: "Aggressive scan with OS, service, script scanning and traceroute",
	Long: `Performs an aggressive scan (-A) that enables:
  • OS detection
  • Version detection
  • Script scanning
  • Traceroute

This provides maximum information but is more detectable.
May require root/administrator privileges for some features.

Examples:
  nmap-toolkit detect aggressive 192.168.1.1
  nmap-toolkit detect aggressive 192.168.1.1 -p 1-1000
  sudo nmap-toolkit detect aggressive 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("⚡ Aggressive Detection Scan")
		utils.PrintInfo("Target: %s", target)
		if ports != "" {
			utils.PrintInfo("Ports: %s", ports)
		}
		utils.PrintWarning("This scan is highly detectable and may take longer")
		utils.PrintInfo("Includes: OS detection, version detection, scripts, traceroute")

		result, err := nmap.AggressiveScan(target, ports, stream)
		if err != nil {
			utils.PrintError("Scan failed: %v", err)
			return
		}

		utils.PrintSuccess("Scan completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
		}
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

	// Add subcommands
	detectCmd.AddCommand(serviceDetectCmd)
	detectCmd.AddCommand(osDetectCmd)
	detectCmd.AddCommand(aggressiveDetectCmd)

	// Flags for service detection
	serviceDetectCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	serviceDetectCmd.Flags().BoolVar(&aggressive, "aggressive", false, "Aggressive version detection")
	serviceDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for OS detection
	osDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for aggressive detection
	aggressiveDetectCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	aggressiveDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

// Package cmd implements scan-related commands
package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	ports      string
	fast       bool
	aggressive bool
	stream     bool
)

// scanCmd represents the scan command group
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform various types of network scans",
	Long:  `Execute different types of network scans including quick scans, port scans, and stealth scans.`,
}

// quickScanCmd performs a quick ping scan
var quickScanCmd = &cobra.Command{
	Use:   "quick [target]",
	Short: "Quick ping scan to discover live hosts",
	Long: `Performs a quick ping scan (-sn) to discover live hosts without port scanning.
This is useful for initial network reconnaissance.

Example:
  nmap-toolkit scan quick 192.168.1.0/24
  nmap-toolkit scan quick 192.168.1.1-50`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🔍 Quick Host Discovery")
		utils.PrintInfo("Target: %s", target)

		result, err := nmap.QuickScan(target)
		if err != nil {
			utils.PrintError("Scan failed: %v", err)
			return
		}

		utils.PrintSuccess("Scan completed")
		fmt.Println()
		utils.PrintData(result.Output)
	},
}

// portScanCmd performs a port scan
var portScanCmd = &cobra.Command{
	Use:   "port [target]",
	Short: "TCP port scan on target",
	Long: `Performs a TCP port scan on the specified target.
You can specify custom ports or use common port ranges.

Examples:
  nmap-toolkit scan port 192.168.1.1
  nmap-toolkit scan port 192.168.1.1 -p 80,443,8080
  nmap-toolkit scan port 192.168.1.1 -p 1-1000
  nmap-toolkit scan port 192.168.1.1 --fast
  nmap-toolkit scan port 192.168.1.1 --aggressive --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🔍 TCP Port Scan")
		utils.PrintInfo("Target: %s", target)
		if ports != "" {
			utils.PrintInfo("Ports: %s", ports)
		}

		options := map[string]bool{
			"verbose":    verbose,
			"fast":       fast,
			"aggressive": aggressive,
			"stream":     stream,
		}

		result, err := nmap.PortScan(target, ports, options)
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

// stealthScanCmd performs a SYN stealth scan
var stealthScanCmd = &cobra.Command{
	Use:   "stealth [target]",
	Short: "SYN stealth scan (requires root/admin)",
	Long: `Performs a SYN stealth scan (-sS) which is less likely to be logged.
This scan type requires root/administrator privileges.

Examples:
  sudo nmap-toolkit scan stealth 192.168.1.1
  sudo nmap-toolkit scan stealth 192.168.1.1 -p 1-65535
  sudo nmap-toolkit scan stealth 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🕵️  SYN Stealth Scan")
		utils.PrintInfo("Target: %s", target)
		utils.PrintWarning("This scan requires root/administrator privileges")

		result, err := nmap.StealthScan(target, ports, stream)
		if err != nil {
			utils.PrintError("Scan failed: %v", err)
			utils.PrintInfo("Try running with sudo/administrator privileges")
			return
		}

		utils.PrintSuccess("Scan completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
		}
	},
}

// udpScanCmd performs a UDP scan
var udpScanCmd = &cobra.Command{
	Use:   "udp [target]",
	Short: "UDP port scan (requires root/admin)",
	Long: `Performs a UDP port scan (-sU) to discover UDP services.
This scan type requires root/administrator privileges and can be slow.

Examples:
  sudo nmap-toolkit scan udp 192.168.1.1
  sudo nmap-toolkit scan udp 192.168.1.1 -p 53,67,68,161
  sudo nmap-toolkit scan udp 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("📡 UDP Port Scan")
		utils.PrintInfo("Target: %s", target)
		utils.PrintWarning("This scan requires root/administrator privileges and may take longer")

		result, err := nmap.UDPScan(target, ports, stream)
		if err != nil {
			utils.PrintError("Scan failed: %v", err)
			utils.PrintInfo("Try running with sudo/administrator privileges")
			return
		}

		utils.PrintSuccess("Scan completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
		}
	},
}

// networkScanCmd scans an entire network
var networkScanCmd = &cobra.Command{
	Use:   "network [network]",
	Short: "Scan entire network for live hosts",
	Long: `Performs a comprehensive network scan to discover all live hosts.

Examples:
  nmap-toolkit scan network 192.168.1.0/24
  nmap-toolkit scan network 10.0.0.0/16 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		utils.PrintHeader("🌐 Network Discovery Scan")
		utils.PrintInfo("Network: %s", network)

		result, err := nmap.ScanNetwork(network, stream)
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
	rootCmd.AddCommand(scanCmd)

	// Add subcommands
	scanCmd.AddCommand(quickScanCmd)
	scanCmd.AddCommand(portScanCmd)
	scanCmd.AddCommand(stealthScanCmd)
	scanCmd.AddCommand(udpScanCmd)
	scanCmd.AddCommand(networkScanCmd)

	// Flags for port scanning
	portScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (e.g., 80,443 or 1-1000)")
	portScanCmd.Flags().BoolVar(&fast, "fast", false, "Fast scan (top 100 ports)")
	portScanCmd.Flags().BoolVar(&aggressive, "aggressive", false, "Aggressive timing (-T4)")
	portScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for stealth scan
	stealthScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (default: top 1000)")
	stealthScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for UDP scan
	udpScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (default: top 100)")
	udpScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for network scan
	networkScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

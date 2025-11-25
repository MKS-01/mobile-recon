// Package cmd implements mobile-specific reconnaissance commands
package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

// mobileCmd represents the mobile reconnaissance command group
var mobileCmd = &cobra.Command{
	Use:   "mobile",
	Short: "Mobile device reconnaissance and security testing",
	Long: `Specialized commands for mobile device reconnaissance including:
  • Android ADB device discovery
  • iOS service detection
  • Mobile app service ports
  • Network-based mobile device identification`,
}

// androidADBCmd scans for Android ADB devices
var androidADBCmd = &cobra.Command{
	Use:   "adb [network]",
	Short: "Scan for Android ADB devices on network",
	Long: `Scans the network for Android devices with ADB (Android Debug Bridge) enabled.
This looks for devices with TCP/IP debugging enabled on common ADB ports (5555-5559).

Examples:
  nmap-toolkit mobile adb 192.168.1.0/24
  nmap-toolkit mobile adb 10.0.0.0/16 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		utils.PrintHeader("🤖 Android ADB Device Discovery")
		utils.PrintInfo("Network: %s", network)
		utils.PrintInfo("Scanning ports: 5555-5559")
		utils.PrintWarning("Only devices with TCP/IP debugging enabled will be detected")

		result, err := nmap.AndroidADBScan(network, stream)
		if err != nil {
			utils.PrintError("Scan failed: %v", err)
			return
		}

		utils.PrintSuccess("Scan completed")
		if !stream {
			fmt.Println()
			utils.PrintData(result.Output)
			fmt.Println()
			utils.PrintInfo("To connect to a discovered device:")
			utils.PrintData("  adb connect <ip>:5555")
		}
	},
}

// mobileScanCmd performs mobile-optimized scanning
var mobileScanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan target for mobile services and ports",
	Long: `Performs a scan optimized for mobile device detection, checking:
  • Common mobile service ports
  • Web services (HTTP/HTTPS)
  • ADB ports (5555, 27042)
  • Mobile debugging ports
  • iOS services (62078 - lockdown)

Examples:
  nmap-toolkit mobile scan 192.168.1.1
  nmap-toolkit mobile scan 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("📱 Mobile Device Scan")
		utils.PrintInfo("Target: %s", target)
		utils.PrintInfo("Scanning mobile-specific ports and services")

		result, err := nmap.MobileScan(target, stream)
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

// iosCmd scans for iOS devices
var iosCmd = &cobra.Command{
	Use:   "ios [network]",
	Short: "Scan for iOS devices on network",
	Long: `Scans the network for iOS devices by looking for:
  • Lockdown service (port 62078)
  • Apple mobile services
  • iTunes sync ports

Examples:
  nmap-toolkit mobile ios 192.168.1.0/24
  nmap-toolkit mobile ios 10.0.0.0/24 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		utils.PrintHeader("🍎 iOS Device Discovery")
		utils.PrintInfo("Network: %s", network)
		utils.PrintInfo("Scanning for iOS lockdown service (port 62078)")

		nmapArgs := []string{
			"-p", "62078",
			"--open",
			"-sV",
			"-v",
			network,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
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

// mitmProxyCmd scans for devices potentially using mitmproxy
var mitmProxyCmd = &cobra.Command{
	Use:   "mitm [network]",
	Short: "Detect potential MITM proxies on network",
	Long: `Scans for common MITM proxy ports used for mobile testing:
  • 8080 - Common HTTP proxy
  • 8443 - Common HTTPS proxy
  • 8888 - Burp Suite default
  • 9090 - Charles Proxy default

Examples:
  nmap-toolkit mobile mitm 192.168.1.0/24
  nmap-toolkit mobile mitm 192.168.1.0/24 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		utils.PrintHeader("🔍 MITM Proxy Detection")
		utils.PrintInfo("Network: %s", network)
		utils.PrintInfo("Scanning for common proxy ports")

		nmapArgs := []string{
			"-p", "8080,8443,8888,9090",
			"--open",
			"-sV",
			"-v",
			network,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
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

// appPortsCmd scans common mobile app backend ports
var appPortsCmd = &cobra.Command{
	Use:   "app-ports [target]",
	Short: "Scan common mobile application backend ports",
	Long: `Scans for ports commonly used by mobile application backends:
  • REST APIs (80, 443, 8080, 8443)
  • WebSocket servers (9000, 3000)
  • Database ports (3306, 5432, 27017)
  • Redis/Cache (6379)
  • Message queues (5672, 9092)

Examples:
  nmap-toolkit mobile app-ports 192.168.1.1
  nmap-toolkit mobile app-ports api.example.com --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🔌 Mobile App Backend Ports")
		utils.PrintInfo("Target: %s", target)

		nmapArgs := []string{
			"-p", "80,443,3000,3306,5432,5672,6379,8080,8443,9000,9092,27017",
			"--open",
			"-sV",
			"-v",
			target,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
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
	rootCmd.AddCommand(mobileCmd)

	// Add subcommands
	mobileCmd.AddCommand(androidADBCmd)
	mobileCmd.AddCommand(mobileScanCmd)
	mobileCmd.AddCommand(iosCmd)
	mobileCmd.AddCommand(mitmProxyCmd)
	mobileCmd.AddCommand(appPortsCmd)

	// Flags
	androidADBCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	mobileScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	iosCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	mitmProxyCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	appPortsCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

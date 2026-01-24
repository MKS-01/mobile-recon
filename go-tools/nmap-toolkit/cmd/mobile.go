package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var mobileCmd = &cobra.Command{
	Use:   "mobile",
	Short: "Mobile device reconnaissance and security testing",
	Long: `Specialized commands for mobile device reconnaissance including:
  - Android ADB device discovery
  - iOS service detection
  - Mobile app service ports
  - Network-based mobile device identification`,
}

var androidADBCmd = &cobra.Command{
	Use:   "adb [network]",
	Short: "Scan for Android ADB devices on network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		output.Header("Android ADB Device Discovery")
		output.Info("Network: %s", network)
		output.Info("Scanning ports: 5555-5559")
		output.Warning("Only devices with TCP/IP debugging enabled will be detected")

		result, err := nmap.AndroidADBScan(network, stream)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
			fmt.Println()
			output.Info("To connect to a discovered device:")
			output.Data("  adb connect <ip>:5555")
		}
	},
}

var mobileScanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan target for mobile services and ports",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Mobile Device Scan")
		output.Info("Target: %s", target)
		output.Info("Scanning mobile-specific ports and services")

		result, err := nmap.MobileScan(target, stream)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

var iosCmd = &cobra.Command{
	Use:   "ios [network]",
	Short: "Scan for iOS devices on network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		output.Header("iOS Device Discovery")
		output.Info("Network: %s", network)
		output.Info("Scanning for iOS lockdown service (port 62078)")

		nmapArgs := []string{
			"-p", "62078",
			"--open",
			"-sV",
			"-v",
			network,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

var mitmProxyCmd = &cobra.Command{
	Use:   "mitm [network]",
	Short: "Detect potential MITM proxies on network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		output.Header("MITM Proxy Detection")
		output.Info("Network: %s", network)
		output.Info("Scanning for common proxy ports")

		nmapArgs := []string{
			"-p", "8080,8443,8888,9090",
			"--open",
			"-sV",
			"-v",
			network,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

var appPortsCmd = &cobra.Command{
	Use:   "app-ports [target]",
	Short: "Scan common mobile application backend ports",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Mobile App Backend Ports")
		output.Info("Target: %s", target)

		nmapArgs := []string{
			"-p", "80,443,3000,3306,5432,5672,6379,8080,8443,9000,9092,27017",
			"--open",
			"-sV",
			"-v",
			target,
		}

		result, err := nmap.CustomScan(nmapArgs, stream)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

func init() {
	RootCmd.AddCommand(mobileCmd)

	mobileCmd.AddCommand(androidADBCmd)
	mobileCmd.AddCommand(mobileScanCmd)
	mobileCmd.AddCommand(iosCmd)
	mobileCmd.AddCommand(mitmProxyCmd)
	mobileCmd.AddCommand(appPortsCmd)

	androidADBCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	mobileScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	iosCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	mitmProxyCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	appPortsCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

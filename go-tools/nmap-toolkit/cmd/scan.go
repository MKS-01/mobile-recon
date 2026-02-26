package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var (
	ports      string
	fast       bool
	aggressive bool
	stream     bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform various types of network scans",
	Long:  `Execute different types of network scans including quick host discovery, port scans, and network-wide scans.`,
}

var quickScanCmd = &cobra.Command{
	Use:   "quick [target]",
	Short: "Quick ping scan to discover live hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Quick Host Discovery")
		output.Info("Target: %s", target)

		result, err := nmap.QuickScan(target)
		if err != nil {
			output.Error("Scan failed: %v", err)
			return
		}

		output.Success("Scan completed")
		fmt.Println()
		output.Data(result.Output)
	},
}

var portScanCmd = &cobra.Command{
	Use:   "port [target]",
	Short: "TCP port scan on target",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("TCP Port Scan")
		output.Info("Target: %s", target)
		if ports != "" {
			output.Info("Ports: %s", ports)
		}

		options := map[string]bool{
			"verbose":    verbose,
			"fast":       fast,
			"aggressive": aggressive,
			"stream":     stream,
		}

		result, err := nmap.PortScan(target, ports, options)
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

var networkScanCmd = &cobra.Command{
	Use:   "network [network]",
	Short: "Scan entire network for live hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		output.Header("Network Discovery Scan")
		output.Info("Network: %s", network)

		result, err := nmap.ScanNetwork(network, stream)
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
	RootCmd.AddCommand(scanCmd)

	scanCmd.AddCommand(quickScanCmd)
	scanCmd.AddCommand(portScanCmd)
	scanCmd.AddCommand(networkScanCmd)

	portScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (e.g., 80,443 or 1-1000)")
	portScanCmd.Flags().BoolVar(&fast, "fast", false, "Fast scan (top 100 ports)")
	portScanCmd.Flags().BoolVar(&aggressive, "aggressive", false, "Aggressive timing (-T4)")
	portScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	networkScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

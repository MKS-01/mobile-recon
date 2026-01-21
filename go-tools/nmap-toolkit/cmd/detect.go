package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Service and OS detection commands",
	Long:  `Perform service version detection and operating system fingerprinting.`,
}

var serviceDetectCmd = &cobra.Command{
	Use:   "service [target]",
	Short: "Detect services and versions running on open ports",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Service Version Detection")
		output.Info("Target: %s", target)
		if ports != "" {
			output.Info("Ports: %s", ports)
		}
		if aggressive {
			output.Info("Mode: Aggressive (intensity 9)")
		}

		result, err := nmap.ServiceVersionScan(target, ports, aggressive, stream)
		if err != nil {
			output.Error("Detection failed: %v", err)
			return
		}

		output.Success("Detection completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

var osDetectCmd = &cobra.Command{
	Use:   "os [target]",
	Short: "Detect operating system (requires root/admin)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Operating System Detection")
		output.Info("Target: %s", target)
		output.Warning("This requires root/administrator privileges")

		result, err := nmap.OSDetection(target, stream)
		if err != nil {
			output.Error("Detection failed: %v", err)
			output.Info("Try running with sudo/administrator privileges")
			return
		}

		output.Success("Detection completed")
		if !stream {
			fmt.Println()
			output.Data(result.Output)
		}
	},
}

var aggressiveDetectCmd = &cobra.Command{
	Use:   "aggressive [target]",
	Short: "Aggressive scan with OS, service, script scanning and traceroute",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Aggressive Detection Scan")
		output.Info("Target: %s", target)
		if ports != "" {
			output.Info("Ports: %s", ports)
		}
		output.Warning("This scan is highly detectable and may take longer")
		output.Info("Includes: OS detection, version detection, scripts, traceroute")

		result, err := nmap.AggressiveScan(target, ports, stream)
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
	rootCmd.AddCommand(detectCmd)

	detectCmd.AddCommand(serviceDetectCmd)
	detectCmd.AddCommand(osDetectCmd)
	detectCmd.AddCommand(aggressiveDetectCmd)

	serviceDetectCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	serviceDetectCmd.Flags().BoolVar(&aggressive, "aggressive", false, "Aggressive version detection")
	serviceDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	osDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	aggressiveDetectCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	aggressiveDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/pkg/output"
	"github.com/MKS-01/mobile-recon/internal/nmap"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Service detection commands",
	Long:  `Perform service version detection on target hosts.`,
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

func init() {
	RootCmd.AddCommand(detectCmd)

	detectCmd.AddCommand(serviceDetectCmd)

	serviceDetectCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	serviceDetectCmd.Flags().BoolVar(&aggressive, "aggressive", false, "Aggressive version detection")
	serviceDetectCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

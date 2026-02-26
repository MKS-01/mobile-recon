package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "SSL/TLS security testing",
	Long:  `Test SSL/TLS configuration and cipher suites on target services.`,
}

var sslScanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "SSL/TLS enumeration and cipher testing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port := "443"
		if ports != "" {
			port = ports
		}

		output.Header("SSL/TLS Enumeration")
		output.Info("Target: %s", target)
		output.Info("Port: %s", port)

		result, err := nmap.SSLScan(target, port, stream)
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
	RootCmd.AddCommand(sslCmd)

	sslCmd.AddCommand(sslScanCmd)

	sslScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Port to scan (default: 443)")
	sslScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
}

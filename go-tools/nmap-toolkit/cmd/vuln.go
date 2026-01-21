package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/common/output"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/spf13/cobra"
)

var scripts string

var vulnCmd = &cobra.Command{
	Use:   "vuln",
	Short: "Vulnerability scanning and NSE scripts",
	Long:  `Perform vulnerability scanning using Nmap Scripting Engine (NSE) scripts.`,
}

var vulnScanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan for common vulnerabilities",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		output.Header("Vulnerability Scan")
		output.Info("Target: %s", target)
		if ports != "" {
			output.Info("Ports: %s", ports)
		}
		output.Warning("This scan may take several minutes")

		result, err := nmap.VulnerabilityScan(target, ports, stream)
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

var sslScanCmd = &cobra.Command{
	Use:   "ssl [target]",
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

var scriptScanCmd = &cobra.Command{
	Use:   "script [target]",
	Short: "Run custom NSE scripts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		if scripts == "" {
			output.Error("Please specify scripts using --scripts flag")
			return
		}

		output.Header("NSE Script Scan")
		output.Info("Target: %s", target)
		output.Info("Scripts: %s", scripts)
		if ports != "" {
			output.Info("Ports: %s", ports)
		}

		result, err := nmap.ScriptScan(target, scripts, ports, stream)
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
	rootCmd.AddCommand(vulnCmd)

	vulnCmd.AddCommand(vulnScanCmd)
	vulnCmd.AddCommand(sslScanCmd)
	vulnCmd.AddCommand(scriptScanCmd)

	vulnScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	vulnScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	sslScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Port to scan (default: 443)")
	sslScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	scriptScanCmd.Flags().StringVar(&scripts, "scripts", "", "NSE scripts to run (required)")
	scriptScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	scriptScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	scriptScanCmd.MarkFlagRequired("scripts")
}

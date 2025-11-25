// Package cmd implements vulnerability scanning commands
package cmd

import (
	"fmt"

	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/nmap"
	"github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	scripts string
)

// vulnCmd represents the vulnerability scanning command group
var vulnCmd = &cobra.Command{
	Use:   "vuln",
	Short: "Vulnerability scanning and NSE scripts",
	Long:  `Perform vulnerability scanning using Nmap Scripting Engine (NSE) scripts.`,
}

// vulnScanCmd performs vulnerability scanning
var vulnScanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "Scan for common vulnerabilities",
	Long: `Performs vulnerability scanning using NSE vuln category scripts.
This can detect common vulnerabilities like:
  • CVE exploits
  • Default credentials
  • Configuration issues
  • Known vulnerabilities

Examples:
  nmap-toolkit vuln scan 192.168.1.1
  nmap-toolkit vuln scan 192.168.1.1 -p 80,443
  nmap-toolkit vuln scan 192.168.1.1 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		utils.PrintHeader("🔓 Vulnerability Scan")
		utils.PrintInfo("Target: %s", target)
		if ports != "" {
			utils.PrintInfo("Ports: %s", ports)
		}
		utils.PrintWarning("This scan may take several minutes")

		result, err := nmap.VulnerabilityScan(target, ports, stream)
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

// sslScanCmd performs SSL/TLS enumeration
var sslScanCmd = &cobra.Command{
	Use:   "ssl [target]",
	Short: "SSL/TLS enumeration and cipher testing",
	Long: `Performs comprehensive SSL/TLS testing including:
  • Supported SSL/TLS versions
  • Cipher suites and strengths
  • Certificate information
  • SSL/TLS vulnerabilities

Examples:
  nmap-toolkit vuln ssl 192.168.1.1
  nmap-toolkit vuln ssl example.com --ports 443
  nmap-toolkit vuln ssl 192.168.1.1 --ports 8443 --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		port := "443"
		if ports != "" {
			port = ports
		}

		utils.PrintHeader("🔒 SSL/TLS Enumeration")
		utils.PrintInfo("Target: %s", target)
		utils.PrintInfo("Port: %s", port)

		result, err := nmap.SSLScan(target, port, stream)
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

// scriptScanCmd runs custom NSE scripts
var scriptScanCmd = &cobra.Command{
	Use:   "script [target]",
	Short: "Run custom NSE scripts",
	Long: `Run specific Nmap Scripting Engine (NSE) scripts against a target.

Common useful scripts:
  • http-enum - Enumerate web directories
  • smb-enum-shares - Enumerate SMB shares
  • ftp-anon - Check for anonymous FTP
  • ssh-brute - SSH brute force
  • http-sql-injection - Test for SQL injection

Examples:
  nmap-toolkit vuln script 192.168.1.1 --scripts "http-enum"
  nmap-toolkit vuln script 192.168.1.1 --scripts "http-*" -p 80,443
  nmap-toolkit vuln script 192.168.1.1 --scripts "default,safe"
  nmap-toolkit vuln script 192.168.1.1 --scripts "smb-*" --stream`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		if scripts == "" {
			utils.PrintError("Please specify scripts using --scripts flag")
			return
		}

		utils.PrintHeader("📜 NSE Script Scan")
		utils.PrintInfo("Target: %s", target)
		utils.PrintInfo("Scripts: %s", scripts)
		if ports != "" {
			utils.PrintInfo("Ports: %s", ports)
		}

		result, err := nmap.ScriptScan(target, scripts, ports, stream)
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
	rootCmd.AddCommand(vulnCmd)

	// Add subcommands
	vulnCmd.AddCommand(vulnScanCmd)
	vulnCmd.AddCommand(sslScanCmd)
	vulnCmd.AddCommand(scriptScanCmd)

	// Flags for vulnerability scan
	vulnScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	vulnScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for SSL scan
	sslScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Port to scan (default: 443)")
	sslScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")

	// Flags for script scan
	scriptScanCmd.Flags().StringVar(&scripts, "scripts", "", "NSE scripts to run (required)")
	scriptScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")
	scriptScanCmd.Flags().BoolVar(&stream, "stream", false, "Stream output in real-time")
	scriptScanCmd.MarkFlagRequired("scripts")
}

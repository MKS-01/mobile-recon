// Package nmap provides a Go wrapper for Nmap network scanning functionality.
package nmap

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ScanResult represents the result of an nmap scan
type ScanResult struct {
	Target  string
	Output  string
	Command string
}

// IsNmapInstalled checks if nmap is installed and available in PATH
func IsNmapInstalled() bool {
	_, err := exec.LookPath("nmap")
	return err == nil
}

// GetVersion returns the installed nmap version
func GetVersion() (string, error) {
	cmd := exec.Command("nmap", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get nmap version: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return string(output), nil
}

// runNmapCommand executes an nmap command and returns the output
func runNmapCommand(args []string, streamOutput bool) (*ScanResult, error) {
	cmd := exec.Command("nmap", args...)

	result := &ScanResult{
		Command: "nmap " + strings.Join(args, " "),
	}

	// Extract target from args
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			result.Target = arg
			break
		}
	}

	if streamOutput {
		// Stream output in real-time
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start nmap: %v", err)
		}

		var outputBuf bytes.Buffer
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			outputBuf.WriteString(line + "\n")
		}

		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			fmt.Println(line)
			outputBuf.WriteString(line + "\n")
		}

		if err := cmd.Wait(); err != nil {
			return nil, fmt.Errorf("nmap scan failed: %v", err)
		}

		result.Output = outputBuf.String()
	} else {
		// Capture all output at once
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("nmap scan failed: %v\nOutput: %s", err, string(output))
		}
		result.Output = string(output)
	}

	return result, nil
}

// QuickScan performs a quick ping scan to discover live hosts
func QuickScan(target string) (*ScanResult, error) {
	args := []string{"-sn", target}
	return runNmapCommand(args, false)
}

// PortScan performs a standard port scan on specified ports
func PortScan(target string, ports string, options map[string]bool) (*ScanResult, error) {
	args := []string{}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	// Add common options
	if options["verbose"] {
		args = append(args, "-v")
	}
	if options["fast"] {
		args = append(args, "-F")
	}
	if options["aggressive"] {
		args = append(args, "-T4")
	}

	args = append(args, target)
	return runNmapCommand(args, options["stream"])
}

// ServiceVersionScan performs service and version detection
func ServiceVersionScan(target string, ports string, aggressive bool, stream bool) (*ScanResult, error) {
	args := []string{"-sV"}

	if aggressive {
		args = append(args, "--version-intensity", "9")
	} else {
		args = append(args, "--version-intensity", "5")
	}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, "-v", target)
	return runNmapCommand(args, stream)
}

// OSDetection performs operating system detection
func OSDetection(target string, stream bool) (*ScanResult, error) {
	args := []string{"-O", "-v", target}
	return runNmapCommand(args, stream)
}

// AggressiveScan performs a comprehensive aggressive scan
func AggressiveScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"-A", "-T4", "-v"}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, target)
	return runNmapCommand(args, stream)
}

// VulnerabilityScan performs vulnerability scanning using NSE scripts
func VulnerabilityScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{
		"--script", "vuln",
		"-sV",
		"-v",
	}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, target)
	return runNmapCommand(args, stream)
}

// SSLScan performs SSL/TLS enumeration and testing
func SSLScan(target string, port string, stream bool) (*ScanResult, error) {
	if port == "" {
		port = "443"
	}

	args := []string{
		"--script", "ssl-enum-ciphers,ssl-cert,ssl-date",
		"-p", port,
		"-v",
		target,
	}

	return runNmapCommand(args, stream)
}

// MobileScan performs a scan optimized for mobile device detection
func MobileScan(target string, stream bool) (*ScanResult, error) {
	// Scan common mobile service ports
	args := []string{
		"-p", "21,22,23,80,443,5555,8080,8443,9000,27042,62078",
		"-sV",
		"--script", "http-headers,http-title,banner",
		"-v",
		target,
	}

	return runNmapCommand(args, stream)
}

// AndroidADBScan scans for Android ADB devices on the network
func AndroidADBScan(network string, stream bool) (*ScanResult, error) {
	args := []string{
		"-p", "5555,5556,5557,5558,5559",
		"--open",
		"-v",
		network,
	}

	return runNmapCommand(args, stream)
}

// CustomScan allows running custom nmap commands
func CustomScan(nmapArgs []string, stream bool) (*ScanResult, error) {
	return runNmapCommand(nmapArgs, stream)
}

// ScanNetwork performs a comprehensive network scan
func ScanNetwork(network string, stream bool) (*ScanResult, error) {
	args := []string{
		"-sn",
		"--min-rate", "1000",
		network,
	}

	return runNmapCommand(args, stream)
}

// StealthScan performs a SYN stealth scan
func StealthScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"-sS"}

	if ports != "" {
		args = append(args, "-p", ports)
	} else {
		args = append(args, "--top-ports", "1000")
	}

	args = append(args, "-v", target)
	return runNmapCommand(args, stream)
}

// UDPScan performs a UDP scan
func UDPScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"-sU"}

	if ports != "" {
		args = append(args, "-p", ports)
	} else {
		args = append(args, "--top-ports", "100")
	}

	args = append(args, "-v", target)
	return runNmapCommand(args, stream)
}

// ScriptScan runs specific NSE scripts
func ScriptScan(target string, scripts string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"--script", scripts}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, "-v", target)
	return runNmapCommand(args, stream)
}

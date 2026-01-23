// Package nmap provides Go wrapper for Nmap network scanning.
package nmap

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ScanResult struct {
	Target  string
	Output  string
	Command string
}

func IsNmapInstalled() bool {
	_, err := exec.LookPath("nmap")
	return err == nil
}

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

func runNmapCommand(args []string, streamOutput bool) (*ScanResult, error) {
	cmd := exec.Command("nmap", args...)

	result := &ScanResult{
		Command: "nmap " + strings.Join(args, " "),
	}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			result.Target = arg
			break
		}
	}

	if streamOutput {
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
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("nmap scan failed: %v\nOutput: %s", err, string(output))
		}
		result.Output = string(output)
	}

	return result, nil
}

func QuickScan(target string) (*ScanResult, error) {
	args := []string{"-sn", target}
	return runNmapCommand(args, false)
}

func PortScan(target string, ports string, options map[string]bool) (*ScanResult, error) {
	args := []string{}

	if ports != "" {
		args = append(args, "-p", ports)
	}

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

func OSDetection(target string, stream bool) (*ScanResult, error) {
	args := []string{"-O", "-v", target}
	return runNmapCommand(args, stream)
}

func AggressiveScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"-A", "-T4", "-v"}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, target)
	return runNmapCommand(args, stream)
}

func VulnerabilityScan(target string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"--script", "vuln", "-sV", "-v"}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, target)
	return runNmapCommand(args, stream)
}

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

func MobileScan(target string, stream bool) (*ScanResult, error) {
	args := []string{
		"-p", "21,22,23,80,443,5555,8080,8443,9000,27042,62078",
		"-sV",
		"--script", "http-headers,http-title,banner",
		"-v",
		target,
	}

	return runNmapCommand(args, stream)
}

func AndroidADBScan(network string, stream bool) (*ScanResult, error) {
	args := []string{
		"-p", "5555,5556,5557,5558,5559",
		"--open",
		"-v",
		network,
	}

	return runNmapCommand(args, stream)
}

func CustomScan(nmapArgs []string, stream bool) (*ScanResult, error) {
	return runNmapCommand(nmapArgs, stream)
}

func ScanNetwork(network string, stream bool) (*ScanResult, error) {
	args := []string{"-sn", "--min-rate", "1000", network}
	return runNmapCommand(args, stream)
}

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

func ScriptScan(target string, scripts string, ports string, stream bool) (*ScanResult, error) {
	args := []string{"--script", scripts}

	if ports != "" {
		args = append(args, "-p", ports)
	}

	args = append(args, "-v", target)
	return runNmapCommand(args, stream)
}

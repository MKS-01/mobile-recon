// Package adb provides core functionality for interacting with Android Debug Bridge.
package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Device struct {
	Serial      string
	State       string
	Product     string
	Model       string
	Device      string
	TransportID string
}

func ExecuteCommand(args ...string) (string, error) {
	cmd := exec.Command("adb", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("adb error: %v - %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

func ExecuteCommandWithDevice(serial string, args ...string) (string, error) {
	cmdArgs := []string{"-s", serial}
	cmdArgs = append(cmdArgs, args...)
	return ExecuteCommand(cmdArgs...)
}

func GetDevices() ([]Device, error) {
	output, err := ExecuteCommand("devices", "-l")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	devices := make([]Device, 0)

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		device := Device{
			Serial: fields[0],
			State:  fields[1],
		}

		for _, field := range fields[2:] {
			parts := strings.Split(field, ":")
			if len(parts) == 2 {
				switch parts[0] {
				case "product":
					device.Product = parts[1]
				case "model":
					device.Model = parts[1]
				case "device":
					device.Device = parts[1]
				case "transport_id":
					device.TransportID = parts[1]
				}
			}
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func GetDefaultDevice() (*Device, error) {
	devices, err := GetDevices()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices connected")
	}

	return &devices[0], nil
}

func IsADBInstalled() bool {
	_, err := exec.LookPath("adb")
	return err == nil
}

func RestartAsRoot(serial string) (bool, error) {
	output, err := ExecuteCommandWithDevice(serial, "root")
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)

	if strings.Contains(output, "restarting adbd as root") ||
		strings.Contains(output, "already running as root") {
		return true, nil
	}

	if strings.Contains(output, "cannot run as root") ||
		strings.Contains(output, "production builds") {
		return false, nil
	}

	return true, nil
}

func GetDeviceArchitecture(serial string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.cpu.abi")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func GetDeviceProperty(serial string, property string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "getprop", property)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func PushFile(serial string, localPath string, remotePath string) error {
	_, err := ExecuteCommandWithDevice(serial, "push", localPath, remotePath)
	return err
}

func ShellCommand(serial string, command string) (string, error) {
	return ExecuteCommandWithDevice(serial, "shell", command)
}

func GetProcessPID(serial string, processName string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "pidof", processName)
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(output), nil
}

func IsDeviceReady(serial string) bool {
	output, err := GetDeviceProperty(serial, "sys.boot_completed")
	if err != nil {
		return false
	}
	return output == "1"
}

// Package adb provides core functionality for interacting with Android Debug Bridge (ADB).
// It includes device management, command execution, and device information retrieval.
package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Device represents an Android device connected via ADB.
// It contains all relevant metadata returned by the 'adb devices -l' command.
type Device struct {
	Serial      string // Unique device identifier (e.g., serial number or IP:port)
	State       string // Connection state (device, offline, unauthorized, etc.)
	Product     string // Product name (ro.product.name)
	Model       string // Device model (ro.product.model)
	Device      string // Device codename (ro.product.device)
	TransportID string // ADB transport layer identifier
}

// ExecuteCommand runs an ADB command with the provided arguments and returns the output.
// It captures both stdout and stderr, returning an error if the command fails.
//
// Parameters:
//   - args: Variable number of arguments to pass to the adb command
//
// Returns:
//   - string: Trimmed stdout output from the command
//   - error: Error with stderr details if command execution fails
//
// Example:
//   output, err := ExecuteCommand("devices", "-l")
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

// ExecuteCommandWithDevice runs an ADB command targeting a specific device by serial number.
// This is useful when multiple devices are connected and you need to target one specifically.
//
// Parameters:
//   - serial: Device serial number or identifier
//   - args: Variable number of arguments to pass to the adb command
//
// Returns:
//   - string: Command output
//   - error: Any error encountered during execution
//
// Example:
//   output, err := ExecuteCommandWithDevice("emulator-5554", "shell", "getprop")
func ExecuteCommandWithDevice(serial string, args ...string) (string, error) {
	cmdArgs := []string{"-s", serial}
	cmdArgs = append(cmdArgs, args...)
	return ExecuteCommand(cmdArgs...)
}

// GetDevices returns a list of all Android devices currently connected via ADB.
// It parses the output of 'adb devices -l' to extract detailed device information.
//
// Returns:
//   - []Device: Slice of Device structs containing device metadata
//   - error: Error if adb command fails or cannot be executed
//
// The function skips the header line and empty lines, parsing each device entry
// to extract serial, state, product, model, device name, and transport ID.
func GetDevices() ([]Device, error) {
	output, err := ExecuteCommand("devices", "-l")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	devices := make([]Device, 0)

	// Iterate through output lines, skipping header (index 0) and empty lines
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// First two fields are always serial and state
		device := Device{
			Serial: fields[0],
			State:  fields[1],
		}

		// Parse additional device metadata from key:value pairs
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

// GetDefaultDevice returns the first connected device from the device list.
// This is useful when only one device is connected or when you want to use
// the first available device as a default target.
//
// Returns:
//   - *Device: Pointer to the first connected device
//   - error: Error if no devices are connected or if GetDevices fails
//
// Note: If multiple devices are connected, this returns the first one in the list.
// For specific device targeting, use ExecuteCommandWithDevice with a serial number.
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

// IsADBInstalled checks if the ADB executable is available in the system PATH.
// This is useful for verifying that ADB is properly installed before attempting
// to execute ADB commands.
//
// Returns:
//   - bool: true if ADB is found in PATH, false otherwise
//
// Example:
//
//	if !IsADBInstalled() {
//	    fmt.Println("Please install Android SDK Platform Tools")
//	}
func IsADBInstalled() bool {
	_, err := exec.LookPath("adb")
	return err == nil
}

// RestartAsRoot restarts the ADB daemon on the device with root permissions.
// This is required for operations that need elevated privileges, such as
// accessing protected directories or running Frida server.
//
// Parameters:
//   - serial: Device serial number to target
//
// Returns:
//   - bool: true if root access was granted, false if device doesn't support root
//   - error: Error if the command fails
//
// Note: Root access is only available on:
//   - Emulators with google_apis images (NOT google_apis_playstore)
//   - Rooted physical devices
//   - userdebug/eng builds
func RestartAsRoot(serial string) (bool, error) {
	output, err := ExecuteCommandWithDevice(serial, "root")
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)

	// Check for success indicators
	if strings.Contains(output, "restarting adbd as root") ||
		strings.Contains(output, "already running as root") {
		return true, nil
	}

	// Check for failure indicators
	if strings.Contains(output, "cannot run as root") ||
		strings.Contains(output, "production builds") {
		return false, nil
	}

	// Unknown response - assume success if no error
	return true, nil
}

// GetDeviceArchitecture returns the CPU architecture of the connected device.
// This is useful for downloading architecture-specific binaries like Frida server.
//
// Parameters:
//   - serial: Device serial number to target
//
// Returns:
//   - string: Architecture string (arm64-v8a, armeabi-v7a, x86_64, x86)
//   - error: Error if the command fails
func GetDeviceArchitecture(serial string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "getprop", "ro.product.cpu.abi")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// GetDeviceProperty retrieves a system property from the device.
//
// Parameters:
//   - serial: Device serial number to target
//   - property: Property name (e.g., "ro.build.version.sdk")
//
// Returns:
//   - string: Property value
//   - error: Error if the command fails
func GetDeviceProperty(serial string, property string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "getprop", property)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// PushFile pushes a local file to the device.
//
// Parameters:
//   - serial: Device serial number to target
//   - localPath: Path to local file
//   - remotePath: Destination path on device
//
// Returns:
//   - error: Error if the push fails
func PushFile(serial string, localPath string, remotePath string) error {
	_, err := ExecuteCommandWithDevice(serial, "push", localPath, remotePath)
	return err
}

// ShellCommand executes a shell command on the device and returns the output.
//
// Parameters:
//   - serial: Device serial number to target
//   - command: Shell command to execute
//
// Returns:
//   - string: Command output
//   - error: Error if the command fails
func ShellCommand(serial string, command string) (string, error) {
	return ExecuteCommandWithDevice(serial, "shell", command)
}

// GetProcessPID returns the PID of a process by name.
//
// Parameters:
//   - serial: Device serial number to target
//   - processName: Name of the process to find
//
// Returns:
//   - string: PID if found, empty string if not running
//   - error: Error if the command fails
func GetProcessPID(serial string, processName string) (string, error) {
	output, err := ExecuteCommandWithDevice(serial, "shell", "pidof", processName)
	if err != nil {
		// pidof returns error if process not found, which is not an error for us
		return "", nil
	}
	return strings.TrimSpace(output), nil
}

// IsDeviceReady checks if the device has fully booted.
//
// Parameters:
//   - serial: Device serial number to target
//
// Returns:
//   - bool: true if device is ready, false otherwise
func IsDeviceReady(serial string) bool {
	output, err := GetDeviceProperty(serial, "sys.boot_completed")
	if err != nil {
		return false
	}
	return output == "1"
}

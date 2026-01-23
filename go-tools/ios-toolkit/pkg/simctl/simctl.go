// Package simctl provides core functionality for interacting with iOS Simulator via xcrun simctl.
package simctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Simulator represents an iOS Simulator device.
type Simulator struct {
	UDID        string `json:"udid"`
	Name        string `json:"name"`
	State       string `json:"state"`
	IsAvailable bool   `json:"isAvailable"`
	DeviceType  string `json:"deviceTypeIdentifier"`
	Runtime     string `json:"runtime"`
	DataPath    string `json:"dataPath"`
}

// App represents an installed application on a simulator.
type App struct {
	BundleID     string
	Name         string
	Path         string
	DataPath     string
	BundleVersion string
}

// Process represents a running process in the simulator.
type Process struct {
	PID  string
	Name string
}

// SimctlDevices represents the JSON structure from `simctl list devices`.
type SimctlDevices struct {
	Devices map[string][]Simulator `json:"devices"`
}

// ExecuteCommand runs a generic command and returns its output.
func ExecuteCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s error: %v - %s", name, err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// SimctlCommand runs an xcrun simctl command.
func SimctlCommand(args ...string) (string, error) {
	cmdArgs := append([]string{"simctl"}, args...)
	return ExecuteCommand("xcrun", cmdArgs...)
}

// SimctlCommandWithDevice runs an xcrun simctl command targeting a specific device.
func SimctlCommandWithDevice(udid string, args ...string) (string, error) {
	return SimctlCommand(append(args, udid)...)
}

// IsXcodeInstalled checks if Xcode command line tools are available.
func IsXcodeInstalled() bool {
	_, err := exec.LookPath("xcrun")
	return err == nil
}

// GetSimulators returns a list of all available simulators.
func GetSimulators() ([]Simulator, error) {
	output, err := SimctlCommand("list", "devices", "-j")
	if err != nil {
		return nil, err
	}

	var result SimctlDevices
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse simctl output: %v", err)
	}

	var simulators []Simulator
	for runtime, devices := range result.Devices {
		for _, device := range devices {
			device.Runtime = runtime
			simulators = append(simulators, device)
		}
	}

	return simulators, nil
}

// GetBootedSimulators returns only simulators that are currently booted.
func GetBootedSimulators() ([]Simulator, error) {
	simulators, err := GetSimulators()
	if err != nil {
		return nil, err
	}

	var booted []Simulator
	for _, sim := range simulators {
		if sim.State == "Booted" {
			booted = append(booted, sim)
		}
	}

	return booted, nil
}

// GetDefaultSimulator returns the first booted simulator or an error if none are booted.
func GetDefaultSimulator() (*Simulator, error) {
	booted, err := GetBootedSimulators()
	if err != nil {
		return nil, err
	}

	if len(booted) == 0 {
		return nil, fmt.Errorf("no booted simulators found. Start one with: xcrun simctl boot <device>")
	}

	return &booted[0], nil
}

// GetSimulatorByUDID finds a simulator by its UDID.
func GetSimulatorByUDID(udid string) (*Simulator, error) {
	simulators, err := GetSimulators()
	if err != nil {
		return nil, err
	}

	for _, sim := range simulators {
		if sim.UDID == udid {
			return &sim, nil
		}
	}

	return nil, fmt.Errorf("simulator with UDID %s not found", udid)
}

// GetInstalledApps returns the list of installed apps on a simulator.
func GetInstalledApps(udid string) ([]App, error) {
	output, err := SimctlCommand("listapps", udid)
	if err != nil {
		// listapps might not be available on older Xcode versions
		// Fall back to parsing the data directory
		return getInstalledAppsFromDataPath(udid)
	}

	return parseAppList(output), nil
}

// getInstalledAppsFromDataPath parses apps from the simulator's data path.
func getInstalledAppsFromDataPath(udid string) ([]App, error) {
	sim, err := GetSimulatorByUDID(udid)
	if err != nil {
		return nil, err
	}

	// Use find to locate Info.plist files in app containers
	output, err := ExecuteCommand("find", sim.DataPath+"/Containers/Bundle/Application", "-name", "Info.plist", "-maxdepth", "3")
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, plistPath := range strings.Split(output, "\n") {
		if plistPath == "" {
			continue
		}

		// Extract bundle ID from Info.plist
		bundleID, _ := ExecuteCommand("defaults", "read", plistPath, "CFBundleIdentifier")
		appName, _ := ExecuteCommand("defaults", "read", plistPath, "CFBundleName")
		version, _ := ExecuteCommand("defaults", "read", plistPath, "CFBundleShortVersionString")

		if bundleID != "" {
			apps = append(apps, App{
				BundleID:      bundleID,
				Name:          appName,
				Path:          strings.TrimSuffix(plistPath, "/Info.plist"),
				BundleVersion: version,
			})
		}
	}

	return apps, nil
}

// parseAppList parses the output of simctl listapps.
func parseAppList(output string) []App {
	var apps []App
	// Parse JSON output from simctl listapps if available
	var appMap map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(output), &appMap); err == nil {
		for bundleID, info := range appMap {
			app := App{BundleID: bundleID}
			if name, ok := info["CFBundleName"].(string); ok {
				app.Name = name
			}
			if path, ok := info["Path"].(string); ok {
				app.Path = path
			}
			if version, ok := info["CFBundleShortVersionString"].(string); ok {
				app.BundleVersion = version
			}
			apps = append(apps, app)
		}
	}
	return apps
}

// LaunchApp launches an app on the simulator.
func LaunchApp(udid, bundleID string) error {
	_, err := SimctlCommand("launch", udid, bundleID)
	return err
}

// TerminateApp terminates a running app on the simulator.
func TerminateApp(udid, bundleID string) error {
	_, err := SimctlCommand("terminate", udid, bundleID)
	return err
}

// SpawnProcess spawns a process on the simulator.
func SpawnProcess(udid string, args ...string) (string, error) {
	cmdArgs := append([]string{"spawn", udid}, args...)
	return SimctlCommand(cmdArgs...)
}

// GetRunningProcesses returns running processes on the simulator by querying the host.
func GetRunningProcesses(udid string) ([]Process, error) {
	// iOS Simulator processes run on macOS, so we query using ps
	// Filter for processes related to the simulator
	output, err := ExecuteCommand("ps", "-ax", "-o", "pid,comm")
	if err != nil {
		return nil, err
	}

	var processes []Process
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			processes = append(processes, Process{
				PID:  fields[0],
				Name: strings.Join(fields[1:], " "),
			})
		}
	}

	return processes, nil
}

// GetSimulatorPID returns the PID of the main Simulator.app process.
func GetSimulatorPID() (string, error) {
	output, err := ExecuteCommand("pgrep", "-f", "Simulator.app")
	if err != nil {
		return "", fmt.Errorf("Simulator.app not running")
	}
	pids := strings.Split(strings.TrimSpace(output), "\n")
	if len(pids) > 0 {
		return pids[0], nil
	}
	return "", fmt.Errorf("Simulator.app not running")
}

// BootSimulator boots a simulator by UDID.
func BootSimulator(udid string) error {
	_, err := SimctlCommand("boot", udid)
	return err
}

// ShutdownSimulator shuts down a simulator by UDID.
func ShutdownSimulator(udid string) error {
	_, err := SimctlCommand("shutdown", udid)
	return err
}

// GetArchitecture returns the architecture of the simulator (arm64 for Apple Silicon, x86_64 for Intel).
func GetArchitecture() (string, error) {
	output, err := ExecuteCommand("uname", "-m")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// OpenURL opens a URL in the simulator's default browser or app.
func OpenURL(udid, url string) error {
	_, err := SimctlCommand("openurl", udid, url)
	return err
}

// AddMedia adds photos/videos to the simulator.
func AddMedia(udid string, paths ...string) error {
	args := append([]string{"addmedia", udid}, paths...)
	_, err := SimctlCommand(args...)
	return err
}

// GetAppDataPath returns the data container path for an app.
func GetAppDataPath(udid, bundleID string) (string, error) {
	output, err := SimctlCommand("get_app_container", udid, bundleID, "data")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// InstallApp installs an app on the simulator.
func InstallApp(udid, appPath string) error {
	_, err := SimctlCommand("install", udid, appPath)
	return err
}

// UninstallApp uninstalls an app from the simulator.
func UninstallApp(udid, bundleID string) error {
	_, err := SimctlCommand("uninstall", udid, bundleID)
	return err
}

// Package frida locates Frida host tooling — the CLI binaries (frida, frida-ps,
// frida-trace) installed via `pip install frida-tools` — across the common
// install locations on macOS and Linux, with a $PATH fallback.
//
// This is host-side tooling shared by the Android and iOS toolkits. Android
// frida-server provisioning (download/extract/push) is device-specific and
// lives in the adb toolkit.
package frida

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// pythonVersions are the CPython minor versions whose pip user-install bin
// directories we probe for Frida tools.
var pythonVersions = []string{"3.9", "3.10", "3.11", "3.12", "3.13"}

// candidatePaths returns well-known absolute locations for a Frida CLI tool,
// covering Homebrew (Intel + Apple Silicon) and pip user installs.
func candidatePaths(tool string) []string {
	paths := []string{
		"/usr/local/bin/" + tool,
		"/opt/homebrew/bin/" + tool,
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".local", "bin", tool))
		for _, v := range pythonVersions {
			paths = append(paths, filepath.Join(home, "Library", "Python", v, "bin", tool))
		}
	}
	return paths
}

// Locate returns the absolute path to a Frida CLI tool (e.g. "frida",
// "frida-ps", "frida-trace"), or "" if it cannot be found. It checks the
// well-known install locations first, then falls back to $PATH.
func Locate(tool string) string {
	for _, p := range candidatePaths(tool) {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	name := tool
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if p, err := exec.LookPath(name); err == nil {
		return p
	}
	return ""
}

// Installed reports whether the Frida host tools appear to be installed.
func Installed() bool {
	return Locate("frida-ps") != ""
}

// Version returns the installed Frida version (from `frida --version`).
func Version() (string, error) {
	frida := Locate("frida")
	if frida == "" {
		return "", fmt.Errorf("frida not found")
	}
	out, err := exec.Command(frida, "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

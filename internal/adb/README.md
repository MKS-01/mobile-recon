# ADB Toolkit

A comprehensive command-line toolkit for Android Debug Bridge (ADB) operations. Perfect for Android development, debugging, and reverse engineering tasks.

## Features

- 🔧 **Device Management**: List devices, get info, reboot, screenshots, screen recording
- 📱 **App Management**: Install, uninstall, list, start/stop apps, pull APKs
- 🔍 **Reconnaissance**: Monitor logs, dump packages, analyze components, inspect databases
- 🎯 **Input Simulation**: Send text, tap, swipe, keyboard events
- 🛠️ **Frida Integration**: Helper commands for dynamic instrumentation
- 🎨 **Beautiful CLI**: Colorful output with progress indicators

## Installation

### Prerequisites

- Go 1.21 or higher
- ADB installed and in PATH
- Connected Android device or emulator

### Install

This toolkit ships as part of the unified `mobile-recon` CLI. From the repo root:

```bash
./scripts/install.sh        # builds and installs mobile-recon to ~/go/bin
```

The ADB commands are then available under `mobile-recon adb`.

## Quick Start

```bash
# List connected devices
mobile-recon adb device list

# Get device info
mobile-recon adb device info

# List installed apps (third-party only)
mobile-recon adb app list -3

# Install an APK
mobile-recon adb app install myapp.apk

# Pull an APK from device
mobile-recon adb app pull com.example.app

# Take a screenshot
mobile-recon adb device screenshot screenshot.png

# Monitor logcat
mobile-recon adb recon logcat
```

## Command Reference

### Device Commands

```bash
# List all connected devices
mobile-recon adb device list

# Get detailed device information
mobile-recon adb device info

# Reboot device
mobile-recon adb device reboot

# Reboot to recovery/bootloader
mobile-recon adb device reboot recovery
mobile-recon adb device reboot bootloader

# Execute shell command
mobile-recon adb device shell "ls -la /sdcard"

# Take screenshot
mobile-recon adb device screenshot [filename.png]

# Record screen
mobile-recon adb device screenrecord [filename.mp4]
```

### App Management

```bash
# List packages
mobile-recon adb app list              # All packages
mobile-recon adb app list -3           # Third-party only
mobile-recon adb app list -s           # System only
mobile-recon adb app list chrome       # Filter by name

# Install/Uninstall
mobile-recon adb app install app.apk
mobile-recon adb app install -r app.apk    # Reinstall keeping data
mobile-recon adb app uninstall com.example.app
mobile-recon adb app uninstall -k com.example.app  # Keep data

# App Operations
mobile-recon adb app clear com.example.app     # Clear data
mobile-recon adb app info com.example.app      # Package info
mobile-recon adb app start com.example.app     # Launch app
mobile-recon adb app stop com.example.app      # Force stop
mobile-recon adb app pull com.example.app      # Pull APK
```

### Reconnaissance & Reverse Engineering

```bash
# Logcat
mobile-recon adb recon logcat                    # Monitor logs
mobile-recon adb recon logcat "ActivityManager"  # Filter
mobile-recon adb recon logcat -s logcat.txt      # Save to file
mobile-recon adb recon logcat -c                 # Clear buffer

# Package Analysis
mobile-recon adb recon dump com.example.app              # Dump info
mobile-recon adb recon dump com.example.app -s dump.txt  # Save
mobile-recon adb recon activities com.example.app        # List activities
mobile-recon adb recon services com.example.app          # List services
mobile-recon adb recon receivers com.example.app         # List receivers

# Data Extraction (requires root)
mobile-recon adb recon files com.example.app             # List files
mobile-recon adb recon db com.example.app                # List databases
mobile-recon adb recon db com.example.app user.db        # Pull database

# System Monitoring
mobile-recon adb recon network                   # Network connections
mobile-recon adb recon processes                 # Running processes
mobile-recon adb recon processes com.example     # Filter processes
```

### Input Simulation

```bash
# Text & Touch
mobile-recon adb input text "Hello World"
mobile-recon adb input tap 500 1000
mobile-recon adb input swipe 500 1500 500 500 300

# Keycodes
mobile-recon adb input key 66        # Enter
mobile-recon adb input home          # Home button
mobile-recon adb input back          # Back button
```

### Frida Integration

```bash
mobile-recon adb frida ps                    # List processes
mobile-recon adb frida server check          # Check server status
mobile-recon adb frida server start          # Start server
mobile-recon adb frida trace com.example.app # Get trace command
```

## Working with Multiple Devices

```bash
# Target specific device
mobile-recon adb -d emulator-5554 device info
mobile-recon adb -d emulator-5554 app list
```

## Common Keycodes

| Code | Action    |
|------|-----------|
| 3    | Home      |
| 4    | Back      |
| 26   | Power     |
| 66   | Enter     |
| 67   | Backspace |
| 82   | Menu      |

## Reverse Engineering Workflow

### 1. Reconnaissance
```bash
mobile-recon adb app info com.target.app
mobile-recon adb recon dump com.target.app -s app_dump.txt
mobile-recon adb recon activities com.target.app
```

### 2. Extract APK
```bash
mobile-recon adb app pull com.target.app target.apk
```

### 3. Monitor Runtime
```bash
mobile-recon adb recon logcat com.target.app -s app_logs.txt
mobile-recon adb recon network
```

### 4. Extract Data (root required)
```bash
mobile-recon adb recon files com.target.app
mobile-recon adb recon db com.target.app user.db
```

### 5. Dynamic Analysis
```bash
mobile-recon adb frida server start
mobile-recon adb frida ps
```

## Project Structure

```
internal/adb/
├── adb.go              # ADB wrapper (package adb)
└── cmd/                # Command implementations (package cmd)
    ├── root.go         # Root command
    ├── device.go       # Device management
    ├── app.go          # App management
    ├── recon.go        # Reconnaissance
    ├── input.go        # Input simulation
    └── frida.go        # Frida integration
```

## Extending the CLI

Add custom commands by creating a new file in `internal/adb/cmd/`:

```go
package cmd

import (
    "github.com/MKS-01/mobile-recon/internal/adb"
    "github.com/MKS-01/mobile-recon/pkg/output"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycmd",
    Short: "My custom command",
    Run: func(cmd *cobra.Command, args []string) {
        serial, _ := getTargetDevice()
        out, _ := adb.ExecuteCommandWithDevice(serial, "shell", "echo", "hello")
        output.Success(out)
    },
}

func init() {
    RootCmd.AddCommand(myCmd)
}
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [color](https://github.com/fatih/color) - Colored output

## License

MIT

---

**Part of the [mobile-recon](../../) toolkit**

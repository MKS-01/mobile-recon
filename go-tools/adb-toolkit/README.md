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

### Build from source

```bash
# Navigate to the adb-toolkit directory
cd go-tools/adb-toolkit

# Install dependencies
go mod download

# Build the binary
go build -o adb-toolkit

# Optional: Install globally
go install
```

## Quick Start

```bash
# List connected devices
./adb-toolkit device list

# Get device info
./adb-toolkit device info

# List installed apps (third-party only)
./adb-toolkit app list -3

# Install an APK
./adb-toolkit app install myapp.apk

# Pull an APK from device
./adb-toolkit app pull com.example.app

# Take a screenshot
./adb-toolkit device screenshot screenshot.png

# Monitor logcat
./adb-toolkit recon logcat
```

## Command Reference

### Device Commands

```bash
# List all connected devices
adb-toolkit device list

# Get detailed device information
adb-toolkit device info

# Reboot device
adb-toolkit device reboot

# Reboot to recovery/bootloader
adb-toolkit device reboot recovery
adb-toolkit device reboot bootloader

# Execute shell command
adb-toolkit device shell "ls -la /sdcard"

# Take screenshot
adb-toolkit device screenshot [filename.png]

# Record screen
adb-toolkit device screenrecord [filename.mp4]
```

### App Management

```bash
# List packages
adb-toolkit app list              # All packages
adb-toolkit app list -3           # Third-party only
adb-toolkit app list -s           # System only
adb-toolkit app list chrome       # Filter by name

# Install/Uninstall
adb-toolkit app install app.apk
adb-toolkit app install -r app.apk    # Reinstall keeping data
adb-toolkit app uninstall com.example.app
adb-toolkit app uninstall -k com.example.app  # Keep data

# App Operations
adb-toolkit app clear com.example.app     # Clear data
adb-toolkit app info com.example.app      # Package info
adb-toolkit app start com.example.app     # Launch app
adb-toolkit app stop com.example.app      # Force stop
adb-toolkit app pull com.example.app      # Pull APK
```

### Reconnaissance & Reverse Engineering

```bash
# Logcat
adb-toolkit recon logcat                    # Monitor logs
adb-toolkit recon logcat "ActivityManager"  # Filter
adb-toolkit recon logcat -s logcat.txt      # Save to file
adb-toolkit recon logcat -c                 # Clear buffer

# Package Analysis
adb-toolkit recon dump com.example.app              # Dump info
adb-toolkit recon dump com.example.app -s dump.txt  # Save
adb-toolkit recon activities com.example.app        # List activities
adb-toolkit recon services com.example.app          # List services
adb-toolkit recon receivers com.example.app         # List receivers

# Data Extraction (requires root)
adb-toolkit recon files com.example.app             # List files
adb-toolkit recon db com.example.app                # List databases
adb-toolkit recon db com.example.app user.db        # Pull database

# System Monitoring
adb-toolkit recon network                   # Network connections
adb-toolkit recon processes                 # Running processes
adb-toolkit recon processes com.example     # Filter processes
```

### Input Simulation

```bash
# Text & Touch
adb-toolkit input text "Hello World"
adb-toolkit input tap 500 1000
adb-toolkit input swipe 500 1500 500 500 300

# Keycodes
adb-toolkit input key 66        # Enter
adb-toolkit input home          # Home button
adb-toolkit input back          # Back button
```

### Frida Integration

```bash
adb-toolkit frida ps                    # List processes
adb-toolkit frida server check          # Check server status
adb-toolkit frida server start          # Start server
adb-toolkit frida trace com.example.app # Get trace command
```

## Working with Multiple Devices

```bash
# Target specific device
adb-toolkit -d emulator-5554 device info
adb-toolkit -d emulator-5554 app list
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
adb-toolkit app info com.target.app
adb-toolkit recon dump com.target.app -s app_dump.txt
adb-toolkit recon activities com.target.app
```

### 2. Extract APK
```bash
adb-toolkit app pull com.target.app target.apk
```

### 3. Monitor Runtime
```bash
adb-toolkit recon logcat com.target.app -s app_logs.txt
adb-toolkit recon network
```

### 4. Extract Data (root required)
```bash
adb-toolkit recon files com.target.app
adb-toolkit recon db com.target.app user.db
```

### 5. Dynamic Analysis
```bash
adb-toolkit frida server start
adb-toolkit frida ps
```

## Project Structure

```
adb-toolkit/
├── main.go              # Entry point
├── cmd/                 # Command implementations
│   ├── root.go         # Root command
│   ├── device.go       # Device management
│   ├── app.go          # App management
│   ├── recon.go        # Reconnaissance
│   ├── input.go        # Input simulation
│   └── frida.go        # Frida integration
└── pkg/                # Shared packages
    ├── adb/            # ADB wrapper
    └── utils/          # Utilities
```

## Extending the CLI

Add custom commands by creating a new file in `cmd/`:

```go
package cmd

import (
    "github.com/mks/adb-toolkit/pkg/adb"
    "github.com/mks/adb-toolkit/pkg/utils"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycmd",
    Short: "My custom command",
    Run: func(cmd *cobra.Command, args []string) {
        serial, _ := getTargetDevice()
        output, _ := adb.ExecuteCommandWithDevice(serial, "shell", "echo", "hello")
        utils.PrintSuccess(output)
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [color](https://github.com/fatih/color) - Colored output

## License

MIT

---

**Part of the [automate-scripts](../../) monorepo**

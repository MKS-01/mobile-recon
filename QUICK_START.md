# Mobile Recon - Quick Start Guide

Get up and running with Mobile Recon in minutes!

## Installation

### Quick Install (Recommended)

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Run the automated installation script
./scripts/install.sh

# Reload your shell
source ~/.zshrc   # for zsh
source ~/.bashrc  # for bash
```

### Manual Build

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Build the unified CLI
cd go-tools/mobile-recon-cli
go build -o mobile-recon

# Build all tools
./mobile-recon build --all
```

## Basic Usage

### List Available Tools

```bash
./mobile-recon list
```

### Run Tools

```bash
# Using the unified CLI
./mobile-recon run adb-toolkit device list
./mobile-recon run nmap-toolkit scan quick 192.168.1.0/24

# Using shortcuts
./mobile-recon adb device list
./mobile-recon nmap scan port 192.168.1.1
```

### Interactive Mode

```bash
./mobile-recon interactive
```

## Tool-Specific Quick Start

### ADB Toolkit (Android Testing)

```bash
# List connected devices
mobile-recon adb device list

# Get device info
mobile-recon adb device info

# List installed apps
mobile-recon adb app list

# Take a screenshot
mobile-recon adb device screenshot

# Analyze an app
mobile-recon adb recon packages com.example.app
```

### Nmap Toolkit (Network Scanning)

```bash
# Quick network scan
mobile-recon nmap scan quick 192.168.1.0/24

# Port scan
mobile-recon nmap scan port 192.168.1.1

# Service detection
mobile-recon nmap detect service 192.168.1.1

# Find Android ADB devices on network
mobile-recon nmap mobile adb 192.168.1.0/24

# Vulnerability scan
mobile-recon nmap vuln scan 192.168.1.1
```

## Common Workflows

### Mobile App Security Testing

```bash
# 1. Find devices on network
mobile-recon nmap mobile adb 192.168.1.0/24

# 2. Connect to device
adb connect 192.168.1.100:5555

# 3. Verify connection
mobile-recon adb device list

# 4. List apps
mobile-recon adb app list

# 5. Analyze target app
mobile-recon adb recon packages com.target.app

# 6. Extract APK
mobile-recon adb app pull com.target.app
```

### Network Reconnaissance

```bash
# 1. Discover hosts
mobile-recon nmap scan network 192.168.1.0/24

# 2. Port scan target
mobile-recon nmap scan port 192.168.1.1 -p 1-1000

# 3. Identify services
mobile-recon nmap detect service 192.168.1.1

# 4. Check for vulnerabilities
mobile-recon nmap vuln scan 192.168.1.1
```

## Global Installation

The installation script (`./scripts/install.sh`) automatically installs all tools globally.

If you want to install manually:

```bash
# Install unified CLI
cd go-tools/mobile-recon-cli
go install

# Install individual tools
cd ../adb-toolkit
go install

cd ../nmap-toolkit
go install

# Now use from anywhere
mobile-recon-cli list
adb-toolkit device list
nmap-toolkit scan quick 192.168.1.0/24
```

**Note**: Make sure `$HOME/go/bin` is in your PATH:
```bash
export PATH="$HOME/go/bin:$PATH"
```

## Useful Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias mr='mobile-recon'
alias mradb='mobile-recon adb'
alias mrnmap='mobile-recon nmap'
alias mri='mobile-recon interactive'
```

Then use:

```bash
mr list
mradb device list
mrnmap scan quick 192.168.1.0/24
mri
```

## Prerequisites

Make sure you have these installed:

- **Go** (1.21+): https://golang.org/dl/
- **ADB** (for Android testing): https://developer.android.com/studio/command-line/adb
- **Nmap** (for network scanning): https://nmap.org/download.html

## Need Help?

```bash
# General help
mobile-recon --help

# Tool-specific help
mobile-recon adb --help
mobile-recon nmap --help

# Command-specific help
mobile-recon adb device --help
mobile-recon nmap scan --help
```

## Documentation

- [Full Documentation](README.md)
- [ADB Toolkit Docs](go-tools/adb-toolkit/README.md)
- [Nmap Toolkit Docs](go-tools/nmap-toolkit/README.md)
- [Unified CLI Docs](go-tools/mobile-recon-cli/README.md)

## Tips

1. **Use Interactive Mode** when learning - it's the easiest way to explore features
2. **Use --stream flag** for real-time output on long scans
3. **Check tool availability** with `mobile-recon list --all`
4. **Run with sudo** when needed for privileged operations (OS detection, stealth scans)

## Common Issues

### "Tool not built"
```bash
mobile-recon build <tool-name>
```

### "Nmap not installed"
```bash
# macOS
brew install nmap

# Linux
sudo apt install nmap
```

### "ADB not installed"
```bash
# macOS
brew install android-platform-tools

# Linux
sudo apt install adb
```

Happy Hacking! 🚀

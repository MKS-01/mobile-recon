# Mobile Recon Toolkit

A comprehensive, unified CLI toolkit for mobile security testing, network reconnaissance, and penetration testing.

## Features

- **Unified CLI** - All tools accessible via single `mobile-recon` command
- **ADB Toolkit** - Android device automation and reverse engineering
- **Nmap Toolkit** - Network scanning and reconnaissance
- **APK Analyzer** - Android APK static analysis
- **iOS Toolkit** - iOS Simulator management and Frida integration

## Quick Start

### Installation (Recommended)

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Run the install script
./scripts/install.sh
```

### Manual Installation

```bash
# Clone and build
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon/go-tools/mobile-recon-cli

# Build the unified binary
go build -o mobile-recon

# Move to your PATH
mv mobile-recon ~/go/bin/
# or
sudo mv mobile-recon /usr/local/bin/
```

### PATH Setup

Add Go bin directory to your PATH. Add this to your `~/.zshrc` or `~/.bashrc`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell:

```bash
source ~/.zshrc   # for zsh
source ~/.bashrc  # for bash
```

## Usage

### Command Structure

```
mobile-recon <tool> <command> [args...]
```

### Available Tools

| Tool | Command | Description |
|------|---------|-------------|
| ADB Toolkit | `mobile-recon adb` | Android device automation |
| Nmap Toolkit | `mobile-recon nmap` | Network reconnaissance |
| APK Analyzer | `mobile-recon apk` | APK static analysis |
| iOS Toolkit | `mobile-recon ios` | iOS Simulator management |

### Examples

```bash
# List all available tools
mobile-recon list

# ADB Toolkit
mobile-recon adb device list
mobile-recon adb app list
mobile-recon adb app pull com.example.app
mobile-recon adb recon logcat

# Nmap Toolkit
mobile-recon nmap scan quick 192.168.1.0/24
mobile-recon nmap scan port 192.168.1.1
mobile-recon nmap detect service 192.168.1.1
mobile-recon nmap vuln scan 192.168.1.1

# APK Analyzer
mobile-recon apk info app.apk
mobile-recon apk manifest app.apk
mobile-recon apk permissions app.apk
mobile-recon apk abuse-permissions app.apk
mobile-recon apk security app.apk
mobile-recon apk strings app.apk

# iOS Toolkit
mobile-recon ios device list
mobile-recon ios device boot <udid>
mobile-recon ios frida ps
```

## Project Structure

```
mobile-recon/
├── go-tools/
│   ├── mobile-recon-cli/    # Unified CLI (main binary)
│   ├── adb-toolkit/         # Android ADB automation
│   ├── nmap-toolkit/        # Network reconnaissance
│   ├── apk-analyzer/        # APK static analysis
│   ├── ios-toolkit/         # iOS Simulator toolkit
│   └── common/              # Shared utilities
├── scripts/
│   └── install.sh           # Installation script
└── README.md
```

## Tools Documentation

### ADB Toolkit

Android Debug Bridge automation toolkit.

```bash
mobile-recon adb --help

# Device Management
mobile-recon adb device list              # List connected devices
mobile-recon adb device info              # Get device information
mobile-recon adb device screenshot        # Take screenshot
mobile-recon adb device screenrecord      # Record screen

# App Management
mobile-recon adb app list                 # List installed apps
mobile-recon adb app install app.apk      # Install APK
mobile-recon adb app uninstall com.app    # Uninstall app
mobile-recon adb app pull com.app         # Extract APK

# Reconnaissance
mobile-recon adb recon logcat             # View logs
mobile-recon adb recon packages com.app   # Analyze package

# Input Automation
mobile-recon adb input tap 500 1000       # Tap screen
mobile-recon adb input text "hello"       # Type text
mobile-recon adb input swipe 100 500 100 100

# Frida Integration
mobile-recon adb frida setup              # Install Frida server
mobile-recon adb frida ps                 # List processes
```

### Nmap Toolkit

Network reconnaissance toolkit.

```bash
mobile-recon nmap --help

# Scanning
mobile-recon nmap scan quick 192.168.1.0/24     # Quick host discovery
mobile-recon nmap scan port 192.168.1.1         # Port scan
mobile-recon nmap scan stealth 192.168.1.1      # Stealth scan (sudo)
mobile-recon nmap scan udp 192.168.1.1          # UDP scan

# Detection
mobile-recon nmap detect service 192.168.1.1    # Service detection
mobile-recon nmap detect os 192.168.1.1         # OS detection (sudo)
mobile-recon nmap detect aggressive 192.168.1.1 # Full detection

# Vulnerability Scanning
mobile-recon nmap vuln scan 192.168.1.1         # Vulnerability scan
mobile-recon nmap vuln ssl 192.168.1.1          # SSL/TLS testing

# Mobile-Specific
mobile-recon nmap mobile adb 192.168.1.0/24     # Find ADB devices
mobile-recon nmap mobile ios 192.168.1.0/24     # Find iOS devices
mobile-recon nmap mobile scan 192.168.1.1       # Mobile port scan
```

### APK Analyzer

Android APK static analysis toolkit.

```bash
mobile-recon apk --help

mobile-recon apk info app.apk                 # APK metadata
mobile-recon apk manifest app.apk             # Parse AndroidManifest.xml
mobile-recon apk permissions app.apk          # List permissions
mobile-recon apk permissions -d app.apk       # Dangerous permissions only
mobile-recon apk abuse-permissions app.apk    # Detect abusive permissions
mobile-recon apk abuse-permissions -m app.apk # Malware permissions only
mobile-recon apk security app.apk             # Security analysis
mobile-recon apk strings app.apk              # Extract strings
mobile-recon apk strings --urls app.apk       # Extract URLs only
mobile-recon apk files app.apk                # List files
mobile-recon apk files --tree app.apk         # Show file tree
```

### iOS Toolkit

iOS Simulator management and Frida integration.

```bash
mobile-recon ios --help

# Device Management
mobile-recon ios device list              # List simulators
mobile-recon ios device list -a           # Include unavailable
mobile-recon ios device boot <udid>       # Boot simulator
mobile-recon ios device shutdown <udid>   # Shutdown simulator
mobile-recon ios device info              # Current device info

# Frida Integration
mobile-recon ios frida setup              # Install Frida
mobile-recon ios frida ps                 # List processes
mobile-recon ios frida apps               # List apps
mobile-recon ios frida attach <pid>       # Attach to process
mobile-recon ios frida spawn <bundle>     # Spawn app
```

## Development

### Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **ADB** - For ADB Toolkit ([Install Guide](https://developer.android.com/studio/command-line/adb))
- **Nmap** - For Nmap Toolkit ([Download](https://nmap.org/download.html))
- **Xcode** - For iOS Toolkit (macOS only)

### Building from Source

```bash
# Clone repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon/go-tools/mobile-recon-cli

# Download dependencies
go mod tidy

# Build
go build -o mobile-recon

# Install to PATH
mv mobile-recon ~/go/bin/
```

### After Making Changes

When you modify the code or add new features:

```bash
# Navigate to the CLI directory
cd go-tools/mobile-recon-cli

# Clean and rebuild
go mod tidy
go build -o mobile-recon

# Reinstall
mv mobile-recon ~/go/bin/

# Or use the install script
cd ../..
./scripts/install.sh
```

### Adding New Features to a Toolkit

1. **Edit the relevant toolkit** in `go-tools/<toolkit>/cmd/`
2. **Rebuild the main binary**:
   ```bash
   cd go-tools/mobile-recon-cli
   go build -o mobile-recon
   mv mobile-recon ~/go/bin/
   ```

### Creating a New Tool

See [NEW_TOOL_SETUP.md](.templates/NEW_TOOL_SETUP.md) for instructions on adding a new toolkit.

## Configuration

### Helpful Aliases

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# Mobile Recon aliases
alias mr='mobile-recon'
alias mradb='mobile-recon adb'
alias mrnmap='mobile-recon nmap'
alias mrapk='mobile-recon apk'
alias mrios='mobile-recon ios'
```

## Troubleshooting

### Command not found: mobile-recon

Ensure `~/go/bin` is in your PATH:

```bash
# Check if in PATH
echo $PATH | grep -o 'go/bin'

# Add to PATH (add to ~/.zshrc or ~/.bashrc)
export PATH="$HOME/go/bin:$PATH"

# Reload shell
source ~/.zshrc
```

### Old alias overriding new binary

If you have an old alias like `alias mobile-recon='mobile-recon-cli'`:

```bash
# Remove the alias from current session
unalias mobile-recon

# Remove from config file
sed -i '' "/alias mobile-recon='mobile-recon-cli'/d" ~/.zshrc
```

### ADB not installed

```bash
# macOS
brew install android-platform-tools

# Ubuntu/Debian
sudo apt install adb

# Verify
adb version
```

### Nmap not installed

```bash
# macOS
brew install nmap

# Ubuntu/Debian
sudo apt install nmap

# Verify
nmap --version
```

## Security Notice

This toolkit is designed for:
- Authorized security testing
- Mobile app development and debugging
- Security research and education
- Network administration

**Always obtain proper authorization** before scanning networks or testing applications you don't own.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Rebuild and test (`go build -o mobile-recon && ./mobile-recon --help`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Nmap](https://nmap.org) - Network mapper
- [ADB](https://developer.android.com/studio/command-line/adb) - Android Debug Bridge
- [Color](https://github.com/fatih/color) - Terminal coloring

---

**Built for security researchers, penetration testers, and mobile developers**

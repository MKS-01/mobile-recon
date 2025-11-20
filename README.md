# Mobile Recon Toolkit

A comprehensive collection of reconnaissance tools and automation scripts for mobile security testing, network analysis, and penetration testing.

## 📁 Project Structure

```
mobile-recon/
├── go-tools/                    # Go-based CLI tools
│   ├── mobile-recon-cli/        # Unified CLI manager
│   ├── adb-toolkit/             # Android ADB automation toolkit
│   └── nmap-toolkit/            # Network reconnaissance toolkit
├── bash-scripts/                # Bash automation scripts
│   └── github-repo.sh           # GitHub repository automation
└── README.md
```

## 🎯 Quick Start

### Using the Unified CLI (Recommended)

The easiest way to use all tools is through the unified CLI manager:

```bash
# 1. Build the unified CLI
cd go-tools/mobile-recon-cli
go build -o mobile-recon

# 2. Build all tools
./mobile-recon build --all

# 3. List available tools
./mobile-recon list

# 4. Run any tool
./mobile-recon run adb-toolkit device list
./mobile-recon run nmap-toolkit scan quick 192.168.1.0/24

# Or use shortcuts
./mobile-recon adb device list
./mobile-recon nmap scan port 192.168.1.1

# Interactive mode
./mobile-recon interactive
```

### Install Globally

```bash
cd go-tools/mobile-recon-cli
go install

# Now use from anywhere
mobile-recon list
mobile-recon adb device list
mobile-recon nmap scan quick 192.168.1.0/24
```

## 🛠️ Tools & Scripts

### Unified CLI Manager

**[Mobile Recon CLI](go-tools/mobile-recon-cli/)** - Centralized tool management interface

**Features:**
- Discover and list all available tools
- Build and install tools
- Run tools with unified interface
- Interactive mode for easy navigation
- Category-based tool browsing
- Shortcuts for common operations

**Quick Start:**
```bash
cd go-tools/mobile-recon-cli
go build -o mobile-recon
./mobile-recon interactive
```

See [detailed documentation](go-tools/mobile-recon-cli/README.md)

---

### Mobile Tools

#### [ADB Toolkit](go-tools/adb-toolkit/)

Comprehensive Android Debug Bridge (ADB) CLI toolkit for Android development and reverse engineering.

**Features:**
- **Device Management**: List, info, reboot, screenshots, screen recording
- **App Management**: Install, uninstall, pull APKs, list packages
- **Reconnaissance**: Logcat analysis, package inspection, database extraction
- **Input Simulation**: Tap, swipe, text input automation
- **Frida Integration**: Helper commands for dynamic instrumentation

**Quick Start:**
```bash
# Via unified CLI
mobile-recon run adb-toolkit device list

# Direct usage
cd go-tools/adb-toolkit
go build -o adb-toolkit
./adb-toolkit device list
```

**Common Commands:**
```bash
adb-toolkit device list              # List connected devices
adb-toolkit device info               # Get device information
adb-toolkit app list                  # List installed apps
adb-toolkit recon packages com.app    # Analyze app package
adb-toolkit frida list-processes      # List running processes
```

See [detailed documentation](go-tools/adb-toolkit/README.md)

---

### Network Tools

#### [Nmap Toolkit](go-tools/nmap-toolkit/)

Advanced network reconnaissance toolkit with mobile-specific scanning capabilities.

**Features:**
- **Network Scanning**: Host discovery, port scanning (TCP/UDP/Stealth)
- **Service Detection**: Version detection, OS fingerprinting
- **Vulnerability Scanning**: Automated vuln detection, SSL/TLS testing
- **Mobile-Specific**: Android ADB discovery, iOS device detection, MITM proxy detection
- **NSE Scripts**: Custom script execution for advanced testing

**Quick Start:**
```bash
# Via unified CLI
mobile-recon run nmap-toolkit scan quick 192.168.1.0/24

# Direct usage
cd go-tools/nmap-toolkit
go build -o nmap-toolkit
./nmap-toolkit scan quick 192.168.1.0/24
```

**Common Commands:**
```bash
nmap-toolkit scan quick 192.168.1.0/24        # Quick host discovery
nmap-toolkit scan port 192.168.1.1            # Port scan
nmap-toolkit detect service 192.168.1.1       # Service detection
nmap-toolkit vuln scan 192.168.1.1            # Vulnerability scan
nmap-toolkit mobile adb 192.168.1.0/24        # Find ADB devices
nmap-toolkit mobile ios 192.168.1.0/24        # Find iOS devices
```

See [detailed documentation](go-tools/nmap-toolkit/README.md)

---

### Bash Scripts

#### [GitHub Repo Script](bash-scripts/github-repo.sh)

Automation script for GitHub repository management.

**Usage:**
```bash
bash bash-scripts/github-repo.sh
```

---

## 🚀 Installation & Setup

### Prerequisites

- **Go**: Version 1.21 or higher ([Download](https://golang.org/dl/))
- **Android ADB**: For ADB Toolkit ([Install Guide](https://developer.android.com/studio/command-line/adb))
- **Nmap**: For Nmap Toolkit ([Download](https://nmap.org/download.html))
- **Bash**: Version 4.0 or higher

### Installation Options

#### Option 1: Quick Install (Recommended)

Use the automated installation script to build and install all tools globally:

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Run the installation script
./scripts/install.sh

# Or with verbose output
./scripts/install.sh --verbose
```

The script will:
- Check prerequisites (Go installation)
- Build all Go CLI tools
- Install them to `$GOPATH/bin` or `$HOME/go/bin`
- Verify PATH configuration
- Optionally set up helpful aliases

After installation, reload your shell:
```bash
source ~/.zshrc   # for zsh
source ~/.bashrc  # for bash
```

#### Option 2: Manual Build (Unified CLI)

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Build unified CLI
cd go-tools/mobile-recon-cli
go build -o mobile-recon

# Build all tools
./mobile-recon build --all

# Install globally (optional)
go install
```

#### Option 3: Individual Tools

```bash
# Clone the repository
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon

# Build ADB Toolkit
cd go-tools/adb-toolkit
go mod download
go build -o adb-toolkit

# Build Nmap Toolkit
cd ../nmap-toolkit
go mod download
go build -o nmap-toolkit
```

---

## 📚 Documentation

Detailed documentation for each tool:

- **[Mobile Recon CLI](go-tools/mobile-recon-cli/README.md)** - Unified CLI manager
- **[ADB Toolkit](go-tools/adb-toolkit/README.md)** - Android automation
- **[Nmap Toolkit](go-tools/nmap-toolkit/README.md)** - Network reconnaissance

---

## 🎯 Use Cases

### Mobile Application Security Testing

```bash
# 1. Discover Android devices on network
mobile-recon nmap mobile adb 192.168.1.0/24

# 2. Connect to device
adb connect 192.168.1.100:5555

# 3. List installed apps
mobile-recon adb app list

# 4. Analyze specific app
mobile-recon adb recon packages com.example.app

# 5. Extract APK for reverse engineering
mobile-recon adb app pull com.example.app
```

### Network Security Audit

```bash
# 1. Discover live hosts
mobile-recon nmap scan network 192.168.1.0/24

# 2. Deep port scan
sudo mobile-recon nmap scan stealth 192.168.1.1

# 3. Service and OS detection
sudo mobile-recon nmap detect aggressive 192.168.1.1

# 4. Vulnerability assessment
mobile-recon nmap vuln scan 192.168.1.1

# 5. SSL/TLS testing
mobile-recon nmap vuln ssl 192.168.1.1
```

### Mobile Development & Testing

```bash
# Device management
mobile-recon adb device list
mobile-recon adb device info
mobile-recon adb device screenshot

# App installation and testing
mobile-recon adb app install app.apk
mobile-recon adb app launch com.example.app

# Log monitoring
mobile-recon adb recon logcat

# Input automation
mobile-recon adb input tap 500 1000
mobile-recon adb input swipe 100 500 100 100
```

### Penetration Testing Workflow

```bash
# 1. Reconnaissance
mobile-recon nmap scan quick 10.0.0.0/24

# 2. Port scanning
sudo mobile-recon nmap scan stealth 10.0.0.50

# 3. Service enumeration
mobile-recon nmap detect service 10.0.0.50

# 4. Vulnerability scanning
mobile-recon nmap vuln scan 10.0.0.50

# 5. Mobile-specific checks
mobile-recon nmap mobile scan 10.0.0.50
mobile-recon nmap mobile app-ports api.example.com
```

---

## 🔧 Development

### Adding New Tools

Use the automated tool generator to create new tools with proper structure:

```bash
# Generate a new tool
./scripts/new-tool.sh <tool-name> <category> <description>

# Example: Create a Frida toolkit
./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation toolkit"

# Example: Create a web testing tool
./scripts/new-tool.sh jwt-analyzer Web "JWT token analysis and manipulation"
```

The generator will:
- ✅ Create complete directory structure
- ✅ Generate boilerplate code with Cobra framework
- ✅ Initialize Go modules
- ✅ Register with mobile-recon-cli automatically
- ✅ Provide step-by-step implementation guide

**Tool Categories:**
- `Mobile` - Android, iOS, mobile app testing
- `Network` - Network scanning, reconnaissance
- `Web` - Web application testing, API testing
- `Other` - General-purpose tools

**Manual Setup:**

If you prefer manual setup, see [NEW_TOOL_SETUP.md](.templates/NEW_TOOL_SETUP.md) for complete instructions

### Project Structure Guidelines

- Each tool should be self-contained with its own README
- Follow consistent command structure using Cobra
- Use the shared color package for output formatting
- Include comprehensive help text and examples

---

## 📋 Project Philosophy

This toolkit follows these principles:

1. **Unified Experience**: Single CLI to access all tools
2. **Modularity**: Each tool is independent and focused
3. **Mobile-First**: Specialized features for mobile security testing
4. **Developer-Friendly**: Clear documentation and intuitive commands
5. **Open Source**: MIT licensed for community contribution

---

## 💡 Tips & Best Practices

### Global Installation

Install the unified CLI globally:
```bash
cd go-tools/mobile-recon-cli
go install

# Now use from anywhere
mobile-recon list
mobile-recon adb device list
```

### Useful Aliases

Add to your `.bashrc` or `.zshrc`:
```bash
alias mr='mobile-recon'
alias mradb='mobile-recon adb'
alias mrnmap='mobile-recon nmap'
alias mri='mobile-recon interactive'
```

### Workflow Integration

```bash
# Create a testing workflow
cat > mobile-test.sh << 'EOF'
#!/bin/bash
echo "Starting Mobile Security Test..."

# Discover devices
mobile-recon nmap mobile adb 192.168.1.0/24

# Check device
mobile-recon adb device list
mobile-recon adb device info

# Analyze app
mobile-recon adb app list | grep $1
mobile-recon adb recon packages $1
EOF

chmod +x mobile-test.sh
./mobile-test.sh com.example.app
```

---

## 🔐 Security Considerations

**Important:** Always obtain proper authorization before:
- Scanning networks or systems
- Testing mobile applications
- Performing penetration testing
- Analyzing devices you don't own

This toolkit is designed for:
- Authorized security testing
- Mobile app development and debugging
- Security research and education
- Network administration

**Legal Disclaimer:** Users are responsible for complying with all applicable laws and regulations. Unauthorized access to computer systems is illegal.

---

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contribution Guidelines

- Follow existing code structure and style
- Add tests for new features
- Update documentation
- Use conventional commit messages

---

## 🎯 Roadmap

### Planned Features

#### Go Tools
- [ ] iOS device automation toolkit
- [ ] Frida scripting toolkit
- [ ] Mobile API testing framework
- [ ] Certificate pinning bypass toolkit
- [ ] Mobile malware analysis tools

#### Network Tools
- [ ] Web application scanner
- [ ] API fuzzing toolkit
- [ ] Wireless network analysis
- [ ] Bluetooth reconnaissance

#### Integrations
- [ ] Burp Suite integration
- [ ] OWASP ZAP integration
- [ ] MobSF integration
- [ ] CI/CD pipeline tools

---

## 📖 Resources

### Learning Resources

- [OWASP Mobile Security Testing Guide](https://owasp.org/www-project-mobile-security-testing-guide/)
- [Android Security Documentation](https://source.android.com/security)
- [Nmap Documentation](https://nmap.org/docs.html)
- [Frida Documentation](https://frida.re/docs/home/)

### Related Projects

- [MobSF](https://github.com/MobSF/Mobile-Security-Framework-MobSF) - Mobile Security Framework
- [Objection](https://github.com/sensepost/objection) - Runtime mobile exploration
- [APKTool](https://github.com/iBotPeaches/Apktool) - APK reverse engineering

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🙏 Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Nmap](https://nmap.org) - Network mapper
- [ADB](https://developer.android.com/studio/command-line/adb) - Android Debug Bridge
- [Color](https://github.com/fatih/color) - Terminal coloring

---

## 📬 Contact & Support

- **GitHub Issues**: [Report bugs or request features](https://github.com/MKS-01/mobile-recon/issues)
- **Discussions**: [Ask questions and share ideas](https://github.com/MKS-01/mobile-recon/discussions)
- **Security Issues**: Please report security vulnerabilities privately

---

**Built for security researchers, penetration testers, and mobile developers** 🚀

*Happy Hacking! (Ethically, of course)* 🔒

# Automate Scripts

A collection of automation tools and scripts for day-to-day development tasks and security research.

## 📁 Project Structure

```
automate-scripts/
├── go-tools/           # Go-based CLI tools
│   └── adb-toolkit/    # Android ADB automation toolkit
├── bash-scripts/       # Bash automation scripts
│   └── github-repo.sh  # GitHub repository automation
├── docs/               # Documentation
└── README.md
```

## 🛠️ Tools & Scripts

### Go Tools

#### [ADB Toolkit](go-tools/adb-toolkit/)
Comprehensive Android Debug Bridge (ADB) CLI toolkit for development and reverse engineering.

**Features:**
- Device management (list, info, reboot, screenshots)
- App management (install, uninstall, pull APKs)
- Reconnaissance tools (logcat, package analysis, database extraction)
- Input simulation (tap, swipe, text input)
- Frida integration helpers

**Quick Start:**
```bash
cd go-tools/adb-toolkit
go build -o adb-toolkit
./adb-toolkit --help
```

See [detailed documentation](go-tools/adb-toolkit/README.md)

### Bash Scripts

#### [GitHub Repo Script](bash-scripts/github-repo.sh)
Automation script for GitHub repository management.

**Usage:**
```bash
bash bash-scripts/github-repo.sh
```

## 🚀 Getting Started

### Prerequisites

- **For Go tools**: Go 1.21 or higher
- **For Bash scripts**: Bash 4.0 or higher
- **For ADB toolkit**: Android Debug Bridge (ADB)

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd automate-scripts
```

2. Navigate to the specific tool you want to use:
```bash
# For Go tools
cd go-tools/adb-toolkit
go mod download
go build -o adb-toolkit

# For Bash scripts
cd bash-scripts
chmod +x *.sh
```

## 📚 Documentation

Each tool has its own detailed documentation:
- [ADB Toolkit Docs](go-tools/adb-toolkit/README.md)

## 🎯 Use Cases

### Mobile Development
- Automate APK installation and testing
- Monitor app logs and performance
- Capture screenshots and recordings
- Simulate user interactions

### Security Research & Reverse Engineering
- Extract APKs from devices
- Analyze app components (activities, services, receivers)
- Pull and inspect databases
- Monitor network traffic
- Dynamic instrumentation with Frida

### DevOps & Automation
- GitHub repository management
- Batch operations across multiple devices
- Automated testing workflows

## 🔧 Adding New Tools

### Go Tools
```bash
# Create new tool directory
mkdir -p go-tools/your-tool-name
cd go-tools/your-tool-name

# Initialize Go module
go mod init github.com/mks/your-tool-name

# Create main.go and build
go build -o your-tool-name
```

### Bash Scripts
```bash
# Create new script
touch bash-scripts/your-script.sh
chmod +x bash-scripts/your-script.sh

# Add shebang and your script
echo '#!/bin/bash' > bash-scripts/your-script.sh
```

## 📋 Project Philosophy

This monorepo follows these principles:

1. **Separation of Concerns**: Go tools and Bash scripts are organized separately
2. **Self-Contained**: Each tool has its own dependencies and documentation
3. **Extensible**: Easy to add new tools and scripts
4. **Practical**: Focus on real-world automation tasks

## 🤝 Contributing

Feel free to add new tools and scripts! Follow the existing structure:
- Add Go tools to `go-tools/`
- Add Bash scripts to `bash-scripts/`
- Include README for complex tools
- Update this main README with new additions

## 📝 License

MIT

## 🎯 Roadmap

### Planned Go Tools
- [ ] iOS device automation toolkit
- [ ] Git automation CLI
- [ ] Docker container management tool
- [ ] API testing framework

### Planned Bash Scripts
- [ ] System setup automation
- [ ] Database backup scripts
- [ ] Log analysis tools
- [ ] Deployment automation

## 💡 Tips

### Global Installation
Install Go tools globally for easy access:
```bash
cd go-tools/adb-toolkit
go install
# Now use from anywhere
adb-toolkit --help
```

### PATH Configuration
Add bash-scripts to your PATH:
```bash
export PATH="$PATH:/path/to/automate-scripts/bash-scripts"
```

### Aliases
Create useful aliases:
```bash
# In your .bashrc or .zshrc
alias adb-tk='~/automate-scripts/go-tools/adb-toolkit/adb-toolkit'
alias gh-repo='bash ~/automate-scripts/bash-scripts/github-repo.sh'
```

---

**Built for developers, by developers** 🚀

For questions or suggestions, open an issue or submit a PR!

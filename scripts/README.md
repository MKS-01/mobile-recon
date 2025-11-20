# Mobile Recon - Installation Scripts

This directory contains automated installation and setup scripts for the Mobile Recon Toolkit.

## Available Scripts

### install.sh

Automated installation script that builds and installs all Go CLI tools globally.

### test-installation.sh

Quick test script to verify that all tools are properly installed and accessible.

### new-tool.sh

Tool generator that creates a new reconnaissance tool from template with proper structure and integration.

**Usage (install.sh):**
```bash
./scripts/install.sh              # Standard installation
./scripts/install.sh --verbose    # Installation with detailed output
./scripts/install.sh --help       # Show help information
```

**Usage (test-installation.sh):**
```bash
./scripts/test-installation.sh    # Verify all tools are installed
```

**Usage (new-tool.sh):**
```bash
./scripts/new-tool.sh <tool-name> <category> <description>

# Examples
./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation toolkit"
./scripts/new-tool.sh jwt-analyzer Web "JWT token analysis and manipulation"
./scripts/new-tool.sh wifi-toolkit Network "WiFi security analysis toolkit"

# Categories: Mobile, Network, Web, Other
```

**What install.sh does:**
1. Checks prerequisites (Go installation, version 1.21+)
2. Builds all Go CLI tools:
   - `mobile-recon-cli` - Unified CLI manager
   - `adb-toolkit` - Android Debug Bridge toolkit
   - `nmap-toolkit` - Network reconnaissance toolkit
3. Installs them to `$GOPATH/bin` (typically `$HOME/go/bin`)
4. Verifies PATH configuration
5. Optionally sets up helpful shell aliases

**What new-tool.sh does:**
1. Validates tool name and category
2. Creates complete directory structure from template
3. Generates all boilerplate code (main.go, cmd/, pkg/)
4. Replaces template variables with your values
5. Initializes Go modules
6. Registers tool with mobile-recon-cli
7. Provides step-by-step implementation guide

See [.templates/NEW_TOOL_SETUP.md](../.templates/NEW_TOOL_SETUP.md) for complete setup and development guide.

**Requirements:**
- Go 1.21 or higher
- Unix-like environment (Linux, macOS, WSL)
- Bash shell

**After installation:**
```bash
# Reload your shell configuration
source ~/.zshrc   # for zsh
source ~/.bashrc  # for bash

# Test installation
./scripts/test-installation.sh

# Verify tools work
mobile-recon-cli list
adb-toolkit --help
nmap-toolkit --help
```

## PATH Configuration

The tools are installed to `$GOPATH/bin` or `$HOME/go/bin` by default. Make sure this directory is in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/go/bin:$PATH"
```

The installation script will warn you if this is not configured correctly.

## Suggested Aliases

For easier access, consider adding these aliases to your shell configuration:

```bash
# Unified CLI shortcuts
alias mobile-recon='mobile-recon-cli'
alias mr='mobile-recon-cli'
alias mradb='mobile-recon-cli adb'
alias mrnmap='mobile-recon-cli nmap'
alias mri='mobile-recon-cli interactive'

# Direct tool shortcuts
alias adb='adb-toolkit'
alias nmap='nmap-toolkit'
```

The installation script can add these automatically if you choose to accept the prompt.

## Troubleshooting

### "command not found" after installation

**Problem:** Tools are installed but not found in PATH.

**Solution:**
```bash
# 1. Check if tools are installed
ls -la ~/go/bin/

# 2. Verify PATH includes Go bin directory
echo $PATH | grep -o "$HOME/go/bin"

# 3. Add to PATH if missing
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### "Go is not installed"

**Problem:** Go is not installed or not in PATH.

**Solution:**
1. Download Go from https://golang.org/dl/
2. Follow installation instructions for your OS
3. Verify: `go version`

### Build failures

**Problem:** Build fails with module errors.

**Solution:**
```bash
# Clean Go module cache
go clean -modcache

# Try installation again
./scripts/install.sh --verbose
```

## Manual Installation

If you prefer to install tools manually:

```bash
# Install unified CLI
cd go-tools/mobile-recon-cli
go install

# Install ADB Toolkit
cd ../adb-toolkit
go install

# Install Nmap Toolkit
cd ../nmap-toolkit
go install
```

## Uninstallation

To remove installed tools:

```bash
rm -f ~/go/bin/mobile-recon-cli
rm -f ~/go/bin/adb-toolkit
rm -f ~/go/bin/nmap-toolkit
```

## Contributing

To add new installation scripts:

1. Create the script in this directory
2. Make it executable: `chmod +x scripts/your-script.sh`
3. Update this README with usage instructions
4. Test on multiple platforms (Linux, macOS, WSL)

## Support

For issues or questions:
- Check the [main README](../README.md)
- See [QUICK_START.md](../QUICK_START.md)
- Open an issue on GitHub

---

**Happy hacking! 🚀**

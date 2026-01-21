# Mobile Recon CLI

Unified command-line interface for managing and running all Mobile Recon tools from a single location.

## Overview

The Mobile Recon CLI provides a centralized way to:
- Discover all available reconnaissance tools
- Build and install tools
- Run tools with a unified interface
- Browse tools by category
- Interactive tool selection

## Features

- **Tool Discovery**: Automatically discovers all available tools in the go-tools directory
- **Unified Interface**: Run any tool from a single CLI
- **Build Management**: Build individual tools or all tools at once
- **Interactive Mode**: User-friendly interactive menu for tool selection
- **Category Browsing**: Browse tools by category (Mobile, Network, etc.)
- **Shortcuts**: Quick shortcuts for commonly used tools

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Build from Source (Development)

For local development and testing:

```bash
cd go-tools/mobile-recon-cli
go mod download
go build -o mobile-recon
```

This builds the binary in the current directory. Run with `./mobile-recon`.

### Install Globally

To install globally and use `mobile-recon` from anywhere, you **must** embed the source path using ldflags:

```bash
cd go-tools/mobile-recon-cli

# Replace <FULL_PATH_TO_GO_TOOLS> with your actual path
go install -ldflags "-s -w -X github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager.SourcePath=<FULL_PATH_TO_GO_TOOLS>"
```

Example:

```bash
go install -ldflags "-s -w -X github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager.SourcePath=/Users/mks/Desktop/C0D3/mobile-recon/go-tools"
```

**Why is this needed?** The CLI needs to know where the tool source directories are located to build them. Without the embedded path, it will look for tools relative to the binary location (e.g., `/Users/mks/go/`) which is incorrect.

### Quick Install Script

You can also use this one-liner (run from the mobile-recon-cli directory):

```bash
go install -ldflags "-s -w -X github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager.SourcePath=$(cd .. && pwd)"
```

## Usage

### List Available Tools

```bash
# List all built tools
mobile-recon list

# List all tools including those not built
mobile-recon list --all
```

### Build Tools

```bash
# Build a specific tool
mobile-recon build adb-toolkit
mobile-recon build nmap-toolkit

# Build all tools
mobile-recon build --all
```

### Run Tools

```bash
# Run a tool with arguments
mobile-recon run adb-toolkit device list
mobile-recon run nmap-toolkit scan quick 192.168.1.0/24

# Show tool help
mobile-recon run adb-toolkit --help
mobile-recon run nmap-toolkit --help
```

### Shortcuts

```bash
# Use shortcuts for quick access
mobile-recon adb device list
mobile-recon nmap scan port 192.168.1.1
```

### Interactive Mode

```bash
# Launch interactive menu
mobile-recon interactive
```

The interactive mode provides:
- Browse all available tools
- Select and build tools
- Run tools with guided input
- Browse by category

### Install Tools Globally

```bash
# Install a tool globally
mobile-recon install adb-toolkit
mobile-recon install nmap-toolkit

# Now you can run directly
adb-toolkit device list
nmap-toolkit scan quick 192.168.1.0/24
```

## Command Reference

### Core Commands

| Command | Description |
|---------|-------------|
| `list` | List all available tools |
| `build [tool]` | Build a specific tool or all tools |
| `run [tool] [args...]` | Run a tool with arguments |
| `install [tool]` | Install a tool globally |
| `interactive` | Launch interactive mode |

### Shortcuts

| Command | Equivalent | Description |
|---------|-----------|-------------|
| `adb [args...]` | `run adb-toolkit [args...]` | Run ADB Toolkit |
| `nmap [args...]` | `run nmap-toolkit [args...]` | Run Nmap Toolkit |

### Flags

| Flag | Description |
|------|-------------|
| `--all`, `-a` | Build all tools / show all tools |
| `--help`, `-h` | Show help information |

## Available Tools

### Mobile Tools

- **ADB Toolkit**: Android Debug Bridge automation toolkit
  - Device management
  - App management
  - Reconnaissance
  - Frida integration

### Network Tools

- **Nmap Toolkit**: Network reconnaissance and scanning toolkit
  - Port scanning
  - Service detection
  - Vulnerability scanning
  - Mobile-specific scans

## Examples

### First Time Setup

```bash
# 1. Build all tools
mobile-recon build --all

# 2. List available tools
mobile-recon list

# 3. Run a tool
mobile-recon run adb-toolkit device list
```

### Daily Usage

```bash
# Quick shortcuts
mobile-recon adb device info
mobile-recon nmap scan quick 192.168.1.0/24

# Or use interactive mode
mobile-recon interactive
```

### Development Workflow

```bash
# Build a tool after making changes
mobile-recon build adb-toolkit

# Test it
mobile-recon run adb-toolkit --help

# Install globally when ready
mobile-recon install adb-toolkit
```

## Tool Directory Structure

```
go-tools/
├── mobile-recon-cli/    # This unified CLI
├── adb-toolkit/         # Android ADB toolkit
├── nmap-toolkit/        # Network scanning toolkit
└── [future tools]/      # Additional tools
```

## Adding New Tools

To add a new tool to the unified CLI:

1. Create the tool in a new directory under `go-tools/`

2. Add the tool to `pkg/toolmanager/tools.yaml`:

```yaml
tools:
  - name: your-tool
    display_name: Your Tool
    dir: your-tool
    binary: your-tool
    description: Description of your tool
    category: Category  # Mobile, Network, Web, etc.
```

3. Rebuild the unified CLI:

```bash
cd go-tools/mobile-recon-cli
go build -o mobile-recon
```

4. Build your new tool:

```bash
./mobile-recon build your-tool
```

## Architecture

The unified CLI uses a tool manager that:

1. **Discovers Tools**: Scans the go-tools directory for known tools
2. **Checks Availability**: Verifies if tools are built
3. **Manages Execution**: Runs tools as subprocesses with proper I/O
4. **Handles Building**: Builds tools using Go commands

Tools are organized by categories for easy browsing and discovery.

## Tips

### Global Installation

Install the unified CLI globally for easy access:

```bash
cd go-tools/mobile-recon-cli
go install

# Now use from anywhere
mobile-recon list
```

### PATH Configuration

Add the built binaries to your PATH:

```bash
export PATH="$PATH:$HOME/go/bin"
```

### Aliases

Create aliases for frequently used commands:

```bash
alias mr='mobile-recon'
alias mrl='mobile-recon list'
alias mrb='mobile-recon build'
alias mri='mobile-recon interactive'
```

## Troubleshooting

### "Tool not found"

Make sure you're using the correct tool name:

```bash
mobile-recon list --all
```

### "Tool not built"

Build the tool first:

```bash
mobile-recon build [tool-name]
```

### "chdir /Users/.../go/[tool-name]: no such file or directory"

This error occurs when the CLI can't find the tool source directories. This happens when you installed globally without embedding the source path.

**Solution:** Reinstall with the correct ldflags:

```bash
cd go-tools/mobile-recon-cli
go install -ldflags "-s -w -X github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli/pkg/toolmanager.SourcePath=$(cd .. && pwd)"
```

### "Failed to initialize tool manager"

Ensure you're running from the correct directory or have proper permissions. If installed globally, make sure the source path was embedded correctly during installation.

## Contributing

This is part of the [mobile-recon](https://github.com/MKS-01/mobile-recon) project.

## License

MIT

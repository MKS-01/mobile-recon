# Tool Templates

This directory contains templates and guides for creating new tools in the Mobile Recon Toolkit.

## Contents

### 📁 tool-template/
Complete template for creating a new reconnaissance tool. Includes:
- Go project structure (main.go, cmd/, pkg/)
- Cobra CLI framework setup
- README template
- Example commands and packages

### 📘 NEW_TOOL_SETUP.md
**Complete single-page setup guide** - Everything you need to create a new tool:
- Quick start (one command)
- Step-by-step instructions
- Code patterns and examples
- Testing guide
- Troubleshooting
- Complete working example (JWT analyzer)

## Quick Start

### Option 1: Use the Generator Script (Recommended)

```bash
# From project root
./scripts/new-tool.sh <tool-name> <category> <description>

# Example
./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation toolkit"
```

The script will:
1. Create tool directory from template
2. Replace all template variables
3. Initialize Go modules
4. Register with mobile-recon-cli
5. Show next steps

### Option 2: Manual Copy

```bash
# Copy template
cp -r .templates/tool-template go-tools/your-tool-name

# Edit files and replace:
# - {{TOOL_NAME}} → your-tool-name
# - {{TOOL_NAME_TITLE}} → Your Tool Name
# - {{TOOL_SHORT_DESCRIPTION}} → Brief description
# - {{TOOL_CATEGORY}} → Mobile, Network, Web, or Other

# Initialize
cd go-tools/your-tool-name
go mod tidy
```

## Template Variables

When using the template manually, replace these placeholders:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{TOOL_NAME}}` | Tool name in kebab-case | `frida-toolkit` |
| `{{TOOL_NAME_TITLE}}` | Tool name in Title Case | `Frida Toolkit` |
| `{{TOOL_SHORT_DESCRIPTION}}` | Brief one-liner | `Frida dynamic instrumentation toolkit` |
| `{{TOOL_LONG_DESCRIPTION}}` | Detailed description | `A comprehensive toolkit for...` |
| `{{TOOL_CATEGORY}}` | Tool category | `Mobile`, `Network`, `Web`, `Other` |
| `{{VERSION}}` | Initial version | `1.0.0` |
| `{{USE_CASE}}` | Primary use case | `security testing and analysis` |

## Tool Categories

Choose the appropriate category for your tool:

- **Mobile** - Android, iOS, mobile app testing
- **Network** - Network scanning, reconnaissance, analysis
- **Web** - Web application testing, API testing
- **Other** - General-purpose or cross-category tools

## File Structure

Every tool should follow this structure:

```
tool-name/
├── cmd/                      # Cobra CLI commands
│   ├── root.go              # Root command
│   └── *.go                 # Subcommands
├── pkg/                     # Core packages
│   └── core/                # Main business logic
│       └── core.go
├── main.go                  # Entry point
├── go.mod                   # Go module
├── README.md                # Documentation
└── .gitignore              # Git ignore
```

## Examples

### Creating a Frida Tool

```bash
./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation toolkit"
cd go-tools/frida-toolkit
```

Customize `cmd/root.go`:
```go
rootCmd = &cobra.Command{
    Use:   "frida-toolkit",
    Short: "Frida dynamic instrumentation toolkit",
    // ... customize as needed
}
```

Add commands:
```bash
cp cmd/example.go cmd/attach.go
# Edit cmd/attach.go with your logic
```

### Creating a Web Tool

```bash
./scripts/new-tool.sh jwt-toolkit Web "JWT token analysis and manipulation"
cd go-tools/jwt-toolkit
```

### Creating a Network Tool

```bash
./scripts/new-tool.sh wifi-analyzer Network "WiFi security analysis toolkit"
cd go-tools/wifi-analyzer
```

## Integration

New tools are automatically integrated with:

1. **mobile-recon-cli** - Unified CLI manager
   - Listed in `mobile-recon-cli list`
   - Buildable with `mobile-recon-cli build <tool>`
   - Runnable with `mobile-recon-cli run <tool>`

2. **Installation Script** - `scripts/install.sh`
   - Automatically detects and builds new tools
   - Installs to `$GOPATH/bin`

3. **Project Documentation**
   - Add to main README.md tools section
   - Link from QUICK_START.md if applicable

## Best Practices

### Naming
- Use kebab-case: `tool-name` ✅
- Not camelCase: `toolName` ❌
- End with purpose: `frida-toolkit`, `jwt-analyzer`

### Commands
- Use noun-verb structure: `process attach`, `token decode`
- Provide examples in help text
- Use consistent flags across commands

### Code
- Keep `cmd/` lean, logic in `pkg/`
- Write unit tests in `pkg/`
- Document public functions
- Handle errors descriptively

### Documentation
- Complete README with real examples
- List all prerequisites
- Include troubleshooting section
- Add usage examples for every command

## Testing Your Tool

```bash
# Build
go build -o tool-name

# Test locally
./tool-name --help
./tool-name command --help

# Install globally
go install

# Test with CLI
mobile-recon-cli list --all
mobile-recon-cli build tool-name
mobile-recon-cli run tool-name
```

## Resources

- **[NEW_TOOL_SETUP.md](NEW_TOOL_SETUP.md)** - Complete setup guide (start here!)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Example Tools](../go-tools/) - Reference existing tools (adb-toolkit, nmap-toolkit)
- [Main README](../README.md) - Project overview

## Need Help?

- Read **[NEW_TOOL_SETUP.md](NEW_TOOL_SETUP.md)** - comprehensive single-page guide
- Check existing tools for examples
- Open an issue for questions

---

**Happy tool building! 🚀**

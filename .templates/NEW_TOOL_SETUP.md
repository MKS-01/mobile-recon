# New Tool Setup Guide

**Complete guide for creating a new tool in the Mobile Recon Toolkit**

---

## 🚀 Quick Setup (One Command)

```bash
cd /path/to/mobile-recon
./scripts/new-tool.sh <tool-name> <category> "<description>"
```

### Examples

```bash
# Frida dynamic instrumentation toolkit
./scripts/new-tool.sh frida-toolkit Mobile "Frida dynamic instrumentation toolkit"

# JWT token analyzer
./scripts/new-tool.sh jwt-analyzer Web "JWT token analysis and manipulation"

# WiFi security analyzer
./scripts/new-tool.sh wifi-toolkit Network "WiFi security analysis toolkit"

# General hash calculator
./scripts/new-tool.sh hash-toolkit Other "File hashing and integrity verification"
```

### Categories
- **Mobile** - Android, iOS, mobile app security testing
- **Network** - Network scanning, reconnaissance, analysis
- **Web** - Web application testing, API testing
- **Other** - General-purpose or cross-category tools

---

## 📋 What the Script Does

When you run `./scripts/new-tool.sh`, it automatically:

1. ✅ Validates tool name (must be kebab-case: `tool-name`)
2. ✅ Creates complete directory structure
3. ✅ Generates all boilerplate code:
   - `main.go` - Entry point
   - `cmd/root.go` - Root command with Cobra
   - `cmd/example.go` - Example subcommand
   - `pkg/core/core.go` - Core business logic
   - `go.mod` - Go module definition
   - `README.md` - Documentation template
   - `.gitignore` - Git ignore rules
4. ✅ Replaces template variables with your values
5. ✅ Initializes Go modules (`go mod tidy`)
6. ✅ Registers tool with `mobile-recon-cli`
7. ✅ Provides next steps for implementation

---

## 📁 Generated Structure

```
go-tools/your-tool/
├── cmd/                       # CLI commands (Cobra framework)
│   ├── root.go               # Root command, flags, helpers
│   └── example.go            # Example subcommand (template)
├── pkg/                      # Core packages
│   └── core/                 # Main business logic
│       └── core.go           # Implementation goes here
├── main.go                   # Entry point (minimal)
├── go.mod                    # Go module definition
├── README.md                 # Tool documentation
└── .gitignore               # Git ignore rules
```

---

## 🛠️ Next Steps After Generation

### 1. Navigate to Your Tool
```bash
cd go-tools/your-tool
```

### 2. Review Generated Files

**cmd/root.go** - Customize your root command:
```go
rootCmd = &cobra.Command{
    Use:   "your-tool",
    Short: "Your short description",
    Long: `Detailed description of what your tool does.

Features:
  • Feature 1
  • Feature 2
  • Feature 3

Perfect for [your use case].`,
    Version: "1.0.0",
}
```

**cmd/example.go** - Your first command template:
```go
var exampleCmd = &cobra.Command{
    Use:   "command [args]",
    Short: "Command description",
    Run: func(cmd *cobra.Command, args []string) {
        runExample(args)
    },
}

func runExample(args []string) {
    printInfo("Running command...")
    // TODO: Implement your logic
    printSuccess("Command completed")
}
```

**pkg/core/core.go** - Your business logic:
```go
package core

// Implement your core functionality here
func YourFunction(input string) (string, error) {
    // TODO: Add your implementation
    return result, nil
}
```

### 3. Add Your Commands

Copy the example command:
```bash
cp cmd/example.go cmd/your-command.go
```

Edit `cmd/your-command.go`:
```go
package cmd

import (
    "github.com/MKS-01/mobile-recon/go-tools/your-tool/pkg/core"
    "github.com/spf13/cobra"
)

var yourCmd = &cobra.Command{
    Use:   "scan [target]",
    Short: "Scan a target",
    Long: `Detailed description of the scan command.

Examples:
  your-tool scan target.com
  your-tool scan target.com --output results.json`,
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        target := args[0]

        printInfo("Scanning %s...", target)

        result, err := core.Scan(target, outputFile)
        if err != nil {
            printError("Scan failed: %v", err)
            return
        }

        printSuccess("Scan completed: %s", result)
    },
}

var outputFile string

func init() {
    rootCmd.AddCommand(yourCmd)
    yourCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
}
```

### 4. Implement Core Logic

Edit `pkg/core/core.go`:
```go
package core

import (
    "fmt"
    "os/exec"
)

// Scan performs a security scan on the target
func Scan(target, output string) (string, error) {
    if target == "" {
        return "", fmt.Errorf("target cannot be empty")
    }

    // Example: Run external tool
    cmd := exec.Command("nmap", "-sV", target)
    result, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("scan failed: %v", err)
    }

    // Process results
    return string(result), nil
}
```

### 5. Add Prerequisites Check (Optional)

If your tool requires external dependencies, add a check in `cmd/root.go`:

```go
import "os/exec"

var rootCmd = &cobra.Command{
    // ... existing config ...
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Check if required tool is installed
        if !isInstalled("nmap") {
            printError("nmap is not installed or not in PATH")
            printInfo("Install from: https://nmap.org/download.html")
            os.Exit(1)
        }
    },
}

func isInstalled(tool string) bool {
    _, err := exec.LookPath(tool)
    return err == nil
}
```

### 6. Build and Test

```bash
# Build the tool
go build -o your-tool

# Test help
./your-tool --help

# Test your command
./your-tool command arg1

# Run tests (if you added any)
go test ./...
```

### 7. Update README

Edit `README.md` with:
- ✅ Real features (remove template placeholders)
- ✅ Actual command examples
- ✅ Prerequisites
- ✅ Use cases
- ✅ Troubleshooting tips

### 8. Install Globally

```bash
# Install to $GOPATH/bin
go install

# Verify
which your-tool
your-tool --help
```

### 9. Test with Mobile Recon CLI

```bash
# List all tools (yours should appear)
mobile-recon-cli list --all

# Build with CLI
mobile-recon-cli build your-tool

# Run with CLI
mobile-recon-cli run your-tool command arg1

# Or use shortcut (if registered)
mobile-recon-cli your-tool command arg1
```

---

## 📚 Code Patterns & Examples

### Pattern 1: Simple Command
```go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List available items",
    Run: func(cmd *cobra.Command, args []string) {
        items, err := core.GetItems()
        if err != nil {
            printError("Failed to list: %v", err)
            return
        }

        for _, item := range items {
            fmt.Println(item)
        }
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

### Pattern 2: Command with Required Arguments
```go
var connectCmd = &cobra.Command{
    Use:   "connect <host> <port>",
    Short: "Connect to a server",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        host := args[0]
        port := args[1]

        if err := core.Connect(host, port); err != nil {
            printError("Connection failed: %v", err)
            return
        }

        printSuccess("Connected to %s:%s", host, port)
    },
}
```

### Pattern 3: Command with Flags
```go
var scanCmd = &cobra.Command{
    Use:   "scan <target>",
    Short: "Scan a target",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        target := args[0]

        opts := &core.ScanOptions{
            Verbose: verboseFlag,
            Output:  outputFlag,
            Timeout: timeoutFlag,
        }

        result, err := core.Scan(target, opts)
        if err != nil {
            printError("Scan failed: %v", err)
            return
        }

        printSuccess("Scan completed")
    },
}

var (
    verboseFlag bool
    outputFlag  string
    timeoutFlag int
)

func init() {
    rootCmd.AddCommand(scanCmd)
    scanCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
    scanCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file")
    scanCmd.Flags().IntVarP(&timeoutFlag, "timeout", "t", 30, "Timeout in seconds")
}
```

### Pattern 4: Running External Commands
```go
package core

import (
    "fmt"
    "os/exec"
)

func RunExternalTool(args ...string) (string, error) {
    cmd := exec.Command("external-tool", args...)

    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("command failed: %v\nOutput: %s", err, output)
    }

    return string(output), nil
}
```

### Pattern 5: File Operations
```go
package core

import (
    "fmt"
    "os"
)

func ProcessFile(path string) error {
    // Check if file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("file not found: %s", path)
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }

    // Process data
    result := string(data)

    // Write output
    outPath := path + ".processed"
    if err := os.WriteFile(outPath, []byte(result), 0644); err != nil {
        return fmt.Errorf("failed to write output: %v", err)
    }

    return nil
}
```

### Pattern 6: Error Handling
```go
// Good error messages with context
func DoOperation(input string) error {
    if input == "" {
        return fmt.Errorf("input cannot be empty")
    }

    result, err := externalCall(input)
    if err != nil {
        return fmt.Errorf("operation failed for input '%s': %v\nTry: check if the service is running", input, err)
    }

    return nil
}
```

---

## 🎨 Output Formatting (Helper Functions)

The template provides these helper functions in `cmd/root.go`:

```go
printSuccess("✓ Operation completed successfully")
printError("✗ Failed to connect: %v", err)
printWarning("⚠ Using default configuration")
printInfo("ℹ Processing 10 items...")
```

### Custom Colors
```go
import "github.com/fatih/color"

// Define custom colors
red := color.New(color.FgRed, color.Bold)
green := color.New(color.FgGreen)
yellow := color.New(color.FgYellow)

// Use them
red.Println("Error message")
green.Printf("Success: %s\n", result)
yellow.Print("Warning: ")
```

---

## 🧪 Testing Your Tool

### Unit Tests
Create `pkg/core/core_test.go`:

```go
package core

import "testing"

func TestYourFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed: test", false},
        {"empty input", "", "", true},
        {"special chars", "test@123", "processed: test@123", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := YourFunction(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("YourFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./pkg/core

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestYourFunction ./pkg/core
```

---

## 📦 Installation Options

### Option 1: Global Install (Recommended)
```bash
cd go-tools/your-tool
go install

# Now use from anywhere
your-tool --help
```

### Option 2: Use install.sh Script
```bash
# From project root
./scripts/install.sh

# This builds and installs ALL tools including yours
```

### Option 3: Manual Build
```bash
cd go-tools/your-tool
go build -o your-tool

# Use locally
./your-tool
```

---

## 🔍 Troubleshooting

### Tool not showing in `mobile-recon-cli list --all`

**Check registration:**
```bash
cat go-tools/mobile-recon-cli/pkg/toolmanager/toolmanager.go | grep your-tool
```

**Should see:**
```go
{
    name:        "your-tool",
    displayName: "Your Tool",
    dir:         "your-tool",
    binary:      "your-tool",
    description: "Your description",
    category:    "Category",
},
```

### Build errors

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod tidy
go mod download

# Try build again
go build
```

### `command not found` after install

```bash
# Check if installed
ls -la ~/go/bin/your-tool

# Check PATH
echo $PATH | grep "go/bin"

# Add to PATH if missing
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Import cycle detected

**Problem:** Circular dependency between packages

**Solution:** Reorganize packages:
- `cmd/` should only import from `pkg/`
- `pkg/` packages should not import from `cmd/`
- Create shared packages in `pkg/` if needed

---

## ✅ Checklist Before Publishing

- [ ] All template variables replaced
- [ ] README.md updated with real content
- [ ] Commands have help text and examples
- [ ] Prerequisites documented
- [ ] Error messages are descriptive
- [ ] Built successfully with `go build`
- [ ] Tested manually with sample inputs
- [ ] Unit tests written and passing
- [ ] Installed globally and tested
- [ ] Works with `mobile-recon-cli`
- [ ] Documentation is complete

---

## 📖 Additional Resources

### Documentation
- [Templates README](README.md) - Template system overview
- This document (NEW_TOOL_SETUP.md) - Complete setup guide

### Code Examples
- [ADB Toolkit](../go-tools/adb-toolkit/) - Mobile tool example
- [Nmap Toolkit](../go-tools/nmap-toolkit/) - Network tool example
- [Mobile Recon CLI](../go-tools/mobile-recon-cli/) - CLI manager

### External Resources
- [Cobra Documentation](https://github.com/spf13/cobra) - CLI framework
- [Effective Go](https://go.dev/doc/effective_go) - Go best practices
- [Go Project Layout](https://github.com/golang-standards/project-layout) - Project structure

---

## 💡 Tips & Best Practices

### Naming
✅ **Good:**
- `frida-toolkit`
- `jwt-analyzer`
- `wifi-scanner`

❌ **Bad:**
- `fridaToolkit` (camelCase)
- `FridaTool` (PascalCase)
- `frida_toolkit` (snake_case)

### Command Structure
✅ **Good:**
```bash
tool-name <noun> <verb> [args] [flags]
tool-name process attach com.app
tool-name token decode jwt-string
```

❌ **Bad:**
```bash
tool-name <verb> <noun> [args]  # Inconsistent
```

### Error Messages
✅ **Good:**
```go
return fmt.Errorf("failed to connect to %s: %v\nCheck if the device is authorized with 'adb devices'", device, err)
```

❌ **Bad:**
```go
return fmt.Errorf("error: %v", err)
```

### Code Organization
✅ **Good:**
- Business logic in `pkg/`
- CLI wrappers in `cmd/`
- Each package has one responsibility

❌ **Bad:**
- All code in `cmd/`
- Mixed responsibilities
- No separation of concerns

---

## 🎯 Example: Complete Tool Setup

Let's create a JWT analyzer tool from start to finish:

```bash
# 1. Generate the tool
./scripts/new-tool.sh jwt-analyzer Web "JWT token analysis and manipulation"

# 2. Navigate to it
cd go-tools/jwt-analyzer

# 3. Add a decode command
cp cmd/example.go cmd/decode.go

# 4. Edit cmd/decode.go
cat > cmd/decode.go << 'EOF'
package cmd

import (
    "github.com/MKS-01/mobile-recon/go-tools/jwt-analyzer/pkg/core"
    "github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
    Use:   "decode <jwt-token>",
    Short: "Decode a JWT token",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        token := args[0]

        printInfo("Decoding JWT token...")

        result, err := core.DecodeJWT(token)
        if err != nil {
            printError("Decode failed: %v", err)
            return
        }

        printSuccess("Decoded successfully:")
        fmt.Println(result)
    },
}

func init() {
    rootCmd.AddCommand(decodeCmd)
}
EOF

# 5. Implement core logic
cat > pkg/core/core.go << 'EOF'
package core

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
)

func DecodeJWT(token string) (string, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return "", fmt.Errorf("invalid JWT format")
    }

    // Decode payload (second part)
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return "", fmt.Errorf("failed to decode payload: %v", err)
    }

    // Pretty print JSON
    var data interface{}
    if err := json.Unmarshal(payload, &data); err != nil {
        return string(payload), nil
    }

    pretty, _ := json.MarshalIndent(data, "", "  ")
    return string(pretty), nil
}
EOF

# 6. Build and test
go mod tidy
go build -o jwt-analyzer
./jwt-analyzer decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"

# 7. Install globally
go install

# 8. Test with mobile-recon-cli
mobile-recon-cli run jwt-analyzer decode "token..."
```

---

**🎉 You're ready to create amazing tools! Happy coding! 🚀**

For questions or issues, refer to:
- This guide (`.templates/NEW_TOOL_SETUP.md`) - comprehensive documentation
- Existing tools in `go-tools/` for working examples
- Open an issue on GitHub for help

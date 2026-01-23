# APK Analyzer

A comprehensive toolkit for static analysis of Android APK files.

## Features

- **APK Information** - Extract metadata, file size, architecture info
- **Manifest Analysis** - Parse and extract AndroidManifest.xml
- **Permission Analysis** - Identify dangerous permissions and their risks
- **String Extraction** - Extract strings from DEX and native libraries
- **Security Analysis** - Detect hardcoded secrets, misconfigurations, and security issues
- **File Management** - List and extract files from APK archives

## Installation

```bash
cd go-tools/apk-analyzer
go build -o apk-analyzer
go install  # Optional: install globally
```

## Usage

### Get APK Information

```bash
# Basic info
apk-analyzer info app.apk

# Verbose output
apk-analyzer info app.apk -v

# JSON output
apk-analyzer info app.apk -o json
```

### Analyze Manifest

```bash
# Analyze manifest
apk-analyzer manifest app.apk

# Extract raw manifest
apk-analyzer manifest app.apk --extract -o AndroidManifest.xml
```

### Analyze Permissions

```bash
# List all permissions
apk-analyzer permissions app.apk

# Show only dangerous permissions
apk-analyzer permissions app.apk --dangerous

# JSON output
apk-analyzer permissions app.apk -o json
```

### Extract Strings

```bash
# Extract all strings (min length 6)
apk-analyzer strings app.apk

# Set minimum string length
apk-analyzer strings app.apk --min 10

# Search for patterns
apk-analyzer strings app.apk --search "api|key|secret"

# Extract only URLs
apk-analyzer strings app.apk --urls

# Extract only email addresses
apk-analyzer strings app.apk --emails

# Extract from specific file
apk-analyzer strings app.apk --file classes.dex
```

### Security Analysis

```bash
# Run security analysis
apk-analyzer security app.apk

# Verbose output with details
apk-analyzer security app.apk -v

# JSON output
apk-analyzer security app.apk -o json
```

### List/Extract Files

```bash
# List all files
apk-analyzer files app.apk

# List as tree
apk-analyzer files app.apk --tree

# Filter by pattern
apk-analyzer files app.apk --filter "*.dex"
apk-analyzer files app.apk --filter "lib/*"

# Extract specific file
apk-analyzer files app.apk --extract classes.dex
```

## Security Checks

The security analysis includes:

### Hardcoded Secrets
- API keys
- Private keys
- AWS credentials
- Passwords and secrets

### Configuration Issues
- Debuggable applications
- Backup enabled
- Insecure HTTP URLs

### Third-Party Detection
- Firebase integration
- Common SDKs

### Anti-Tampering Detection
- Root detection code
- Frida detection
- Xposed detection

## Dangerous Permissions

The tool identifies Android dangerous permissions including:

| Permission | Risk Category |
|------------|---------------|
| CAMERA | Privacy |
| RECORD_AUDIO | Privacy |
| ACCESS_FINE_LOCATION | Privacy/Tracking |
| READ_CONTACTS | Privacy |
| READ_SMS | Privacy |
| SEND_SMS | Financial |
| CALL_PHONE | Financial |
| READ_EXTERNAL_STORAGE | Privacy |

## Output Formats

All commands support two output formats:

- `text` (default) - Human-readable colored output
- `json` - Machine-parseable JSON output

Use `-o json` flag for JSON output.

## Use Cases

### Mobile Security Testing
```bash
# Quick security overview
apk-analyzer security app.apk

# Check permissions
apk-analyzer permissions app.apk -d
```

### Reverse Engineering
```bash
# List all files
apk-analyzer files app.apk --tree

# Extract strings for analysis
apk-analyzer strings app.apk --min 8 -o json > strings.json
```

### Malware Analysis
```bash
# Find suspicious URLs
apk-analyzer strings app.apk --urls

# Security analysis
apk-analyzer security app.apk -v
```

### API Key Discovery
```bash
# Search for API keys
apk-analyzer strings app.apk --search "api[_-]?key|secret|token"
```

## Limitations

- **Binary Manifest**: APK files contain binary XML. For full manifest parsing, use `apktool` or `aapt2`
- **Obfuscation**: String extraction may be limited for heavily obfuscated apps
- **Static Analysis**: This tool performs static analysis only; runtime behavior is not analyzed

## Related Tools

For complete APK analysis, consider using these tools alongside apk-analyzer:

- **apktool** - Full APK decompilation
- **jadx** - DEX to Java decompiler
- **aapt2** - Android Asset Packaging Tool
- **dex2jar** - DEX to JAR conversion

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please submit pull requests with:
- Clear description of changes
- Tests for new functionality
- Updated documentation

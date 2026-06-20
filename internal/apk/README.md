# APK Analyzer

A comprehensive toolkit for static analysis of Android APK files.

## Features

- **APK Information** - Extract metadata, file size, architecture info
- **Manifest Analysis** - Parse and extract AndroidManifest.xml
- **Permission Analysis** - Identify dangerous permissions and their risks
- **Abusive Permission Detection** - Detect permissions commonly abused by malware
- **String Extraction** - Extract strings from DEX and native libraries
- **Security Analysis** - Detect hardcoded secrets, misconfigurations, and security issues
- **File Management** - List and extract files from APK archives

## Installation

This toolkit ships as part of the unified `mobile-recon` CLI:

```bash
./scripts/install.sh        # from the repo root
```

The APK commands are then available under `mobile-recon apk`.

## Usage

### Get APK Information

```bash
# Basic info
mobile-recon apk info app.apk

# Verbose output
mobile-recon apk info app.apk -v

# JSON output
mobile-recon apk info app.apk -o json
```

### Analyze Manifest

```bash
# Analyze manifest
mobile-recon apk manifest app.apk

# Extract raw manifest
mobile-recon apk manifest app.apk --extract -o AndroidManifest.xml
```

### Analyze Permissions

```bash
# List all permissions
mobile-recon apk permissions app.apk

# Show only dangerous permissions
mobile-recon apk permissions app.apk --dangerous

# JSON output
mobile-recon apk permissions app.apk -o json
```

### Detect Abusive Permissions

```bash
# Scan for permissions commonly abused by malware
mobile-recon apk abuse-permissions app.apk

# Show only malware-associated permissions (high priority)
mobile-recon apk abuse-permissions app.apk --malware

# Verbose output with detailed descriptions
mobile-recon apk abuse-permissions app.apk -v

# JSON output for integration
mobile-recon apk abuse-permissions app.apk -o json
```

### Extract Strings

```bash
# Extract all strings (min length 6)
mobile-recon apk strings app.apk

# Set minimum string length
mobile-recon apk strings app.apk --min 10

# Search for patterns
mobile-recon apk strings app.apk --search "api|key|secret"

# Extract only URLs
mobile-recon apk strings app.apk --urls

# Extract only email addresses
mobile-recon apk strings app.apk --emails

# Extract from specific file
mobile-recon apk strings app.apk --file classes.dex
```

### Security Analysis

```bash
# Run security analysis
mobile-recon apk security app.apk

# Verbose output with details
mobile-recon apk security app.apk -v

# JSON output
mobile-recon apk security app.apk -o json
```

### List/Extract Files

```bash
# List all files
mobile-recon apk files app.apk

# List as tree
mobile-recon apk files app.apk --tree

# Filter by pattern
mobile-recon apk files app.apk --filter "*.dex"
mobile-recon apk files app.apk --filter "lib/*"

# Extract specific file
mobile-recon apk files app.apk --extract classes.dex
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

## Abusive Permission Detection

The `abuse-permissions` command detects two categories of potentially malicious permissions:

### Malware Permissions
Top permissions widely abused by known malware (25+ indicators):
- INTERNET, WAKE_LOCK - Persistence and data exfiltration
- CAMERA, RECORD_AUDIO - Surveillance capabilities
- READ/WRITE_EXTERNAL_STORAGE - Data theft
- READ_SMS, RECEIVE_SMS - OTP interception
- SYSTEM_ALERT_WINDOW - Overlay attacks/clickjacking
- RECEIVE_BOOT_COMPLETED - Persistence across reboots
- REQUEST_INSTALL_PACKAGES - Dropper functionality

### Other Common Abused Permissions
Permissions commonly misused but also found in legitimate apps (20+ indicators):
- FOREGROUND_SERVICE variants - Background execution
- BIND_ACCESSIBILITY_SERVICE - Keylogging, credential theft
- BIND_NOTIFICATION_LISTENER_SERVICE - Notification interception
- QUERY_ALL_PACKAGES - App reconnaissance

Each permission includes:
- **Status**: normal, dangerous, or unknown
- **Info**: Brief capability description
- **Description**: Detailed explanation of abuse potential

## Output Formats

All commands support two output formats:

- `text` (default) - Human-readable colored output
- `json` - Machine-parseable JSON output

Use `-o json` flag for JSON output.

## Use Cases

### Mobile Security Testing
```bash
# Quick security overview
mobile-recon apk security app.apk

# Check permissions
mobile-recon apk permissions app.apk -d
```

### Reverse Engineering
```bash
# List all files
mobile-recon apk files app.apk --tree

# Extract strings for analysis
mobile-recon apk strings app.apk --min 8 -o json > strings.json
```

### Malware Analysis
```bash
# Detect abusive permissions (primary malware indicator)
mobile-recon apk abuse-permissions app.apk -v

# Find suspicious URLs
mobile-recon apk strings app.apk --urls

# Security analysis
mobile-recon apk security app.apk -v
```

### API Key Discovery
```bash
# Search for API keys
mobile-recon apk strings app.apk --search "api[_-]?key|secret|token"
```

## Limitations

- **Binary Manifest**: APK files contain binary XML. For full manifest parsing, use `apktool` or `aapt2`
- **Obfuscation**: String extraction may be limited for heavily obfuscated apps
- **Static Analysis**: This tool performs static analysis only; runtime behavior is not analyzed

## Related Tools

For complete APK analysis, consider using these tools alongside `mobile-recon apk`:

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

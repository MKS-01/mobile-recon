# Nmap Toolkit

A comprehensive Go-based CLI wrapper for Nmap that provides powerful network reconnaissance capabilities with a focus on mobile security testing and penetration testing.

## Features

### Network Scanning
- **Quick Host Discovery**: Fast ping scans to identify live hosts
- **Port Scanning**: TCP, UDP, and SYN stealth scans
- **Network-wide Scanning**: Discover all devices on a network segment

### Detection & Fingerprinting
- **Service Version Detection**: Identify services and their versions
- **OS Fingerprinting**: Detect operating systems
- **Aggressive Scanning**: Comprehensive detection with scripts and traceroute

### Vulnerability Assessment
- **Vulnerability Scanning**: Automated vulnerability detection using NSE scripts
- **SSL/TLS Testing**: Comprehensive SSL/TLS cipher and certificate analysis
- **Custom NSE Scripts**: Run any Nmap Scripting Engine (NSE) scripts

### Mobile-Specific Reconnaissance
- **Android ADB Discovery**: Find devices with ADB debugging enabled
- **iOS Device Detection**: Locate iOS devices on the network
- **Mobile Service Scanning**: Scan common mobile app backend ports
- **MITM Proxy Detection**: Identify testing proxies (Burp, Charles, mitmproxy)

## Installation

### Prerequisites

1. **Install Nmap**:
   ```bash
   # macOS
   brew install nmap

   # Linux (Debian/Ubuntu)
   sudo apt install nmap

   # Linux (RHEL/CentOS)
   sudo yum install nmap

   # Windows
   # Download from https://nmap.org/download.html
   ```

2. **Install Go** (1.21 or higher):
   ```bash
   # Visit https://golang.org/dl/
   ```

### Build from Source

```bash
cd go-tools/nmap-toolkit
go mod download
go build -o nmap-toolkit
```

### Install Globally

```bash
cd go-tools/nmap-toolkit
go install
```

## Usage

### Basic Commands

```bash
# Show help
nmap-toolkit --help

# Quick host discovery
nmap-toolkit scan quick 192.168.1.0/24

# Basic port scan
nmap-toolkit scan port 192.168.1.1

# Scan specific ports
nmap-toolkit scan port 192.168.1.1 -p 80,443,8080

# Fast scan (top 100 ports)
nmap-toolkit scan port 192.168.1.1 --fast

# Aggressive timing with real-time output
nmap-toolkit scan port 192.168.1.1 --aggressive --stream
```

### Network Discovery

```bash
# Scan entire network for live hosts
nmap-toolkit scan network 192.168.1.0/24

# SYN stealth scan (requires root)
sudo nmap-toolkit scan stealth 192.168.1.1

# UDP port scan (requires root)
sudo nmap-toolkit scan udp 192.168.1.1 -p 53,67,161
```

### Service & OS Detection

```bash
# Detect services and versions
nmap-toolkit detect service 192.168.1.1

# Aggressive service detection
nmap-toolkit detect service 192.168.1.1 --aggressive

# OS detection (requires root)
sudo nmap-toolkit detect os 192.168.1.1

# Full aggressive scan (OS, service, scripts, traceroute)
sudo nmap-toolkit detect aggressive 192.168.1.1
```

### Vulnerability Scanning

```bash
# Scan for vulnerabilities
nmap-toolkit vuln scan 192.168.1.1

# Scan specific ports for vulnerabilities
nmap-toolkit vuln scan 192.168.1.1 -p 80,443,8080

# SSL/TLS enumeration
nmap-toolkit vuln ssl example.com

# Custom SSL port
nmap-toolkit vuln ssl 192.168.1.1 --ports 8443

# Run custom NSE scripts
nmap-toolkit vuln script 192.168.1.1 --scripts "http-enum,http-title"
nmap-toolkit vuln script 192.168.1.1 --scripts "smb-*"
```

### Mobile Security Testing

```bash
# Discover Android ADB devices on network
nmap-toolkit mobile adb 192.168.1.0/24

# Scan for iOS devices
nmap-toolkit mobile ios 192.168.1.0/24

# Mobile-optimized service scan
nmap-toolkit mobile scan 192.168.1.100

# Detect MITM proxies (Burp, Charles, mitmproxy)
nmap-toolkit mobile mitm 192.168.1.0/24

# Scan mobile app backend ports
nmap-toolkit mobile app-ports api.example.com
```

## Command Reference

### Scan Commands

| Command | Description | Requires Root |
|---------|-------------|---------------|
| `scan quick [target]` | Quick ping scan for host discovery | No |
| `scan port [target]` | TCP port scan | No |
| `scan stealth [target]` | SYN stealth scan | Yes |
| `scan udp [target]` | UDP port scan | Yes |
| `scan network [network]` | Network-wide host discovery | No |

### Detection Commands

| Command | Description | Requires Root |
|---------|-------------|---------------|
| `detect service [target]` | Service version detection | No |
| `detect os [target]` | Operating system fingerprinting | Yes |
| `detect aggressive [target]` | Full aggressive scan | Partial |

### Vulnerability Commands

| Command | Description | Requires Root |
|---------|-------------|---------------|
| `vuln scan [target]` | Vulnerability scanning with NSE | No |
| `vuln ssl [target]` | SSL/TLS enumeration | No |
| `vuln script [target]` | Run custom NSE scripts | No |

### Mobile Commands

| Command | Description | Requires Root |
|---------|-------------|---------------|
| `mobile adb [network]` | Find Android ADB devices | No |
| `mobile ios [network]` | Find iOS devices | No |
| `mobile scan [target]` | Mobile service scan | No |
| `mobile mitm [network]` | Detect MITM proxies | No |
| `mobile app-ports [target]` | Scan mobile backend ports | No |

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--ports` | `-p` | Specify ports (e.g., 80,443 or 1-1000) |
| `--verbose` | `-v` | Verbose output |
| `--stream` | | Stream output in real-time |
| `--fast` | | Fast scan mode (top 100 ports) |
| `--aggressive` | | Aggressive timing/detection |
| `--scripts` | | Specify NSE scripts to run |

## Examples

### Comprehensive Mobile App Security Test

```bash
# 1. Discover all devices on network
nmap-toolkit scan network 192.168.1.0/24 --stream

# 2. Find Android ADB devices
nmap-toolkit mobile adb 192.168.1.0/24

# 3. Comprehensive service scan on target
nmap-toolkit detect service 192.168.1.100 --aggressive --stream

# 4. Check for vulnerabilities
nmap-toolkit vuln scan 192.168.1.100 --stream

# 5. Test SSL/TLS configuration
nmap-toolkit vuln ssl 192.168.1.100 --ports 443

# 6. Scan mobile app backend
nmap-toolkit mobile app-ports api.example.com --stream
```

### Penetration Testing Workflow

```bash
# 1. Initial reconnaissance
nmap-toolkit scan quick 10.0.0.0/24

# 2. Deep port scan on discovered hosts
sudo nmap-toolkit scan stealth 10.0.0.50 --stream

# 3. Service and OS detection
sudo nmap-toolkit detect aggressive 10.0.0.50 --stream

# 4. Vulnerability assessment
nmap-toolkit vuln scan 10.0.0.50 -p 80,443,8080,8443 --stream

# 5. Web application enumeration
nmap-toolkit vuln script 10.0.0.50 --scripts "http-*" -p 80,443
```

### Network Security Audit

```bash
# Find all live hosts
nmap-toolkit scan network 10.0.0.0/8 --stream

# Identify exposed services
nmap-toolkit detect service 10.0.0.0/8 --stream

# Check for SSL/TLS issues
nmap-toolkit vuln ssl 10.0.0.50

# Look for common vulnerabilities
nmap-toolkit vuln scan 10.0.0.0/8 --stream
```

## Tips & Best Practices

### Performance

- Use `--stream` flag for real-time output on long scans
- Use `--fast` for quick reconnaissance (scans top 100 ports)
- Use `--aggressive` for faster timing (may be less accurate)
- Limit scan scope with specific port ranges when possible

### Stealth

- Use `scan stealth` for SYN scans that are harder to detect
- Avoid aggressive scans on production systems
- Be aware that vulnerability scans may trigger IDS/IPS systems

### Mobile Testing

- Ensure devices are on the same network segment for discovery
- Android ADB discovery only works if TCP/IP debugging is enabled
- iOS discovery uses the lockdown service (port 62078)
- MITM proxy detection helps identify test infrastructure

### Authorization

- Always obtain proper authorization before scanning networks
- Many scan types require root/administrator privileges
- Respect rate limits and avoid DoS conditions
- Document all scan activities for compliance

## Common NSE Script Categories

```bash
# Run all default scripts
nmap-toolkit vuln script [target] --scripts "default"

# Safe scripts only
nmap-toolkit vuln script [target] --scripts "safe"

# Vulnerability detection
nmap-toolkit vuln script [target] --scripts "vuln"

# Authentication testing
nmap-toolkit vuln script [target] --scripts "auth"

# Brute force attacks (use with caution)
nmap-toolkit vuln script [target] --scripts "brute"

# Discovery scripts
nmap-toolkit vuln script [target] --scripts "discovery"
```

## Troubleshooting

### "Nmap is not installed or not in PATH"

Install nmap using your package manager or download from https://nmap.org

### "Scan failed" with permission errors

Some scan types require root/administrator privileges:
```bash
sudo nmap-toolkit scan stealth [target]
sudo nmap-toolkit detect os [target]
sudo nmap-toolkit scan udp [target]
```

### Slow scans

- Use `--fast` flag for quicker results
- Specify exact ports with `-p` instead of scanning all ports
- Use `--aggressive` for faster timing
- Reduce the network range

### No devices found in mobile scans

- Ensure devices are on the same network
- For Android: Enable ADB over TCP/IP on the device
- For iOS: Ensure the device is not in airplane mode
- Check firewall rules aren't blocking scan traffic

## Security Considerations

This tool is designed for:
- Authorized penetration testing
- Network security auditing
- Mobile application security testing
- Research and educational purposes

**Always obtain proper authorization before scanning any networks or systems you don't own.**

## Contributing

This is part of the [mobile-recon](https://github.com/MKS-01/mobile-recon) project.

## License

MIT

## Credits

Built on top of [Nmap](https://nmap.org) - the industry standard network mapper.

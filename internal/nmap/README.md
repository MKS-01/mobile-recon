# Nmap Toolkit

A lightweight Go CLI wrapper for Nmap focused on local network discovery and mobile device reconnaissance.

## Features

- **Network Discovery**: Quickly find all devices on your local network
- **Port Scanning**: TCP port scanning with configurable options
- **Service Detection**: Identify services running on discovered hosts
- **SSL/TLS Testing**: Analyze SSL/TLS configuration of services
- **Mobile Device Discovery**: Find Android (ADB) and iOS devices on the network
- **MITM Proxy Detection**: Locate testing proxies (Burp, Charles, mitmproxy)

## Installation

### Prerequisites

**Nmap** must be installed:

```bash
# macOS
brew install nmap

# Linux (Debian/Ubuntu)
sudo apt install nmap

# Linux (RHEL/CentOS)
sudo yum install nmap
```

### Install

This toolkit ships as part of the unified `mobile-recon` CLI:

```bash
./scripts/install.sh        # from the repo root
```

The Nmap commands are then available under `mobile-recon nmap`.

## Usage

### Network Discovery

```bash
# Quick ping scan to find live hosts
mobile-recon nmap scan quick 192.168.1.0/24

# Discover all devices on network
mobile-recon nmap scan network 192.168.1.0/24 --stream
```

### Port Scanning

```bash
# Basic port scan
mobile-recon nmap scan port 192.168.1.100

# Scan specific ports
mobile-recon nmap scan port 192.168.1.100 -p 80,443,8080

# Fast scan (top 100 ports)
mobile-recon nmap scan port 192.168.1.100 --fast

# Aggressive timing with real-time output
mobile-recon nmap scan port 192.168.1.100 --aggressive --stream
```

### Service Detection

```bash
# Detect services and versions
mobile-recon nmap detect service 192.168.1.100

# Aggressive service detection
mobile-recon nmap detect service 192.168.1.100 --aggressive --stream
```

### SSL/TLS Testing

```bash
# Test SSL/TLS on default port 443
mobile-recon nmap ssl scan example.com

# Test SSL/TLS on custom port
mobile-recon nmap ssl scan 192.168.1.100 -p 8443
```

### Mobile Device Discovery

```bash
# Find Android devices with ADB enabled
mobile-recon nmap mobile adb 192.168.1.0/24

# Find iOS devices
mobile-recon nmap mobile ios 192.168.1.0/24

# Mobile-optimized service scan
mobile-recon nmap mobile scan 192.168.1.100

# Detect MITM proxies (Burp, Charles, mitmproxy)
mobile-recon nmap mobile mitm 192.168.1.0/24
```

## Command Reference

| Command | Description |
|---------|-------------|
| `scan quick [target]` | Quick ping scan for host discovery |
| `scan port [target]` | TCP port scan |
| `scan network [network]` | Network-wide host discovery |
| `detect service [target]` | Service version detection |
| `ssl scan [target]` | SSL/TLS enumeration |
| `mobile adb [network]` | Find Android ADB devices |
| `mobile ios [network]` | Find iOS devices |
| `mobile scan [target]` | Mobile service scan |
| `mobile mitm [network]` | Detect MITM proxies |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--ports` | `-p` | Specify ports (e.g., 80,443 or 1-1000) |
| `--verbose` | `-v` | Verbose output |
| `--stream` | | Stream output in real-time |
| `--fast` | | Fast scan mode (top 100 ports) |
| `--aggressive` | | Aggressive timing/detection |

## Example Workflow

```bash
# 1. Discover all devices on network
mobile-recon nmap scan network 192.168.1.0/24 --stream

# 2. Find mobile devices
mobile-recon nmap mobile adb 192.168.1.0/24
mobile-recon nmap mobile ios 192.168.1.0/24

# 3. Scan a specific device
mobile-recon nmap detect service 192.168.1.100 --stream

# 4. Test SSL/TLS if applicable
mobile-recon nmap ssl scan 192.168.1.100 -p 443
```

## Tips

- Use `--stream` for real-time output on longer scans
- Use `--fast` for quick reconnaissance
- Ensure mobile devices are on the same network segment
- Android ADB discovery only works if TCP/IP debugging is enabled
- iOS discovery uses the lockdown service (port 62078)

## Authorization

Always obtain proper authorization before scanning networks or systems you don't own.

## License

MIT

## Credits

Built on [Nmap](https://nmap.org) - the network mapper.

# Mobile Recon Toolkit

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)
![Status](https://img.shields.io/badge/Status-Work%20In%20Progress-orange)
![Built with Claude](https://img.shields.io/badge/Built%20with-Claude-blueviolet)

A unified CLI toolkit for mobile security testing — Android, iOS, and the networks they run on.

> **Educational project.** Built for learning mobile security concepts, understanding tooling, and experimenting in lab environments. Still in active development and requires thorough testing before use in any real scenario.

> Built with [Claude](https://claude.ai) as an AI-assisted development experiment.

## Preview

![mobile-recon CLI](screenshot/cli-preview.png)

## Tools

| Command | Description |
|---------|-------------|
| `mobile-recon adb` | Android device automation & reverse engineering via ADB |
| `mobile-recon nmap` | Local network discovery & mobile device reconnaissance |
| `mobile-recon apk` | Android APK static analysis |
| `mobile-recon ios` | iOS Simulator management & Frida integration |

## Requirements

| Dependency | Required For | Install |
|------------|-------------|---------|
| Go 1.21+ | Building | [golang.org/dl](https://golang.org/dl/) |
| ADB | ADB Toolkit | `brew install android-platform-tools` |
| Nmap | Nmap Toolkit | `brew install nmap` |
| Xcode | iOS Toolkit | Mac App Store |

## Installation

```bash
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon
./scripts/install.sh
```

The install script builds the binary and installs it to `~/go/bin`. If `mobile-recon` is not found after install, ensure Go's bin is in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"  # add to ~/.zshrc or ~/.bashrc
```

## Usage

```
mobile-recon <tool> <command> [flags]
```

### ADB Toolkit

```bash
mobile-recon adb device list
mobile-recon adb device info
mobile-recon adb app list
mobile-recon adb app pull com.example.app
mobile-recon adb recon logcat
mobile-recon adb frida setup
mobile-recon adb frida ps
```

### Nmap Toolkit

```bash
mobile-recon nmap scan quick 192.168.1.0/24
mobile-recon nmap scan port 192.168.1.1
mobile-recon nmap detect service 192.168.1.1
mobile-recon nmap ssl scan 192.168.1.1
mobile-recon nmap mobile adb 192.168.1.0/24
mobile-recon nmap mobile ios 192.168.1.0/24
mobile-recon nmap mobile mitm 192.168.1.0/24
```

### APK Analyzer

```bash
mobile-recon apk info app.apk
mobile-recon apk manifest app.apk
mobile-recon apk permissions app.apk
mobile-recon apk security app.apk
mobile-recon apk strings --urls app.apk
```

### iOS Toolkit

```bash
mobile-recon ios device list
mobile-recon ios device boot <udid>
mobile-recon ios frida ps
mobile-recon ios frida attach <pid>
mobile-recon ios frida spawn <bundle-id>
```

## Contributing

1. Fork and create a branch: `git checkout -b feature/your-feature`
2. Make changes and verify: `go build ./...`
3. Rebuild: `cd go-tools/mobile-recon-cli && go build -o mobile-recon && mv mobile-recon ~/go/bin/`
4. Open a Pull Request

For adding a new toolkit, see [.templates/NEW_TOOL_SETUP.md](.templates/NEW_TOOL_SETUP.md).

## Disclaimer

This project is **educational and experimental** — built to learn how mobile security tooling works. It is still in progress and has not been fully tested. Use it in controlled lab environments only.

- Not production-ready
- Always obtain proper authorization before scanning networks or testing applications you don't own

**Use responsibly. The author is not liable for any misuse.**

## License

MIT — see [LICENSE](LICENSE).

## Built With

- [Claude](https://claude.ai) — AI pair programmer
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Nmap](https://nmap.org) — Network mapper
- [ADB](https://developer.android.com/studio/command-line/adb) — Android Debug Bridge
- [color](https://github.com/fatih/color) — Terminal colors

# Mobile Recon Toolkit

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey)
![Status](https://img.shields.io/badge/Status-Work%20In%20Progress-orange)

A unified CLI toolkit for mobile security testing — Android, iOS, and the networks they run on.

> **Educational project** for learning mobile security in lab environments. Still in active development; built with [Claude](https://claude.ai) as an AI-assisted experiment. Use in controlled environments only.

---

<img src="screenshot/cli-frida.png" width="500">

*One-command Frida server setup on Android*

---

## Tools

| Command | What it does |
|---------|--------------|
| `mobile-recon adb`  | Android device automation over ADB — devices, APKs, logs, input, and one-command Frida server setup |
| `mobile-recon apk`  | Static analysis for APKs — manifest, permissions, security checks, strings/URLs |
| `mobile-recon ios`  | iOS Simulator management with Frida — boot, attach, spawn, instrument |
| `mobile-recon nmap` | Local network discovery — live hosts, services, SSL/TLS, ADB/MITM detection |

Run `mobile-recon list` for the full tree, or `mobile-recon <tool> --help` for a tool's commands.

---

## Getting Started

### Requirements

| Dependency | Required For | Install |
|------------|-------------|---------|
| Go 1.21+ | Building | [golang.org/dl](https://golang.org/dl/) |
| ADB | `adb` | `brew install android-platform-tools` |
| Nmap | `nmap` | `brew install nmap` |
| Xcode | `ios` | Mac App Store |

### Install

```bash
git clone https://github.com/MKS-01/mobile-recon.git
cd mobile-recon
./scripts/install.sh
```

Builds the binary into `~/go/bin`. If `mobile-recon` is not found, add Go's bin to your PATH:
`export PATH="$HOME/go/bin:$PATH"`.

---

## Project Structure

Single Go module (`github.com/MKS-01/mobile-recon`); each toolkit is a package under `internal/`, compiled into one binary.

```
internal/
  cli/          # unified root command + `list`
  adb/  apk/  ios/  nmap/   # each: <tool>.go + cmd/
pkg/
  output/       # shared console output
  frida/        # Frida host-tooling lookup
```

Build from the repo root: `go build ./...` · `go vet ./...`.

---

## Disclaimer

Educational and experimental — not production-ready. Always obtain authorization before scanning networks or testing apps you don't own. Use responsibly; the author is not liable for misuse.

## License

MIT — see [LICENSE](LICENSE).

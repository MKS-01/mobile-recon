# Mobile Recon Toolkit

Experimental, AI-assisted CLI for mobile security testing — Android, iOS, and the networks they run on. Built for learning and lab use; not production-ready, so test only what you're authorized to.

## Install

```bash
git clone https://github.com/MKS-01/mobile-recon.git && cd mobile-recon && ./scripts/install.sh
```

## Usage

```bash
mobile-recon adb device list                # Android devices over ADB
mobile-recon apk security app.apk           # static APK security audit
mobile-recon nmap scan quick 192.168.1.0/24 # discover hosts on a network
mobile-recon ios device list                # iOS simulators
```

Add `--json` for machine-readable output. Run `mobile-recon list` to see all tools.

MIT licensed.

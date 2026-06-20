# iOS Toolkit

A comprehensive Go-based toolkit for iOS Simulator operations with Frida integration for dynamic instrumentation.

## Key Feature: No Jailbreak Required!

Unlike physical iOS devices, **iOS Simulators do NOT need jailbreaking** to use Frida. This is because:

1. Simulator apps run as regular macOS processes
2. Frida attaches directly via the macOS kernel
3. No frida-server needed (unlike Android)

## Installation

This toolkit ships as part of the unified `mobile-recon` CLI:

```bash
./scripts/install.sh        # from the repo root
```

The iOS commands are then available under `mobile-recon ios`.

## Prerequisites

- Xcode command line tools: `xcode-select --install`
- Frida tools: `pip3 install frida-tools`
- A booted iOS Simulator

## Usage

### Device Management

```bash
# List all simulators
mobile-recon ios device list

# Boot a simulator
mobile-recon ios device boot "iPhone 16 Pro"

# Shutdown simulators
mobile-recon ios device shutdown

# Get info about current simulator
mobile-recon ios device info
```

### Frida Commands

```bash
# Verify Frida setup
mobile-recon ios frida setup

# List running processes
mobile-recon ios frida ps

# List installed apps
mobile-recon ios frida apps

# Attach to a running app
mobile-recon ios frida attach com.example.app

# Spawn and attach (for early hooking)
mobile-recon ios frida spawn com.example.app

# Spawn with a script
mobile-recon ios frida spawn -s hook.js com.example.app

# Trace method calls
mobile-recon ios frida trace com.example.app

# Kill an app
mobile-recon ios frida kill com.example.app
```

### Target Specific Simulator

Use the `-u` flag to target a specific simulator by UDID:

```bash
mobile-recon ios -u 12345678-ABCD-1234-ABCD-123456789ABC frida ps
```

## Frida Scripts for iOS

Common Frida script patterns for iOS:

```javascript
// Hook NSURLSession
Interceptor.attach(ObjC.classes.NSURLSession['- dataTaskWithRequest:completionHandler:'].implementation, {
    onEnter: function(args) {
        var request = ObjC.Object(args[2]);
        console.log('URL: ' + request.URL().absoluteString());
    }
});

// Hook Keychain access
Interceptor.attach(Module.findExportByName(null, 'SecItemCopyMatching'), {
    onEnter: function(args) {
        console.log('SecItemCopyMatching called');
    }
});

// Bypass SSL pinning
var sslSetPeerDomainName = Module.findExportByName(null, 'SSLSetPeerDomainName');
if (sslSetPeerDomainName) {
    Interceptor.replace(sslSetPeerDomainName, new NativeCallback(function(context, name, len) {
        return 0;
    }, 'int', ['pointer', 'pointer', 'uint32']));
}
```

## Comparison with Android Frida Setup

| Feature | Android | iOS Simulator |
|---------|---------|---------------|
| Jailbreak/Root | Required | Not needed |
| frida-server | Must be pushed & started | Not needed |
| Connection | USB (`-U`) | Device ID (`-D <udid>`) |
| Setup complexity | Multiple steps | Just install frida-tools |

## Troubleshooting

### "No booted simulator found"
Boot a simulator first:
```bash
mobile-recon ios device boot "iPhone 16 Pro"
```

### "Frida not found"
Install Frida tools:
```bash
pip3 install frida-tools
# or
brew install frida
```

### "Permission denied" when attaching
This shouldn't happen on simulators. If it does, try:
```bash
# Check if System Integrity Protection is blocking
csrutil status
```

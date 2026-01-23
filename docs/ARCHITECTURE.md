# Mobile Recon Architecture

## Overview

Mobile Recon is a modular, extensible toolkit for mobile security testing and network reconnaissance. It provides a unified CLI interface to manage and run multiple specialized tools.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACE                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌────────────┐  │
│   │   CLI Mode  │    │ Interactive │    │  Shortcuts  │    │   Direct   │  │
│   │             │    │    Mode     │    │  (adb/nmap) │    │   Tools    │  │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └─────┬──────┘  │
│          │                  │                  │                  │         │
│          └──────────────────┴─────────┬────────┴──────────────────┘         │
│                                       │                                     │
└───────────────────────────────────────┼─────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MOBILE-RECON-CLI                                   │
│                         (Unified Entry Point)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        TOOL MANAGER                                  │   │
│   │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │   │
│   │  │   Discover   │  │    Build     │  │     Run      │               │   │
│   │  │    Tools     │  │    Tools     │  │    Tools     │               │   │
│   │  └──────────────┘  └──────────────┘  └──────────────┘               │   │
│   │                                                                      │   │
│   │  ┌──────────────────────────────────────────────────────────────┐   │   │
│   │  │                    tools.yaml (Registry)                      │   │   │
│   │  │  - Tool definitions                                           │   │   │
│   │  │  - Categories                                                 │   │   │
│   │  │  - Build configurations                                       │   │   │
│   │  └──────────────────────────────────────────────────────────────┘   │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└───────────────────────────────────────┬─────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            COMMON PACKAGES                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────┐         ┌─────────────────────┐                   │
│   │   common/output     │         │    common/exec      │                   │
│   │                     │         │    (planned)        │                   │
│   │  - Colored output   │         │                     │                   │
│   │  - Headers/Sections │         │  - Command runner   │                   │
│   │  - Tables           │         │  - Output capture   │                   │
│   │  - Progress         │         │  - Error handling   │                   │
│   └─────────────────────┘         └─────────────────────┘                   │
│                                                                             │
└───────────────────────────────────────┬─────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           SPECIALIZED TOOLS                                  │
├──────────────────────────────┬──────────────────────────────────────────────┤
│                              │                                              │
│   ┌────────────────────┐     │     ┌────────────────────┐                   │
│   │    ADB TOOLKIT     │     │     │   NMAP TOOLKIT     │                   │
│   │    (Mobile)        │     │     │    (Network)       │                   │
│   ├────────────────────┤     │     ├────────────────────┤                   │
│   │ Commands:          │     │     │ Commands:          │                   │
│   │  - device          │     │     │  - scan            │                   │
│   │  - app             │     │     │  - detect          │                   │
│   │  - recon           │     │     │  - vuln            │                   │
│   │  - input           │     │     │  - mobile          │                   │
│   │  - frida           │     │     │                    │                   │
│   └─────────┬──────────┘     │     └─────────┬──────────┘                   │
│             │                │               │                              │
│             ▼                │               ▼                              │
│   ┌────────────────────┐     │     ┌────────────────────┐                   │
│   │    pkg/adb         │     │     │    pkg/nmap        │                   │
│   │  (ADB wrapper)     │     │     │  (Nmap wrapper)    │                   │
│   └────────────────────┘     │     └────────────────────┘                   │
│                              │                                              │
└──────────────────────────────┴──────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          EXTERNAL DEPENDENCIES                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐                 │
│   │   ADB    │   │   Nmap   │   │  Frida   │   │  Others  │                 │
│   │ (Android │   │ (Network │   │ (Dynamic │   │          │                 │
│   │  Bridge) │   │ Scanner) │   │  Instr.) │   │          │                 │
│   └──────────┘   └──────────┘   └──────────┘   └──────────┘                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              COMMAND FLOW                                     │
└──────────────────────────────────────────────────────────────────────────────┘

  User Input                    Processing                         Output
  ──────────                    ──────────                         ──────

  $ mobile-recon          ┌─────────────────┐              ┌─────────────────┐
       │                  │  Parse Command  │              │  Formatted      │
       ▼                  │  (Cobra CLI)    │              │  Terminal       │
  ┌─────────┐             └────────┬────────┘              │  Output         │
  │  list   │────────────────────▶ │                       └────────▲────────┘
  │  build  │                      ▼                                │
  │  run    │             ┌─────────────────┐                       │
  │  adb    │             │  Tool Manager   │                       │
  │  nmap   │             │                 │                       │
  └─────────┘             │  - Discover     │                       │
                          │  - Validate     │                       │
                          │  - Execute      │                       │
                          └────────┬────────┘                       │
                                   │                                │
                                   ▼                                │
                          ┌─────────────────┐                       │
                          │  Tool Binary    │───────────────────────┘
                          │                 │
                          │  adb-toolkit    │
                          │  nmap-toolkit   │
                          └─────────────────┘


┌──────────────────────────────────────────────────────────────────────────────┐
│                              BUILD FLOW                                       │
└──────────────────────────────────────────────────────────────────────────────┘

  $ mobile-recon build --all
           │
           ▼
  ┌─────────────────┐
  │  Load tools.yaml │
  └────────┬────────┘
           │
           ▼
  ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
  │  For each tool  │────▶│  cd tool.Path   │────▶│  go build -o    │
  └─────────────────┘     └─────────────────┘     │  tool.Binary    │
           │                                       └────────┬────────┘
           │                                                │
           ▼                                                ▼
  ┌─────────────────┐                             ┌─────────────────┐
  │  Update status  │◀────────────────────────────│  Binary created │
  │  tool.Available │                             │  in tool dir    │
  └─────────────────┘                             └─────────────────┘
```

## Current Directory Structure

```
mobile-recon/
├── go-tools/                           # Go-based CLI tools
│   │
│   ├── mobile-recon-cli/               # Unified CLI (entry point)
│   │   ├── cmd/                        # Cobra commands
│   │   │   ├── root.go                 # Root command & init
│   │   │   ├── build.go                # Build commands
│   │   │   ├── run.go                  # Run commands & shortcuts
│   │   │   ├── list.go                 # List commands
│   │   │   └── interactive.go          # Interactive mode
│   │   ├── pkg/
│   │   │   └── toolmanager/            # Tool management
│   │   │       ├── toolmanager.go      # Core logic
│   │   │       └── tools.yaml          # Tool registry
│   │   ├── main.go                     # Entry point
│   │   ├── go.mod
│   │   └── README.md
│   │
│   ├── adb-toolkit/                    # Android ADB tool
│   │   ├── cmd/                        # Cobra commands
│   │   │   ├── root.go
│   │   │   ├── device.go               # Device management
│   │   │   ├── app.go                  # App management
│   │   │   ├── recon.go                # Reconnaissance
│   │   │   ├── input.go                # Input simulation
│   │   │   └── frida.go                # Frida integration
│   │   ├── pkg/
│   │   │   └── adb/                    # ADB wrapper
│   │   │       └── adb.go
│   │   ├── main.go
│   │   └── go.mod
│   │
│   ├── nmap-toolkit/                   # Network scanning tool
│   │   ├── cmd/                        # Cobra commands
│   │   │   ├── root.go
│   │   │   ├── scan.go                 # Scanning commands
│   │   │   ├── detect.go               # Detection commands
│   │   │   ├── vuln.go                 # Vulnerability scans
│   │   │   └── mobile.go               # Mobile-specific scans
│   │   ├── pkg/
│   │   │   └── nmap/                   # Nmap wrapper
│   │   │       └── nmap.go
│   │   ├── main.go
│   │   └── go.mod
│   │
│   └── common/                         # Shared packages
│       └── output/                     # Output formatting
│           └── output.go
│
├── bash-scripts/                       # Legacy bash scripts
├── scripts/                            # Installation scripts
├── docs/                               # Documentation
└── README.md
```

## Tool Registry (tools.yaml)

```yaml
tools:
  - name: adb-toolkit
    display_name: ADB Toolkit
    dir: adb-toolkit
    binary: adb-toolkit
    description: Android Debug Bridge automation toolkit
    category: Mobile

  - name: nmap-toolkit
    display_name: Nmap Toolkit
    dir: nmap-toolkit
    binary: nmap-toolkit
    description: Network reconnaissance and scanning toolkit
    category: Network
```

## Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| `mobile-recon-cli` | Unified entry point, tool discovery, build management |
| `toolmanager` | Tool registration, discovery, building, and execution |
| `tools.yaml` | Declarative tool configuration |
| `common/output` | Consistent terminal output formatting |
| `adb-toolkit` | Android device interaction and reconnaissance |
| `nmap-toolkit` | Network scanning and vulnerability assessment |

---

## Future Scope

### Phase 1: Core Enhancements

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PLANNED ENHANCEMENTS                                 │
└─────────────────────────────────────────────────────────────────────────────┘

  ┌─────────────────────┐
  │  Configuration      │
  │  Management         │
  │                     │
  │  - Global config    │
  │  - Per-tool config  │
  │  - Profiles         │
  └─────────────────────┘
           │
           ▼
  ┌─────────────────────┐     ┌─────────────────────┐
  │  Plugin System      │────▶│  Third-party        │
  │                     │     │  Tool Integration   │
  │  - Dynamic loading  │     │                     │
  │  - Hot reload       │     │  - Burp Suite       │
  │  - Versioning       │     │  - JADX             │
  └─────────────────────┘     │  - APKTool          │
                              └─────────────────────┘
```

#### 1.1 Configuration System
- Global configuration file (`~/.mobile-recon/config.yaml`)
- Per-tool configuration overrides
- Environment profiles (dev, staging, production)
- Credential management

#### 1.2 Enhanced Tool Manager
- Tool versioning and updates
- Dependency management between tools
- Parallel tool execution
- Tool health checks

#### 1.3 Common Packages Expansion
```
common/
├── output/          # Terminal output (existing)
├── exec/            # Command execution utilities
├── config/          # Configuration management
├── network/         # Network utilities
├── crypto/          # Cryptographic utilities
└── storage/         # Data persistence
```

### Phase 2: New Tools

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           NEW TOOL CATEGORIES                                │
└─────────────────────────────────────────────────────────────────────────────┘

  Mobile                    Network                   Web
  ──────                    ───────                   ───
  ┌──────────────┐          ┌──────────────┐          ┌──────────────┐
  │ adb-toolkit  │          │ nmap-toolkit │          │ web-toolkit  │
  │ (existing)   │          │ (existing)   │          │ (planned)    │
  └──────────────┘          └──────────────┘          └──────────────┘
  ┌──────────────┐          ┌──────────────┐          ┌──────────────┐
  │ ios-toolkit  │          │ wifi-toolkit │          │ api-toolkit  │
  │ (planned)    │          │ (planned)    │          │ (planned)    │
  └──────────────┘          └──────────────┘          └──────────────┘
  ┌──────────────┐          ┌──────────────┐
  │ apk-toolkit  │          │ proxy-toolkit│          Forensics
  │ (planned)    │          │ (planned)    │          ─────────
  └──────────────┘          └──────────────┘          ┌──────────────┐
                                                      │ forensic-    │
                                                      │ toolkit      │
                                                      └──────────────┘
```

#### Planned Tools

| Tool | Category | Description |
|------|----------|-------------|
| `ios-toolkit` | Mobile | iOS device interaction (libimobiledevice) |
| `apk-toolkit` | Mobile | APK analysis, decompilation, repackaging |
| `wifi-toolkit` | Network | WiFi auditing and analysis |
| `proxy-toolkit` | Network | MITM proxy management (mitmproxy, Burp) |
| `web-toolkit` | Web | Web application testing utilities |
| `api-toolkit` | Web | API security testing |
| `forensic-toolkit` | Forensics | Mobile forensics and data extraction |

### Phase 3: Advanced Features

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ADVANCED ARCHITECTURE                               │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌─────────────────┐
                              │   Web UI        │
                              │   (Optional)    │
                              └────────┬────────┘
                                       │
                                       ▼
  ┌─────────────┐             ┌─────────────────┐             ┌─────────────┐
  │   CLI       │────────────▶│   Core Engine   │◀────────────│   API       │
  │   Interface │             │                 │             │   Server    │
  └─────────────┘             │  - Orchestrator │             └─────────────┘
                              │  - Scheduler    │
                              │  - Event Bus    │
                              └────────┬────────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
                    ▼                  ▼                  ▼
           ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
           │   Workflow   │   │   Report     │   │   Plugin     │
           │   Engine     │   │   Generator  │   │   Manager    │
           └──────────────┘   └──────────────┘   └──────────────┘
```

#### 3.1 Workflow Engine
- Define recon workflows in YAML
- Chain multiple tools together
- Conditional execution
- Parallel execution

Example workflow:
```yaml
name: full-mobile-recon
steps:
  - tool: nmap-toolkit
    command: mobile adb ${TARGET_NETWORK}
    output: devices

  - tool: adb-toolkit
    command: device list
    foreach: devices
    output: connected_devices

  - tool: adb-toolkit
    command: app list -3
    foreach: connected_devices
    output: apps

  - tool: adb-toolkit
    command: recon dump ${app}
    foreach: apps
    parallel: true
```

#### 3.2 Report Generation
- HTML/PDF reports
- JSON export for automation
- Markdown summaries
- Vulnerability tracking

#### 3.3 API Server Mode
- REST API for tool execution
- WebSocket for real-time output
- Authentication & authorization
- Rate limiting

---

## Scaling Strategy

### Horizontal Scaling (More Tools)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         TOOL ADDITION WORKFLOW                               │
└─────────────────────────────────────────────────────────────────────────────┘

  1. Create Tool                2. Register                   3. Build & Use
  ──────────────                ────────                      ──────────────

  go-tools/                     tools.yaml                    $ mobile-recon
  └── new-toolkit/              ┌─────────────┐                  build new-toolkit
      ├── cmd/                  │ - name: ... │
      │   └── root.go           │   dir: ...  │               $ mobile-recon
      ├── pkg/                  │   binary:.. │                  run new-toolkit
      ├── main.go               └─────────────┘
      └── go.mod
```

**Steps to add a new tool:**

1. **Create directory structure:**
   ```bash
   mkdir -p go-tools/new-toolkit/{cmd,pkg}
   ```

2. **Initialize Go module:**
   ```bash
   cd go-tools/new-toolkit
   go mod init github.com/MKS-01/mobile-recon/go-tools/new-toolkit
   ```

3. **Create main.go and cmd/root.go** (follow existing patterns)

4. **Register in tools.yaml:**
   ```yaml
   - name: new-toolkit
     display_name: New Toolkit
     dir: new-toolkit
     binary: new-toolkit
     description: Description here
     category: Category
   ```

5. **Rebuild mobile-recon-cli** (to embed updated tools.yaml)

### Vertical Scaling (More Features per Tool)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         FEATURE ADDITION PATTERN                             │
└─────────────────────────────────────────────────────────────────────────────┘

  adb-toolkit/
  └── cmd/
      ├── root.go
      ├── device.go          # Existing
      ├── app.go             # Existing
      ├── recon.go           # Existing
      │
      ├── ssl.go             # NEW: SSL pinning bypass
      ├── backup.go          # NEW: Backup/restore
      └── automation.go      # NEW: UI automation
```

**Pattern for adding commands:**

1. Create new command file in `cmd/`
2. Define Cobra command with subcommands
3. Implement logic in `pkg/` if complex
4. Register in `root.go` init()

### Team Scaling

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         TEAM ORGANIZATION                                    │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌─────────────────┐
                              │   Core Team     │
                              │                 │
                              │ - mobile-recon- │
                              │   cli           │
                              │ - common/       │
                              └────────┬────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ▼                        ▼                        ▼
     ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
     │  Mobile Team    │      │  Network Team   │      │  Web Team       │
     │                 │      │                 │      │                 │
     │  - adb-toolkit  │      │  - nmap-toolkit │      │  - web-toolkit  │
     │  - ios-toolkit  │      │  - wifi-toolkit │      │  - api-toolkit  │
     │  - apk-toolkit  │      │  - proxy-toolkit│      │                 │
     └─────────────────┘      └─────────────────┘      └─────────────────┘
```

**Benefits of modular architecture:**
- Teams can work independently on different tools
- Common packages ensure consistency
- Tools can be released independently
- Easy onboarding for new contributors

---

## Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| CLI Framework | Cobra | Command parsing, help generation |
| Terminal UI | promptui | Interactive menus |
| Output | fatih/color | Colored terminal output |
| Config | YAML (gopkg.in/yaml.v3) | Tool registry, configuration |
| Build | Go modules | Dependency management |

## Best Practices

### Code Organization
- Each tool is a standalone Go module
- Shared code goes in `common/`
- Commands in `cmd/`, business logic in `pkg/`
- Consistent error handling and output formatting

### Adding Features
1. Check if it belongs in `common/` (reusable)
2. Follow existing patterns in similar tools
3. Use `common/output` for consistent formatting
4. Add tests for new functionality
5. Update documentation

### Tool Design Principles
- Single responsibility per command
- Consistent flag naming across tools
- Helpful error messages with suggestions
- Support both interactive and scripted usage

---

## Roadmap

| Phase | Timeline | Deliverables |
|-------|----------|--------------|
| 1.0 | Current | ADB Toolkit, Nmap Toolkit, Unified CLI |
| 1.1 | Next | Configuration system, enhanced common packages |
| 1.2 | Future | iOS Toolkit, APK Toolkit |
| 2.0 | Future | Workflow engine, report generation |
| 3.0 | Future | API server, Web UI |

---

## Contributing

See the main [README.md](../README.md) for contribution guidelines.

## License

MIT

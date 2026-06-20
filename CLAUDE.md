# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build

This is a multi-module Go repo — there is no `go.work`, so `go build ./...` from the repo root does **not** work. Each directory under `go-tools/` is a separate module linked by `replace` directives.

Build and vet all modules:
```bash
for mod in go-tools/*/; do (cd "$mod" && go build ./... && go vet ./...); done
```

Or use the `/build-all` skill which does this with error reporting.

Install the unified binary to `~/go/bin/mobile-recon`:
```bash
./scripts/install.sh
```

The install script is interactive on first run (PATH/alias prompts); subsequent runs are clean.

## Lint

golangci-lint is configured in `.golangci.yml` (standard defaults + misspell, unconvert). Run per-module:
```bash
(cd go-tools/<module> && golangci-lint run)
```

Or use `/lint-all` to run across all modules.

## Project Structure

- `go-tools/mobile-recon-cli/` — unified CLI entry point, imports all toolkits via `replace` directives
- `go-tools/adb-toolkit/` — Android device automation (ADB, Frida setup)
- `go-tools/nmap-toolkit/` — network discovery and scanning
- `go-tools/apk-analyzer/` — APK static analysis
- `go-tools/ios-toolkit/` — iOS Simulator management with Frida
- `go-tools/common/` — shared output/formatting utilities
- `scripts/` — install and scaffolding scripts

All toolkits use Cobra for CLI structure. When adding commands, follow the existing `cmd/root.go` + per-subcommand file pattern.

## Adding New Toolkits

Use `/new-tool <name>` to scaffold. Note: `scripts/new-tool.sh` depends on a missing `.templates/` directory — the skill scaffolds manually by mirroring `ios-toolkit` (the smallest reference).

New modules must be wired into `mobile-recon-cli` via `require` + `replace` in its `go.mod` and registered in `cmd/root.go` and `cmd/list.go`.

## Testing

No test files exist yet. When adding tests, run them per-module: `(cd go-tools/<module> && go test ./...)`.

## Key Dependencies

- Go 1.21+, Cobra CLI framework, fatih/color for output
- Runtime: ADB (android-platform-tools), Nmap, Xcode (for iOS Simulator)

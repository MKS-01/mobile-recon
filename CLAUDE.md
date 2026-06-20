# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build

This is a **single Go module** rooted at the repo (`github.com/MKS-01/mobile-recon`).
Everything compiles into one binary, so the standard one-shot commands work from the
repo root:

```bash
go build ./...
go vet ./...
```

Or use the `/build-all` skill which does this with error reporting.

Install the unified binary to `~/go/bin/mobile-recon`:
```bash
./scripts/install.sh
```

The install script is interactive on first run (PATH/alias prompts); subsequent runs are clean.

## Lint

golangci-lint is configured in `.golangci.yml` (standard defaults + misspell, unconvert):
```bash
golangci-lint run ./...
```

Or use `/lint-all`.

## Project Structure

- `main.go` — thin entry point, calls `internal/cli`
- `internal/cli/` — unified root command, banner, `list` (registers each toolkit and tags it with a cobra group)
- `internal/adb/` — Android device automation (`adb.go` helper + `cmd/` cobra commands)
- `internal/apk/` — APK static analysis (`apk.go` + `cmd/`)
- `internal/nmap/` — network discovery and scanning (`nmap.go` + `cmd/`)
- `internal/ios/` — iOS Simulator management with Frida (`ios.go` + `cmd/`)
- `pkg/output/` — shared console output/formatting utilities
- `scripts/` — install and scaffolding scripts

Each toolkit is `internal/<tool>/` with a helper package (`<tool>.go`, package `<tool>`)
and a `cmd/` subpackage (package `cmd`) holding the Cobra command tree. The toolkit
command packages are imported into `internal/cli` and registered in-process — there is
**one** binary, no per-tool build/install step. When adding commands, follow the existing
`cmd/root.go` + per-subcommand file pattern.

## Adding New Toolkits

Use `/new-tool <name>` to scaffold. Mirror `internal/ios` (the smallest reference):
create `internal/<name>/<name>.go` + `internal/<name>/cmd/`, then register the new
`cmd.RootCmd` in `internal/cli/root.go` (via `register(...)` with a group ID). No
`go.mod`/`replace` wiring is needed anymore.

## Output & JSON

`pkg/output` is the single console-output layer with a global text/JSON mode, wired to
the root command's `--json`, `--no-color`, and `--quiet` flags (applied in a
`PersistentPreRun`; `cobra.EnableTraverseRunHooks` is on so toolkit pre-runs still fire).

When adding a command, use `output.Success/Info/Section/...` for human output (these
route to stderr in JSON mode) and gate machine output on `output.IsJSON()`:

```go
if output.IsJSON() {
    if err := output.JSON(payload); err != nil {
        output.Error("Failed to generate JSON: %v", err)
    }
    return
}
// ... text rendering ...
```

All four toolkits emit JSON for their data/listing/scan commands (`apk`; `adb` device/app/recon
listings; `nmap` scans; `ios` device/app listings). Pure side-effecting actions (reboot, shell,
install, tap, boot, frida attach/spawn/trace, logcat stream) remain text-only — they have no
structured result to emit.

## Testing

No test files exist yet. When adding tests, run them from the repo root: `go test ./...`.

## Key Dependencies

- Go 1.21+, Cobra CLI framework, fatih/color for output
- Runtime: ADB (android-platform-tools), Nmap, Xcode (for iOS Simulator)

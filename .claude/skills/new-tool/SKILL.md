---
name: new-tool
description: Scaffold a new toolkit under internal/ and wire it into the mobile-recon CLI. Use when the user asks to add a new tool, toolkit, or top-level subcommand.
disable-model-invocation: true
---

Add a new toolkit named `$ARGUMENTS` (kebab-case stem, e.g. `frida`). If no name was given, ask for one plus a one-line description.

This repo is a **single Go module** (`github.com/MKS-01/mobile-recon`); toolkits are plain packages under `internal/`. There is no per-tool `go.mod`, no `replace` wiring, and no build/install step — everything compiles into the one `mobile-recon` binary.

Note: `scripts/new-tool.sh` predates this layout and is obsolete — scaffold manually by mirroring `internal/ios/` (the smallest reference).

1. Create `internal/<name>/` with:
   - `<name>.go` — helper package (`package <name>`), if the tool needs shared logic (wrappers around an external binary, parsing, etc.). Mirror `internal/ios/ios.go`.
   - `cmd/` — the Cobra command tree (`package cmd`): `root.go` exporting `RootCmd`, plus one file per subcommand. Mirror `internal/ios/cmd/`.
   - `README.md` — mirror an existing toolkit README (point install at `./scripts/install.sh`, examples as `mobile-recon <name> ...`).

2. Use the shared output package `github.com/MKS-01/mobile-recon/pkg/output` for console output. Frida host-tooling lookups go through `github.com/MKS-01/mobile-recon/pkg/frida`.

3. Register the toolkit in `internal/cli/root.go`: import its `cmd` package (aliased, e.g. `<name>cmd`) and add a `register(<name>cmd.RootCmd, "<group>")` call in `init()`, using `"mobile"` or `"network"` (or add a new `cobra.Group` to `toolGroups` for a new category). `list` reads the command tree automatically — no separate listing to update.

4. Run `go build ./...` and `go vet ./...` from the repo root (or `/build-all`) to verify.

5. Update the Tools section of the root `README.md` with the new toolkit and example commands.

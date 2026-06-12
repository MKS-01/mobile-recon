---
name: new-tool
description: Scaffold a new toolkit under go-tools/ and wire it into the mobile-recon CLI. Use when the user asks to add a new tool, toolkit, or top-level subcommand.
disable-model-invocation: true
---

Add a new toolkit named `$ARGUMENTS` (kebab-case, e.g. `frida-toolkit`). If no name was given, ask for one plus a one-line description.

Note: `scripts/new-tool.sh` exists but depends on `.templates/tool-template/`, which is missing from the repo — it will fail. Scaffold manually by mirroring an existing toolkit instead (ios-toolkit is the smallest reference).

1. Create `go-tools/<name>/` with the standard layout: `main.go`, `cmd/root.go` (+ one file per subcommand), `go.mod`, `README.md`. Copy the structure and style from `go-tools/ios-toolkit/`.

2. `go.mod` module path: `github.com/MKS-01/mobile-recon/go-tools/<name>`, `go 1.21`, cobra dependency. If it uses shared output helpers, require `.../go-tools/common` with a `replace ... => ../common` directive (see adb-toolkit's go.mod for the pattern).

3. Wire into the unified CLI in `go-tools/mobile-recon-cli`:
   - Add `require` + `replace` entries for the new module in its `go.mod`.
   - Register the toolkit's root command in `cmd/root.go`, following how the existing four toolkits are added.
   - Add it to the tool listing in `cmd/list.go`.

4. Run `go mod tidy` in the new module and in `mobile-recon-cli`, then verify with `/build-all`.

5. Update the Tools section of the root `README.md` with the new toolkit and example commands.

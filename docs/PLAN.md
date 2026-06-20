# Mobile Recon — Refactor Plan

> Status: proposal · Author: review pass 2026-06-20 · Branch: `refactor`

This document reviews the current code/CLI organization, identifies the structural
problems, and lays out a phased refactor plus a feature roadmap.

---

## 1. Current state review

### 1.1 How it's wired today

```
mobile-recon-cli ──imports──> adb-toolkit/cmd.RootCmd   (compiled IN-PROCESS)
                 ──imports──> nmap-toolkit/cmd.RootCmd
                 ──imports──> apk-analyzer/cmd.RootCmd
                 ──imports──> ios-toolkit/cmd.RootCmd
                 ──also has─> toolmanager{ build / install / RunTool(exec) }  ← parallel model
```

`cmd/root.go` registers every toolkit's `RootCmd` directly with `AddCommand`, so
`mobile-recon adb …` runs **compiled-in** code in a single binary. That is the real
execution path and it works.

### 1.2 The central problem: two contradictory execution models

Alongside the in-process model, `pkg/toolmanager` implements a *second*, separate
model that treats each toolkit as a standalone external binary to be built, marked
`Available`, and run via `exec.Command` (`RunTool`). **`RunTool` is never called.**
The consequences:

- **`mobile-recon list` lies.** It prints `✓ Available` / `✗ Not built` per tool, but
  the subcommands are *always* available because they're linked into the unified
  binary. The status reflects whether a stray standalone binary happens to exist on
  disk/PATH — irrelevant to whether the command runs.
- **`mobile-recon build` / `install` are vestigial.** They compile per-tool binaries
  that the unified CLI never invokes. The real install path is `scripts/install.sh`,
  which builds only `mobile-recon-cli`.
- New contributors can't tell which model is "true." This is the biggest source of
  the "everything feels nested / tangled" feeling.

### 1.3 Duplication

- **Output helpers exist twice.** `common/output` (`Success`, `Error`, `Header`, …)
  and `apk-analyzer/pkg/utils/output.go` (`PrintSuccess`, `PrintError`, …) are the
  same code with divergent APIs. apk-analyzer ignores `common` entirely (8 files
  import the local `utils`).
- **Banner / divider glyphs** are hand-duplicated between `cmd/root.go` and
  `scripts/install.sh`.
- **Frida host-tooling helpers** partially overlap: adb-toolkit and ios-toolkit each
  re-implement "locate the `frida`/`frida-ps`/`frida-trace` binary, get version, check
  installed." (The frida-*server* download/xz/arch logic is genuinely Android-only and
  should stay in adb.)

### 1.4 Module / dependency drift

- 6 separate `go.mod` modules glued by `replace` directives, **no `go.work`**. So
  `go build ./...` from root doesn't work — CLAUDE.md documents a `for mod in …` loop
  as the workaround. This is daily friction for every build, vet, and lint.
- **Version skew:** `nmap-toolkit` pins `cobra v1.8.0` while everything else is
  `v1.10.1`; `pflag` similarly split (`v1.0.5` vs `v1.0.9`). Each module re-lists the
  same transitive deps.

### 1.5 Layout inconsistency & nesting

- Path depth: `go-tools/<tool>/cmd/*.go` + `go-tools/<tool>/pkg/<x>/*.go` + `main.go`
  *per module*. The `go-tools/` prefix + per-tool module + `cmd` + `pkg` is 3–4 levels
  before you reach real code.
- apk-analyzer has both `pkg/apk` and `pkg/utils`; the others have a single `pkg/<x>`.
- **`apk-analyzer/pkg/apk/apk.go` is 1,621 lines** — a god-file mixing manifest
  parsing, file walking, string extraction, and security heuristics.

### 1.6 Smaller correctness / UX issues

- `list.go` builds categories from a `map[string][]Tool`; **Go randomizes map
  iteration**, so the tool list prints in non-deterministic order.
- `getShortName()` derives the CLI command (`adb`, `nmap`, …) by string-trimming
  `-toolkit`/`-analyzer` suffixes at runtime — fragile. The command alias should be
  explicit config, not inferred.
- `toolmanager.New()` ignores the error from `DiscoverTools()`.
- No tests anywhere; no CI.

---

## 2. Target architecture

Decision: **commit fully to the single-binary, in-process model** and delete the
parallel exec model. Rationale below.

> **Tradeoff — single binary vs. launcher-of-binaries.** A launcher (keep `RunTool`,
> drop the in-process imports) would let tools version and ship independently and keep
> the core binary tiny. But this is one repo, one release, one author; the tools share
> `common`; and a single static Go binary is the whole appeal for a recon tool you
> `scp` onto a box. The launcher model buys isolation we don't need and costs us a
> build/install dance we already see causing confusion. **Recommend single binary.**
> (If independent distribution ever matters, revisit — see roadmap R5.)

### 2.1 Proposed layout

Collapse the redundant `go-tools/` nesting and converge on one module with internal
packages. Cobra command trees stay per-domain but become plain packages, not modules.

```
mobile-recon/
├── go.mod                      # ONE module: github.com/MKS-01/mobile-recon
├── main.go                     # thin entrypoint -> internal/cli.Execute()
├── internal/
│   ├── cli/                    # root cmd, banner, version, list
│   │   ├── root.go
│   │   └── list.go
│   ├── adb/                    # was adb-toolkit
│   │   ├── command/            # cobra commands (device, app, input, recon, frida)
│   │   └── adb.go              # was pkg/adb
│   ├── nmap/
│   ├── apk/                    # apk.go split: manifest.go, files.go, strings.go, security.go
│   └── ios/
├── pkg/
│   ├── output/                 # the ONE output package (was common/output)
│   └── frida/                  # shared frida host-tooling helpers
├── docs/
├── scripts/
└── tools.yaml                  # if kept, drives `list` metadata only (no build/exec)
```

Why `internal/` for toolkits: Go's `internal/` makes them un-importable outside this
module — correct, since they're not meant to be reused as libraries — and removes the
`replace`-directive bookkeeping entirely. `pkg/` holds the genuinely shared, stable
surface (`output`, `frida`).

If keeping modules separate is preferred for now, the minimal-disruption alternative
is: **add a `go.work`** at the root and **dedupe `common`** — that alone kills the
build-loop friction and version skew without moving files. (See Phase 0.)

---

## 3. Refactor plan (phased, each phase compiles & is shippable)

### Phase 0 — Stop the bleeding (low risk, high payoff) ✅ done
- Add `go.work` at repo root so `go build ./...` / `go vet ./...` / lint work in one
  shot. Update CLAUDE.md to drop the for-loop workaround.
- Align dependency versions: bump nmap-toolkit to `cobra v1.10.1`; run `go mod tidy`
  everywhere.
- Delete `apk-analyzer/pkg/utils/output.go`; repoint its 8 importers at `common/output`
  (map `PrintSuccess`→`Success`, etc.). Removes the duplicate output API.
- Fix `list` ordering: sort categories and tools deterministically.

### Phase 1 — Resolve the dual-model confusion ✅ done
- ~~Delete the dead exec path: `RunTool`, `Available`, `Binary`, the `build`/`install`
  cobra commands, and the on-disk binary probing in `DiscoverTools`.~~ Deleted the
  entire `toolmanager` package, `build.go`, and `tools.yaml` (all dead) — dropped the
  now-unused `gopkg.in/yaml.v3` dependency.
- ~~Drop `tools.yaml` and make `list` introspect `rootCmd.Commands()` directly.~~ `list`
  and `--help` now read from the live cobra command tree; categories are expressed as
  cobra `Group`s (`Mobile Tools` / `Network Tools`) assigned at registration in
  `root.go`. Single source of truth — `list` can no longer disagree with what runs, and
  the misleading "✓ Available / ✗ Not built" status is gone.

### Phase 2 — Collapse the module sprawl (the big move) ✅ done
- Merged all six modules into one `go.mod` at the repo root
  (`github.com/MKS-01/mobile-recon`); deleted every per-module `go.mod`/`go.sum`, the
  `replace` directives, the per-tool `main.go`, and `go.work`. `go build ./...` /
  `go vet ./...` / `golangci-lint run ./...` now work from the root.
- Moved `go-tools/<tool>/cmd` → `internal/<tool>/cmd` and `go-tools/<tool>/pkg/<x>` →
  `internal/<tool>/<tool>.go`. **As-built deviations from the sketch above:** kept the
  command subpackage named `cmd` (not `command`) to avoid touching every `package`
  clause; the iOS helper package was renamed `simctl` → `ios` so dir and package match.
- Promoted `common/output` → `pkg/output`.
- Pointed `scripts/install.sh` at the root module (one `go build`); refreshed the
  `/build-all` and `/lint-all` skills and CLAUDE.md for the single-module layout.
- Per-tool `README.md`s preserved under `internal/<tool>/` (their build snippets still
  reference old paths — cleaned up alongside Phase 3 docs work).

### Phase 3 — Break up the god-files & extract shared frida ✅ done
- Split `internal/apk/apk.go` (1,621 LOC) along its seams: `apk.go` (open/zip core,
  98 LOC), `files.go`, `strings.go`, `manifest.go`, `security.go` (analysis logic), and
  `permissions_data.go` (the ~680-line dangerous/malware permission data tables).
  Removed dead `parseAXML` and `stringContains`.
- Extracted shared frida host-tooling into `pkg/frida` (`Locate(tool)`, `Installed()`,
  `Version()`), collapsing the duplicated path-finders in both toolkits; left Android
  frida-*server* provisioning (download/extract/push) in `internal/adb`. Also removed
  the unused `spawnApp` var in the iOS frida command.

### Phase 4 — Tests & CI (no tests exist today)
- Unit tests for the pure logic: arch mapping (`mapArchToFrida`), apk parsing, nmap
  output parsing, permission abuse heuristics.
- A GitHub Actions workflow: `go build ./... && go vet ./... && golangci-lint run &&
  go test ./...` on push/PR.

**Sequencing note:** Phases 0–1 are safe and independently mergeable. Phase 2 is the
disruptive one (touches every import) and should be its own PR. 3–4 follow.

---

## 4. New features to add now (cross-cutting, high value)

These are small, broadly useful, and become easy once Phase 0/1 land:

1. **Global `--json` output.** ✅ done (foundation). `pkg/output` now has a text/JSON
   mode: in JSON mode status messages go to stderr and a command emits its payload via
   `output.JSON(v)` on stdout. Global `--json` flag on the root command, applied in a
   `PersistentPreRun` (with `cobra.EnableTraverseRunHooks`). Migrated across **all four
   toolkits**: `apk` (overloaded `-o` gone; `files --extract` now uses `--dest`), `adb`
   (device/app/recon listings), `nmap` (all scans — streaming auto-disabled in JSON mode),
   and `ios` (device/app listings). Pure side-effecting actions (reboot, shell, install,
   tap, boot, frida attach/spawn/trace, logcat stream) stay text-only by design.
2. **Global `--quiet` / `--no-color` / `NO_COLOR` support.** ✅ done. `--no-color` sets
   `color.NoColor`; `--quiet` suppresses Info/Header/Section/Divider and the banner.
   (`fatih/color` already honors the `NO_COLOR` env var.)
3. **`mobile-recon doctor`.** One command that checks the environment: `adb`, `nmap`,
   `frida`, Xcode/`simctl` presence + versions, connected devices/simulators. Replaces
   the scattered `IsADBInstalled` / `isFridaToolsInstalled` checks with one report.
4. **`mobile-recon version`** with build metadata (already have `Version` ldflag — make
   it a real subcommand and stamp git SHA + build date).
5. **Shell completions.** Cobra gives this nearly free:
   `mobile-recon completion {bash|zsh|fish}`. Big day-to-day UX win.
6. **Consistent device targeting.** `--device/-d` exists for adb; add an analogous
   `--udid` for ios and document the convention so target selection is uniform.

---

## 5. Roadmap (later)

- **R1 — Report export.** Aggregate a recon session (device info + apk findings + nmap
  results) into a single Markdown/HTML report. Natural once `--json` exists.
- **R2 — Session/output directory.** `--out ./run-<timestamp>/` that captures raw tool
  output, pulled APKs, and the report in one artifact folder.
- **R3 — APK deep analysis.** Secret/credential scanning, certificate & signing-scheme
  inspection, native `.so` listing, third-party SDK fingerprinting.
- **R4 — Frida script library.** Bundle common instrumentation scripts (SSL pinning
  bypass, root/jailbreak detection bypass) behind `mobile-recon adb frida run <script>`.
- **R5 — Plugin model (revisit the launcher question).** *If* third parties ever need to
  ship tools independently, reintroduce an exec-based plugin discovery — but as a
  deliberate, documented extension point, not the accidental dual model we have now.
- **R6 — iOS device (not just simulator) support** via `frida-server` over USB, mirroring
  the Android flow.
- **R7 — TUI dashboard** (bubbletea) for live device/scan status.
- **R8 — Cross-platform release pipeline.** GoReleaser: tagged multi-arch binaries +
  Homebrew tap.

---

## 6. Suggested first PR

Phase 0 + Phase 1 together: kills the build-loop friction, the duplicate output
package, the misleading `list` status, and the dead `build`/`install`/`RunTool` code —
all low-risk, no file moves — and immediately makes the project legible. Phase 2 (the
module merge) follows as a dedicated PR.

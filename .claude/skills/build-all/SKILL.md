---
name: build-all
description: Build and vet every Go module in this repo. Use after making code changes to verify everything compiles, or when the user asks to build, rebuild, or check the build.
---

This repo is multi-module: each directory under `go-tools/` (adb-toolkit, apk-analyzer, common, ios-toolkit, mobile-recon-cli, nmap-toolkit) has its own `go.mod`, linked by `replace` directives. There is no `go.work` (it's gitignored), so `go build ./...` from the repo root does NOT work — you must build per module.

1. Build and vet each module:

```bash
for mod in go-tools/*/; do
  echo "=== $mod ==="
  (cd "$mod" && go build ./... && go vet ./...) || echo "FAILED: $mod"
done
```

2. `mobile-recon-cli` imports all the other modules, so building it last catches cross-module breakage. If only it fails, the problem is usually a signature change in another toolkit it consumes.

3. If the user wants the installed binary updated (so `mobile-recon` on their PATH reflects the change), run `./scripts/install.sh` — it builds `mobile-recon-cli` and installs to `~/go/bin/mobile-recon`. Note: the script is interactive (PATH/alias prompts) only on first install; on rebuilds it usually runs clean.

Report which modules passed and show full compiler output for any failure.

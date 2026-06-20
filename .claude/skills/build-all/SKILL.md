---
name: build-all
description: Build and vet the Go module in this repo. Use after making code changes to verify everything compiles, or when the user asks to build, rebuild, or check the build.
---

This repo is a single Go module rooted at the repo (`github.com/MKS-01/mobile-recon`); all toolkits live under `internal/` and compile into one binary.

1. Build and vet everything from the repo root:

```bash
go build ./... && go vet ./...
```

2. The toolkit command packages under `internal/<tool>/cmd` are imported by `internal/cli`, so a break in any toolkit surfaces when the whole module builds. Show full compiler output for any failure.

3. If the user wants the installed binary updated (so `mobile-recon` on their PATH reflects the change), run `./scripts/install.sh` — it builds the root module and installs to `~/go/bin/mobile-recon`. The script is interactive (PATH/alias prompts) only on first install; on rebuilds it usually runs clean.

Report pass/fail and show full compiler output for any failure.

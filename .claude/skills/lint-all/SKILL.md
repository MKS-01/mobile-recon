---
name: lint-all
description: Run golangci-lint across every Go module in this repo. Use after making code changes to catch lint issues, or when the user asks to lint or check code quality.
---

This repo is multi-module: each directory under `go-tools/` has its own `go.mod`. golangci-lint must be run per module.

1. Run golangci-lint on each module:

```bash
for mod in go-tools/*/; do
  echo "=== $mod ==="
  (cd "$mod" && golangci-lint run ./...) || echo "FAILED: $mod"
done
```

2. If golangci-lint is not installed, install it first:

```bash
brew install golangci-lint
```

3. Report which modules passed and show full linter output for any findings.

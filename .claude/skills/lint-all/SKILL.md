---
name: lint-all
description: Run golangci-lint across the Go module in this repo. Use after making code changes to catch lint issues, or when the user asks to lint or check code quality.
---

This repo is a single Go module rooted at the repo. golangci-lint runs once from the repo root.

1. Run golangci-lint:

```bash
golangci-lint run ./...
```

2. If golangci-lint is not installed, install it first:

```bash
brew install golangci-lint
```

3. Report pass/fail and show full linter output for any findings.

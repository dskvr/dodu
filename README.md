# dodu

> **Docker disk atlas** — an `ncdu`-style, fast, navigable TUI/CLI that shows
> what consumes disk in your Docker environment (images, containers, volumes,
> build cache, logs) and suggests safe, dry-run cleanup plans.

[![CI](https://github.com/tyutyutyu/dodu/actions/workflows/ci.yml/badge.svg)](https://github.com/tyutyutyu/dodu/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> Status: **early development** — phase 00 (bootstrap) complete. See
> [PLAN.md](PLAN.md) and [.github/prompts/](.github/prompts/) for the roadmap.

## Why

`docker system df -v` tells you how much disk is used, but not **what** to
clean or **how much** you'd actually reclaim. `dodu` answers:

- What is eating my disk? (images vs. volumes vs. build cache vs. logs)
- Which images/volumes are *exclusive* vs. *shared*?
- What can I safely prune, and how much would I get back?

## Install

```bash
go install github.com/tyutyutyu/dodu/cmd/dodu@latest
```

Pre-built binaries (Linux, macOS — amd64 & arm64) and Homebrew/Scoop packages
arrive with the `v0.1.0` release.

## Quick start

```bash
dodu version    # show build info
dodu            # launch TUI (placeholder until phase 06)
```

## Development

Requires **Go 1.23+** and [Task](https://taskfile.dev).
`golangci-lint` is auto-installed by `task lint`.

```bash
task verify     # lint + test + build (CI gate)
task test       # unit tests
task build      # ./bin/dodu
task run        # build + run
```

### Project layout

| Path | Purpose |
|---|---|
| `cmd/dodu/` | Binary entry point |
| `internal/cli/` | Cobra command tree |
| `internal/tui/` | Bubbletea TUI (phases 06–07) |
| `pkg/docker/` | Docker SDK wrapper (phase 01) |
| `pkg/scan/` | Snapshot collector (phase 02) |
| `pkg/size/` | Layer accounting (phase 03) |
| `pkg/group/` | Tree builders (phase 04) |
| `pkg/cache/` | Persistent cache (phase 05) |
| `pkg/plan/` | Cleanup planner (phase 08) |
| `pkg/export/` | JSON/CSV export (phase 10) |
| `test/integration/` | E2E tests (`-tags=integration`) |

## Contributing

1. Pick a phase from [`.github/prompts/`](.github/prompts/).
2. Run `task verify` before committing.
3. Use Conventional Commits (`feat:`, `fix:`, `chore:`, …).

## License

[MIT](LICENSE) © dodu contributors

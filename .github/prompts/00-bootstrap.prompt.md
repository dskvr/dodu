---
mode: agent
description: Project skeleton — Go module, Taskfile, golangci-lint, GitHub Actions CI, mappastruktúra
---

# Phase 00 — Bootstrap

## Goal
Üres repóból working `dodu` Go projektet csinálni: build, lint, test parancsokkal, CI-vel és a `PLAN.md`-ben rögzített mappastruktúrával.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **4. Architektúra** és **8. Tech stack**.
- Stack: **Go 1.23+**, Cobra (CLI), Bubbletea (TUI), Docker SDK, bbolt, slog, testify, golangci-lint, Taskfile, goreleaser.
- Bin neve: `dodu`. Modul: `github.com/<org>/dodu` — kérdezd meg az org nevét, ha nem ismert.

## Deliverables
1. `go.mod` Go 1.23+, modul név egyeztetve.
2. Mappák létrehozva (mind üres `.gitkeep`-pel vagy `doc.go`-val):
   - `cmd/dodu/main.go` — minimal `cobra.Command` `dodu --version` és `dodu` (TUI placeholder).
   - `pkg/docker/`, `pkg/scan/`, `pkg/size/`, `pkg/group/`, `pkg/plan/`, `pkg/cache/`, `pkg/export/`
   - `internal/tui/`, `internal/cli/`
   - `test/integration/` (build tag: `integration`)
3. `Taskfile.yml` targetekkel: `verify` (lint+test+build), `lint`, `test`, `build`, `run`, `tidy`, `version-bump`.
4. `.golangci.yml` ésszerű alapokkal (govet, staticcheck, errcheck, revive, gosec, gofumpt, gci).
5. `.github/workflows/ci.yml`: matrix (linux/macos, amd64/arm64), `task verify` futtatás, cache.
6. `README.md` — telepítés, használat, contributing skeleton.
7. `CHANGELOG.md` — Keep a Changelog formátum, `[Unreleased]` szekcióval.
8. `.gitignore` Go + IDE alapok + `dist/`, `coverage.out`.
9. `LICENSE` — MIT (vagy kérdezd meg).
10. `cmd/dodu/main.go` working: `dodu --version` printeli `dev` vagy ldflags-ből injektált verziót.

## Acceptance Criteria
- [ ] `task verify` lokálisan zöld (lint + test + build).
- [ ] `./bin/dodu --version` futtatható.
- [ ] `go vet ./...` és `golangci-lint run` hibamentes.
- [ ] CI workflow átmegy push-on (legalább lint+test).
- [ ] `README.md` tartalmazza: install (go install), gyors start, status badge.

## Anti-goals
- **Ne** írj még TUI-t, scan-t, vagy üzleti logikát. Csak a skeleton.
- **Ne** húzz be felesleges dependency-t (pl. logger framework slog helyett).
- **Ne** generálj boilerplate-et minden csomagba — üres `doc.go` elég.

## Verification
```bash
task verify
./bin/dodu --version
golangci-lint run
```

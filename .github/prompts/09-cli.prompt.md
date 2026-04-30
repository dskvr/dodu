---
mode: agent
description: Cobra CLI — scan, prune, export, cache, version alparancsok
---

# Phase 09 — CLI Commands

## Goal
Headless használat CI/script környezetben: minden TUI funkció elérhető legyen non-interaktív parancsként is.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **4 (CLI side)**.
- Cobra root: `dodu`. Default action (no subcommand) = TUI indítás.

## Deliverables
1. `internal/cli/root.go` — root command, persistent flag-ek: `--host`, `--no-cache`, `--log-level`, `--readonly`.
2. `internal/cli/scan.go` — `dodu scan [--json|--csv|--text] [--group type|project]`.
3. `internal/cli/prune.go` — `dodu prune [--dry-run|--apply] [--yes] [--filter kind=...]`.
4. `internal/cli/export.go` — `dodu export <file.{json,csv}>`.
5. `internal/cli/cache.go` — `dodu cache info|purge`.
6. `internal/cli/version.go` — verzió + git SHA + build dátum (ldflags).
7. Exit kódok dokumentálva: 0 ok, 1 általános hiba, 2 daemon nem elérhető, 3 readonly violation.
8. Tesztek: minden subcommand legalább `--help` smoke + happy path mock-kal.

## Acceptance Criteria
- [ ] `dodu --help` minden parancsot felsorol, példákkal.
- [ ] `dodu scan --json` parsable JSON-t ad (validál sémával).
- [ ] `dodu prune --apply` konfirmáció nélkül **megtagad**, ha nincs `--yes`.
- [ ] Non-tty stdin-en `--yes` automatikus elvárás (CI mód).
- [ ] Exit kódok minden ágon helyesek.

## Anti-goals
- **Ne** duplikálj logikát — minden subcommand a `pkg/*` használja.
- **Ne** csinálj interaktív promptot a CLI módban (kivéve TUI-t).

## Verification
```bash
go test ./internal/cli/...
dodu --help
dodu scan --json | jq .
```

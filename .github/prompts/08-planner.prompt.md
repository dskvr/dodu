---
mode: agent
description: Dry-run prune planner — guardrail-ek, reclaim becslés, két lépéses confirm, audit log
---

# Phase 08 — Cleanup Planner

## Goal
Biztonságos törlés-tervező: a felhasználó node-okat jelöl törlésre (`d`), a planner megmutatja a teljes tervet és **becsült visszanyerhető helyet**, majd két lépéses confirm után végrehajtja prune API-kkal.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **6 (Guardrail-ek)**.
- Default mindenhol **dry-run**. Execute csak explicit `x` + `yes` után.
- Audit log: `~/.local/state/dodu/audit.log` (XDG_STATE_HOME).

## Deliverables
1. `pkg/plan/plan.go`:
   ```go
   type Item struct { Kind; ID; Name; EstReclaim int64; Reason string }
   type Plan struct {
       Items []Item
       EstReclaim int64
       Warnings []string  // pl. "image referenced by running container"
       Blocked []Item     // soha nem törölhetők
   }
   func Build(snap *scan.Snapshot, marks []Mark) *Plan
   ```
2. `pkg/plan/guards.go`:
   - Futó konténerhez tartozó volumes/images → `Blocked` (kemény tiltás).
   - In-use volume → `Blocked`.
   - Tagged image, ami egy másik tagged image parent-je → `Warning`.
3. `pkg/plan/exec.go`:
   - `func (p *Plan) Execute(ctx, client docker.Client) (Report, error)` — csoportosítva hívja a megfelelő prune API-t (filterekkel ID-kre szűkítve, ahol megy).
   - Audit log írás minden tételről (timestamp, kind, id, reclaimed_bytes, success).
4. `internal/tui/views/planner.go` — overlay modal:
   - Lista: blokkolt (piros), warning (sárga), törlés (zöld).
   - Footer: "Reclaim ~1.2 GiB · Press `x` to execute, `Esc` to cancel".
   - Execute után második megerősítés: `yes` beírása.
5. CLI: `dodu prune --dry-run` (default), `dodu prune --apply` (még akkor is megerősít, ha `--yes` nincs).
6. Unit tesztek: guardok minden esete + execute mock-kal.

## Acceptance Criteria
- [ ] `Plan` soha nem tartalmaz `Blocked` itemet az `Items` listában.
- [ ] Execute után `Report.ActualReclaim` ≥ 80% × `EstReclaim` (integration test, lokális daemon).
- [ ] Audit log strukturált (JSON Lines), append-only.
- [ ] `--apply --yes` flag CI-friendly, de TUI-ban mindig kell interaktív confirm.
- [ ] `DODU_READONLY=1` env letiltja az execute-ot mindenhol.

## Anti-goals
- **Ne** írj saját graph-walker-t a layer eltávolításra — bízd a Docker prune-ra.
- **Ne** törölj semmit guardrail-megkerülés nélkül.

## Verification
```bash
go test ./pkg/plan/...
go test -tags=integration ./test/integration/plan_test.go
```

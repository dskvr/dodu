---
mode: agent
description: TUI nézetek — Atlas (tree), Details, Breakdown bar, csoportosítás váltás
---

# Phase 07 — TUI Views

## Goal
Renderelhető nézetek a `pkg/group.Node` fán: navigálható Atlas (ncdu-szerű), jobb oldali Details panel, alul Breakdown bar.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **5**.
- Két csoportosítás: `g` váltás by-type ↔ by-project.
- Sortby `s`: size (default), name, count.

## Deliverables
1. `internal/tui/views/atlas.go` — bal oldali tree:
   - Sorok: `[size bar] [size text] [shared marker] [name] [count]`
   - Size bar: lipgloss `█` / `░`, hossz a parent legnagyobb child-jához viszonyítva.
   - Shared marker: `~` ha `Estimated` vagy a node-ban van shared layer.
2. `internal/tui/views/details.go` — jobb panel a kiválasztott node-ra:
   - Image: tag-ek, ID, created, used by containers (Refs).
   - Container: image, status, log path+size, mounts.
   - Volume: driver, mountpoint, used by.
   - BuildCache: type, parent, last used.
3. `internal/tui/views/breakdown.go` — alul vízszintes stacked bar a totals-szal (Images/Containers/Volumes/BuildCache/Logs), szín szerint.
4. Layout: `lipgloss.JoinHorizontal/Vertical`, reszponzív (min 80×24, ideális 120×40).
5. State: belépés child-ba → stack push; vissza → pop.
6. Csoportosítás váltás: őrizze meg a kijelölést, ha lehet (kind+id alapján).

## Acceptance Criteria
- [ ] Navigáció (j/k/h/l) működik, fókusz mindig látható.
- [ ] `g` váltja a csoportosítást, `s` rendezi újra a current szintet.
- [ ] Size bar arányos és nem lóg túl.
- [ ] Details panel kis terminálban (80 oszlop) elrejthető (`Tab`).
- [ ] Snapshot szerinti rendering — semmi blokkoló SDK hívás a render közben.

## Anti-goals
- **Ne** implementáld még a planner overlay-t (08).
- **Ne** írj animációkat (transitions) — minimalista marad.

## Verification
```bash
go test ./internal/tui/views/...
go run ./cmd/dodu  # manuális UX
```

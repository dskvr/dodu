---
mode: agent
description: Bubbletea app skeleton — main loop, keymap, status bar, help overlay
---

# Phase 06 — TUI Shell

## Goal
Bubbletea-alapú TUI alapváz: app loop, központi `Model`, globális keymap, status bar, help overlay. Még nem renderel adatot — placeholder view.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **5 (UX/Billentyűk)**.
- Lib: `github.com/charmbracelet/bubbletea`, `lipgloss`, `bubbles/key`, `bubbles/help`, `bubbles/spinner`.

## Deliverables
1. `internal/tui/app.go` — `type App struct { ...; sub Model }`; `tea.Program` indítás.
2. `internal/tui/keymap.go` — minden bind a PLAN szerint (`?`, `q`, `g`, `s`, `r`, `d`, `D`, `x`, `e`, navigáció).
3. `internal/tui/help.go` — bubbles/help overlay, full + short.
4. `internal/tui/status.go` — alsó sávban: daemon ID rövidítve, snapshot kor (pl. "12s ago"), aktív csoportosítás, sortby.
5. `internal/tui/styles.go` — Lipgloss színek, dark/light auto-detect.
6. `internal/tui/messages.go` — `tea.Msg` típusok: `SnapshotLoadedMsg`, `ScanStartedMsg`, `ErrorMsg`, `RefreshMsg`.
7. `cmd/dodu/main.go`-ba integráció: ha nincs alparancs, indítsa a TUI-t.
8. Loading state: spinner + "Scanning Docker daemon...".

## Acceptance Criteria
- [ ] `dodu` indít, `?` help-et nyit, `q` kilép.
- [ ] Resize-re reszponzív (no overflow).
- [ ] Globális keymap dokumentálva a help overlay-ben.
- [ ] Nincs panic ha a daemon nem elérhető — error screen + retry.

## Anti-goals
- **Ne** implementáld még az atlas/details view-kat (07).
- **Ne** rakj üzleti logikát a TUI-ba — csak hívás `pkg/scan`, `pkg/group` felé.

## Verification
```bash
go run ./cmd/dodu
# manuálisan: ?, q, resize
go test ./internal/tui/...
```

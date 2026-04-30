---
mode: agent
description: Párhuzamos scanner — images, containers, volumes, build cache, log méretek
---

# Phase 02 — Scan Engine

## Goal
Egy párhuzamos collector (`pkg/scan`), ami a Docker daemon-ból lehúzza az összes méret-releváns adatot, és egységes `Snapshot` struktúrába rendezi.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **4** és **2 (KPI)**.
- Cél: cold scan **< 8 s**, közepes daemon ellen (~200 image, ~50 container, ~100 volume).
- `pkg/docker.Client` interface használata; **soha** ne hívd közvetlenül az SDK-t.

## Deliverables
1. `pkg/scan/snapshot.go`:
   ```go
   type Snapshot struct {
       Daemon      DaemonInfo
       Images      []Image
       Containers  []Container
       Volumes     []Volume
       BuildCache  []BuildCacheEntry
       CapturedAt  time.Time
       Duration    time.Duration
   }
   ```
2. `pkg/scan/scanner.go`:
   - `type Scanner struct { client docker.Client; concurrency int }`
   - `func (s *Scanner) Scan(ctx) (*Snapshot, error)` — `errgroup`-pal párhuzamosan futtatja a 4 fő collectort + log méretek lekérése.
   - Részleges hiba: ha pl. volumes lekérés hibázik, a többi adat menjen, és tegyél `Snapshot.Errors []error`-t.
3. Container log méret: párhuzamos `os.Stat(LogPath)` minden konténerre, **nem blokkoló** (timeout per fájl).
4. Strukturált logging (`slog`) minden lépésnél, debug szinten timing-gal.
5. Unit tesztek mock client-tel: happy path + részleges hiba.

## Acceptance Criteria
- [ ] `Scan()` egyetlen hívásban visszaad konzisztens `Snapshot`-ot.
- [ ] Részleges hibák **nem buktatják** az egész scan-t.
- [ ] Tesztek lefedik: happy, részleges hiba, ctx cancel.
- [ ] Benchmarks: `go test -bench=. ./pkg/scan/...` mock client-tel <100 ms.

## Anti-goals
- **Ne** számolj még shared/exclusive layer-eket itt — az a `pkg/size` dolga.
- **Ne** csoportosíts — az a `pkg/group` dolga.
- **Ne** cache-elj — az a `pkg/cache` dolga.

## Verification
```bash
go test -race ./pkg/scan/...
go test -bench=. ./pkg/scan/...
```

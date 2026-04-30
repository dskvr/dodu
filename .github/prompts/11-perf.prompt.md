---
mode: agent
description: Profiling, benchmarking, KPI-ok teljesítése
---

# Phase 11 — Performance

## Goal
A `PLAN.md` KPI-jainak validálása és bug-ok elhárítása: cold < 8s, warm < 2s, RAM ≤ 100 MB, binary ≤ 15 MB.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **2 (KPI)**.
- Eszközök: `pprof`, `go test -bench`, `hyperfine`, `time`.

## Deliverables
1. `test/perf/scan_bench_test.go` — bench mock-kal (200 image / 50 container / 100 volume) + nagy (2000/500/1000).
2. `test/perf/scan_real.sh` — `hyperfine` script lokális daemon ellen, riport markdownba.
3. CPU/heap pprof recipe a `docs/profiling.md`-ben.
4. Optimalizációk: `errgroup` concurrency tuning, alloc csökkentés (pre-size slices), lipgloss render cache.
5. Binary méret: `task build:release` `-ldflags="-s -w"` + UPX opcionális (default off).
6. KPI riport: `docs/perf-report.md` aktuális mérésekkel + lokális gép specifikációval.

## Acceptance Criteria
- [ ] Bench result rögzítve `docs/perf-report.md`-ben — cold scan < 8s 2000 image-ig.
- [ ] Warm cache hit < 200 ms.
- [ ] `dodu` binary < 15 MB stripped, Linux amd64.
- [ ] Heap < 100 MB tipikus snapshot esetén (`pprof` snapshot mellékelve).

## Anti-goals
- **Ne** mikro-optimalizálj olvashatatlanra — csak ha benchmarks indokolja.
- **Ne** rakj UPX-et default-tá (vírusirtó issue).

## Verification
```bash
go test -bench=. -benchmem ./test/perf/...
hyperfine --warmup 1 'dodu scan --no-cache --json'
ls -lh bin/dodu
```

# Profiling dodu

dodu ships small benchmark and profiling recipes for the scanner and grouping pipeline.
The mocked Docker client makes synthetic large fixtures trivial, so most performance work
can be done without a real daemon.

## Bench (synthetic large daemon)

```bash
go test -tags perf -bench=. -benchmem -benchtime=5x ./test/perf/...
```

Tweak `mock.NewLargeClient(images, containers, volumes)` (see `pkg/docker/mock`) to size
the fixture for your scenario.

## CPU profile

```bash
go test -tags perf -bench=BenchmarkScanLarge -benchtime=10x \
    -cpuprofile cpu.out ./test/perf/...
go tool pprof -http=: cpu.out
```

## Memory profile

```bash
go test -tags perf -bench=BenchmarkScanLarge -benchtime=10x \
    -memprofile mem.out ./test/perf/...
go tool pprof -http=: mem.out
```

## Live binary trace

To collect a runtime trace from a real `dodu scan` run:

```bash
GODEBUG=trace=1 dodu scan
go tool trace trace.out
```

## Targets to watch

| Metric                              | Goal                |
|-------------------------------------|---------------------|
| Scan, 2k img / 500 ctr / 1k vol     | < 800 ms p95        |
| Memory after scan                   | < 64 MiB resident   |
| TUI navigation key → repaint        | < 16 ms p95         |
| `dodu prune --apply` 100 items      | < 5 s wall          |

If a benchmark exceeds the goal, capture a flamegraph (`pprof -http`), file an issue
with the trace attached, and link it from `docs/perf-report.md`.

# Performance report — dodu 0.3.0

> Reference run on Linux/amd64, Go 1.23, against the mock daemon used by
> `test/perf/scan_bench_test.go`. Numbers are indicative, not normative.

## Scan throughput (synthetic, 2000 images / 500 containers / 1000 volumes)

| Metric        | Goal     | Result        |
|---------------|----------|---------------|
| Wall          | < 800 ms | ~210 ms       |
| Allocs/op     | < 200k   | ~95k          |
| Bytes/op      | < 32 MiB | ~14 MiB       |

Run via:

```bash
go test -tags perf -bench=BenchmarkScanLarge -benchmem -benchtime=10x ./test/perf/...
```

## Hot paths observed

1. `pkg/scan` — JSON unmarshalling of the Docker API responses dominates allocations.
2. `pkg/group/grouping.go` — sorting children at every nesting level. Bounded by the
   structure of the snapshot; no further wins without caching.
3. `pkg/docker/mock` — string interning of synthetic IDs (only relevant in tests).

## Anti-patterns intentionally avoided

- No `time.Sleep` in production paths.
- No reflection-heavy mapping; explicit DTO struct copies in `pkg/docker/sdk.go`.
- Bounded log-stat fan-out (`Scanner.LogConcurrency = 16` default).

## How to file a regression

If you observe a regression of >20% in any goal above:

1. Capture `go test ... -cpuprofile cpu.out` per `docs/profiling.md`.
2. Open an issue tagged `perf` with the profile and `go version` output.
3. Bisect using `git bisect run task test`.

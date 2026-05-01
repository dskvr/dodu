# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-05-01

### Added
- `pkg/scan`: parallel snapshot collector (`Scanner.Scan`) gathering images,
  containers, volumes, build cache, and per-container log file sizes via
  `errgroup`. Partial collector failures are recorded in `Snapshot.Errors`
  rather than aborting; only ping or context cancellation aborts.
  Configurable `LogConcurrency` (default 16) and `LogStatTimeout` (default 1s).
- `pkg/size`: per-image accounting (`ImageSize{Total,Shared,Exclusive,
  Estimated}`) honouring missing `SharedSize`; per-container accounting
  (writable layer + log file); `ComputeTotals` aggregating images, containers,
  volumes, build cache, logs, and a conservative reclaimable estimate; IEC/SI
  human-readable size formatter.
- `pkg/group`: navigable `Node` tree with `ByType` and `ByProject` builders
  (compose project/service grouping with `<orphan>` bucket), bidirectional
  refs (image↔containers, volume↔containers), and `Sort{BySize,ByName,
  ByCount}`.
- `pkg/cache`: bbolt-backed snapshot cache with gob encoding, daemon-keyed
  hashing (`Key`), versioned entries with corruption fallback, and
  `Policy.Fresh` TTL check (default 60s). Default path honours
  `XDG_CACHE_HOME` then `~/.cache/dodu/snapshots.db`.

### Fixed
- `pkg/docker/mock`: race in `bump`/`CallCount` writes when invoked from
  parallel goroutines (added mutex).

## [0.1.0] - 2026-05-01

### Added
- `pkg/docker`: thin, mockable Docker Engine SDK wrapper exposing only the
  surface dodu needs (`Ping`, `ListImages`, `ListContainers`, `ListVolumes`,
  `BuildCacheUsage`, `DiskUsage`, `LogFileSize`, `Prune{Images,Containers,
  Volumes,BuildCache}`, `Close`).
- Internal types decouple consumers from the Docker SDK; SDK is imported
  only from `pkg/docker`.
- Sentinel errors (`ErrDaemonUnreachable`, `ErrPermissionDenied`,
  `ErrReadOnly`) with `mapError` translation for socket / network failures.
- `pkg/docker/mock`: programmable test double implementing `docker.Client`
  with per-method `*Func` overrides and call counters.
- Unit tests for error mapping, options, mock invariants, and interface
  conformance (compile-time assertion).

### Changed
- Pinned transitive deps (`grpc`, `genproto`, `otel*`, `x/sys`, `x/time`,
  `go-connections`) to versions compatible with Go 1.23 toolchain.

## [0.0.1] - 2026-04-30

### Added
- Project skeleton: Go module, package layout (`pkg/{docker,scan,size,group,plan,cache,export}`, `internal/{cli,tui}`, `cmd/dodu`).
- Cobra-based CLI root with `version` subcommand.
- Taskfile targets: `verify`, `lint`, `test`, `build`, `run`, `tidy`, `clean`, `version-bump`.
- `golangci-lint` v2 configuration with formatters (gofumpt, gci).
- GitHub Actions CI workflow (Linux + macOS) running `task verify`.
- Phase-by-phase development prompt files under `.github/prompts/`.
- `PLAN.md` with vision, KPIs, MVP scope, architecture, UX, and roadmap.
- MIT license, `.gitignore`, `CHANGELOG.md`.

[Unreleased]: https://github.com/tyutyutyu/dodu/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/tyutyutyu/dodu/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/tyutyutyu/dodu/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/tyutyutyu/dodu/releases/tag/v0.0.1

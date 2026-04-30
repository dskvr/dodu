# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/tyutyutyu/dodu/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/tyutyutyu/dodu/releases/tag/v0.0.1

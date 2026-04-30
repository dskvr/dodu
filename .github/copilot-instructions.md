# Copilot Instructions for `dodu`

> Project: ncdu-style TUI/CLI for Docker disk usage. See [PLAN.md](../PLAN.md)
> for vision, KPIs, and architecture; phase-by-phase prompts in
> [.github/prompts/](./prompts/).

## Conventions

- **Language:** Go 1.23+ (CI pinned to 1.23). Module `github.com/tyutyutyu/dodu`.
- **Layout:** `cmd/dodu` (binary), `internal/{cli,tui}`, `pkg/{docker,scan,size,group,plan,cache,export}`, `test/integration` (`-tags=integration`).
- **SDK isolation:** the Docker Engine SDK is imported **only** by `pkg/docker`. Everything else uses the `docker.Client` interface for testability.
- **Disk-first focus:** features that don't help answer "what eats my disk?" or "what can I safely reclaim?" do not belong in the MVP.
- **Honest sizing:** always distinguish `shared` vs. `exclusive` and surface an `Estimated` flag when uncertain.
- **Safety:** every destructive operation defaults to dry-run; execute requires two-step confirmation; respect `DODU_READONLY=1`.

## Tooling

- Build / test / lint: `task verify` (gate). Other targets: `task build`, `task test`, `task test-cover`, `task run -- <args>`, `task tidy`.
- Lint: `golangci-lint` v2 (auto-installed by `task lint-install`); config in [`.golangci.yml`](../.golangci.yml).
- Formatters: `gofumpt` + `gci` (run via `golangci-lint`).
- Version: stored in [`VERSION`](../VERSION); injected via `-ldflags` in `Taskfile.yml`. Bump with `task version-bump -- patch|minor|major`.

## Workflow rules

1. Pick a phase from [`.github/prompts/`](./prompts/) and follow its Acceptance Criteria.
2. Run `task verify` **before** committing — must be green.
3. Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `ci:`).
4. Update [`CHANGELOG.md`](../CHANGELOG.md) under `[Unreleased]` for any user-visible change.
5. Update this file when introducing new conventions or architectural decisions.

## Anti-goals (do not do)

- Do not add management features (start/stop/exec containers) — that's `lazydocker`'s job.
- Do not import the Docker SDK outside `pkg/docker`.
- Do not bypass guardrails for "convenience".
- Do not add dependencies casually — prefer stdlib (`log/slog`, `errors`, `errgroup`).

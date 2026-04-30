---
mode: agent
description: goreleaser, Homebrew tap, apt/scoop, dokumentáció site, v0.1.0 release
---

# Phase 12 — Release

## Goal
v0.1.0 MVP release: cross-platform binárik, Homebrew tap, scoop bucket, GitHub Release notes, dokumentáció.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **10 (Roadmap)**.
- Tag formátum: SemVer, `v0.1.0`.
- `CHANGELOG.md` Keep a Changelog formátum (a 00-bootstrap-ban inicializálva).

## Deliverables
1. `.goreleaser.yaml`:
   - Builds: linux/macos × amd64/arm64, ldflags verzió-injektálás.
   - Archives: tar.gz + zip.
   - Checksums (sha256).
   - Homebrew tap (`homebrew-tap` repo) formula.
   - Scoop bucket (Windows opcionális).
   - GitHub Release notes a CHANGELOG-ból.
2. `.github/workflows/release.yml`:
   - Trigger: `tag push v*`.
   - Steps: checkout, setup-go, `goreleaser release`.
   - Permissions: contents: write.
3. `docs/install.md` — minden install útvonal (go install, brew, scoop, manual binary).
4. `docs/quickstart.md` — 5 perces tutorial.
5. `docs/screenshots/` — atlas, planner, details (asciinema/PNG).
6. `CHANGELOG.md` `[Unreleased]` → `[0.1.0] - <date>` mozgatás.
7. README badges: build, version, license, go report card.

## Acceptance Criteria
- [ ] `goreleaser release --snapshot --clean` lokálisan sikeres.
- [ ] Tag `v0.1.0` push → GitHub Release artifacts feltöltve.
- [ ] `brew install <org>/tap/dodu` működik (ha tap repo kész).
- [ ] `dodu --version` az injektált verziót mutatja.
- [ ] Quickstart 10 perc alatt végigvihető zero context-tel.

## Anti-goals
- **Ne** publikálj release-t törött `task verify` mellett.
- **Ne** csinálj manuális build-et release-hez — csak goreleaser.

## Verification
```bash
goreleaser check
goreleaser release --snapshot --clean
git tag v0.1.0 && git push origin v0.1.0
```

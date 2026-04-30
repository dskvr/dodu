---
mode: agent
description: Docker SDK wrapper — mockolható interface, connection handling, error mapping
---

# Phase 01 — Docker Client Wrapper

## Goal
Egy szűk, **mockolható** absztrakció a Docker daemon felett (`pkg/docker`), ami csak a `dodu` által használt műveleteket exponálja. Ez izolálja az SDK-t és lehetővé teszi a unit teszteket.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **4. Architektúra**.
- Hivatalos SDK: `github.com/docker/docker/client`.
- Connection: `DOCKER_HOST` env, alapértelmezetten `unix:///var/run/docker.sock`.
- MVP-ben **csak lokális** daemon. Remote v0.2-re halasztva.

## Deliverables
1. `pkg/docker/client.go`:
   - `type Client interface { ... }` az alábbi metódusokkal:
     - `Ping(ctx) (Info, error)`
     - `ListImages(ctx) ([]Image, error)` — kiterjesztett: SharedSize, VirtualSize, Containers, RepoTags.
     - `ListContainers(ctx, all bool) ([]Container, error)` — SizeRw, SizeRootFs, LogPath, Labels, State.
     - `ListVolumes(ctx) ([]Volume, error)` — UsageData (méret), Scope, Driver, Labels.
     - `BuildCacheUsage(ctx) ([]BuildCache, error)` — `docker system df -v` ekvivalens.
     - `DiskUsage(ctx) (DiskUsage, error)` — összesítő.
     - `Events(ctx) (<-chan Event, error)` — cache invalidációhoz (post-MVP-ben elég stub).
     - `PruneImages/Containers/Volumes/BuildCache(ctx, filters) (PruneReport, error)` — **csak ezek írhatnak**.
     - `LogFileSize(ctx, containerID) (int64, error)` — `os.Stat(LogPath)`, ha elérhető.
   - `New(opts ...Option) (Client, error)` — env-ből konfigurál, `WithHost`, `WithTimeout` opciók.
2. `pkg/docker/types.go` — saját típusok (ne szivárogtasd az SDK típusait kifelé).
3. `pkg/docker/errors.go` — `ErrDaemonUnreachable`, `ErrPermissionDenied`, error mapping.
4. `pkg/docker/mock/mock.go` — `Client` interface mock implementáció (testify/mock vagy kézi).
5. Unit tesztek: legalább happy path + error mapping.

## Acceptance Criteria
- [ ] `pkg/docker` fordít, tesztek zöldek.
- [ ] **Egyetlen** SDK import a `pkg/docker`-ben; máshol nem.
- [ ] `Client` interface 100%-ban mockolható; nincs `*client.Client` exponálva.
- [ ] `LogFileSize` graceful, ha `LogPath` üres vagy nem elérhető (nem panic).
- [ ] Error mapping: `connect: permission denied` → `ErrPermissionDenied` actionable üzenettel.

## Anti-goals
- **Ne** wrap-eld az összes SDK metódust — csak amit a `dodu` használ.
- **Ne** rakj retry logikát ide; majd a hívó dönt.
- **Ne** írj TUI-specifikus dolgot ebbe a csomagba.

## Verification
```bash
go test ./pkg/docker/...
# Manuális smoke (ha van docker):
go run ./cmd/dodu debug ping  # ha van ilyen alparancs
```

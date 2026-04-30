# dodu — Docker Disk Atlas (TUI/CLI)

> Egy **ncdu-szerű, gyors, navigálható TUI** Docker környezethez. Megmutatja, **mi eszi a lemezt** (images, containers, volumes, build cache, logok), és **biztonságos takarítási javaslatokat** ad dry-runnal.

---

## 1. Termékvízió

- **One-liner:** `dodu` — *disk-first* Docker tárhely atlasz, ncdu UX-szel.
- **Elsődleges felhasználó:** DevOps/SRE, backend devek, CI karbantartók.
- **Másodlagos:** Docker Desktop / homelab userek.
- **Differenciálás vs. létező eszközök:**
  - `lazydocker`/`ctop`: menedzsment-fókusz → mi: **disk-fókusz**.
  - `dive`: layer-elemzés egy image-re → mi: **teljes daemon** + kapcsolatok.
  - `docker system df -v`: nyers szöveg → mi: **navigálható atlasz** + cleanup planner.

---

## 2. Sikerkritériumok (KPI)

| KPI | Cél |
|---|---|
| Indulás → első képernyő (warm cache) | < 2 s |
| Indulás → első képernyő (cold) | < 8 s |
| Méret-pontosság (vs. `docker system df -v`) | eltérés ≤ 5% a fő metrikáknál |
| Dry-run javaslatból tényleges reclaim | ≥ 80% (a felhasználó által elfogadott terv esetén) |
| Alap billentyűk megtanulhatósága | < 10 perc, on-screen help-pel |
| Binary méret | ≤ 15 MB (statikus, stripped) |
| Memória használat | ≤ 100 MB tipikus daemon ellen |

---

## 3. MVP Scope

**Bent (v0.1):**
- Images: shared vs. exclusive layer méret megkülönböztetése + becslés-jelölés.
- Containers: writable layer + log fájl méret.
- Volumes: named + anonymous, használt/árva detektálás.
- Build cache: típus szerinti bontás.
- Csoportosítás: **by type** és **by compose project** (label `com.docker.compose.project`).
- Dry-run prune planner (images / containers / volumes / build-cache) reclaim becsléssel.
- Export: JSON és CSV.
- Linux + macOS (amd64, arm64).

**Kint (post-MVP):**
- Remote Docker host (DOCKER_HOST, SSH).
- Networks részletes nézet.
- Swarm / Kubernetes integráció.
- Folyamatos watch / TUI auto-refresh telemetriával.
- Webes dashboard.

---

## 4. Architektúra (high-level)

```
+----------------------+       +-----------------------+
|        TUI           |       |        CLI            |
|  (Bubbletea/Lipgloss)|       | (cobra: scan, prune,  |
|                      |       |  export, version)     |
+----------+-----------+       +-----------+-----------+
           |                               |
           v                               v
     +-----+-------------------------------+-----+
     |              Core Engine                  |
     |  - Scanner (images/containers/vols/cache) |
     |  - Sizer   (shared/exclusive accounting)  |
     |  - Grouper (by-type, by-compose)          |
     |  - Planner (dry-run prune + reclaim calc) |
     |  - Cache   (BoltDB / file, TTL, hash key) |
     +---------------------+---------------------+
                           |
                           v
                +----------+----------+
                |  Docker SDK (Go)    |
                |  /var/run/docker.sock|
                +---------------------+
```

**Modulok (pkg/ alatt):**
- `pkg/docker` — SDK wrapper, mockolható interface.
- `pkg/scan` — párhuzamos collector (images/containers/volumes/buildcache/logs).
- `pkg/size` — shared layer accounting; exclusive vs. shared/cumulative.
- `pkg/group` — by-type és by-compose-project node-fa építés.
- `pkg/plan` — prune tervező, reclaim számítás, guardrail-ek.
- `pkg/cache` — perzisztens cache (`~/.cache/dodu/`), `daemon-id+last-event-id` kulccsal.
- `pkg/export` — JSON/CSV serializer.
- `internal/tui` — Bubbletea modellek, view-k, keymap.
- `internal/cli` — Cobra parancsok.

---

## 5. UX / Billentyűk (ncdu-inspired)

| Key | Action |
|---|---|
| `↑/↓ j/k` | navigáció |
| `→ l Enter` | belépés node-ba |
| `← h Esc` | vissza |
| `g` | csoportosítás váltás (type ↔ compose) |
| `s` | rendezés (size, name, count) |
| `r` | rescan |
| `d` | törlés-jelölés (planner kosár) |
| `D` | planner megnyitása (dry-run preview) |
| `x` | execute prune (megerősítéssel) |
| `e` | export (JSON/CSV, file selector) |
| `?` | help overlay |
| `q` | quit |

**Vizualizáció:** size bar (lipgloss), shared méret szürkével, exclusive színesen.

---

## 6. Biztonság / Guardrail-ek

1. **Default dry-run** minden destruktív műveletnél.
2. Futó konténerhez kötött volumes / images **soha nem törölhetők** automatikusan.
3. Execute előtt **két lépéses megerősítés** (típus szerinti összegzés + Enter "yes").
4. Audit log: `~/.local/state/dodu/audit.log` minden végrehajtott prune-ról.
5. `--no-prune` env (`DODU_READONLY=1`) globális kapcsoló CI/szerver használatra.
6. Soha nem írunk a Docker socketbe a `prune` parancsokon kívül a planner alapján.

---

## 7. Fejlesztési fázisok (prompt fájlok)

A `.github/prompts/` mappában minden fázishoz egy `*.prompt.md`. Sorrend:

| # | Prompt | Cél |
|---|---|---|
| 00 | `00-bootstrap.prompt.md` | Go modul, Taskfile, lint, CI, mappastruktúra |
| 01 | `01-docker-client.prompt.md` | Docker SDK wrapper + mock interface |
| 02 | `02-scan-engine.prompt.md` | Párhuzamos collector minden típusra |
| 03 | `03-sizing.prompt.md` | Shared/exclusive layer accounting |
| 04 | `04-grouping.prompt.md` | By-type + by-compose node-fa |
| 05 | `05-cache.prompt.md` | Perzisztens cache + invalidáció |
| 06 | `06-tui-shell.prompt.md` | Bubbletea app skeleton + keymap |
| 07 | `07-tui-views.prompt.md` | Atlas, details, breakdown panek |
| 08 | `08-planner.prompt.md` | Dry-run prune tervező + guardrail-ek |
| 09 | `09-cli.prompt.md` | Cobra parancsok (scan/prune/export) |
| 10 | `10-export.prompt.md` | JSON/CSV export |
| 11 | `11-perf.prompt.md` | Profiling, target KPI-okhoz illesztés |
| 12 | `12-release.prompt.md` | goreleaser, brew/apt/scoop, docs |

Minden prompt önállóan futtatható és tartalmazza: **goal, context, deliverables, acceptance criteria, anti-goals**.

---

## 8. Tech stack

- **Nyelv:** Go 1.23+
- **TUI:** [Bubbletea](https://github.com/charmbracelet/bubbletea) + Lipgloss + Bubbles
- **CLI:** [Cobra](https://github.com/spf13/cobra)
- **Docker:** hivatalos `github.com/docker/docker/client`
- **Cache:** `bbolt` (single-file embed)
- **Logging:** `log/slog`
- **Test:** `testify` + testcontainers-go integration tesztekhez
- **Build:** Taskfile (`task verify`, `task build`, `task release`)
- **CI:** GitHub Actions (lint, test, build matrix, goreleaser)
- **Lint:** `golangci-lint`

---

## 9. Kockázatok

| Kockázat | Mitigáció |
|---|---|
| Shared layer méret félrevezető | Külön "shared" / "exclusive" oszlop + tooltip |
| Docker API változás | SDK pinning + integration tesztek több verzió ellen |
| Lassú daemon (sok image) | Párhuzamos collector + perzisztens cache |
| Téves törlés | Dry-run default, két lépéses confirm, audit log |
| Permission (rootless / Docker Desktop) | Egyértelmű hibaüzenet + dokumentáció |

---

## 10. Roadmap

- **v0.1 (MVP):** scan + TUI atlas + dry-run prune + export.
- **v0.2:** remote host, watch mode, networks.
- **v0.3:** plugin API, Kubernetes adapter (kubelet image GC view).
- **v1.0:** stabil API, brew/apt/scoop, dokumentációs site.

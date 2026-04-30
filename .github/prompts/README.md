# dodu — Development Prompts

Ezek a fájlok a `dodu` fejlesztését **fázisokra bontva** vezérlik. Minden prompt egy önálló, futtatható szakaszt definiál: cél, kontextus, deliverables, acceptance criteria, anti-goals, verifikáció.

## Használat (VS Code Copilot Chat)

```
/00-bootstrap
/01-docker-client
...
```

vagy kézzel: nyisd meg a kívánt `*.prompt.md` fájlt, és futtasd `Copilot: Run Prompt File` paranccsal.

## Sorrend

A prompt-ok **lineáris függőségben** vannak; az egyikre épül a következő. Csak akkor lépj tovább, ha az aktuális **Acceptance Criteria** zöld.

| # | Fájl | Cél röviden |
|---|---|---|
| 00 | [00-bootstrap](./00-bootstrap.prompt.md) | Repo skeleton |
| 01 | [01-docker-client](./01-docker-client.prompt.md) | Docker SDK wrapper |
| 02 | [02-scan-engine](./02-scan-engine.prompt.md) | Párhuzamos collector |
| 03 | [03-sizing](./03-sizing.prompt.md) | Shared/exclusive accounting |
| 04 | [04-grouping](./04-grouping.prompt.md) | by-type / by-compose fa |
| 05 | [05-cache](./05-cache.prompt.md) | Perzisztens snapshot cache |
| 06 | [06-tui-shell](./06-tui-shell.prompt.md) | Bubbletea váz |
| 07 | [07-tui-views](./07-tui-views.prompt.md) | Atlas / Details / Breakdown |
| 08 | [08-planner](./08-planner.prompt.md) | Dry-run prune planner |
| 09 | [09-cli](./09-cli.prompt.md) | Cobra alparancsok |
| 10 | [10-export](./10-export.prompt.md) | JSON / CSV |
| 11 | [11-perf](./11-perf.prompt.md) | KPI validáció |
| 12 | [12-release](./12-release.prompt.md) | v0.1.0 release |

## Globális elvek (minden fázisban)

- **Test-first** ahol értelmes (legalább kontraktus + happy path).
- **Mockolható** interfészek minden külső függőséghez.
- **Disk-first**: ne menjünk el menedzsment-feature irányba.
- **Őszinte méret**: shared vs. exclusive mindig látszik.
- **Biztonság**: dry-run default, két lépéses confirm.
- **`task verify` zöld** minden phase végén.

A teljes terv: [../../PLAN.md](../../PLAN.md).

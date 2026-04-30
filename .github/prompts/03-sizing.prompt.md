---
mode: agent
description: Shared/exclusive layer accounting + méret-aggregáció
---

# Phase 03 — Sizing & Layer Accounting

## Goal
Pontos és **őszinte** méretezés: különbség `shared` (más image-ekkel közös) és `exclusive` (csak ezé az image-é) layer méret között. Ez a `dodu` egyik kulcs differenciátora.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **2 (KPI: pontosság ≤ 5%)** és **6 (őszinte jelölés)**.
- Forrás: `Image.SharedSize`, `Image.Size`, `Image.VirtualSize` az SDK-ból; layer-ek lekérése szükség esetén `ImageInspect`-tel.
- `docker system df -v` jelenti a baseline-t a validációhoz.

## Deliverables
1. `pkg/size/sizing.go`:
   - `type ImageSize struct { Total, Shared, Exclusive int64; Estimated bool }`
   - `func ImageSizes(images []docker.Image) map[string]ImageSize`
   - Dokumentáld a számítást: `Exclusive = Size - SharedSize` (az SDK definíciója szerint).
2. `pkg/size/container.go`:
   - `type ContainerSize struct { WritableLayer, LogFile, Total int64 }`
   - `func ContainerSizes(containers []docker.Container) map[string]ContainerSize`
3. `pkg/size/totals.go`:
   - `type Totals struct { Images, Containers, Volumes, BuildCache, Logs, Reclaimable int64 }`
   - `func ComputeTotals(snap *scan.Snapshot) Totals`
4. `pkg/size/format.go` — human-readable formatter (`1.2 GB`, IEC vs SI flag).
5. Unit tesztek táblázatosan; edge case: 0 byte, nincs SharedSize info → `Estimated=true`.
6. **Validációs script** (`test/manual/validate-sizes.sh`): összehasonlítja `dodu scan --json` totals-t `docker system df` outputtal, eltérés %.

## Acceptance Criteria
- [ ] Lokális daemon ellen az `Images` total ≤ 5%-on belül van `docker system df` Images sorához képest.
- [ ] `Estimated=true` minden olyan image-en, ahol az SDK nem adott vissza shared infót.
- [ ] Formatter helyes: `1024 -> "1.0 KiB"` (IEC default), `1000 -> "1.0 kB"` (SI flag).

## Anti-goals
- **Ne** próbálj per-layer dedup-ot graphdriver szinten — az SDK adatokra építünk.
- **Ne** mutass hamis pontosságot — `Estimated` mindig látszódjon.

## Verification
```bash
go test ./pkg/size/...
./test/manual/validate-sizes.sh
```

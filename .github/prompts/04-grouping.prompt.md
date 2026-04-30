---
mode: agent
description: Csoportosítás — by type és by compose project, navigálható node-fa
---

# Phase 04 — Grouping (by-type & by-compose)

## Goal
A `Snapshot`-ot navigálható **node-fává** alakítani két nézetben:
1. **By Type:** Images / Containers / Volumes / BuildCache / Logs root node-ok.
2. **By Compose Project:** `com.docker.compose.project` label szerint csoportosítva, plusz `<orphan>` bucket.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **3** és **5 (`g` keybind)**.
- Compose label: `com.docker.compose.project`, service: `com.docker.compose.service`.
- A node-fát a TUI fogja renderelni — az API legyen view-agnosztikus.

## Deliverables
1. `pkg/group/node.go`:
   ```go
   type Node struct {
       Name      string
       Kind      Kind   // RootType, RootProject, Image, Container, Volume, BuildCache, LogFile, Service, Composite
       Size      Size   // Total, Shared, Exclusive (Estimated bool)
       Children  []*Node
       Refs      Refs   // backlinks: ImageID, ContainerIDs, VolumeNames
       Meta      map[string]string
   }
   ```
2. `pkg/group/by_type.go` — `func ByType(snap, sizes) *Node`.
3. `pkg/group/by_project.go` — `func ByProject(snap, sizes) *Node`; `<orphan>` minden olyannak, ami nem compose-managed.
4. `pkg/group/refs.go` — kapcsolatok: image → containers, volume → containers, container → image+volumes+log.
5. Sortable: `func (n *Node) Sort(by SortKey)` (Size, Name, Count).
6. Unit tesztek: szintetikus snapshot → várt fa.

## Acceptance Criteria
- [ ] Egy konténer mindig megtalálható **mindkét** nézetben (by-type és by-project).
- [ ] Parent `Size` = Σ children `Size` (Exclusive összegre, Shared esetén dedup).
- [ ] `Refs` bidirekcionális — pl. image node listázza a hivatkozó container ID-ket.
- [ ] `<orphan>` bucket létezik project nézetben, ha van non-compose container.

## Anti-goals
- **Ne** rendelj UI-t a node-okhoz (icon, color) — az TUI dolga.
- **Ne** módosítsd a `Snapshot`-ot.

## Verification
```bash
go test ./pkg/group/...
```

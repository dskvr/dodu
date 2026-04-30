---
mode: agent
description: JSON / CSV export — stabil séma, dokumentált, validálható
---

# Phase 10 — Export

## Goal
Exportáld a snapshot-ot JSON-be (gépnek) és CSV-be (embernek/Excelnek), **stabil sémával**, hogy más eszközök is használhassák.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **3 (MVP scope)**.
- Séma verziózva: `schema_version: "1.0"`.

## Deliverables
1. `pkg/export/json.go` — `func ToJSON(w io.Writer, snap *Snapshot, sizes Totals) error`.
2. `pkg/export/csv.go` — több CSV (images.csv, containers.csv, volumes.csv, build_cache.csv) zip-be vagy stdout esetén egy file egy típusra.
3. `schemas/snapshot.schema.json` — JSON Schema Draft-07.
4. `pkg/export/validate.go` — `func ValidateJSON([]byte) error` (`gojsonschema` vagy hasonló).
5. Unit tesztek: roundtrip + schema validáció + golden files (`testdata/`).
6. Dokumentáció: `docs/export-schema.md`.

## Acceptance Criteria
- [ ] `dodu export out.json` érvényes a schema ellen.
- [ ] CSV oszlopok stabilak, nevek dokumentálva.
- [ ] Backward compat: minor verzió frissítésnél (1.x) parser nem törhet.
- [ ] Számértékek byte-ban (nem human-readable string).

## Anti-goals
- **Ne** rakj a JSON-be UI-specifikus mezőt (color, icon).
- **Ne** csonkítsd az ID-ket — full SHA256.

## Verification
```bash
go test ./pkg/export/...
dodu export /tmp/dodu.json && jq -e . /tmp/dodu.json
```

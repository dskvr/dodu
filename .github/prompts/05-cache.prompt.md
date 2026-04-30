---
mode: agent
description: Perzisztens snapshot cache — bbolt, TTL, daemon-events alapú invalidáció
---

# Phase 05 — Persistent Cache

## Goal
Warm-start < 2 s legyen: a legutóbbi `Snapshot` perzisztált, és érvényesnek számít, ha a daemon nem változott.

## Context
- Lásd: [../../PLAN.md](../../PLAN.md) szakasz **2 (KPI: warm < 2s)**.
- Cache helye: `~/.cache/dodu/snapshots.db` (XDG compliant).
- Backend: `bbolt`.
- Invalidáció: TTL (default 60s) + opcionális Docker `Events` figyelése (post-MVP a watcher).

## Deliverables
1. `pkg/cache/cache.go`:
   - `type Cache interface { Load(key) (*Snapshot, bool, error); Save(key, *Snapshot) error; Purge() error }`
   - `key = sha256(daemon.ID + daemon.ServerVersion)`.
2. `pkg/cache/bolt.go` — bbolt implementáció, gob encoding.
3. `pkg/cache/policy.go`:
   - `type Policy struct { TTL time.Duration }`
   - `func (p Policy) Fresh(snap *Snapshot, now time.Time) bool`
4. CLI integráció: `dodu scan --no-cache`, `dodu cache purge`, `dodu cache info`.
5. Unit tesztek: save→load roundtrip, TTL expiry, purge.

## Acceptance Criteria
- [ ] Második futtatás (warm) `< 200 ms`-ig visszaad cached snapshot-ot mock daemon-nal.
- [ ] Cache file korrupció esetén graceful fallback (full scan + log warning).
- [ ] XDG: `XDG_CACHE_HOME` figyelembe van véve.
- [ ] `dodu cache info` mutatja a méretet, last update-et, daemon ID-t.

## Anti-goals
- **Ne** használj nehéz DB-t (sqlite). Bbolt single-file, embed.
- **Ne** cache-elj prune eredményt — destruktív műveletek mindig friss adaton dolgozzanak.

## Verification
```bash
go test ./pkg/cache/...
dodu scan && time dodu scan  # második < 200 ms
```

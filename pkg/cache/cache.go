// Package cache persists scan.Snapshots so subsequent dodu invocations can
// render instantly while a fresh scan runs in the background.
package cache

import "github.com/tyutyutyu/dodu/pkg/scan"

// Cache stores and retrieves snapshots keyed by daemon identity.
type Cache interface {
	Load(key string) (*scan.Snapshot, error)
	Save(key string, snap *scan.Snapshot) error
	Purge() error
	Close() error
}

package cache

import (
	"time"

	"github.com/tyutyutyu/dodu/pkg/scan"
)

// Policy decides whether a cached snapshot is fresh enough to display.
type Policy struct {
	TTL time.Duration
}

// DefaultPolicy is the built-in TTL: 60 seconds.
var DefaultPolicy = Policy{TTL: 60 * time.Second}

// Fresh reports whether snap is within the policy's TTL of now.
func (p Policy) Fresh(snap *scan.Snapshot, now time.Time) bool {
	if snap == nil || p.TTL <= 0 {
		return false
	}
	if snap.CapturedAt.IsZero() {
		return false
	}
	return now.Sub(snap.CapturedAt) <= p.TTL
}

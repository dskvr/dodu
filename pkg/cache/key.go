package cache

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/tyutyutyu/dodu/pkg/docker"
)

// Key derives a stable cache key from daemon identity.
//
// We hash daemon ID + ServerVersion so that talking to a different daemon (or
// after an upgrade) produces a different key without leaking the raw ID into
// filenames.
func Key(d docker.DaemonInfo) string {
	h := sha256.Sum256([]byte(d.ID + "|" + d.ServerVersion))
	return hex.EncodeToString(h[:16])
}

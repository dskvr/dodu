package size

import (
	"fmt"
	"math"
)

// Unit selects the formatting style.
type Unit int

const (
	// IEC selects 1024-based units (B, KiB, MiB, GiB, TiB).
	IEC Unit = iota
	// SI selects 1000-based units (B, kB, MB, GB, TB).
	SI
)

// Format returns a short human-readable size string.
func Format(bytes int64, u Unit) string {
	if bytes < 0 {
		return "-" + Format(-bytes, u)
	}
	if bytes < 1024 && u == IEC {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1000 && u == SI {
		return fmt.Sprintf("%d B", bytes)
	}
	var (
		base     float64
		suffixes []string
	)
	switch u {
	case SI:
		base = 1000
		suffixes = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	default:
		base = 1024
		suffixes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	}
	v := float64(bytes)
	exp := int(math.Floor(math.Log(v) / math.Log(base)))
	if exp >= len(suffixes) {
		exp = len(suffixes) - 1
	}
	if exp < 0 {
		exp = 0
	}
	scaled := v / math.Pow(base, float64(exp))
	return fmt.Sprintf("%.1f %s", scaled, suffixes[exp])
}

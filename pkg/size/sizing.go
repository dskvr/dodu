// Package size computes shared/exclusive layer accounting and totals from a
// scan.Snapshot. It is the source of truth for "how big is this thing on disk".
package size

import "github.com/tyutyutyu/dodu/pkg/docker"

// ImageSize is the per-image disk accounting.
//
//   - Total: full apparent size (image.Size).
//   - Shared: bytes shared with other images (image.SharedSize).
//   - Exclusive: bytes that disappear if this image is deleted (Total-Shared).
//   - Estimated: true if the daemon did not provide SharedSize, in which case
//     Shared is 0 and Exclusive equals Total but the value is uncertain.
type ImageSize struct {
	Total     int64
	Shared    int64
	Exclusive int64
	Estimated bool
}

// ImageSizes returns per-image accounting indexed by image ID.
func ImageSizes(images []docker.Image) map[string]ImageSize {
	out := make(map[string]ImageSize, len(images))
	for _, im := range images {
		s := ImageSize{Total: im.Size}
		switch {
		case im.SharedSize < 0:
			// Daemon did not report shared info.
			s.Estimated = true
			s.Exclusive = im.Size
		default:
			s.Shared = im.SharedSize
			s.Exclusive = im.Size - im.SharedSize
			if s.Exclusive < 0 {
				s.Exclusive = 0
			}
		}
		out[im.ID] = s
	}
	return out
}

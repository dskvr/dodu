package size

import "github.com/tyutyutyu/dodu/pkg/docker"

// ContainerSize splits per-container disk usage into writable layer + log file.
type ContainerSize struct {
	WritableLayer int64
	LogFile       int64
	Total         int64
}

// ContainerSizes returns per-container accounting indexed by container ID.
//
// logSizes maps container ID → log file size; missing entries count as zero.
func ContainerSizes(containers []docker.Container, logSizes map[string]int64) map[string]ContainerSize {
	out := make(map[string]ContainerSize, len(containers))
	for _, c := range containers {
		writable := c.SizeRw
		if writable < 0 {
			writable = 0
		}
		log := logSizes[c.ID]
		out[c.ID] = ContainerSize{
			WritableLayer: writable,
			LogFile:       log,
			Total:         writable + log,
		}
	}
	return out
}

package mock

import (
	"fmt"

	"github.com/tyutyutyu/dodu/pkg/docker"
)

// NewLargeClient builds a mock daemon with synthetic images, containers and
// volumes, useful for benchmarks and stress tests. All sizes are deterministic
// pseudo-values derived from the index.
func NewLargeClient(images, containers, volumes int) *Client {
	c := New()
	c.Daemon = docker.DaemonInfo{ID: "MOCK:LARGE", ServerVersion: "0.0.0-bench"}

	c.Images = make([]docker.Image, images)
	for i := 0; i < images; i++ {
		c.Images[i] = docker.Image{
			ID:         fmt.Sprintf("sha256:img%08d", i),
			RepoTags:   []string{fmt.Sprintf("repo/img%d:v1", i)},
			Size:       int64(1<<20) * int64(10+i%200),
			SharedSize: int64(1<<20) * int64(i%5),
			Containers: int64(i % 4),
		}
	}

	c.Containers = make([]docker.Container, containers)
	for i := 0; i < containers; i++ {
		state := "exited"
		if i%5 == 0 {
			state = "running"
		}
		c.Containers[i] = docker.Container{
			ID:         fmt.Sprintf("ctr%08d", i),
			Names:      []string{fmt.Sprintf("/c-%d", i)},
			Image:      fmt.Sprintf("repo/img%d:v1", i%images),
			ImageID:    fmt.Sprintf("sha256:img%08d", i%images),
			State:      state,
			SizeRw:     int64(1<<20) * int64(i%50),
			SizeRootFs: int64(1<<20) * int64(100+i%200),
		}
		c.LogSizes[c.Containers[i].ID] = int64(1<<10) * int64(i%1024)
	}

	c.Volumes = make([]docker.Volume, volumes)
	for i := 0; i < volumes; i++ {
		c.Volumes[i] = docker.Volume{
			Name:       fmt.Sprintf("vol%08d", i),
			Driver:     "local",
			Mountpoint: fmt.Sprintf("/var/lib/docker/volumes/vol%08d/_data", i),
			UsageBytes: int64(1<<20) * int64(i%500),
			RefCount:   int64(i % 3),
		}
	}
	return c
}

package docker

import (
	"context"
	"time"
)

// Client is the narrow surface dodu needs from the Docker daemon.
//
// All methods are safe for concurrent use by multiple goroutines unless noted.
// Implementations MUST translate transport / authn errors via the sentinel
// errors in errors.go where possible.
type Client interface {
	// Ping returns daemon info or an error if unreachable.
	Ping(ctx context.Context) (DaemonInfo, error)

	// ListImages returns all images with size accounting fields populated.
	ListImages(ctx context.Context) ([]Image, error)

	// ListContainers returns containers. If all is false, only running ones.
	// Returned items have SizeRw / SizeRootFs populated.
	ListContainers(ctx context.Context, all bool) ([]Container, error)

	// ListVolumes returns all volumes. UsageBytes / RefCount are populated
	// when the daemon reports them; otherwise -1.
	ListVolumes(ctx context.Context) ([]Volume, error)

	// BuildCacheUsage returns build cache records analogous to
	// `docker system df -v` build cache section.
	BuildCacheUsage(ctx context.Context) ([]BuildCacheEntry, error)

	// DiskUsage returns the aggregate usage breakdown.
	DiskUsage(ctx context.Context) (DiskUsage, error)

	// LogFileSize returns the size of the container's JSON log file in bytes.
	// Returns 0 (no error) if the log path is empty or unreadable.
	LogFileSize(ctx context.Context, containerID string) (int64, error)

	// PruneImages removes dangling / unused images per filters.
	PruneImages(ctx context.Context, f PruneFilters) (PruneReport, error)
	// PruneContainers removes stopped containers per filters.
	PruneContainers(ctx context.Context, f PruneFilters) (PruneReport, error)
	// PruneVolumes removes unused volumes per filters.
	PruneVolumes(ctx context.Context, f PruneFilters) (PruneReport, error)
	// PruneBuildCache removes build cache entries per filters.
	PruneBuildCache(ctx context.Context, f PruneFilters) (PruneReport, error)

	// RemoveImage deletes a single image by ID. force allows removal of images
	// referenced by stopped containers; pruneChildren removes untagged parents.
	RemoveImage(ctx context.Context, id string, force, pruneChildren bool) (int64, error)
	// RemoveContainer deletes a single container by ID. force kills running ones
	// (callers should not pass true unless the user explicitly opted in).
	RemoveContainer(ctx context.Context, id string, force, removeVolumes bool) error
	// RemoveVolume deletes a single volume by name. force is required if the
	// daemon thinks the volume is in use.
	RemoveVolume(ctx context.Context, name string, force bool) error

	// Close releases any underlying resources.
	Close() error
}

// Option configures a Client constructor.
type Option func(*options)

type options struct {
	host    string
	timeout time.Duration
}

// WithHost overrides the daemon endpoint (default: env / unix socket).
func WithHost(host string) Option {
	return func(o *options) { o.host = host }
}

// WithTimeout sets the request timeout for SDK calls.
func WithTimeout(d time.Duration) Option {
	return func(o *options) { o.timeout = d }
}

func defaultOptions() options {
	return options{
		timeout: 30 * time.Second,
	}
}

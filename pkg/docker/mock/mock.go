// Package mock provides a Client implementation for unit tests.
package mock

import (
	"context"
	"errors"

	"github.com/tyutyutyu/dodu/pkg/docker"
)

// Client is a programmable mock implementing docker.Client.
//
// Each method consults its corresponding `*Func` field; if nil, a zero value
// is returned with no error. Tests may also use the canned slices for the
// simple happy-path case.
type Client struct {
	PingFunc            func(ctx context.Context) (docker.DaemonInfo, error)
	ListImagesFunc      func(ctx context.Context) ([]docker.Image, error)
	ListContainersFunc  func(ctx context.Context, all bool) ([]docker.Container, error)
	ListVolumesFunc     func(ctx context.Context) ([]docker.Volume, error)
	BuildCacheFunc      func(ctx context.Context) ([]docker.BuildCacheEntry, error)
	DiskUsageFunc       func(ctx context.Context) (docker.DiskUsage, error)
	LogFileSizeFunc     func(ctx context.Context, id string) (int64, error)
	PruneImagesFunc     func(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error)
	PruneContainersFunc func(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error)
	PruneVolumesFunc    func(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error)
	PruneBuildCacheFunc func(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error)
	CloseFunc           func() error

	// Canned data used when the corresponding *Func is nil.
	Daemon     docker.DaemonInfo
	Images     []docker.Image
	Containers []docker.Container
	Volumes    []docker.Volume
	BuildCache []docker.BuildCacheEntry
	LogSizes   map[string]int64

	// CallCount tracks how many times each method has been invoked.
	CallCount map[string]int
}

// New returns a mock Client with empty defaults.
func New() *Client {
	return &Client{CallCount: map[string]int{}, LogSizes: map[string]int64{}}
}

func (c *Client) bump(name string) {
	if c.CallCount == nil {
		c.CallCount = map[string]int{}
	}
	c.CallCount[name]++
}

func (c *Client) Ping(ctx context.Context) (docker.DaemonInfo, error) {
	c.bump("Ping")
	if c.PingFunc != nil {
		return c.PingFunc(ctx)
	}
	return c.Daemon, nil
}

func (c *Client) ListImages(ctx context.Context) ([]docker.Image, error) {
	c.bump("ListImages")
	if c.ListImagesFunc != nil {
		return c.ListImagesFunc(ctx)
	}
	return c.Images, nil
}

func (c *Client) ListContainers(ctx context.Context, all bool) ([]docker.Container, error) {
	c.bump("ListContainers")
	if c.ListContainersFunc != nil {
		return c.ListContainersFunc(ctx, all)
	}
	return c.Containers, nil
}

func (c *Client) ListVolumes(ctx context.Context) ([]docker.Volume, error) {
	c.bump("ListVolumes")
	if c.ListVolumesFunc != nil {
		return c.ListVolumesFunc(ctx)
	}
	return c.Volumes, nil
}

func (c *Client) BuildCacheUsage(ctx context.Context) ([]docker.BuildCacheEntry, error) {
	c.bump("BuildCacheUsage")
	if c.BuildCacheFunc != nil {
		return c.BuildCacheFunc(ctx)
	}
	return c.BuildCache, nil
}

func (c *Client) DiskUsage(ctx context.Context) (docker.DiskUsage, error) {
	c.bump("DiskUsage")
	if c.DiskUsageFunc != nil {
		return c.DiskUsageFunc(ctx)
	}
	return docker.DiskUsage{
		Images:     c.Images,
		Containers: c.Containers,
		Volumes:    c.Volumes,
		BuildCache: c.BuildCache,
	}, nil
}

func (c *Client) LogFileSize(ctx context.Context, id string) (int64, error) {
	c.bump("LogFileSize")
	if c.LogFileSizeFunc != nil {
		return c.LogFileSizeFunc(ctx, id)
	}
	return c.LogSizes[id], nil
}

func (c *Client) PruneImages(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error) {
	c.bump("PruneImages")
	if c.PruneImagesFunc != nil {
		return c.PruneImagesFunc(ctx, f)
	}
	return docker.PruneReport{}, nil
}

func (c *Client) PruneContainers(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error) {
	c.bump("PruneContainers")
	if c.PruneContainersFunc != nil {
		return c.PruneContainersFunc(ctx, f)
	}
	return docker.PruneReport{}, nil
}

func (c *Client) PruneVolumes(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error) {
	c.bump("PruneVolumes")
	if c.PruneVolumesFunc != nil {
		return c.PruneVolumesFunc(ctx, f)
	}
	return docker.PruneReport{}, nil
}

func (c *Client) PruneBuildCache(ctx context.Context, f docker.PruneFilters) (docker.PruneReport, error) {
	c.bump("PruneBuildCache")
	if c.PruneBuildCacheFunc != nil {
		return c.PruneBuildCacheFunc(ctx, f)
	}
	return docker.PruneReport{}, nil
}

func (c *Client) Close() error {
	c.bump("Close")
	if c.CloseFunc != nil {
		return c.CloseFunc()
	}
	return nil
}

// ErrCanned is a convenience error tests can return from *Func fields.
var ErrCanned = errors.New("mock: canned error")

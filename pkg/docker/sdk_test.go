package docker

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

func testSDK(t *testing.T, handler http.HandlerFunc) *sdkClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	cli, err := dockerclient.NewClientWithOpts(dockerclient.WithHost(server.URL), dockerclient.WithVersion("1.47"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cli.Close() })
	return &sdkClient{cli: cli}
}

func TestListVolumesUsesDaemonDiskUsage(t *testing.T) {
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1.47/system/df" || r.URL.Query().Get("type") != "volume" {
			t.Errorf("unexpected request %s", r.URL.String())
			http.Error(w, "wrong endpoint", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"Volumes":[{"Name":"data","Driver":"local","UsageData":{"Size":12345,"RefCount":2}},{"Name":"external","Driver":"plugin"}]}`)
	})
	volumes, err := c.ListVolumes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(volumes) != 2 || volumes[0].UsageBytes != 12345 || volumes[0].RefCount != 2 || volumes[1].UsageBytes != -1 {
		t.Fatalf("volumes = %+v", volumes)
	}
}

func TestRemoteLogPathDoesNotReadLocalFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.json")
	if err := os.WriteFile(path, []byte("unrelated local data"), 0o600); err != nil {
		t.Fatal(err)
	}
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"Id":"container","LogPath":%q}`, path)
	})
	size, err := c.LogFileSize(context.Background(), "container")
	if err != nil {
		t.Fatal(err)
	}
	if size != 0 {
		t.Fatalf("remote log returned unrelated local size %d", size)
	}
}

func TestBuildCachePruneScopesExactIDs(t *testing.T) {
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1.47/build/prune" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
		}
		if r.URL.Query().Get("all") != "1" {
			t.Errorf("selected cache prune must include all record types: %s", r.URL)
		}
		args, err := filters.FromJSON(r.URL.Query().Get("filters"))
		if err != nil {
			t.Fatal(err)
		}
		ids := args.Get("id")
		if len(ids) != 1 {
			t.Fatalf("id filter = %v", ids)
		}
		re, err := regexp.Compile(ids[0])
		if err != nil {
			t.Fatal(err)
		}
		if !re.MatchString("cache.1") || !re.MatchString("cache2") || re.MatchString("cacheX1") || re.MatchString("prefix-cache2") {
			t.Errorf("unsafe ID filter %q", ids[0])
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"CachesDeleted":["cache.1"],"SpaceReclaimed":10}`)
	})
	report, err := c.PruneBuildCache(context.Background(), PruneFilters{IDs: []string{"cache.1", "cache2"}})
	if err != nil || report.SpaceReclaimed != 10 {
		t.Fatalf("prune = %+v, %v", report, err)
	}
}

func TestPruneRejectsUnsupportedIDFiltersWithoutRequest(t *testing.T) {
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unsafe prune request %s", r.URL)
		http.Error(w, "unexpected", 500)
	})
	for _, prune := range []func(context.Context, PruneFilters) (PruneReport, error){c.PruneImages, c.PruneContainers, c.PruneVolumes} {
		if _, err := prune(context.Background(), PruneFilters{IDs: []string{"selected"}}); err == nil {
			t.Error("unsupported IDs accepted")
		}
	}
	if _, err := c.PruneBuildCache(context.Background(), PruneFilters{IDs: []string{""}}); err == nil {
		t.Error("empty ID accepted")
	}
}

func TestDiskUsagePreservesGroupingMetadata(t *testing.T) {
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"LayersSize":42,"Images":[{"Id":"image","RepoDigests":["app@sha256:digest"]}],"Containers":[{"Id":"container","Mounts":[{"Type":"volume","Name":"data","Source":"/arbitrary/storage/data","Destination":"/data","RW":true}]}]}`)
	})
	usage, err := c.DiskUsage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if usage.LayersSize != 42 || len(usage.Images) != 1 || len(usage.Images[0].RepoDigests) != 1 || len(usage.Containers) != 1 || len(usage.Containers[0].Mounts) != 1 || usage.Containers[0].Mounts[0].Name != "data" {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestUnscopedBuildCachePrunePreservesDefaultRetention(t *testing.T) {
	c := testSDK(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("all") == "1" {
			t.Error("unscoped prune enabled all")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"CachesDeleted":[],"SpaceReclaimed":0}`)
	})
	if _, err := c.PruneBuildCache(context.Background(), PruneFilters{}); err != nil {
		t.Fatal(err)
	}
}

package cache_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tyutyutyu/dodu/pkg/cache"
	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/scan"
)

func tempDB(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshots.db")
}

func TestBoltRoundTrip(t *testing.T) {
	c, err := cache.OpenBolt(tempDB(t))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	snap := &scan.Snapshot{
		Daemon:     docker.DaemonInfo{ID: "d1", ServerVersion: "26.0"},
		Images:     []docker.Image{{ID: "img", RepoTags: []string{"foo:1"}, Size: 100, SharedSize: 10}},
		Containers: []docker.Container{{ID: "c1", Names: []string{"/n"}, ImageID: "img", State: "running", SizeRw: 5}},
		Volumes:    []docker.Volume{{Name: "v", UsageBytes: 1024, RefCount: 1}},
		BuildCache: []docker.BuildCacheEntry{{ID: "bc", Size: 64, InUse: true}},
		LogSizes:   map[string]int64{"c1": 32},
		CapturedAt: time.Now(),
		Duration:   time.Second,
		Errors:     []error{errors.New("partial")},
	}

	key := cache.Key(snap.Daemon)
	if err := c.Save(key, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := c.Load(key)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Daemon.ID != "d1" || len(got.Images) != 1 || got.Images[0].ID != "img" {
		t.Errorf("snapshot mismatch: %+v", got)
	}
	if got.LogSizes["c1"] != 32 {
		t.Errorf("log size lost: %d", got.LogSizes["c1"])
	}
	if len(got.Errors) != 1 || got.Errors[0].Error() != "partial" {
		t.Errorf("errors lost: %v", got.Errors)
	}
}

func TestBoltLoadMiss(t *testing.T) {
	c, _ := cache.OpenBolt(tempDB(t))
	t.Cleanup(func() { _ = c.Close() })
	_, err := c.Load("nope")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestBoltCorruptEntry(t *testing.T) {
	dbPath := tempDB(t)
	c, _ := cache.OpenBolt(dbPath)
	// Inject a corrupt blob via Save then mutate file? Simpler: Save valid
	// then corrupt by writing junk through a fresh handle.
	snap := &scan.Snapshot{Daemon: docker.DaemonInfo{ID: "d"}, CapturedAt: time.Now()}
	if err := c.Save("k", snap); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()

	// Truncate file to corrupt it.
	if err := os.WriteFile(dbPath, []byte("not a bolt db"), 0o600); err != nil {
		t.Fatal(err)
	}

	c2, err := cache.OpenBolt(dbPath)
	if err == nil {
		// bolt may reject; either way the cache must be unusable, not crash.
		_, err2 := c2.Load("k")
		if err2 == nil {
			t.Errorf("expected load failure on corrupt db")
		}
		_ = c2.Close()
	}
}

func TestPurge(t *testing.T) {
	c, _ := cache.OpenBolt(tempDB(t))
	t.Cleanup(func() { _ = c.Close() })
	snap := &scan.Snapshot{Daemon: docker.DaemonInfo{ID: "d"}, CapturedAt: time.Now()}
	_ = c.Save("k", snap)
	if err := c.Purge(); err != nil {
		t.Fatalf("purge: %v", err)
	}
	if _, err := c.Load("k"); !errors.Is(err, cache.ErrNotFound) {
		t.Errorf("after purge, want ErrNotFound, got %v", err)
	}
}

func TestPolicyFresh(t *testing.T) {
	now := time.Now()
	p := cache.Policy{TTL: time.Minute}
	if !p.Fresh(&scan.Snapshot{CapturedAt: now.Add(-30 * time.Second)}, now) {
		t.Error("30s old should be fresh")
	}
	if p.Fresh(&scan.Snapshot{CapturedAt: now.Add(-2 * time.Minute)}, now) {
		t.Error("2min old should be stale")
	}
	if p.Fresh(nil, now) {
		t.Error("nil snapshot should not be fresh")
	}
	if p.Fresh(&scan.Snapshot{}, now) {
		t.Error("zero capture time should not be fresh")
	}
}

func TestKeyDeterministic(t *testing.T) {
	a := cache.Key(docker.DaemonInfo{ID: "x", ServerVersion: "26"})
	b := cache.Key(docker.DaemonInfo{ID: "x", ServerVersion: "26"})
	if a != b {
		t.Error("Key not deterministic")
	}
	c := cache.Key(docker.DaemonInfo{ID: "y", ServerVersion: "26"})
	if a == c {
		t.Error("Key not differentiating daemons")
	}
}

package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/export"
	"github.com/tyutyutyu/dodu/pkg/scan"
)

func sampleSnap() *scan.Snapshot {
	return &scan.Snapshot{
		Daemon: docker.DaemonInfo{ID: "d1", ServerVersion: "26.0"},
		Images: []docker.Image{
			{ID: "img1", RepoTags: []string{"foo:1"}, Size: 1000, SharedSize: 100},
		},
		Containers: []docker.Container{
			{ID: "c1", Names: []string{"/a"}, Image: "foo:1", ImageID: "img1", State: "running", SizeRw: 50},
		},
		Volumes: []docker.Volume{
			{Name: "v1", Driver: "local", UsageBytes: 4096, RefCount: 1},
		},
		BuildCache: []docker.BuildCacheEntry{
			{ID: "bc1", Type: "regular", Size: 64, InUse: false},
		},
		LogSizes:   map[string]int64{"c1": 32},
		CapturedAt: time.Unix(1700000000, 0),
		Duration:   500 * time.Millisecond,
	}
}

func TestToJSONRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	if err := export.ToJSON(&buf, sampleSnap(), export.ToolInfo{Name: "dodu", Version: "test"}); err != nil {
		t.Fatalf("ToJSON: %v", err)
	}
	var doc export.Document
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if doc.SchemaVersion != export.SchemaVersion {
		t.Errorf("schema version = %q", doc.SchemaVersion)
	}
	if len(doc.Snapshot.Images) != 1 || doc.Snapshot.Images[0].ID != "img1" {
		t.Errorf("images lost: %+v", doc.Snapshot.Images)
	}
	if doc.Snapshot.LogSizes["c1"] != 32 {
		t.Errorf("log size lost")
	}
	if doc.Totals.Images == 0 {
		t.Error("totals not computed")
	}
}

func TestToJSONNilSnapshot(t *testing.T) {
	if err := export.ToJSON(&bytes.Buffer{}, nil, export.ToolInfo{}); err == nil {
		t.Error("expected error on nil snapshot")
	}
}

func TestToCSVImages(t *testing.T) {
	var buf bytes.Buffer
	if err := export.ToCSV(&buf, sampleSnap(), export.CSVImages); err != nil {
		t.Fatalf("ToCSV: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "id,repo_tags,size_bytes") {
		t.Errorf("missing header: %s", out)
	}
	if !strings.Contains(out, "img1") || !strings.Contains(out, "1000") {
		t.Errorf("missing data: %s", out)
	}
}

func TestToCSVAllKinds(t *testing.T) {
	for _, k := range export.AllCSVKinds {
		var buf bytes.Buffer
		if err := export.ToCSV(&buf, sampleSnap(), k); err != nil {
			t.Errorf("kind %s: %v", k, err)
		}
		if buf.Len() == 0 {
			t.Errorf("kind %s: empty output", k)
		}
	}
}

func TestToCSVUnknownKind(t *testing.T) {
	if err := export.ToCSV(&bytes.Buffer{}, sampleSnap(), export.CSVKind("nope")); err == nil {
		t.Error("expected error on unknown kind")
	}
}

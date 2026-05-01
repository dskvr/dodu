package plan_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/docker/mock"
	"github.com/tyutyutyu/dodu/pkg/group"
	"github.com/tyutyutyu/dodu/pkg/plan"
	"github.com/tyutyutyu/dodu/pkg/scan"
)

func snap() *scan.Snapshot {
	return &scan.Snapshot{
		Images: []docker.Image{
			{ID: "img-running", RepoTags: []string{"app:1"}, Size: 1000, SharedSize: 100, Containers: 1},
			{ID: "img-free", RepoTags: nil, Size: 500},
		},
		Containers: []docker.Container{
			{
				ID: "c-run", Names: []string{"/r"}, ImageID: "img-running", State: "running", SizeRw: 50,
				Mounts: []docker.ContainerMount{{Type: "volume", Name: "v-used"}},
			},
			{ID: "c-stop", Names: []string{"/s"}, ImageID: "img-free", State: "exited", SizeRw: 30},
		},
		Volumes: []docker.Volume{
			{Name: "v-used", UsageBytes: 4096, RefCount: 1},
			{Name: "v-orphan", UsageBytes: 2048, RefCount: 0},
		},
		BuildCache: []docker.BuildCacheEntry{
			{ID: "bc-used", Size: 100, InUse: true},
			{ID: "bc-free", Size: 200, InUse: false},
		},
		LogSizes: map[string]int64{"c-stop": 10},
	}
}

func TestBuildBlocksRunningImage(t *testing.T) {
	p := plan.Build(snap(), []plan.Mark{
		{Kind: group.KindImage, ID: "img-running"},
	})
	if len(p.Items) != 0 {
		t.Errorf("running-image must be blocked, got items %+v", p.Items)
	}
	if len(p.Blocked) != 1 {
		t.Fatalf("expected 1 blocked, got %d", len(p.Blocked))
	}
}

func TestBuildAllowsFreeImage(t *testing.T) {
	p := plan.Build(snap(), []plan.Mark{
		{Kind: group.KindImage, ID: "img-free"},
	})
	if len(p.Items) != 1 || p.Items[0].ID != "img-free" {
		t.Fatalf("expected 1 item img-free, got %+v", p.Items)
	}
	if p.EstReclaim != 500 {
		t.Errorf("est reclaim = %d", p.EstReclaim)
	}
}

func TestBuildBlocksRunningContainerAndUsedVolume(t *testing.T) {
	p := plan.Build(snap(), []plan.Mark{
		{Kind: group.KindContainer, ID: "c-run"},
		{Kind: group.KindVolume, ID: "v-used"},
		{Kind: group.KindVolume, ID: "v-orphan"},
		{Kind: group.KindBuildCache, ID: "bc-used"},
		{Kind: group.KindBuildCache, ID: "bc-free"},
		{Kind: group.KindContainer, ID: "c-stop"},
	})
	// Allowed: c-stop, v-orphan, bc-free
	if len(p.Items) != 3 {
		t.Errorf("items = %d, want 3 (%+v)", len(p.Items), p.Items)
	}
	// Blocked: c-run, v-used, bc-used
	if len(p.Blocked) != 3 {
		t.Errorf("blocked = %d, want 3", len(p.Blocked))
	}
}

func TestExecuteHappy(t *testing.T) {
	m := mock.New()
	p := &plan.Plan{
		Items: []plan.Item{
			{Kind: group.KindContainer, ID: "c1", EstReclaim: 100},
			{Kind: group.KindImage, ID: "img1", EstReclaim: 500},
			{Kind: group.KindVolume, ID: "v1", EstReclaim: 1000},
		},
		EstReclaim: 1600,
	}
	auditPath := filepath.Join(t.TempDir(), "audit.log")
	rep, err := p.Execute(context.Background(), m, plan.ExecOptions{AuditLog: auditPath})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if rep.ActualReclaim != 1600 {
		t.Errorf("reclaim = %d, want 1600", rep.ActualReclaim)
	}
	if m.CallCount["RemoveImage"] != 1 || m.CallCount["RemoveContainer"] != 1 || m.CallCount["RemoveVolume"] != 1 {
		t.Errorf("call counts: %+v", m.CallCount)
	}
	data, err := os.ReadFile(auditPath)
	if err != nil {
		t.Fatalf("read audit: %v", err)
	}
	if len(data) == 0 {
		t.Error("audit log empty")
	}
}

func TestExecuteReadOnly(t *testing.T) {
	p := &plan.Plan{Items: []plan.Item{{Kind: group.KindImage, ID: "x"}}}
	_, err := p.Execute(context.Background(), mock.New(), plan.ExecOptions{ReadOnly: true})
	if !errors.Is(err, plan.ErrReadOnly) {
		t.Errorf("want ErrReadOnly, got %v", err)
	}
}

func TestExecuteRecordsFailure(t *testing.T) {
	m := mock.New()
	m.RemoveImageFunc = func(_ context.Context, _ string, _, _ bool) (int64, error) {
		return 0, errors.New("boom")
	}
	p := &plan.Plan{Items: []plan.Item{{Kind: group.KindImage, ID: "x", EstReclaim: 99}}}
	rep, err := p.Execute(context.Background(), m, plan.ExecOptions{})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(rep.Results) != 1 || rep.Results[0].Success {
		t.Errorf("expected one failed result, got %+v", rep.Results)
	}
	if rep.ActualReclaim != 0 {
		t.Errorf("reclaim must be 0 on failure, got %d", rep.ActualReclaim)
	}
}

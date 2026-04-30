package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/docker/mock"
)

// Compile-time assertion: mock.Client satisfies docker.Client.
var _ docker.Client = (*mock.Client)(nil)

func TestMockReturnsCannedData(t *testing.T) {
	m := mock.New()
	m.Images = []docker.Image{{ID: "img-1", Size: 100}}
	m.Containers = []docker.Container{{ID: "c-1", SizeRw: 10}}

	imgs, err := m.ListImages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(imgs) != 1 || imgs[0].ID != "img-1" {
		t.Errorf("unexpected images: %+v", imgs)
	}

	if c := m.CallCount["ListImages"]; c != 1 {
		t.Errorf("expected 1 ListImages call, got %d", c)
	}
}

func TestMockFuncOverride(t *testing.T) {
	m := mock.New()
	m.ListVolumesFunc = func(_ context.Context) ([]docker.Volume, error) {
		return nil, mock.ErrCanned
	}
	_, err := m.ListVolumes(context.Background())
	if !errors.Is(err, mock.ErrCanned) {
		t.Fatalf("expected ErrCanned, got %v", err)
	}
}

func TestMockLogFileSizeMap(t *testing.T) {
	m := mock.New()
	m.LogSizes["abc"] = 4096
	got, err := m.LogFileSize(context.Background(), "abc")
	if err != nil || got != 4096 {
		t.Fatalf("expected 4096, got %d (err=%v)", got, err)
	}
	missing, _ := m.LogFileSize(context.Background(), "missing")
	if missing != 0 {
		t.Errorf("expected 0 for missing id, got %d", missing)
	}
}

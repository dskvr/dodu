//go:build perf

package perf

import (
	"context"
	"log/slog"
	"testing"

	"github.com/tyutyutyu/dodu/pkg/docker/mock"
	"github.com/tyutyutyu/dodu/pkg/scan"
)

// BenchmarkScanLarge measures Scanner throughput against a synthetic mock with
// 2k images, 500 containers, 1k volumes. Run with:
//
//	go test -tags perf -bench=. -benchmem ./test/perf/...
func BenchmarkScanLarge(b *testing.B) {
	c := mock.NewLargeClient(2000, 500, 1000)
	s := scan.New(c, slog.Default())
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Scan(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

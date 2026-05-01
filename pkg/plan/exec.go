package plan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/group"
)

// ErrReadOnly is returned when DODU_READONLY=1 blocks an execute.
var ErrReadOnly = errors.New("plan: execution disabled (DODU_READONLY=1 or readonly option)")

// Report is the outcome of executing a Plan.
type Report struct {
	Results       []ItemResult
	ActualReclaim int64
	StartedAt     time.Time
	FinishedAt    time.Time
}

// ItemResult records what happened to a single planned item.
type ItemResult struct {
	Item    Item
	Success bool
	Err     string
	Reclaim int64
}

// ExecOptions tunes Execute.
type ExecOptions struct {
	// ReadOnly aborts immediately. Defaults to true if env DODU_READONLY=1.
	ReadOnly bool
	// AuditLog is opened in append mode for JSON-Lines records. Empty disables.
	AuditLog string
	// Force passes force=true to remove calls. The TUI/CLI must guard this.
	Force bool
}

// Execute carries out the Plan via the docker.Client. Blocked items are never
// touched. All actions are appended to the audit log when configured.
func (p *Plan) Execute(ctx context.Context, client docker.Client, opts ExecOptions) (Report, error) {
	if p == nil {
		return Report{}, errors.New("plan: nil")
	}
	if opts.ReadOnly || os.Getenv("DODU_READONLY") == "1" {
		return Report{}, ErrReadOnly
	}
	rep := Report{StartedAt: time.Now()}

	var (
		audit io.WriteCloser
		amu   sync.Mutex
	)
	if opts.AuditLog != "" {
		if err := os.MkdirAll(filepath.Dir(opts.AuditLog), 0o750); err == nil {
			f, err := os.OpenFile(opts.AuditLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
			if err == nil {
				audit = f
				defer func() { _ = f.Close() }()
			}
		}
	}
	writeAudit := func(rec map[string]any) {
		if audit == nil {
			return
		}
		amu.Lock()
		defer amu.Unlock()
		_ = json.NewEncoder(audit).Encode(rec)
	}

	for _, it := range p.Items {
		select {
		case <-ctx.Done():
			return rep, ctx.Err()
		default:
		}
		res := ItemResult{Item: it}
		var err error
		switch it.Kind {
		case group.KindImage:
			_, err = client.RemoveImage(ctx, it.ID, opts.Force, true)
		case group.KindContainer:
			err = client.RemoveContainer(ctx, it.ID, opts.Force, false)
		case group.KindVolume:
			err = client.RemoveVolume(ctx, it.ID, opts.Force)
		case group.KindBuildCache:
			// No per-ID API in the daemon; use prune as fallback for this kind.
			report, perr := client.PruneBuildCache(ctx, docker.PruneFilters{})
			err = perr
			if err == nil {
				res.Reclaim = report.SpaceReclaimed
			}
		default:
			err = fmt.Errorf("unsupported kind %s", it.Kind)
		}
		if err != nil {
			res.Err = err.Error()
		} else {
			res.Success = true
			if res.Reclaim == 0 {
				res.Reclaim = it.EstReclaim
			}
			rep.ActualReclaim += res.Reclaim
		}
		rep.Results = append(rep.Results, res)
		writeAudit(map[string]any{
			"ts":       time.Now().UTC().Format(time.RFC3339Nano),
			"kind":     string(it.Kind),
			"id":       it.ID,
			"name":     it.Name,
			"success":  res.Success,
			"err":      res.Err,
			"reclaim":  res.Reclaim,
			"estimate": it.EstReclaim,
		})
	}

	rep.FinishedAt = time.Now()
	return rep, nil
}

// DefaultAuditLogPath returns $XDG_STATE_HOME/dodu/audit.log or
// $HOME/.local/state/dodu/audit.log.
func DefaultAuditLogPath() (string, error) {
	if p := os.Getenv("XDG_STATE_HOME"); p != "" {
		return filepath.Join(p, "dodu", "audit.log"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state", "dodu", "audit.log"), nil
}

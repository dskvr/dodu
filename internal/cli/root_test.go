package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	var buf bytes.Buffer
	info := BuildInfo{Version: "1.2.3", Commit: "abc", Date: "2026-01-01"}
	if err := printVersion(&buf, info); err != nil {
		t.Fatalf("printVersion returned error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "dodu 1.2.3") {
		t.Errorf("expected version in output, got %q", got)
	}
	if !strings.Contains(got, "abc") || !strings.Contains(got, "2026-01-01") {
		t.Errorf("expected commit and date in output, got %q", got)
	}
}

func TestRootCommandHelp(t *testing.T) {
	root := newRootCmd(BuildInfo{Version: "test"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("help should not error: %v", err)
	}
	if !strings.Contains(buf.String(), "navigable TUI/CLI") {
		t.Errorf("help missing description: %q", buf.String())
	}
}

func TestVersionSubcommand(t *testing.T) {
	root := newRootCmd(BuildInfo{Version: "9.9.9", Commit: "deadbeef", Date: "today"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("version subcommand failed: %v", err)
	}
	if !strings.Contains(buf.String(), "9.9.9") {
		t.Errorf("expected version, got %q", buf.String())
	}
}

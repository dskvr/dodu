// Command dodu is the entry point for the dodu CLI/TUI binary.
package main

import (
	"os"

	"github.com/tyutyutyu/dodu/internal/cli"
)

// Build-time injected via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(cli.Execute(cli.BuildInfo{Version: version, Commit: commit, Date: date}))
}

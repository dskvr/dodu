// Command dodu is the entry point for the dodu CLI/TUI binary.
package main

import (
	"fmt"
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
	if err := cli.Execute(cli.BuildInfo{Version: version, Commit: commit, Date: date}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

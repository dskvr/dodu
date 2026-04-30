// Package cli wires up the dodu Cobra command tree.
//
// The skeleton command exposes only `version` and a placeholder root that will
// later launch the TUI. Subcommands (scan, prune, export, cache) are added in
// later phases per docs at .github/prompts/09-cli.prompt.md.
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// BuildInfo holds version metadata injected at build time.
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

// Execute runs the root command.
func Execute(info BuildInfo) error {
	root := newRootCmd(info)
	return root.Execute()
}

func newRootCmd(info BuildInfo) *cobra.Command {
	root := &cobra.Command{
		Use:   "dodu",
		Short: "Docker disk atlas — ncdu-style TUI for Docker storage",
		Long: "dodu is a fast, navigable TUI/CLI that shows what consumes disk in your\n" +
			"Docker environment (images, containers, volumes, build cache, logs) and\n" +
			"suggests safe, dry-run cleanup plans.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// TUI launch lands here in phase 06. For now, print a hint.
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "dodu: TUI not implemented yet — see `dodu --help`.")
			return err
		},
	}

	root.AddCommand(newVersionCmd(info))
	return root
}

func newVersionCmd(info BuildInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return printVersion(cmd.OutOrStdout(), info)
		},
	}
}

func printVersion(w io.Writer, info BuildInfo) error {
	_, err := fmt.Fprintf(w, "dodu %s (commit %s, built %s)\n", info.Version, info.Commit, info.Date)
	return err
}

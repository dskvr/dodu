package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tyutyutyu/dodu/pkg/cache"
)

func newCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage the snapshot cache",
	}
	cmd.AddCommand(newCacheInfoCmd(), newCacheClearCmd())
	return cmd
}

func newCacheInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print cache location and entries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := cache.DefaultPath()
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Path: %s\n", path)
			c, err := cache.OpenBolt(path)
			if err != nil {
				if errors.Is(err, cache.ErrNotFound) {
					return nil
				}
				return err
			}
			_ = c.Close()
			return nil
		},
	}
}

func newCacheClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "purge",
		Short: "Delete all cached snapshots",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path, err := cache.DefaultPath()
			if err != nil {
				return err
			}
			c, err := cache.OpenBolt(path)
			if err != nil {
				return err
			}
			defer func() { _ = c.Close() }()
			if err := c.Purge(); err != nil {
				return err
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Cache purged.")
			return nil
		},
	}
}

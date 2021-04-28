package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root rklotz command
func NewRootCmd(version string) *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:     "rklotz",
		Short:   "rKlotz is a simple one-user file-based blog engine",
		Version: version,
	}

	cmd.AddCommand(NewServerCmd(ctx, version))

	return cmd
}

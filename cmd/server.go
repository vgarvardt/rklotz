package cmd

import (
	"context"

	wErrors "github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/vgarvardt/rklotz/pkg/server"
)

// NewServerCmd creates new server command
func NewServerCmd(ctx context.Context, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Runs rKlotz server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := server.LoadConfig()
			if err != nil {
				return wErrors.Wrap(err, "failed to load config")
			}

			return server.Run(cfg, version)
		},
	}
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vgarvardt/rklotz/pkg/server"
)

// NewServerCmd creates new server command
func NewServerCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Runs rKlotz server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := server.LoadConfig(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			return server.Run(cmd.Context(), cfg, version)
		},
	}
}

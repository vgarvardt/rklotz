package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates a new version command
func NewVersionCmd(ctx context.Context, version string) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version information",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rklotz %s\n", version)
		},
	}
}

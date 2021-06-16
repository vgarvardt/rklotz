package cmd

import "github.com/spf13/cobra"

// NewRootCmd creates the root rklotz command
func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rklotz",
		Short:   "rKlotz is a simple one-user file-based blog engine",
		Version: version,
	}

	cmd.AddCommand(NewServerCmd(version))

	return cmd
}

package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// NewFreshRSSCommand creates the FreshRSS command.
func NewFreshRSSCommand(ctx context.Context, opts *Options) *cobra.Command {
	freshRSSCommand := &cobra.Command{
		Use:   "freshrss",
		Short: "FreshRSS service commands",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	freshRSSCommand.AddCommand(
		newFreshRSSCapabilitiesCommand(opts),
		newFreshRSSConfigCommand(opts),
		newFreshRSSAuthCommand(ctx, opts),
	)

	return freshRSSCommand
}

func newFreshRSSCapabilitiesCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "describe FreshRSS capabilities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return capabilities(cmd.OutOrStdout(), opts, "freshrss")
		},
	}
}

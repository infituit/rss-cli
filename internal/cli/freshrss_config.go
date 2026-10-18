package cli

import (
	"errors"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/spf13/cobra"
)

func newFreshRSSConfigCommand(opts *Options) *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "manage FreshRSS configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configCommand.AddCommand(newFreshRSSConfigInitCommand(opts), newFreshRSSConfigShowCommand(opts))
	return configCommand
}

func newFreshRSSConfigInitCommand(opts *Options) *cobra.Command {
	var (
		url   string
		force bool
	)

	initCommand := &cobra.Command{
		Use:   "init",
		Short: "configure the FreshRSS endpoint",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if url == "" {
				return errors.New("--url is required")
			}
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			if err := config.SaveService(name, "freshrss", config.ServiceConfig{Version: 1, Provider: "freshrss", BaseURL: url}, force); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"service": "freshrss", "status": "configured"})
		},
	}

	initCommand.Flags().StringVar(&url, "url", "", "FreshRSS base URL or Fever endpoint")
	initCommand.Flags().BoolVar(&force, "force", false, "overwrite existing service configuration")
	return initCommand
}

func newFreshRSSConfigShowCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show FreshRSS configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			cfg, err := config.LoadService(name, "freshrss")
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, cfg)
		},
	}
}

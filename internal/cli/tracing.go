package cli

import (
	"fmt"
	"net/url"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewTracingCommand creates the tracing command.
func NewTracingCommand(opts *Options) *cobra.Command {
	tracingCommand := &cobra.Command{
		Use:   "tracing",
		Short: "manage tracing",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	tracingCommand.AddCommand(newTracingConfigCommand(opts))
	return tracingCommand
}

func newTracingConfigCommand(opts *Options) *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "manage tracing configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configCommand.AddCommand(
		newTracingConfigInitCommand(opts),
		newTracingConfigShowCommand(opts),
		newTracingConfigDisableCommand(opts),
	)
	return configCommand
}

func newTracingConfigInitCommand(opts *Options) *cobra.Command {
	var endpoint string
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "configure the tracing endpoint",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := validateTracingEndpoint(endpoint); err != nil {
				return err
			}
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			contextConfig, err := config.LoadContext(name)
			if err != nil {
				return err
			}
			contextConfig.Tracing = &config.TracingConfig{Enabled: true, Endpoint: endpoint}
			if err := config.UpdateContext(name, contextConfig); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]any{
				"context":  name,
				"endpoint": endpoint,
				"enabled":  true,
				"status":   "configured",
			})
		},
	}
	initCommand.Flags().StringVar(&endpoint, "endpoint", "", "OTLP/HTTP traces endpoint")
	_ = initCommand.MarkFlagRequired("endpoint")
	return initCommand
}

func newTracingConfigShowCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show tracing configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			contextConfig, err := config.LoadContext(name)
			if err != nil {
				return err
			}
			tracingConfig := config.TracingConfig{}
			if contextConfig.Tracing != nil {
				tracingConfig = *contextConfig.Tracing
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]any{
				"context":  name,
				"endpoint": tracingConfig.Endpoint,
				"enabled":  tracingConfig.Enabled,
			})
		},
	}
}

func newTracingConfigDisableCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "disable tracing",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			contextConfig, err := config.LoadContext(name)
			if err != nil {
				return err
			}
			if contextConfig.Tracing == nil {
				contextConfig.Tracing = &config.TracingConfig{}
			} else {
				contextConfig.Tracing.Enabled = false
			}
			if err := config.UpdateContext(name, contextConfig); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]any{
				"context": name,
				"enabled": false,
				"status":  "disabled",
			})
		},
	}
}

func validateTracingEndpoint(endpoint string) error {
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("invalid tracing endpoint %q", endpoint)
	}
	return nil
}

package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/infituit/rss-cli/internal/output"
	"github.com/infituit/rss-cli/pkg/version"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command.
func NewRootCommand(ctx context.Context, stdout, stderr io.Writer) *cobra.Command {
	opts := &Options{Output: output.JSON}
	showVersion := false
	root := &cobra.Command{
		Use:           "rss-cli",
		Short:         "RSS service command line client",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if showVersion {
				return writeVersion(cmd, opts)
			}
			return cmd.Help()
		},
	}

	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetIn(os.Stdin)
	root.PersistentFlags().StringVar(&opts.Context, "context", "", "context name")
	root.PersistentFlags().Var(newOutputFormatValue(&opts.Output), "output", "output format: json or table")
	root.PersistentFlags().BoolVar(&showVersion, "version", false, "show build metadata")

	root.AddCommand(
		NewVersionCommand(opts),
		NewCapabilitiesCommand(opts),
		NewContextCommand(opts),
		NewFreshRSSCommand(ctx, opts),
	)

	return root
}

// NewVersionCommand creates the version command.
func NewVersionCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show build metadata",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeVersion(cmd, opts)
		},
	}
}

// NewCapabilitiesCommand creates the capabilities command.
func NewCapabilitiesCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "describe supported capabilities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return capabilities(cmd.OutOrStdout(), opts, "rss-cli")
		},
	}
}

func writeVersion(cmd *cobra.Command, opts *Options) error {
	return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{
		"version":    version.VERSION,
		"build_time": version.BUILDTIME,
		"git_hash":   version.GITHASH,
	})
}

type outputFormatValue struct {
	target *output.Format
}

func newOutputFormatValue(target *output.Format) *outputFormatValue {
	return &outputFormatValue{target: target}
}

func (v *outputFormatValue) String() string {
	return string(*v.target)
}

func (v *outputFormatValue) Set(value string) error {
	format := output.Format(value)
	if format != output.JSON && format != output.Table {
		return fmt.Errorf("unsupported output format %q", value)
	}
	*v.target = format
	return nil
}

func (v *outputFormatValue) Type() string {
	return "string"
}

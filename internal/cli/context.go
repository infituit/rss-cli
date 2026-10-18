package cli

import (
	"errors"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewContextCommand creates the context command.
func NewContextCommand(opts *Options) *cobra.Command {
	contextCommand := &cobra.Command{
		Use:   "context",
		Short: "manage contexts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	contextCommand.AddCommand(
		newCreateContextCommand(opts),
		newUseContextCommand(opts),
		newListContextCommand(opts),
		newShowContextCommand(opts),
		newDeleteContextCommand(opts),
	)

	return contextCommand
}

func newCreateContextCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "create NAME",
		Short: "create a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SaveContext(args[0], config.ContextConfig{Version: 1}); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"context": args[0], "status": "created"})
		},
	}
}

func newUseContextCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "use NAME",
		Short: "select the current context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetCurrent(args[0]); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"context": args[0], "status": "current"})
		},
	}
}

func newListContextCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list contexts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return listContexts(cmd.OutOrStdout(), opts.Output)
		},
	}
}

func newShowContextCommand(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show the selected context",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			return showContext(name, cmd.OutOrStdout(), opts.Output)
		},
	}
}

func newDeleteContextCommand(opts *Options) *cobra.Command {
	confirm := false
	deleteCommand := &cobra.Command{
		Use:   "delete NAME",
		Short: "delete a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return errors.New("usage: rss-cli context delete NAME --confirm")
			}
			if err := config.DeleteContext(args[0]); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"context": args[0], "status": "deleted"})
		},
	}
	deleteCommand.Flags().BoolVar(&confirm, "confirm", false, "confirm context deletion")
	return deleteCommand
}

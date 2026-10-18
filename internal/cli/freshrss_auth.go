package cli

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/credential"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/infituit/rss-cli/pkg/thirdparty/freshrss"
	"github.com/spf13/cobra"
)

func newFreshRSSAuthCommand(ctx context.Context, opts *Options) *cobra.Command {
	authCommand := &cobra.Command{
		Use:   "auth",
		Short: "manage FreshRSS credentials",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	authCommand.AddCommand(
		newFreshRSSAuthSetCommand(opts),
		newFreshRSSAuthCheckCommand(ctx, opts),
		newFreshRSSAuthClearCommand(opts),
	)
	return authCommand
}

func newFreshRSSAuthSetCommand(opts *Options) *cobra.Command {
	var (
		user             string
		apiPasswordStdin bool
		apiKeyStdin      bool
		force            bool
	)

	setCommand := &cobra.Command{
		Use:   "set",
		Short: "save FreshRSS credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			credentialsPath, err := config.CredentialsPath(name, "freshrss")
			if err != nil {
				return err
			}
			values := credential.Values{Username: user}
			if apiPasswordStdin {
				data, readErr := io.ReadAll(cmd.InOrStdin())
				if readErr != nil {
					return readErr
				}
				values.APIPassword = string(data)
			}
			if apiKeyStdin {
				data, readErr := io.ReadAll(cmd.InOrStdin())
				if readErr != nil {
					return readErr
				}
				values.APIKey = string(data)
			}
			if values.Username == "" && values.APIPassword == "" && values.APIKey == "" {
				return errors.New("--user and --api-password-stdin, or --api-key-stdin are required")
			}
			if err := credential.Save(credentialsPath, values, force); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"service": "freshrss", "status": "credentials_saved"})
		},
	}

	setCommand.Flags().StringVar(&user, "user", "", "FreshRSS username")
	setCommand.Flags().BoolVar(&apiPasswordStdin, "api-password-stdin", false, "read API password from stdin")
	setCommand.Flags().BoolVar(&apiKeyStdin, "api-key-stdin", false, "read API key from stdin")
	setCommand.Flags().BoolVar(&force, "force", false, "replace existing credentials")
	return setCommand
}

func newFreshRSSAuthCheckCommand(ctx context.Context, opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "check FreshRSS authentication",
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
			credentialsPath, err := config.CredentialsPath(name, "freshrss")
			if err != nil {
				return err
			}
			values, err := credential.Load(credentialsPath)
			if err != nil {
				return err
			}
			handler, err := freshrss.NewHandler(freshrss.Config{BaseURL: cfg.BaseURL, Username: values.Username, APIPassword: values.APIPassword, APIKey: values.APIKey}, http.DefaultClient)
			if err != nil {
				return err
			}
			response, err := handler.AuthCheck(ctx)
			if err != nil {
				return err
			}
			if !response.Authenticated {
				return errors.New("authentication failed")
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, response)
		},
	}
}

func newFreshRSSAuthClearCommand(opts *Options) *cobra.Command {
	confirm := false
	clearCommand := &cobra.Command{
		Use:   "clear",
		Short: "remove FreshRSS credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !confirm {
				return errors.New("usage: rss-cli freshrss auth clear --confirm")
			}
			name, err := config.ResolveContext(opts.Context)
			if err != nil {
				return err
			}
			credentialsPath, err := config.CredentialsPath(name, "freshrss")
			if err != nil {
				return err
			}
			if err := credential.Delete(credentialsPath); err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), opts.Output, map[string]string{"service": "freshrss", "status": "credentials_cleared"})
		},
	}
	clearCommand.Flags().BoolVar(&confirm, "confirm", false, "confirm credential removal")
	return clearCommand
}

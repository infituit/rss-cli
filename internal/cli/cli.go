package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/output"
)

// Options represents the global command options.
type Options struct {
	Context string
	Output  output.Format
}

// Run executes the CLI with the provided context.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	root := NewRootCommand(ctx, stdout, stderr)
	root.SetArgs(args)
	return root.Execute()
}

func listContexts(stdout io.Writer, format output.Format) error {
	names, err := config.ListContexts()
	if err != nil {
		return err
	}

	current, _ := config.Current()
	result := make([]map[string]any, 0, len(names))
	for _, name := range names {
		services, _ := config.ListServices(name)
		if services == nil {
			services = []string{}
		}
		result = append(result, map[string]any{"name": name, "current": name == current, "services": services})
	}

	return output.Write(stdout, format, result)
}

func showContext(name string, stdout io.Writer, format output.Format) error {
	services, err := config.ListServices(name)
	if err != nil {
		return err
	}

	result := map[string]any{"name": name, "services": []any{}}
	for _, service := range services {
		serviceCfg, serviceErr := config.LoadService(name, service)
		hasCreds, credErr := config.HasCredentials(name, service)
		if serviceErr != nil {
			return serviceErr
		}
		if credErr != nil {
			return credErr
		}
		result["services"] = append(result["services"].([]any), map[string]any{
			"name":                service,
			"provider":            serviceCfg.Provider,
			"base_url":            serviceCfg.BaseURL,
			"credentials_present": hasCreds,
		})
	}

	return output.Write(stdout, format, result)
}

func capabilities(stdout io.Writer, opts *Options, provider string) error {
	name, err := config.ResolveContext(opts.Context)
	configured := err == nil
	result := map[string]any{
		"provider": provider,
		"capabilities": []string{
			"config.init",
			"config.show",
			"auth.set",
			"auth.check",
			"auth.clear",
		},
		"context": map[string]any{"configured": configured},
	}
	if configured {
		has, hasErr := serviceConfigured(name, provider)
		if hasErr != nil {
			return hasErr
		}
		result["context"].(map[string]any)["name"] = name
		result["context"].(map[string]any)["service_configured"] = has
	}

	return output.Write(stdout, opts.Output, result)
}

func serviceConfigured(contextName, provider string) (bool, error) {
	_, err := config.LoadService(contextName, provider)
	if errors.Is(err, os.ErrNotExist) || (err != nil && strings.HasPrefix(err.Error(), "read config:")) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

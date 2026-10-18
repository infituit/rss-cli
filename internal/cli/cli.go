package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/credential"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/infituit/rss-cli/pkg/thirdparty/freshrss"
	"github.com/infituit/rss-cli/pkg/version"
)

type Options struct {
	Context string
	Output  output.Format
}

func Run(args []string, stdout, stderr io.Writer) error {
	return RunContext(context.Background(), args, stdout, stderr)
}

func RunContext(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	opts, rest, err := parseGlobal(args)
	if err != nil {
		return err
	}
	if len(rest) == 0 || rest[0] == "help" || rest[0] == "--help" || rest[0] == "-h" {
		printHelp(stdout)
		return nil
	}
	if rest[0] == "version" || rest[0] == "--version" {
		return output.Write(stdout, opts.Output, map[string]string{"version": version.VERSION, "build_time": version.BUILDTIME, "git_hash": version.GITHASH})
	}
	if rest[0] == "capabilities" {
		return capabilities(stdout, opts, "rss-cli")
	}
	switch rest[0] {
	case "context":
		return runContext(rest[1:], opts, stdout)
	case "freshrss":
		return runFreshRSS(ctx, rest[1:], opts, stdout)
	default:
		return fmt.Errorf("unknown command %q", rest[0])
	}
}

func parseGlobal(args []string) (Options, []string, error) {
	opts := Options{Output: output.JSON}
	rest := append([]string(nil), args...)
	for len(rest) > 0 {
		switch rest[0] {
		case "--context":
			if len(rest) < 2 {
				return opts, nil, errors.New("--context requires a value")
			}
			opts.Context, rest = rest[1], rest[2:]
		case "--output":
			if len(rest) < 2 {
				return opts, nil, errors.New("--output requires a value")
			}
			opts.Output = output.Format(rest[1])
			if opts.Output != output.JSON && opts.Output != output.Table {
				return opts, nil, fmt.Errorf("unsupported output format %q", opts.Output)
			}
			rest = rest[2:]
		default:
			return opts, rest, nil
		}
	}
	return opts, rest, nil
}

func runContext(args []string, opts Options, stdout io.Writer) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(stdout, "context commands: create, list, show, use, delete")
		return err
	}
	switch args[0] {
	case "create":
		if len(args) != 2 {
			return errors.New("usage: rss-cli context create NAME")
		}
		if err := config.SaveContext(args[1], config.ContextConfig{Version: 1}); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"context": args[1], "status": "created"})
	case "use":
		if len(args) != 2 {
			return errors.New("usage: rss-cli context use NAME")
		}
		if err := config.SetCurrent(args[1]); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"context": args[1], "status": "current"})
	case "list":
		return listContexts(stdout, opts.Output)
	case "show":
		name, err := config.ResolveContext(opts.Context)
		if err != nil {
			return err
		}
		return showContext(name, stdout, opts.Output)
	case "delete":
		if err := deleteContext(args[1:]); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"context": args[1], "status": "deleted"})
	default:
		return fmt.Errorf("unknown context command %q", args[0])
	}
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
		result["services"] = append(result["services"].([]any), map[string]any{"name": service, "provider": serviceCfg.Provider, "base_url": serviceCfg.BaseURL, "credentials_present": hasCreds})
	}
	return output.Write(stdout, format, result)
}
func deleteContext(args []string) error {
	if len(args) != 2 || args[0] == "" || args[1] != "--confirm" {
		return errors.New("usage: rss-cli context delete NAME --confirm")
	}
	if err := config.DeleteContext(args[0]); err != nil {
		return err
	}
	return nil
}

func runFreshRSS(ctx context.Context, args []string, opts Options, stdout io.Writer) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(stdout, "freshrss commands: config, auth, capabilities")
		return err
	}
	switch args[0] {
	case "capabilities":
		return capabilities(stdout, opts, "freshrss")
	case "config":
		return runFreshRSSConfig(args[1:], opts, stdout)
	case "auth":
		return runFreshRSSAuth(ctx, args[1:], opts, stdout)
	default:
		return fmt.Errorf("unknown freshrss command %q", args[0])
	}
}
func capabilities(stdout io.Writer, opts Options, provider string) error {
	name, err := config.ResolveContext(opts.Context)
	configured := err == nil
	result := map[string]any{"provider": provider, "capabilities": []string{"config.init", "config.show", "auth.set", "auth.check", "auth.clear"}, "context": map[string]any{"configured": configured}}
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

func runFreshRSSConfig(args []string, opts Options, stdout io.Writer) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(stdout, "freshrss config commands: init, show")
		return err
	}
	name, err := config.ResolveContext(opts.Context)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		fs := flag.NewFlagSet("freshrss config init", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		url := fs.String("url", "", "FreshRSS base URL or Fever endpoint")
		force := fs.Bool("force", false, "overwrite existing service configuration")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *url == "" {
			return errors.New("--url is required")
		}
		if err := config.SaveService(name, "freshrss", config.ServiceConfig{Version: 1, Provider: "freshrss", BaseURL: *url}, *force); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"service": "freshrss", "status": "configured"})
	case "show":
		cfg, err := config.LoadService(name, "freshrss")
		if err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, cfg)
	default:
		return fmt.Errorf("unknown freshrss config command %q", args[0])
	}
}
func runFreshRSSAuth(ctx context.Context, args []string, opts Options, stdout io.Writer) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(stdout, "freshrss auth commands: set, check, clear")
		return err
	}
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
	switch args[0] {
	case "set":
		fs := flag.NewFlagSet("freshrss auth set", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		user := fs.String("user", "", "FreshRSS username")
		apiPassword := fs.String("api-password-stdin", "", "read API password from stdin when set to true")
		apiKey := fs.String("api-key-stdin", "", "read API key from stdin when set to true")
		force := fs.Bool("force", false, "replace existing credentials")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		values := credential.Values{Username: *user}
		if *apiPassword == "true" {
			data, readErr := io.ReadAll(os.Stdin)
			if readErr != nil {
				return readErr
			}
			values.APIPassword = string(data)
		}
		if *apiKey == "true" {
			data, readErr := io.ReadAll(os.Stdin)
			if readErr != nil {
				return readErr
			}
			values.APIKey = string(data)
		}
		if values.Username == "" && values.APIKey == "" {
			return errors.New("--user and --api-password-stdin=true, or --api-key-stdin=true are required")
		}
		if err := credential.Save(credentialsPath, values, *force); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"service": "freshrss", "status": "credentials_saved"})
	case "check":
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
		return output.Write(stdout, opts.Output, response)
	case "clear":
		if len(args) != 2 || args[1] != "--confirm" {
			return errors.New("usage: rss-cli freshrss auth clear --confirm")
		}
		if err := credential.Delete(credentialsPath); err != nil {
			return err
		}
		return output.Write(stdout, opts.Output, map[string]string{"service": "freshrss", "status": "credentials_cleared"})
	default:
		return fmt.Errorf("unknown freshrss auth command %q", args[0])
	}
}

func printHelp(w io.Writer) {
	_, _ = fmt.Fprintln(w, "rss-cli - RSS service command line client\n\nUsage:\n  rss-cli [--context NAME] [--output json|table] <command>\n\nCommands:\n  context    manage contexts\n  freshrss   FreshRSS service commands\n  capabilities  describe supported capabilities\n  version    show build metadata")
}

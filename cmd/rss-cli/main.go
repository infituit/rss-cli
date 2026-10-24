package main

import (
	"context"
	"os"

	"github.com/infituit/rss-cli/internal/cli"
	"github.com/infituit/rss-cli/internal/config"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/infituit/rss-cli/internal/telemetry"
)

func main() {
	tracingConfig := telemetry.Config{}
	current, err := config.Current()
	if err != nil {
		_ = output.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	if current != "" {
		contextConfig, loadErr := config.LoadContext(current)
		if loadErr != nil {
			_ = output.WriteError(os.Stderr, loadErr)
			os.Exit(1)
		}
		if contextConfig.Tracing != nil {
			tracingConfig = telemetry.Config{
				Enabled:  contextConfig.Tracing.Enabled,
				Endpoint: contextConfig.Tracing.Endpoint,
			}
		}
	}
	shutdown, err := telemetry.Setup(context.Background(), tracingConfig)
	if err != nil {
		_ = output.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	ctx, span := telemetry.Tracer().Start(context.Background(), telemetry.CommandSpanName())
	err = cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	span.End()
	shutdownCtx, cancel := telemetry.ShutdownContext()
	defer cancel()
	if shutdownErr := shutdown(shutdownCtx); shutdownErr != nil {
		_ = output.WriteError(os.Stderr, shutdownErr)
	}
	if err != nil {
		_ = output.WriteError(os.Stderr, err)
		os.Exit(1)
	}
}

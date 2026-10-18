package main

import (
	"context"
	"os"

	"github.com/infituit/rss-cli/internal/cli"
	"github.com/infituit/rss-cli/internal/output"
	"github.com/infituit/rss-cli/internal/telemetry"
)

func main() {
	shutdown, err := telemetry.Setup(context.Background())
	if err != nil {
		_ = output.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	ctx, span := telemetry.Tracer().Start(context.Background(), telemetry.CommandSpanName())
	err = cli.RunContext(ctx, os.Args[1:], os.Stdout, os.Stderr)
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

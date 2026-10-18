// Package main describes the logger example usages.
// nolint: mnd
package main

import (
	"context"
	"errors"

	"github.com/infituit/rss-cli/pkg/contextx"
	"github.com/infituit/rss-cli/pkg/logger"
)

func main() {
	logger.Init(logger.Config{
		LogDir:       "./logs",
		LogMaxSizeMB: 1,
		LogMaxNum:    3,

		ToStdErr:     false,
		AlsoToStdErr: true,
		Level:        logger.LevelDebug,
	})
	defer logger.G.Flush()

	// generates a request context.
	ctx := contextx.New(context.Background(), contextx.WithValues(map[string]any{"a": "ok", "b": 3, "c": false}))

	// business logs.
	logger.G.Biz(ctx).Debug("debug message")
	logger.G.Biz(ctx).Info("received a request %s", "install/agent")
	logger.G.Biz(ctx).WithErr(errors.New("invalid params")).Error("failed to decode params")
	logger.G.Biz(ctx).
		With("host-id", 123).With("state", "unknown").
		Warn("got unexpected data")
	logger.G.Biz(ctx).With("host-id", 123, "state", "unknown").
		Warn("got unexpected data")
	logger.G.Biz(ctx).With("wrong-key").Info("try to print a wrong key-value pair")

	// system logs.
	logger.G.Sys().Debug("debug message")
	logger.G.Sys().Info("brings up a backend thread %d", 123)
	logger.G.Sys().WithErr(errors.New("internal error")).Error("failed to launch task")
	logger.G.Sys().Ctx(ctx).With("host-id", 123).With("state", "unknown").
		Warn("got unexpected data")
	logger.G.Sys().Ctx(ctx).With("host-id", 123, "state", "unknown").
		Warn("got unexpected data")
	logger.G.Sys().Ctx(ctx).With("wrong-key").Info("try to print a wrong key-value pair")
}

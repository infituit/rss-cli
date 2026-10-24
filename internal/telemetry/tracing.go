package telemetry

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/infituit/rss-cli/pkg/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName             = "rss-cli"
	commandSpanName         = "rss-cli.command"
	httpTracerName          = "rss-cli.http"
	shutdownTimeout         = 5 * time.Second
	serviceNameAttribute    = "service.name"
	serviceVersionAttribute = "service.version"
)

// Config controls OTLP tracing for a process.
type Config struct {
	Enabled  bool
	Endpoint string
}

// Setup configures OTLP tracing from the selected context configuration.
func Setup(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	endpoint, err := url.Parse(cfg.Endpoint)
	if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") {
		return nil, fmt.Errorf("invalid tracing endpoint %q", cfg.Endpoint)
	}
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(cfg.Endpoint))
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String(serviceNameAttribute, serviceName),
		attribute.String(serviceVersionAttribute, version.VERSION),
	))
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}

func ShutdownContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), shutdownTimeout)
}

func Tracer() trace.Tracer { return otel.Tracer(serviceName) }

func CommandSpanName() string { return commandSpanName }

func HTTPTracerName() string { return httpTracerName }

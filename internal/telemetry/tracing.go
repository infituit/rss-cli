package telemetry

import (
	"context"
	"os"
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
	serviceName                    = "rss-cli"
	commandSpanName                = "rss-cli.command"
	httpTracerName                 = "rss-cli.http"
	shutdownTimeout                = 5 * time.Second
	otelSDKDisabledEnv             = "OTEL_SDK_DISABLED"
	otelExporterOTLPEndpoint       = "OTEL_EXPORTER_OTLP_ENDPOINT"
	otelExporterOTLPTracesEndpoint = "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"
	otelSDKDisabledValue           = "true"
	serviceNameAttribute           = "service.name"
	serviceVersionAttribute        = "service.version"
)

// Setup configures OTLP tracing when an OTLP endpoint is present.
func Setup(ctx context.Context) (func(context.Context) error, error) {
	if os.Getenv(otelSDKDisabledEnv) == otelSDKDisabledValue ||
		(os.Getenv(otelExporterOTLPEndpoint) == "" && os.Getenv(otelExporterOTLPTracesEndpoint) == "") {
		return func(context.Context) error { return nil }, nil
	}
	exporter, err := otlptracehttp.New(ctx)
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

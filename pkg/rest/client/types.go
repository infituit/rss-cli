// Package client provides the shared HTTP REST client used by providers.
package client

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Capability struct {
	Name                 string
	BaseURL              string
	HTTPClient           HTTPClient
	Logger               *slog.Logger
	ToleranceLatencyTime time.Duration
	MaxRetryCycle        int
	MetricObserver       MetricObserver
	HeaderMaskers        map[string]func(string) string
	URLQueryMaskers      map[string]func(string) string
}

type MetricObserver interface {
	Observe(name, method, path string, statusCode int, duration time.Duration)
}

const ToleranceLatencyTimeDefault = 500 * time.Millisecond

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NopLogger struct{}

func (NopLogger) Debug(string, ...any) {}
func (NopLogger) Info(string, ...any)  {}
func (NopLogger) Warn(string, ...any)  {}
func (NopLogger) Error(string, ...any) {}

type Result struct {
	FullURL    string
	Body       io.ReadCloser
	Err        error
	StatusCode int
	Status     string
	Header     http.Header

	logResponse  bool
	queryMaskers map[string]func(string) string
}

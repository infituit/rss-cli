// Package logger provides the logger interface.
package logger

import (
	"io"
	"time"

	"github.com/infituit/rss-cli/pkg/contextx"
)

// ILogger is the logger interface.
type ILogger interface {
	// Sys returns the system logger.
	Sys() ILoggerOption

	// Biz returns the business logger.
	Biz(ctx contextx.IContext) ILoggerOption

	// Flush flushes all logs into files.
	Flush()
}

// ILoggerOption is the logger option interface.
type ILoggerOption interface {
	// Ctx add context into logger.
	Ctx(ctx contextx.IContext) ILoggerOption

	// With add key-value pairs into logger.
	With(kvs ...interface{}) ILoggerOption

	// WithErr add error into logger.
	WithErr(err error) ILoggerOption

	// WithDuration add cost into logger.
	WithDuration(duration time.Duration) ILoggerOption

	// AssignWhenLogging assigns the logger message to str when logging.
	AssignWhenLogging(str *string) ILoggerOption

	// Debug logs debug message.
	Debug(format string, args ...interface{})

	// Info logs info message.
	Info(format string, args ...interface{})

	// Warn logs warn message.
	Warn(format string, args ...interface{})

	// Error logs error message.
	Error(format string, args ...interface{})

	// DebugWriter returns the writer for debug message.
	DebugWriter() io.Writer

	// InfoWriter returns the writer for info message.
	InfoWriter() io.Writer

	// WarnWriter returns the writer for warn message.
	WarnWriter() io.Writer

	// ErrorWriter returns the writer for error message.
	ErrorWriter() io.Writer
}

// ILoggerPrinter is the logger printer interface.
type ILoggerPrinter interface {
	// Log prints log with category, level, format and args.
	Log(category Category, level Level, depth int, format string, args ...interface{})
}

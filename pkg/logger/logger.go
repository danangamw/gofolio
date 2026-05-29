// Package logger provides a structured logger based on slog that is integrated
// with OpenTelemetry. Logs are automatically sent via OTLP to Grafana Loki
// (via Grafana Alloy), and each log entry includes trace_id and
// span_id for correlation with traces in Grafana Tempo.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/trace"
)

// logger is the global instance of slog.Logger used throughout the application.
var logger *slog.Logger

// Options configures the logger behavior during initialization.
type Options struct {
	// ServiceName is used as the "service" field in each log entry.
	ServiceName string
	// Level determines the minimum log level: "debug", "info", "warn", "error".
	// Default: "info".
	Level string
	// Environment is displayed as the "env" field (e.g. "development", "production").
	Environment string
}

// Init initializes the global logger. Must be called after telemetry.Init()
// so that the OTel bridge can connect to the already configured LoggerProvider.
//
// Contoh:
//
//	logger.Init(logger.Options{
//	    ServiceName: "go-cms",
//	    Level:       "info",
//	    Environment: "development",
//	})
func Init(opts Options) {
	level := parseLevel(opts.Level)

	// Handler for stdout (development: text, production: JSON)
	var stdoutHandler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug, // display file:line only when debug
	}

	if strings.EqualFold(opts.Environment, "production") {
		stdoutHandler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		// Development: colored text easy to read in the terminal
		stdoutHandler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	// OTel bridge: logs are automatically sent to Grafana Loki via Alloy
	otelHandler := otelslog.NewHandler(opts.ServiceName, otelslog.WithLoggerProvider(nil))
	// nil → uses the global LoggerProvider already configured by telemetry.Init()

	// Fan-out: log to stdout AND to OTLP simultaneously
	multiHandler := &multiHandler{
		handlers: []slog.Handler{stdoutHandler, otelHandler},
	}

	// Add standard fields that are always present in every log entry
	logger = slog.New(multiHandler).With(
		slog.String("service", opts.ServiceName),
		slog.String("env", opts.Environment),
	)

	// Set as the global slog logger so other packages can use standard slog.Info()
	slog.SetDefault(logger)
}

// --- Public Helper Functions ---

// Info logs at INFO level with automatic trace correlation from the context.
func Info(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelInfo, msg, args...)
}

// Debug mencatat log level DEBUG.
func Debug(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelDebug, msg, args...)
}

// Warn mencatat log level WARN.
func Warn(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at ERROR level. To include an error:
//
//	logger.Error(ctx, "failed to process", "error", err)
func Error(ctx context.Context, msg string, args ...any) {
	logWithTrace(ctx, slog.LevelError, msg, args...)
}

// --- Internal Helpers ---

// logWithTrace extracts trace_id and span_id from OpenTelemetry context
// then includes them in each log entry for correlation in Grafana.
func logWithTrace(ctx context.Context, level slog.Level, msg string, args ...any) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		args = append(args,
			"trace_id", spanCtx.TraceID().String(),
			"span_id", spanCtx.SpanID().String(),
		)
	}
	logger.Log(ctx, level, msg, args...)
}

// parseLevel converts a string level to slog.Level.
func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// multiHandler is a slog.Handler that delegates to multiple handlers simultaneously.
// This enables logs to be sent to stdout and OTLP simultaneously.
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				// Continue to the next handler even if one fails
				continue
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/logger"
)

// gormSlogLogger is a custom GORM logger that:
//  1. Forwards all query logs to global slog (→ stdout + Loki via OTLP)
//  2. Includes trace_id & span_id so that query logs can be correlated to HTTP requests in Grafana
//  3. Marks slow queries with WARN level
type gormSlogLogger struct {
	level                     logger.LogLevel
	slowQueryThreshold        time.Duration
	ignoreRecordNotFoundError bool
}

// newGormLogger creates an instance of gormSlogLogger.
// slowThreshold: query duration considered "slow" (default: 200ms).
func newGormLogger(level logger.LogLevel, slowThreshold time.Duration) logger.Interface {
	return &gormSlogLogger{
		level:                     level,
		slowQueryThreshold:        slowThreshold,
		ignoreRecordNotFoundError: true, // record not found is not an error that needs to be logged
	}
}

// LogMode implements logger.Interface — called by GORM to change the log level.
func (l *gormSlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

// Info implements logger.Interface for GORM internal INFO log level.
func (l *gormSlogLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Info {
		slog.InfoContext(ctx, fmt.Sprintf(msg, args...), "source", "gorm")
	}
}

// Warn implements logger.Interface for GORM internal WARN log level.
func (l *gormSlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Warn {
		slog.WarnContext(ctx, fmt.Sprintf(msg, args...), "source", "gorm")
	}
}

// Error implements logger.Interface for GORM internal ERROR log level.
func (l *gormSlogLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= logger.Error {
		slog.ErrorContext(ctx, fmt.Sprintf(msg, args...), "source", "gorm")
	}
}

// Trace is called by GORM after each query completes execution.
// Here is where we log SQL, duration, and rows affected.
func (l *gormSlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Collect standard log attributes
	attrs := []any{
		"source", "gorm",
		"duration_ms", elapsed.Milliseconds(),
		"rows", rows,
	}

	// Include trace_id + span_id if there is an active span in the context
	// → query logs can be clicked in Loki to jump directly to the trace in Tempo
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		attrs = append(attrs,
			"trace_id", span.SpanContext().TraceID().String(),
			"span_id", span.SpanContext().SpanID().String(),
		)
	}

	switch {
	case err != nil && !(l.ignoreRecordNotFoundError && isNotFound(err)):
		// Query error → ERROR level + include SQL for debugging
		slog.ErrorContext(ctx, "gorm query error",
			append(attrs, "error", err, "sql", sql)...,
		)

	case elapsed >= l.slowQueryThreshold:
		// Slow query → WARN level + include SQL
		slog.WarnContext(ctx, "gorm slow query",
			append(attrs, "sql", sql, "threshold_ms", l.slowQueryThreshold.Milliseconds())...,
		)

	case l.level >= logger.Info:
		// Query normal → INFO level + include SQL (only in development)
		slog.InfoContext(ctx, "gorm query",
			append(attrs, "sql", sql)...,
		)
	}
}

// isNotFound checks whether the error is gorm.ErrRecordNotFound
// without importing gorm directly to avoid circular dependency.
func isNotFound(err error) bool {
	return err != nil && err.Error() == "record not found"
}

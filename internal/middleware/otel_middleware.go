package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// OtelHTTP wraps the handler with OpenTelemetry instrumentation via otelhttp.
// Each HTTP request will automatically:
//   - Create a new span (stored in Grafana Tempo)
//   - Record HTTP metrics: duration, status code, bytes in/out (to Prometheus)
//   - Inject trace context into the request context
//
// How to use in RegisterRoutes():
//
//	r.Use(middleware.OtelHTTP("go-cms"))
func OtelHTTP(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName,
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		)
	}
}

// TraceIDHeader injects the trace ID from the active span into the X-Trace-ID response header.
// Useful for debugging: clients can see the trace ID in the response header
// then search for it directly in Grafana Tempo.
func TraceIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		if span.SpanContext().IsValid() {
			w.Header().Set("X-Trace-ID", span.SpanContext().TraceID().String())
		}
		next.ServeHTTP(w, r)
	})
}

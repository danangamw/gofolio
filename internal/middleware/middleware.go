package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

// client holds rate limit info for an IP address.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter returns an IP-based token bucket rate limiting middleware.
// r: limit of requests per second.
// b: burst size (max requests allowed at once).
func RateLimiter(r rate.Limit, b int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*client)

	// Background cleanup routine to prevent memory leaks from one-off IP entries.
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Extract IP address (handles direct and proxy-forwarded requests).
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}

			if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
				ip = xff
			}

			mu.Lock()
			c, exists := clients[ip]
			if !exists {
				c = &client{
					limiter: rate.NewLimiter(r, b),
				}
				clients[ip] = c
			}
			c.lastSeen = time.Now()
			mu.Unlock()

			if !c.limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// OtelHTTP wraps a handler with OpenTelemetry instrumentation.
// Each request creates a span, records HTTP metrics, and propagates trace context.
func OtelHTTP(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName,
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		)
	}
}

// TraceIDHeader injects X-Trace-ID into every response so clients can look
// up the request trace directly in Grafana Tempo.
func TraceIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		if span.SpanContext().IsValid() {
			w.Header().Set("X-Trace-ID", span.SpanContext().TraceID().String())
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds standard security-related HTTP response headers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

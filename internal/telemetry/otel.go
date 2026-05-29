package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config stores the telemetry configuration.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	CollectorAddr  string // contoh: "localhost:4317"
}

// ShutdownFunc is the function called when the application stops to
// clean up all telemetry resources gracefully.
type ShutdownFunc func(ctx context.Context) error

// Init initializes Tracing, Metrics, and Logging via OTLP to Grafana Alloy.
// Returns ShutdownFunc which must be called before the application stops.
func Init(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	res, err := newResource(cfg)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create resource: %w", err)
	}

	// Dial to Grafana Alloy via gRPC with connection timeout
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		cfg.CollectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: dial collector %q: %w", cfg.CollectorAddr, err)
	}

	// Verify the connection is available within the timeout
	if err := waitForConnection(dialCtx, conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: collector not reachable: %w", err)
	}

	// 1. Setup Trace Provider
	tp, err := newTraceProvider(ctx, conn, res, cfg.Environment)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup trace provider: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// 2. Setup Metric Provider
	mp, err := newMetricProvider(ctx, conn, res)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup metric provider: %w", err)
	}
	otel.SetMeterProvider(mp)

	// 3. Setup Log Provider
	lp, err := newLogProvider(ctx, conn, res)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup log provider: %w", err)
	}

	// Shutdown function — use a new context so it is not affected by an expired ctx
	shutdown := func(ctx context.Context) error {
		var errs []error
		if err := tp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
		if err := mp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("metric provider shutdown: %w", err))
		}
		if err := lp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("log provider shutdown: %w", err))
		}
		conn.Close()
		if len(errs) > 0 {
			return fmt.Errorf("telemetry shutdown errors: %v", errs)
		}
		return nil
	}

	return shutdown, nil
}

// newResource builds an OTel Resource with complete service attributes.
func newResource(cfg Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
}

// newTraceProvider creates a TracerProvider with a sampler appropriate for the environment.
func newTraceProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource, env string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	sampler := sdktrace.AlwaysSample()
	if env == "production" {
		// In production, sample 10% of traces to save costs
		sampler = sdktrace.TraceIDRatioBased(0.1)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

// newMetricProvider creates a MeterProvider with a periodic reader of 15 seconds.
func newMetricProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second)),
		),
	)
	return mp, nil
}

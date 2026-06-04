package telemetry

import (
	"context"
	"fmt"
	"time"

	"log"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds the telemetry initialization parameters.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	CollectorAddr  string // e.g. "localhost:4317"
}

var globalLoggerProvider *sdklog.LoggerProvider

// GetLoggerProvider returns the active OTel LoggerProvider.
func GetLoggerProvider() *sdklog.LoggerProvider {
	return globalLoggerProvider
}

// ShutdownFunc must be deferred in main() to flush and close all providers.
type ShutdownFunc func(ctx context.Context) error

// Init initializes Tracing, Metrics, and Logging via OTLP to Grafana Alloy.
// Returns a ShutdownFunc that must be called before the application exits.
//
// Usage in main():
//
//	shutdown, err := telemetry.Init(ctx, telemetry.Config{...})
//	if err != nil { log.Fatal(err) }
//	defer shutdown(ctx)
func Init(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	res, err := newResource(cfg)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create resource: %w", err)
	}

	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		cfg.CollectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: dial collector %q: %w", cfg.CollectorAddr, err)
	}

	if err := waitForConnection(dialCtx, conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: collector not reachable: %w", err)
	}

	// 1. Traces
	tp, err := newTraceProvider(ctx, conn, res, cfg.Environment)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup trace provider: %w", err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 2. Metrics
	mp, err := newMetricProvider(ctx, conn, res)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup metric provider: %w", err)
	}
	otel.SetMeterProvider(mp)

	// Start Go runtime instrumentation
	if err := runtime.Start(runtime.WithMeterProvider(mp)); err != nil {
		log.Printf("WARN: failed to start OTel Go runtime metrics: %v", err)
	}

	// 3. Logs
	lp, err := newLogProvider(ctx, conn, res)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("telemetry: setup log provider: %w", err)
	}
	global.SetLoggerProvider(lp)
	globalLoggerProvider = lp

	shutdown := func(ctx context.Context) error {
		var errs []error
		for _, fn := range []func(context.Context) error{tp.Shutdown, mp.Shutdown, lp.Shutdown} {
			if err := fn(ctx); err != nil {
				errs = append(errs, err)
			}
		}
		conn.Close()
		if len(errs) > 0 {
			return fmt.Errorf("telemetry shutdown errors: %v", errs)
		}
		return nil
	}
	return shutdown, nil
}

func newResource(cfg Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
			attribute.String("loki.resource.labels", "service.name"),
		),
	)
}

func newTraceProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource, env string) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	sampler := sdktrace.AlwaysSample()
	if env == "production" {
		sampler = sdktrace.TraceIDRatioBased(0.1) // sample 10% in prod
	}
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	), nil
}

func newMetricProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*metric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(15*time.Second))),
	), nil
}

func newLogProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	exp, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	return sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
	), nil
}

// waitForConnection waits for the gRPC connection to be ready or context timeout.
func waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	conn.Connect()
	for {
		state := conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return nil
		}
		if !conn.WaitForStateChange(ctx, state) {
			return ctx.Err()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}

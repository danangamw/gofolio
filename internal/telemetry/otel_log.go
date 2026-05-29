package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// newLogProvider creates a LoggerProvider that exports logs via OTLP to Alloy.
func newLogProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(exporter),
		),
	)
	return lp, nil
}

// waitForConnection waits for the gRPC connection to be ready, or timeouts if the collector is unavailable.
// This prevents the application from hanging forever when Alloy is not ready yet.
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

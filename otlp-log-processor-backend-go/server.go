package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	listenAddr            = flag.String("listenAddr", "localhost:4317", "The listen address")
	maxReceiveMessageSize = flag.Int("maxReceiveMessageSize", 16777216, "The max message size in bytes the server can receive")
	attributeKey          = flag.String("attributeKey", "service.name", "The attributeKey to count")
	durationWindow        = flag.Duration("duration", time.Second*10, "The duration between the output of the stats of the attributeKey")
	bufferSize            = flag.Uint("bufferSize", 1000, "The size of the buffer for log ingestion")
)

const name = "dash0.com/otlp-log-processor-backend"

var (
	tracer                = otel.Tracer(name)
	meter                 = otel.Meter(name)
	logger                = otelslog.NewLogger(name)
	logsReceivedCounter   metric.Int64Counter
	resourceLogHitCounter metric.Int64Counter
	scopeLogHitCounter    metric.Int64Counter
)

func init() {
	var err error
	logsReceivedCounter, err = meter.Int64Counter("com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"))
	if err != nil {
		panic(err)
	}
	resourceLogHitCounter, err = meter.Int64Counter("com.dash0.homeexercise.logs.resourcehit",
		metric.WithDescription("The number of hits of the attibute match from the resource"),
		metric.WithUnit("{log}"))
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	slog.SetDefault(logger)
	logger.Info("Starting application")

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(context.Background())
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	flag.Parse()

	slog.Debug("Starting listener", slog.String("listenAddr", *listenAddr))
	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.MaxRecvMsgSize(*maxReceiveMessageSize),
		grpc.Creds(insecure.NewCredentials()),
	)
	collogspb.RegisterLogsServiceServer(grpcServer, newServer(*listenAddr, *attributeKey, *durationWindow, *bufferSize))

	slog.Debug("Starting gRPC server", "attributeKey", *attributeKey, "durationWindow", durationWindow)

	return grpcServer.Serve(listener)
}

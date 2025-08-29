package main

import (
	"context"
	"log/slog"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	v1 "go.opentelemetry.io/proto/otlp/common/v1"
)

type dash0LogsServiceServer struct {
	addr           string
	attributeKey   string
	durationWindow time.Duration
	logIntake      chan v1.AnyValue
	logStats       map[string]uint64

	collogspb.UnimplementedLogsServiceServer
}

func newServer(addr string, attributeKey string, durationWindow time.Duration, bufferSize uint) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{
		addr:           addr,
		attributeKey:   attributeKey,
		durationWindow: durationWindow,
		logIntake:      make(chan v1.AnyValue, bufferSize),
	}

	go s.Start()
	return s
}

func (l *dash0LogsServiceServer) Start() {
	// Read and convert Values from logIntake channel
	// Start summary output every durationWindow
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")
	logsReceivedCounter.Add(ctx, 1)

	// The ResourceLogs typically only contain a single entry, but for propagated logs they might be bundled
	for _, logs := range request.ResourceLogs {
		for _, attrs := range logs.Resource.Attributes {
			if attrs.Key == l.attributeKey {
				resourceLogHitCounter.Add(ctx, 1)
				l.logIntake <- *attrs.Value
			}
		}
		for _, scopes := range logs.ScopeLogs {
			for _, logRecords := range scopes.LogRecords {
				for _, logRecordAttribute := range logRecords.Attributes {
					if logRecordAttribute.Key == l.attributeKey {
						l.logStats[logRecordAttribute.Value.GetStringValue()] += 1
					}
				}
			}
		}
	}
	// TODO: if the key did not match, there is no "unknown" metric

	return &collogspb.ExportLogsServiceResponse{}, nil
}

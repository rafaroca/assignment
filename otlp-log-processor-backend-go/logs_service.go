package main

import (
	"context"
	"log/slog"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

// TODO: make attributeKey settable via env
const attributeKey = "service.name"

type dash0LogsServiceServer struct {
	addr     string
	logStats map[string]uint64

	collogspb.UnimplementedLogsServiceServer
}

func newServer(addr string) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{
		addr:     addr,
		logStats: make(map[string]uint64),
	}
	return s
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")
	logsReceivedCounter.Add(ctx, 1)

	// Do something with the logs

	// The ResourceLogs typically only contain a single entry, but for propagated logs they might be bundled
	for _, logs := range request.ResourceLogs {
		for _, attrs := range logs.Resource.Attributes {
			if attrs.Key == attributeKey {
				l.logStats[attrs.Value.GetStringValue()] += 1
				resourceLogHitCounter.Add(ctx, 1)
				// TODO: is this only ever strings? Could be different types as well
			}
		}
		for _, scopes := range logs.ScopeLogs {
			for _, logRecords := range scopes.LogRecords {
				for _, logRecordAttribute := range logRecords.Attributes {
					if logRecordAttribute.Key == attributeKey {
						l.logStats[logRecordAttribute.Value.GetStringValue()] += 1
					}
				}
			}
		}

		// TODO: if the key did not match, there is no "unknown" metric
	}

	logger.Info("Current mappings", "logStats", l.logStats)
	return &collogspb.ExportLogsServiceResponse{}, nil
}

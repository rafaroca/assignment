package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otelcommon "go.opentelemetry.io/proto/otlp/common/v1"
)

type dash0LogsServiceServer struct {
	addr           string
	attributeKey   string
	durationWindow time.Duration
	logIntake      chan string
	logStats       map[string]uint64

	collogspb.UnimplementedLogsServiceServer
}

func newServer(addr string, attributeKey string, durationWindow time.Duration, bufferSize uint) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{
		addr:           addr,
		attributeKey:   attributeKey,
		durationWindow: durationWindow,
		logIntake:      make(chan string, bufferSize),
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
	if request.ResourceLogs != nil {
		for _, logs := range request.ResourceLogs {
			if logs.Resource != nil && logs.Resource.Attributes != nil {
				for _, attrs := range logs.Resource.Attributes {
					if attrs.Key == l.attributeKey {
						resourceLogHitCounter.Add(ctx, 1)
						l.logIntake <- extractStringValue(attrs.Value)
					}
				}
			}
			if logs.ScopeLogs != nil {
				for _, scopes := range logs.ScopeLogs {
					for _, logRecords := range scopes.LogRecords {
						for _, logRecordAttribute := range logRecords.Attributes {
							if logRecordAttribute.Key == l.attributeKey {
								l.logIntake <- extractStringValue(logRecordAttribute.Value)
							}
						}
					}
				}
			}
		}
	}
	// TODO: if the key did not match, there is no "unknown" metric

	return &collogspb.ExportLogsServiceResponse{}, nil
}

func extractStringValue(value *otelcommon.AnyValue) (strValue string) {
	switch v := value.GetValue().(type) {
	case *otelcommon.AnyValue_StringValue:
		strValue = v.StringValue
	case *otelcommon.AnyValue_BoolValue:
		strValue = fmt.Sprintf("%t", v.BoolValue)
	case *otelcommon.AnyValue_IntValue:
		strValue = fmt.Sprintf("%d", v.IntValue)
	case *otelcommon.AnyValue_DoubleValue:
		strValue = fmt.Sprintf("%f", v.DoubleValue)
	case *otelcommon.AnyValue_ArrayValue:
		strs := make([]string, len(v.ArrayValue.Values))
		for i, elem := range v.ArrayValue.Values {
			strs[i] = extractStringValue(elem)
		}
		strValue = "[" + strings.Join(strs, ",") + "]"
	case *otelcommon.AnyValue_KvlistValue:
		kvPairs := make([]string, len(v.KvlistValue.Values))
		for i, kv := range v.KvlistValue.Values {
			kvPairs[i] = kv.Key + ":" + extractStringValue(kv.Value)
		}
		strValue = "{" + strings.Join(kvPairs, ",") + "}"
	case nil:
		strValue = ""
	default:
		strValue = ""
	}
	return
}

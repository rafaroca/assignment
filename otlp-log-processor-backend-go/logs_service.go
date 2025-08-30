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
	addr         string
	attributeKey string
	processor    *dash0LogsProcessor
	logExport    chan<- string

	collogspb.UnimplementedLogsServiceServer
}

func newServer(addr string, attributeKey string, durationWindow time.Duration, bufferSize uint) collogspb.LogsServiceServer {
	logIntakeChannel := make(chan string, bufferSize)

	processor := &dash0LogsProcessor{
		durationWindow: durationWindow,
		logStats:       make(map[string]uint64),
		logIntake:      logIntakeChannel,
	}

	s := &dash0LogsServiceServer{
		addr:         addr,
		attributeKey: attributeKey,
		processor:    processor,
		logExport:    logIntakeChannel,
	}

	go processor.StartLogProcessing()

	return s
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")
	logsReceivedCounter.Add(ctx, 1)

	if request.ResourceLogs != nil {
		for _, resourceLog := range request.ResourceLogs {
			if resourceLog.Resource != nil && resourceLog.Resource.Attributes != nil {
				for _, attributes := range resourceLog.Resource.Attributes {
					if attributes.Key == l.attributeKey {
						resourceAttributeHitCounter.Add(ctx, 1)
						l.logExport <- extractStringValue(attributes.Value)
					}
				}
			}
			if resourceLog.ScopeLogs != nil {
				for _, scopeLog := range resourceLog.ScopeLogs {
					if scopeLog.LogRecords != nil {
						for _, logRecord := range scopeLog.LogRecords {
							if logRecord.Attributes != nil {
								for _, logRecordAttribute := range logRecord.Attributes {
									if logRecordAttribute.Key == l.attributeKey {
										logAttributeHitCounter.Add(ctx, 1)
										l.logExport <- extractStringValue(logRecordAttribute.Value)
									}
								}
							}
						}
					}
					if scopeLog.Scope != nil && scopeLog.Scope.Attributes != nil {
						for _, scopeAttribute := range scopeLog.Scope.Attributes {
							if scopeAttribute.Key == l.attributeKey {
								scopeAttributeHitCounter.Add(ctx, 1)
								l.logExport <- extractStringValue(scopeAttribute.Value)
							}
						}
					}
				}
			}
		}
	}

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

package main

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otelcommon "go.opentelemetry.io/proto/otlp/common/v1"
	otellogs "go.opentelemetry.io/proto/otlp/logs/v1"
	otelresource "go.opentelemetry.io/proto/otlp/resource/v1"
)

func TestExtractStringValue(t *testing.T) {
	tests := map[string]struct {
		input    *otelcommon.AnyValue
		expected string
	}{
		"StringValue": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_StringValue{
					StringValue: "test string",
				},
			},
			expected: "test string",
		},
		"BoolValue_True": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_BoolValue{
					BoolValue: true,
				},
			},
			expected: "true",
		},
		"BoolValue_False": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_BoolValue{
					BoolValue: false,
				},
			},
			expected: "false",
		},
		"IntValue_Positive": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_IntValue{
					IntValue: 42,
				},
			},
			expected: "42",
		},
		"IntValue_Negative": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_IntValue{
					IntValue: -123,
				},
			},
			expected: "-123",
		},
		"IntValue_Zero": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_IntValue{
					IntValue: 0,
				},
			},
			expected: "0",
		},
		"DoubleValue_Positive": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_DoubleValue{
					DoubleValue: 3.14159,
				},
			},
			expected: "3.141590",
		},
		"DoubleValue_Negative": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_DoubleValue{
					DoubleValue: -2.718,
				},
			},
			expected: "-2.718000",
		},
		"DoubleValue_Zero": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_DoubleValue{
					DoubleValue: 0.0,
				},
			},
			expected: "0.000000",
		},
		"ArrayValue_Empty": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_ArrayValue{
					ArrayValue: &otelcommon.ArrayValue{
						Values: []*otelcommon.AnyValue{},
					},
				},
			},
			expected: "[]",
		},
		"ArrayValue_SingleString": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_ArrayValue{
					ArrayValue: &otelcommon.ArrayValue{
						Values: []*otelcommon.AnyValue{
							{
								Value: &otelcommon.AnyValue_StringValue{
									StringValue: "item1",
								},
							},
						},
					},
				},
			},
			expected: "[item1]",
		},
		"ArrayValue_Multiple": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_ArrayValue{
					ArrayValue: &otelcommon.ArrayValue{
						Values: []*otelcommon.AnyValue{
							{
								Value: &otelcommon.AnyValue_StringValue{
									StringValue: "item1",
								},
							},
							{
								Value: &otelcommon.AnyValue_IntValue{
									IntValue: 42,
								},
							},
							{
								Value: &otelcommon.AnyValue_BoolValue{
									BoolValue: true,
								},
							},
						},
					},
				},
			},
			expected: "[item1,42,true]",
		},
		"KvlistValue_Empty": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_KvlistValue{
					KvlistValue: &otelcommon.KeyValueList{
						Values: []*otelcommon.KeyValue{},
					},
				},
			},
			expected: "{}",
		},
		"KvlistValue_SinglePair": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_KvlistValue{
					KvlistValue: &otelcommon.KeyValueList{
						Values: []*otelcommon.KeyValue{
							{
								Key: "key1",
								Value: &otelcommon.AnyValue{
									Value: &otelcommon.AnyValue_StringValue{
										StringValue: "value1",
									},
								},
							},
						},
					},
				},
			},
			expected: "{key1:value1}",
		},
		"KvlistValue_MultiplePairs": {
			input: &otelcommon.AnyValue{
				Value: &otelcommon.AnyValue_KvlistValue{
					KvlistValue: &otelcommon.KeyValueList{
						Values: []*otelcommon.KeyValue{
							{
								Key: "key1",
								Value: &otelcommon.AnyValue{
									Value: &otelcommon.AnyValue_StringValue{
										StringValue: "value1",
									},
								},
							},
							{
								Key: "key2",
								Value: &otelcommon.AnyValue{
									Value: &otelcommon.AnyValue_IntValue{
										IntValue: 123,
									},
								},
							},
						},
					},
				},
			},
			expected: "{key1:value1,key2:123}",
		},
		"NilValue": {
			input: &otelcommon.AnyValue{
				Value: nil,
			},
			expected: "",
		},
		"NilAnyValue": {
			input:    nil,
			expected: "",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := extractStringValue(test.input)
			if result != test.expected {
				t.Errorf("Test %s failed: expected '%s', got '%s'", testName, test.expected, result)
			}
		})
	}
}

func createResourceAttributesRequest() *collogspb.ExportLogsServiceRequest {
	return &collogspb.ExportLogsServiceRequest{
		ResourceLogs: []*otellogs.ResourceLogs{
			{
				Resource: &otelresource.Resource{
					Attributes: []*otelcommon.KeyValue{
						{
							Key: "service.name",
							Value: &otelcommon.AnyValue{
								Value: &otelcommon.AnyValue_StringValue{
									StringValue: "test-service",
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestLogsServiceServer_Export_ResourceAttributes(t *testing.T) {
	ctx := context.Background()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createResourceAttributesRequest()

	_, err := server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	select {
	case exported := <-logExportChannel:
		if exported != "test-service" {
			t.Errorf("Expected 'test-service', got '%s'", exported)
		}
	case <-time.After(time.Second):
		t.Error("Expected value to be sent to logExport channel but nothing was received")
	}
}

func TestLogsServiceServer_Export_ResourceAttributes_CounterIncrement(t *testing.T) {
	ctx := context.Background()

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	meter := provider.Meter("test")

	counter, err := meter.Int64Counter("com.dash0.homeexercise.logs.resourceattributehit")
	if err != nil {
		t.Fatalf("Failed to create counter: %v", err)
	}

	originalCounter := resourceAttributeHitCounter
	resourceAttributeHitCounter = counter
	defer func() { resourceAttributeHitCounter = originalCounter }()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createResourceAttributesRequest()

	_, err = server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	var resourceMetrics metricdata.ResourceMetrics
	err = reader.Collect(ctx, &resourceMetrics)
	if err != nil {
		t.Errorf("Failed to collect metrics: %v", err)
	}

	found := false
	for _, sm := range resourceMetrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "com.dash0.homeexercise.logs.resourceattributehit" {
				for _, dp := range m.Data.(metricdata.Sum[int64]).DataPoints {
					if dp.Value == 1 {
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		t.Error("Expected resourceAttributeHitCounter to be incremented by 1")
	}
}

func createLogRecordAttributesRequest() *collogspb.ExportLogsServiceRequest {
	return &collogspb.ExportLogsServiceRequest{
		ResourceLogs: []*otellogs.ResourceLogs{
			{
				ScopeLogs: []*otellogs.ScopeLogs{
					{
						LogRecords: []*otellogs.LogRecord{
							{
								Attributes: []*otelcommon.KeyValue{
									{
										Key: "service.name",
										Value: &otelcommon.AnyValue{
											Value: &otelcommon.AnyValue_StringValue{
												StringValue: "test-log-service",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestLogsServiceServer_Export_LogRecordAttributes(t *testing.T) {
	ctx := context.Background()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createLogRecordAttributesRequest()

	_, err := server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	select {
	case exported := <-logExportChannel:
		if exported != "test-log-service" {
			t.Errorf("Expected 'test-log-service', got '%s'", exported)
		}
	case <-time.After(time.Second):
		t.Error("Expected value to be sent to logExport channel but nothing was received")
	}
}

func TestLogsServiceServer_Export_LogRecordAttributes_CounterIncrement(t *testing.T) {
	ctx := context.Background()

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	meter := provider.Meter("test")

	counter, err := meter.Int64Counter("com.dash0.homeexercise.logs.logattributehit")
	if err != nil {
		t.Fatalf("Failed to create counter: %v", err)
	}

	originalCounter := logAttributeHitCounter
	logAttributeHitCounter = counter
	defer func() { logAttributeHitCounter = originalCounter }()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createLogRecordAttributesRequest()

	_, err = server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	var resourceMetrics metricdata.ResourceMetrics
	err = reader.Collect(ctx, &resourceMetrics)
	if err != nil {
		t.Errorf("Failed to collect metrics: %v", err)
	}

	found := false
	for _, sm := range resourceMetrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "com.dash0.homeexercise.logs.logattributehit" {
				for _, dp := range m.Data.(metricdata.Sum[int64]).DataPoints {
					if dp.Value == 1 {
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		t.Error("Expected logAttributeHitCounter to be incremented by 1")
	}
}

func createScopeAttributesRequest() *collogspb.ExportLogsServiceRequest {
	return &collogspb.ExportLogsServiceRequest{
		ResourceLogs: []*otellogs.ResourceLogs{
			{
				ScopeLogs: []*otellogs.ScopeLogs{
					{
						Scope: &otelcommon.InstrumentationScope{
							Attributes: []*otelcommon.KeyValue{
								{
									Key: "service.name",
									Value: &otelcommon.AnyValue{
										Value: &otelcommon.AnyValue_StringValue{
											StringValue: "test-scope-service",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestLogsServiceServer_Export_ScopeAttributes(t *testing.T) {
	ctx := context.Background()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createScopeAttributesRequest()

	_, err := server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	select {
	case exported := <-logExportChannel:
		if exported != "test-scope-service" {
			t.Errorf("Expected 'test-scope-service', got '%s'", exported)
		}
	case <-time.After(time.Second):
		t.Error("Expected value to be sent to logExport channel but nothing was received")
	}
}

func TestLogsServiceServer_Export_ScopeAttributes_CounterIncrement(t *testing.T) {
	ctx := context.Background()

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	meter := provider.Meter("test")

	counter, err := meter.Int64Counter("com.dash0.homeexercise.logs.scopeattributehit")
	if err != nil {
		t.Fatalf("Failed to create counter: %v", err)
	}

	originalCounter := scopeAttributeHitCounter
	scopeAttributeHitCounter = counter
	defer func() { scopeAttributeHitCounter = originalCounter }()

	logExportChannel := make(chan string, 10)
	server := &dash0LogsServiceServer{
		addr:         "localhost:4317",
		attributeKey: "service.name",
		logExport:    logExportChannel,
	}

	request := createScopeAttributesRequest()

	_, err = server.Export(ctx, request)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	var resourceMetrics metricdata.ResourceMetrics
	err = reader.Collect(ctx, &resourceMetrics)
	if err != nil {
		t.Errorf("Failed to collect metrics: %v", err)
	}

	found := false
	for _, sm := range resourceMetrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "com.dash0.homeexercise.logs.scopeattributehit" {
				for _, dp := range m.Data.(metricdata.Sum[int64]).DataPoints {
					if dp.Value == 1 {
						found = true
						break
					}
				}
			}
		}
	}

	if !found {
		t.Error("Expected scopeAttributeHitCounter to be incremented by 1")
	}
}

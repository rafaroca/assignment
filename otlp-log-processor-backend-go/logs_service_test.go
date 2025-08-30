package main

import (
	"testing"

	v1 "go.opentelemetry.io/proto/otlp/common/v1"
)

func TestExtractStringValue(t *testing.T) {
	tests := map[string]struct {
		input    *v1.AnyValue
		expected string
	}{
		"StringValue": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_StringValue{
					StringValue: "test string",
				},
			},
			expected: "test string",
		},
		"BoolValue_True": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_BoolValue{
					BoolValue: true,
				},
			},
			expected: "true",
		},
		"BoolValue_False": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_BoolValue{
					BoolValue: false,
				},
			},
			expected: "false",
		},
		"IntValue_Positive": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_IntValue{
					IntValue: 42,
				},
			},
			expected: "42",
		},
		"IntValue_Negative": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_IntValue{
					IntValue: -123,
				},
			},
			expected: "-123",
		},
		"IntValue_Zero": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_IntValue{
					IntValue: 0,
				},
			},
			expected: "0",
		},
		"DoubleValue_Positive": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_DoubleValue{
					DoubleValue: 3.14159,
				},
			},
			expected: "3.141590",
		},
		"DoubleValue_Negative": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_DoubleValue{
					DoubleValue: -2.718,
				},
			},
			expected: "-2.718000",
		},
		"DoubleValue_Zero": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_DoubleValue{
					DoubleValue: 0.0,
				},
			},
			expected: "0.000000",
		},
		"ArrayValue_Empty": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_ArrayValue{
					ArrayValue: &v1.ArrayValue{
						Values: []*v1.AnyValue{},
					},
				},
			},
			expected: "[]",
		},
		"ArrayValue_SingleString": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_ArrayValue{
					ArrayValue: &v1.ArrayValue{
						Values: []*v1.AnyValue{
							{
								Value: &v1.AnyValue_StringValue{
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
			input: &v1.AnyValue{
				Value: &v1.AnyValue_ArrayValue{
					ArrayValue: &v1.ArrayValue{
						Values: []*v1.AnyValue{
							{
								Value: &v1.AnyValue_StringValue{
									StringValue: "item1",
								},
							},
							{
								Value: &v1.AnyValue_IntValue{
									IntValue: 42,
								},
							},
							{
								Value: &v1.AnyValue_BoolValue{
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
			input: &v1.AnyValue{
				Value: &v1.AnyValue_KvlistValue{
					KvlistValue: &v1.KeyValueList{
						Values: []*v1.KeyValue{},
					},
				},
			},
			expected: "{}",
		},
		"KvlistValue_SinglePair": {
			input: &v1.AnyValue{
				Value: &v1.AnyValue_KvlistValue{
					KvlistValue: &v1.KeyValueList{
						Values: []*v1.KeyValue{
							{
								Key: "key1",
								Value: &v1.AnyValue{
									Value: &v1.AnyValue_StringValue{
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
			input: &v1.AnyValue{
				Value: &v1.AnyValue_KvlistValue{
					KvlistValue: &v1.KeyValueList{
						Values: []*v1.KeyValue{
							{
								Key: "key1",
								Value: &v1.AnyValue{
									Value: &v1.AnyValue_StringValue{
										StringValue: "value1",
									},
								},
							},
							{
								Key: "key2",
								Value: &v1.AnyValue{
									Value: &v1.AnyValue_IntValue{
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
			input: &v1.AnyValue{
				Value: nil,
			},
			expected: "",
		},
		"NilAnyValue": {
			input: nil,
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
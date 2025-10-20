package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

var (
	colTests = []struct {
		collection CollectionName
		yamlStr    string
	}{
		{ARGS, "ARGS"},
		{ARGS_GET, "ARGS_GET"},
		{ARGS_GET_NAMES, "ARGS_GET_NAMES"},
		{ARGS_NAMES, "ARGS_NAMES"},
		{ARGS_POST_NAMES, "ARGS_POST_NAMES"},
		{ARGS_POST, "ARGS_POST"},
		{ENV, "ENV"},
		{FILES, "FILES"},
		{GEO, "GEO"},
		{GLOBAL, "GLOBAL"},
		{IP, "IP"},
		{MATCHED_VARS_NAMES, "MATCHED_VARS_NAMES"},
		{MATCHED_VARS, "MATCHED_VARS"},
		{MULTIPART_PART_HEADERS, "MULTIPART_PART_HEADERS"},
		{PERF_RULES, "PERF_RULES"},
		{REQUEST_COOKIES_NAMES, "REQUEST_COOKIES_NAMES"},
		{REQUEST_COOKIES, "REQUEST_COOKIES"},
		{REQUEST_HEADERS_NAMES, "REQUEST_HEADERS_NAMES"},
		{REQUEST_HEADERS, "REQUEST_HEADERS"},
		{RESOURCE, "RESOURCE"},
		{RESPONSE_HEADERS_NAMES, "RESPONSE_HEADERS_NAMES"},
		{RESPONSE_HEADERS, "RESPONSE_HEADERS"},
		{RULE, "RULE"},
		{SESSION, "SESSION"},
		{TX, "TX"},
		{USER, "USER"},
		{XML, "XML"},
	}
)

func TestCollectionNameToString(t *testing.T) {
	for _, tt := range colTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			if tt.collection.String() != tt.yamlStr {
				t.Errorf("Expected %q, got %q", tt.yamlStr, tt.collection.String())
			}
		})
	}
}

func TestStringToCollectionName(t *testing.T) {
	for _, tt := range colTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			collection := stringToCollectionName(tt.yamlStr)
			if collection != tt.collection {
				t.Errorf("Expected %q, got %q", tt.collection, collection)
			}
		})
	}
}

func TestMarshalCollectionName(t *testing.T) {
	for _, tt := range colTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			data, err := yaml.Marshal(tt.collection)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			if string(data) != tt.yamlStr+"\n" {
				t.Errorf("Expected %q, got %q", tt.yamlStr+"\n", data)
			}
		})
	}
}

func TestUnknownCollectionName(t *testing.T) {
	t.Run("marshal unknown", func(t *testing.T) {
		unknown := UNKNOWN_COLLECTION
		_, err := unknown.MarshalYAML()
		assert.Error(t, err)
		assert.Equal(t, "Unknown collection name", err.Error())
	})
}

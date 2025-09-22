package types

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

var (
	stringTests = []struct {
		operator OperatorType
		yamlStr  string
	}{
		{BeginsWith, "beginsWith"},
		{Contains, "contains"},
		{ContainsWord, "containsWord"},
		{DetectSQLi, "detectSQLi"},
		{DetectXSS, "detectXSS"},
		{EndsWith, "endsWith"},
		{Eq, "eq"},
		{FuzzyHash, "fuzzyHash"},
		{Ge, "ge"},
		{GeoLookup, "geoLookup"},
		{GsbLookup, "gsbLookup"},
		{Gt, "gt"},
		{InspectFile, "inspectFile"},
		{IpMatchF, "ipMatchF"},
		{IpMatchFromFile, "ipMatchFromFile"},
		{IpMatch, "ipMatch"},
		{Le, "le"},
		{Lt, "lt"},
		{Pmf, "pmf"},
		{PmFromFile, "pmFromFile"},
		{Pm, "pm"},
		{Rbl, "rbl"},
		{Rsub, "rsub"},
		{Rx, "rx"},
		{RxGlobal, "rxGlobal"},
		{StrEq, "streq"},
		{StrMatch, "strmatch"},
		{UnconditionalMatch, "unconditionalMatch"},
		{ValidateByteRange, "validateByteRange"},
		{ValidateDTD, "validateDTD"},
		{ValidateHash, "validateHash"},
		{ValidateSchema, "validateSchema"},
		{ValidateUrlEncoding, "validateUrlEncoding"},
		{ValidateUtf8Encoding, "validateUtf8Encoding"},
		{VerifyCC, "verifyCC"},
		{VerifyCPF, "verifyCPF"},
		{VerifySSN, "verifySSN"},
		{VerifySVNR, "verifySVNR"},
		{Within, "within"},
	}
)

func TestOperatorTypeToString(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			if tt.operator.String() != tt.yamlStr {
				t.Errorf("Expected %q, got %q", tt.yamlStr, tt.operator.String())
			}
		})
	}
}

func TestStringToOperatorType(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			op := stringToOperatorType(tt.yamlStr)
			if op != tt.operator {
				t.Errorf("Expected %q, got %q", tt.operator, op)
			}
		})
	}
}

func TestMarshalOperatorType(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			data, err := yaml.Marshal(tt.operator)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			if string(data) != tt.yamlStr+"\n" {
				t.Errorf("Expected %q, got %q", tt.yamlStr+"\n", data)
			}
		})
	}
}

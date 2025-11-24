package types

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	operatorUnmarshalTests = []struct {
		name           string
		input          string
		expected       Operator
		expectedString string
	}{
		{
			name: "Standard format",
			input: `name: rx
value: ^.*$`,
			expected: Operator{
				Name:  Rx,
				Value: "^.*$",
			},
			expectedString: "@rx ^.*$",
		},
		{
			name: "Standard format, negated",
			input: `name: rx
value: ^.*$
negate: true`,
			expected: Operator{
				Name:   Rx,
				Value:  "^.*$",
				Negate: true,
			},
			expectedString: "!@rx ^.*$",
		},
		{
			name:  "Compact format",
			input: `rx: ^.*$`,
			expected: Operator{
				Name:  Rx,
				Value: "^.*$",
			},
			expectedString: "@rx ^.*$",
		},
		{
			name: "Compact format, negated",
			input: `rx: ^.*$
negate: true`,
			expected: Operator{
				Name:   Rx,
				Value:  "^.*$",
				Negate: true,
			},
			expectedString: "!@rx ^.*$",
		},
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

func TestUnmarshalOperator(t *testing.T) {
	for _, tt := range operatorUnmarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			var result Operator
			err := yaml.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			require.Equal(t, tt.expected, result)
			require.Equal(t, tt.expectedString, result.ToString())
		})
	}
}

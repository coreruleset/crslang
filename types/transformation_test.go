package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

var (
	stringTests = []struct {
		transformation Transformation
		yamlStr        string
	}{
		{Base64Decode, "base64Decode"},
		{Base64DecodeExt, "base64DecodeExt"},
		{Base64Encode, "base64Encode"},
		{CmdLine, "cmdLine"},
		{CompressWhitespace, "compressWhitespace"},
		{EscapeSeqDecode, "escapeSeqDecode"},
		{CssDecode, "cssDecode"},
		{HexEncode, "hexEncode"},
		{HexDecode, "hexDecode"},
		{HtmlEntityDecode, "htmlEntityDecode"},
		{JsDecode, "jsDecode"},
		{Length, "length"},
		{Lowercase, "lowercase"},
		{Md5, "md5"},
		{None, "none"},
		{NormalisePath, "normalisePath"},
		{normalizePath, "normalizePath"},
		{NormalisePathWin, "normalisePathWin"},
		{normalizePathWin, "normalizePathWin"},
		{ParityEven7bit, "parityEven7bit"},
		{ParityOdd7bit, "parityOdd7bit"},
		{ParityZero7bit, "parityZero7bit"},
		{RemoveComments, "removeComments"},
		{RemoveCommentsChar, "removeCommentsChar"},
		{RemoveNulls, "removeNulls"},
		{RemoveWhitespace, "removeWhitespace"},
		{ReplaceComments, "replaceComments"},
		{ReplaceNulls, "replaceNulls"},
		{Sha1, "sha1"},
		{SqlHexDecode, "sqlHexDecode"},
		{Trim, "trim"},
		{TrimLeft, "trimLeft"},
		{TrimRight, "trimRight"},
		{Uppercase, "uppercase"},
		{UrlEncode, "urlEncode"},
		{UrlDecode, "urlDecode"},
		{UrlDecodeUni, "urlDecodeUni"},
		{Utf8toUnicode, "utf8toUnicode"},
	}
)

func TestTransformationToString(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			if tt.transformation.String() != tt.yamlStr {
				t.Errorf("Expected %q, got %q", tt.yamlStr, tt.transformation.String())
			}
		})
	}
}

func TestStringToOperatorType(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			op := stringToTransformation(tt.yamlStr)
			if op != tt.transformation {
				t.Errorf("Expected %q, got %q", tt.transformation, op)
			}
		})
	}
}

func TestMarshalTransformation(t *testing.T) {
	for _, tt := range stringTests {
		t.Run(tt.yamlStr, func(t *testing.T) {
			data, err := yaml.Marshal(tt.transformation)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			if string(data) != tt.yamlStr+"\n" {
				t.Errorf("Expected %q, got %q", tt.yamlStr+"\n", data)
			}
		})
	}
}

func TestUnknownTransformation(t *testing.T) {
	t.Run("marshal unknown", func(t *testing.T) {
		unknown := UnknownTransformation
		_, err := unknown.MarshalYAML()
		assert.Error(t, err)
		assert.Equal(t, "Invalid transformation value", err.Error())
	})

	t.Run("add unknown transformation", func(t *testing.T) {
		var ts Transformations
		err := ts.AddTransformation("invalidTransformation")
		assert.Error(t, err)
		assert.Equal(t, "Invalid transformation value: invalidTransformation", err.Error())
		assert.Empty(t, ts.Transformations)
	})
}

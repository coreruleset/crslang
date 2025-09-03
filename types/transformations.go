package types

import "fmt"

type Transformations struct {
	Transformations []Transformation `yaml:"transformations,omitempty"`
}

type Transformation string

const (
	Base64Decode       Transformation = "base64Decode"
	Base64DecodeExt    Transformation = "base64DecodeExt"
	Base64Encode       Transformation = "base64Encode"
	CmdLine            Transformation = "cmdLine"
	CompressWhitespace Transformation = "compressWhitespace"
	EscapeSeqDecode    Transformation = "escapeSeqDecode"
	CssDecode          Transformation = "cssDecode"
	HexEncode          Transformation = "hexEncode"
	HexDecode          Transformation = "hexDecode"
	HtmlEntityDecode   Transformation = "htmlEntityDecode"
	JsDecode           Transformation = "jsDecode"
	Length             Transformation = "length"
	Lowercase          Transformation = "lowercase"
	Md5                Transformation = "md5"
	None               Transformation = "none"
	NormalisePath      Transformation = "normalisePath"
	normalizePath      Transformation = "normalizePath"
	NormalisePathWin   Transformation = "normalisePathWin"
	normalizePathWin   Transformation = "normalizePathWin"
	ParityEven7bit     Transformation = "parityEven7bit"
	ParityOdd7bit      Transformation = "parityOdd7bit"
	ParityZero7bit     Transformation = "parityZero7bit"
	RemoveComments     Transformation = "removeComments"
	RemoveCommentsChar Transformation = "removeCommentsChar"
	RemoveNulls        Transformation = "removeNulls"
	RemoveWhitespace   Transformation = "removeWhitespace"
	ReplaceComments    Transformation = "replaceComments"
	ReplaceNulls       Transformation = "replaceNulls"
	Sha1               Transformation = "sha1"
	SqlHexDecode       Transformation = "sqlHexDecode"
	Trim               Transformation = "trim"
	TrimLeft           Transformation = "trimLeft"
	TrimRight          Transformation = "trimRight"
	Uppercase          Transformation = "uppercase"
	UrlEncode          Transformation = "urlEncode"
	UrlDecode          Transformation = "urlDecode"
	UrlDecodeUni       Transformation = "urlDecodeUni"
	Utf8toUnicode      Transformation = "utf8toUnicode"
)

var (
	allTransformation = map[string]Transformation{
		"base64Decode":       Base64Decode,
		"base64DecodeExt":    Base64DecodeExt,
		"base64Encode":       Base64Encode,
		"cmdLine":            CmdLine,
		"compressWhitespace": CompressWhitespace,
		"escapeSeqDecode":    EscapeSeqDecode,
		"cssDecode":          CssDecode,
		"hexEncode":          HexEncode,
		"hexDecode":          HexDecode,
		"htmlEntityDecode":   HtmlEntityDecode,
		"jsDecode":           JsDecode,
		"length":             Length,
		"lowercase":          Lowercase,
		"md5":                Md5,
		"none":               None,
		"normalisePath":      NormalisePath,
		"normalizePath":      normalizePath,
		"normalisePathWin":   NormalisePathWin,
		"normalizePathWin":   normalizePathWin,
		"parityEven7bit":     ParityEven7bit,
		"parityOdd7bit":      ParityOdd7bit,
		"parityZero7bit":     ParityZero7bit,
		"removeComments":     RemoveComments,
		"removeCommentsChar": RemoveCommentsChar,
		"removeNulls":        RemoveNulls,
		"removeWhitespace":   RemoveWhitespace,
		"replaceComments":    ReplaceComments,
		"replaceNulls":       ReplaceNulls,
		"sha1":               Sha1,
		"sqlHexDecode":       SqlHexDecode,
		"trim":               Trim,
		"trimLeft":           TrimLeft,
		"trimRight":          TrimRight,
		"uppercase":          Uppercase,
		"urlEncode":          UrlEncode,
		"urlDecode":          UrlDecode,
		"urlDecodeUni":       UrlDecodeUni,
		"utf8toUnicode":      Utf8toUnicode,
	}
)

func (t *Transformations) AddTransformation(transformation string) error {
	constValue, exists := allTransformation[transformation]
	if !exists {
		return fmt.Errorf("Invalid transformation value: %s", transformation)
	}
	t.Transformations = append(t.Transformations, constValue)
	return nil
}

func (t Transformations) ToString() string {
	results := []string{}
	for _, transformation := range t.Transformations {
		results = append(results, string(transformation))
	}
	result := ""
	for i, value := range results {
		if i == 0 {
			result += "t:" + value
		} else {
			result += ",t:" + value
		}
	}
	return result
}

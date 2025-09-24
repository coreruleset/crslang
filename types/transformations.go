package types

import "fmt"

type Transformations struct {
	Transformations []Transformation `yaml:"transformations,omitempty"`
}

type Transformation int

const (
	UnknownTransformation Transformation = iota
	Base64Decode
	Base64DecodeExt
	Base64Encode
	CmdLine
	CompressWhitespace
	EscapeSeqDecode
	CssDecode
	HexEncode
	HexDecode
	HtmlEntityDecode
	JsDecode
	Length
	Lowercase
	Md5
	None
	NormalisePath
	normalizePath
	NormalisePathWin
	normalizePathWin
	ParityEven7bit
	ParityOdd7bit
	ParityZero7bit
	RemoveComments
	RemoveCommentsChar
	RemoveNulls
	RemoveWhitespace
	ReplaceComments
	ReplaceNulls
	Sha1
	SqlHexDecode
	Trim
	TrimLeft
	TrimRight
	Uppercase
	UrlEncode
	UrlDecode
	UrlDecodeUni
	Utf8toUnicode
)

func (t Transformation) String() string {
	switch t {
	case Base64Decode:
		return "base64Decode"
	case Base64DecodeExt:
		return "base64DecodeExt"
	case Base64Encode:
		return "base64Encode"
	case CmdLine:
		return "cmdLine"
	case CompressWhitespace:
		return "compressWhitespace"
	case EscapeSeqDecode:
		return "escapeSeqDecode"
	case CssDecode:
		return "cssDecode"
	case HexEncode:
		return "hexEncode"
	case HexDecode:
		return "hexDecode"
	case HtmlEntityDecode:
		return "htmlEntityDecode"
	case JsDecode:
		return "jsDecode"
	case Length:
		return "length"
	case Lowercase:
		return "lowercase"
	case Md5:
		return "md5"
	case None:
		return "none"
	case NormalisePath:
		return "normalisePath"
	case normalizePath:
		return "normalizePath"
	case NormalisePathWin:
		return "normalisePathWin"
	case normalizePathWin:
		return "normalizePathWin"
	case ParityEven7bit:
		return "parityEven7bit"
	case ParityOdd7bit:
		return "parityOdd7bit"
	case ParityZero7bit:
		return "parityZero7bit"
	case RemoveComments:
		return "removeComments"
	case RemoveCommentsChar:
		return "removeCommentsChar"
	case RemoveNulls:
		return "removeNulls"
	case RemoveWhitespace:
		return "removeWhitespace"
	case ReplaceComments:
		return "replaceComments"
	case ReplaceNulls:
		return "replaceNulls"
	case Sha1:
		return "sha1"
	case SqlHexDecode:
		return "sqlHexDecode"
	case Trim:
		return "trim"
	case TrimLeft:
		return "trimLeft"
	case TrimRight:
		return "trimRight"
	case Uppercase:
		return "uppercase"
	case UrlEncode:
		return "urlEncode"
	case UrlDecode:
		return "urlDecode"
	case UrlDecodeUni:
		return "urlDecodeUni"
	case Utf8toUnicode:
		return "utf8toUnicode"
	default:
		return "unknown"
	}
}

func stringToTransformation(value string) Transformation {
	switch value {
	case "base64Decode":
		return Base64Decode
	case "base64DecodeExt":
		return Base64DecodeExt
	case "base64Encode":
		return Base64Encode
	case "cmdLine":
		return CmdLine
	case "compressWhitespace":
		return CompressWhitespace
	case "escapeSeqDecode":
		return EscapeSeqDecode
	case "cssDecode":
		return CssDecode
	case "hexEncode":
		return HexEncode
	case "hexDecode":
		return HexDecode
	case "htmlEntityDecode":
		return HtmlEntityDecode
	case "jsDecode":
		return JsDecode
	case "length":
		return Length
	case "lowercase":
		return Lowercase
	case "md5":
		return Md5
	case "none":
		return None
	case "normalisePath":
		return NormalisePath
	case "normalizePath":
		return normalizePath
	case "normalisePathWin":
		return NormalisePathWin
	case "normalizePathWin":
		return normalizePathWin
	case "parityEven7bit":
		return ParityEven7bit
	case "parityOdd7bit":
		return ParityOdd7bit
	case "parityZero7bit":
		return ParityZero7bit
	case "removeComments":
		return RemoveComments
	case "removeCommentsChar":
		return RemoveCommentsChar
	case "removeNulls":
		return RemoveNulls
	case "removeWhitespace":
		return RemoveWhitespace
	case "replaceComments":
		return ReplaceComments
	case "replaceNulls":
		return ReplaceNulls
	case "sha1":
		return Sha1
	case "sqlHexDecode":
		return SqlHexDecode
	case "trim":
		return Trim
	case "trimLeft":
		return TrimLeft
	case "trimRight":
		return TrimRight
	case "uppercase":
		return Uppercase
	case "urlEncode":
		return UrlEncode
	case "urlDecode":
		return UrlDecode
	case "urlDecodeUni":
		return UrlDecodeUni
	case "utf8toUnicode":
		return Utf8toUnicode
	default:
		return UnknownTransformation
	}
}

func (t Transformation) MarshalYAML() (interface{}, error) {
	if t == UnknownTransformation {
		return nil, fmt.Errorf("Invalid transformation value")
	}
	return t.String(), nil
}

func (t *Transformation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	if err := unmarshal(&name); err != nil {
		return err
	}
	tConst := stringToTransformation(name)
	if tConst == UnknownTransformation {
		return fmt.Errorf("Transformation %s is not valid", name)
	}
	*t = tConst
	return nil
}

func (t *Transformations) AddTransformation(transformation string) error {
	tConst := stringToTransformation(transformation)
	if tConst == UnknownTransformation {
		return fmt.Errorf("Invalid transformation value: %s", transformation)
	}
	t.Transformations = append(t.Transformations, tConst)
	return nil
}

func (t Transformations) ToString() string {
	results := []string{}
	for _, transformation := range t.Transformations {
		results = append(results, transformation.String())
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

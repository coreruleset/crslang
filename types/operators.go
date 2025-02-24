package types

import "fmt"

type Operator struct {
	Name  OperatorType `yaml:"name"`
	Value string       `yaml:",omitempty"`
}

type OperatorType string

const (
	BeginsWith           OperatorType = "beginsWith"
	Contains             OperatorType = "contains"
	ContainsWord         OperatorType = "containsWord"
	DetectSQLi           OperatorType = "detectSQLi"
	DetectXSS            OperatorType = "detectXSS"
	EndsWith             OperatorType = "endsWith"
	Eq                   OperatorType = "eq"
	FuzzyHash            OperatorType = "fuzzyHash"
	Ge                   OperatorType = "ge"
	GeoLookup            OperatorType = "geoLookup"
	GsbLookup            OperatorType = "gsbLookup"
	Gt                   OperatorType = "gt"
	InspectFile          OperatorType = "inspectFile"
	IpMatchF             OperatorType = "ipMatchF"
	IpMatchFromFile      OperatorType = "ipMatchFromFile"
	IpMatch              OperatorType = "ipMatch"
	Le                   OperatorType = "le"
	Lt                   OperatorType = "lt"
	Pmf                  OperatorType = "pmf"
	PmFromFile           OperatorType = "pmFromFile"
	Pm                   OperatorType = "pm"
	Rbl                  OperatorType = "rbl"
	Rsub                 OperatorType = "rsub"
	Rx                   OperatorType = "rx"
	RxGlobal             OperatorType = "rxGlobal"
	StrEq                OperatorType = "streq"
	StrMatch             OperatorType = "strmatch"
	UnconditionalMatch   OperatorType = "unconditionalMatch"
	ValidateByteRange    OperatorType = "validateByteRange"
	ValidateDTD          OperatorType = "validateDTD"
	ValidateHash         OperatorType = "validateHash"
	ValidateSchema       OperatorType = "validateSchema"
	ValidateUrlEncoding  OperatorType = "validateUrlEncoding"
	ValidateUtf8Encoding OperatorType = "validateUtf8Encoding"
	VerifyCC             OperatorType = "verifyCC"
	VerifyCPF            OperatorType = "verifyCPF"
	VerifySSN            OperatorType = "verifySSN"
	VerifySVNR           OperatorType = "verifySVNR"
	Within               OperatorType = "within"
)

var (
	allOperators = map[string]OperatorType{
		"beginsWith":           BeginsWith,
		"contains":             Contains,
		"containsWord":         ContainsWord,
		"detectSQLi":           DetectSQLi,
		"detectXSS":            DetectXSS,
		"endsWith":             EndsWith,
		"eq":                   Eq,
		"fuzzyHash":            FuzzyHash,
		"ge":                   Ge,
		"geoLookup":            GeoLookup,
		"gsbLookup":            GsbLookup,
		"gt":                   Gt,
		"inspectFile":          InspectFile,
		"ipMatchF":             IpMatchF,
		"ipMatchFromFile":      IpMatchFromFile,
		"ipMatch":              IpMatch,
		"le":                   Le,
		"lt":                   Lt,
		"pmf":                  Pmf,
		"pmFromFile":           PmFromFile,
		"pm":                   Pm,
		"rbl":                  Rbl,
		"rsub":                 Rsub,
		"rx":                   Rx,
		"rxGlobal":             RxGlobal,
		"streq":                StrEq,
		"strmatch":             StrMatch,
		"unconditionalMatch":   UnconditionalMatch,
		"validateByteRange":    ValidateByteRange,
		"validateDTD":          ValidateDTD,
		"validateHash":         ValidateHash,
		"validateSchema":       ValidateSchema,
		"validateUrlEncoding":  ValidateUrlEncoding,
		"validateUtf8Encoding": ValidateUtf8Encoding,
		"verifyCC":             VerifyCC,
		"verifyCPF":            VerifyCPF,
		"verifySSN":            VerifySSN,
		"verifySVNR":           VerifySVNR,
		"within":               Within,
	}
)

func (o *Operator) SetOperatorName(name string) error {
	operatorConst, ok := allOperators[name]
	if !ok {
		return fmt.Errorf("Operator name %s not found", name)
	}

	o.Name = operatorConst
	return nil
}

func (o *Operator) SetOperatorValue(value string) {
	o.Value = value
}

func (o *Operator) ToString() string {
	return "@" + string(o.Name) + " " + o.Value
}

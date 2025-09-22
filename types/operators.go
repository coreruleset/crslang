package types

import (
	"fmt"
)

type Operator struct {
	Negate bool         `yaml:"negate,omitempty"`
	Name   OperatorType `yaml:"name"`
	Value  string       `yaml:",omitempty"`
}

type OperatorType int

const (
	UnknownOperator OperatorType = iota
	BeginsWith
	Contains
	ContainsWord
	DetectSQLi
	DetectXSS
	EndsWith
	Eq
	FuzzyHash
	Ge
	GeoLookup
	GsbLookup
	Gt
	InspectFile
	IpMatchF
	IpMatchFromFile
	IpMatch
	Le
	Lt
	Pmf
	PmFromFile
	Pm
	Rbl
	Rsub
	Rx
	RxGlobal
	StrEq
	StrMatch
	UnconditionalMatch
	ValidateByteRange
	ValidateDTD
	ValidateHash
	ValidateSchema
	ValidateUrlEncoding
	ValidateUtf8Encoding
	VerifyCC
	VerifyCPF
	VerifySSN
	VerifySVNR
	Within
)

func (o OperatorType) String() string {
	switch o {
	case BeginsWith:
		return "beginsWith"
	case Contains:
		return "contains"
	case ContainsWord:
		return "containsWord"
	case DetectSQLi:
		return "detectSQLi"
	case DetectXSS:
		return "detectXSS"
	case EndsWith:
		return "endsWith"
	case Eq:
		return "eq"
	case FuzzyHash:
		return "fuzzyHash"
	case Ge:
		return "ge"
	case GeoLookup:
		return "geoLookup"
	case GsbLookup:
		return "gsbLookup"
	case Gt:
		return "gt"
	case InspectFile:
		return "inspectFile"
	case IpMatchF:
		return "ipMatchF"
	case IpMatchFromFile:
		return "ipMatchFromFile"
	case IpMatch:
		return "ipMatch"
	case Le:
		return "le"
	case Lt:
		return "lt"
	case Pmf:
		return "pmf"
	case PmFromFile:
		return "pmFromFile"
	case Pm:
		return "pm"
	case Rbl:
		return "rbl"
	case Rsub:
		return "rsub"
	case Rx:
		return "rx"
	case RxGlobal:
		return "rxGlobal"
	case StrEq:
		return "streq"
	case StrMatch:
		return "strmatch"
	case UnconditionalMatch:
		return "unconditionalMatch"
	case ValidateByteRange:
		return "validateByteRange"
	case ValidateDTD:
		return "validateDTD"
	case ValidateHash:
		return "validateHash"
	case ValidateSchema:
		return "validateSchema"
	case ValidateUrlEncoding:
		return "validateUrlEncoding"
	case ValidateUtf8Encoding:
		return "validateUtf8Encoding"
	case VerifyCC:
		return "verifyCC"
	case VerifyCPF:
		return "verifyCPF"
	case VerifySSN:
		return "verifySSN"
	case VerifySVNR:
		return "verifySVNR"
	case Within:
		return "within"
	default:
		return "unknownOperator"
	}
}

func stringToOperatorType(name string) OperatorType {
	switch name {
	case "beginsWith":
		return BeginsWith
	case "contains":
		return Contains
	case "containsWord":
		return ContainsWord
	case "detectSQLi":
		return DetectSQLi
	case "detectXSS":
		return DetectXSS
	case "endsWith":
		return EndsWith
	case "eq":
		return Eq
	case "fuzzyHash":
		return FuzzyHash
	case "ge":
		return Ge
	case "geoLookup":
		return GeoLookup
	case "gsbLookup":
		return GsbLookup
	case "gt":
		return Gt
	case "inspectFile":
		return InspectFile
	case "ipMatchF":
		return IpMatchF
	case "ipMatchFromFile":
		return IpMatchFromFile
	case "ipMatch":
		return IpMatch
	case "le":
		return Le
	case "lt":
		return Lt
	case "pmf":
		return Pmf
	case "pmFromFile":
		return PmFromFile
	case "pm":
		return Pm
	case "rbl":
		return Rbl
	case "rsub":
		return Rsub
	case "rx":
		return Rx
	case "rxGlobal":
		return RxGlobal
	case "streq":
		return StrEq
	case "strmatch":
		return StrMatch
	case "unconditionalMatch":
		return UnconditionalMatch
	case "validateByteRange":
		return ValidateByteRange
	case "validateDTD":
		return ValidateDTD
	case "validateHash":
		return ValidateHash
	case "validateSchema":
		return ValidateSchema
	case "validateUrlEncoding":
		return ValidateUrlEncoding
	case "validateUtf8Encoding":
		return ValidateUtf8Encoding
	case "verifyCC":
		return VerifyCC
	case "verifyCPF":
		return VerifyCPF
	case "verifySSN":
		return VerifySSN
	case "verifySVNR":
		return VerifySVNR
	case "within":
		return Within
	default:
		return UnknownOperator
	}
}

func (o OperatorType) MarshalYAML() (interface{}, error) {
	return o.String(), nil
}

func (o *OperatorType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	if err := unmarshal(&name); err != nil {
		return err
	}
	operatorConst := stringToOperatorType(name)
	if operatorConst == UnknownOperator {
		return fmt.Errorf("Operator name %s is not valid", name)
	}
	*o = operatorConst
	return nil
}

func (o *Operator) SetOperatorName(name string) error {
	operatorConst := stringToOperatorType(name)
	if operatorConst == UnknownOperator {
		return fmt.Errorf("Operator name %s is not valid", name)
	}

	o.Name = operatorConst
	return nil
}

func (o *Operator) SetOperatorValue(value string) {
	o.Value = value
}

func (o *Operator) SetOperatorNot(not bool) {
	o.Negate = not
}

func (o *Operator) ToString() string {
	if o.Value != "" {
		return "@" + o.Name.String() + " " + o.Value
	} else {
		return "@" + o.Name.String()
	}
}

package types

import (
	"strconv"
	"strings"
)

type Metadata interface {
	SetComment(value string)
	SetId(value string)
	SetPhase(value string)
	SetMsg(value string)
	SetMaturity(value string)
	SetRev(value string)
	SetSeverity(value string)
	AddTag(value string)
	SetVer(value string)
}

type CommentMetadata struct {
	Comment string `yaml:"comment,omitempty"`
}

func (c CommentMetadata) GetKind() Kind {
	return UnknownKind
}

func (c *CommentMetadata) SetComment(value string) {
	c.Comment = value
}

func (c *CommentMetadata) SetId(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetPhase(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetMsg(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetMaturity(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetRev(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetSeverity(value string) {
	// Do nothing
}

func (c *CommentMetadata) AddTag(value string) {
	// Do nothing
}

func (c *CommentMetadata) SetVer(value string) {
	// Do nothing
}

func (c CommentMetadata) ToSeclang() string {
	return commentToSeclang(c.Comment)
}

// commentToSeclang converts comment strings to SecLang format by replacing newlines with newline + #
func commentToSeclang(comment string) string {
	lines := strings.Split(comment, "\n")
	res := ""
	for i, line := range lines {
		if i != len(lines)-1 && line != "" {
			res += "# " + line + "\n"
		} else if i != len(lines)-1 {
			res += "#\n"
		}
	}
	return res
}

type SecRuleMetadata struct {
	OnlyPhaseMetadata `yaml:",inline"`
	Id                int      `yaml:"id,omitempty"`
	Msg               string   `yaml:"message,omitempty"`
	Maturity          string   `yaml:"maturity,omitempty"`
	Rev               string   `yaml:"revision,omitempty"`
	Severity          string   `yaml:"severity,omitempty"`
	Tags              []string `yaml:"tags,omitempty"`
	Ver               string   `yaml:"version,omitempty"`
}

type OnlyPhaseMetadata struct {
	CommentMetadata `yaml:",inline"`
	Phase           string `yaml:"phase,omitempty"`
}

func (m *OnlyPhaseMetadata) ToString() string {
	return "phase:" + m.Phase
}

func (m *OnlyPhaseMetadata) SetId(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) SetPhase(value string) {
	m.Phase = value
}

func (m *OnlyPhaseMetadata) SetMsg(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) SetMaturity(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) SetRev(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) SetSeverity(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) AddTag(value string) {
	// Do nothing
}

func (m *OnlyPhaseMetadata) SetVer(value string) {
	// Do nothing
}

func (s *SecRuleMetadata) ToString() string {
	result := ""
	result += s.OnlyPhaseMetadata.ToString() + ", id:" + strconv.Itoa(s.Id)
	if s.Msg != "" {
		result += ", msg:'" + s.Msg + "'"
	}
	if s.Maturity != "" {
		result += ", maturity:'" + s.Maturity + "'"
	}
	if s.Rev != "" {
		result += ", rev:'" + s.Rev + "'"
	}
	if s.Severity != "" {
		result += ", severity:'" + s.Severity + "'"
	}
	if s.Ver != "" {
		result += ", ver:'" + s.Ver + "'"
	}
	return result
}

func CopySecRuleMetadata(s SecRuleMetadata) *SecRuleMetadata {
	copy := new(SecRuleMetadata)
	copy.OnlyPhaseMetadata = s.OnlyPhaseMetadata
	copy.Id = s.Id
	copy.Maturity = s.Maturity
	copy.Msg = s.Msg
	copy.Rev = s.Rev
	copy.Severity = s.Severity
	copy.Ver = s.Ver
	copy.Tags = append(copy.Tags, s.Tags...)
	return copy
}

func (s *SecRuleMetadata) SetId(value string) {
	if s.Id != 0 {
		panic("Id already set")
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	s.Id = id
}

// func (s *SecRuleMetadata) SetPhase(value string) {
// 	s.Phase = value
// }

func (s *SecRuleMetadata) SetMsg(value string) {
	s.Msg = value
}

func (s *SecRuleMetadata) SetMaturity(value string) {
	s.Maturity = value
}

func (s *SecRuleMetadata) SetRev(value string) {
	s.Rev = value
}

func (s *SecRuleMetadata) SetSeverity(value string) {
	s.Severity = value
}

func (s *SecRuleMetadata) AddTag(value string) {
	s.Tags = append(s.Tags, value)
}

func (s *SecRuleMetadata) SetVer(value string) {
	s.Ver = value
}

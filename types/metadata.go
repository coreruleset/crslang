package types

import (
	"strconv"
)

type CommentMetadata struct {
	Comment string `yaml:"comment,omitempty"`
}

func (c *CommentMetadata) SetComment(value string) {
	c.Comment = value
}

func (c CommentMetadata) ToSeclang() string {
	return c.Comment
}

type SecRuleMetadata struct {
	OnlyPhaseMetadata	`yaml:",inline"`
	Id       int    `yaml:"id,omitempty"`
	Msg      string `yaml:"message,omitempty"`
	Maturity string `yaml:"maturity,omitempty"`
	Rev      string `yaml:"revision,omitempty"`
	Severity string `yaml:"severity,omitempty"`
	Ver      string `yaml:"version,omitempty"`
}

type OnlyPhaseMetadata struct{
	CommentMetadata `yaml:",inline"`
	Phase             string `yaml:"phase"`
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

func (s *SecRuleMetadata) SetVer(value string) {
	s.Ver = value
}
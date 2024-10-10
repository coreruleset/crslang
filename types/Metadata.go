package types

import (
	"strconv"
)

type SecRuleMetada struct {
	OnlyPhaseMetadata	`yaml:"-"`
	Id       int    `yaml:"id,omitempty"`
	Msg      string `yaml:"message,omitempty"`
	Maturity string `yaml:"maturity,omitempty"`
	Rev      string `yaml:"revision,omitempty"`
	Severity string `yaml:"severity,omitempty"`
	Ver      string `yaml:"version,omitempty"`
}

type OnlyPhaseMetadata struct{
	Phase             string `yaml:"phase"`
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



func (s *SecRuleMetada) String() string {
	return "Id: " + strconv.Itoa(s.Id) + ", Phase: " + s.Phase + ", Msg: " + s.Msg + ", Maturity: " + s.Maturity + ", Rev: " + s.Rev + ", Severity: " + s.Severity + ", Ver: " + s.Ver
}

func (s *SecRuleMetada) SetId(value string) {
	if s.Id != 0 {
		panic("Id already set")
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	s.Id = id
}

func (s *SecRuleMetada) SetPhase(value string) {
	s.Phase = value
}

func (s *SecRuleMetada) SetMsg(value string) {
	s.Msg = value
}

func (s *SecRuleMetada) SetMaturity(value string) {
	s.Maturity = value
}

func (s *SecRuleMetada) SetRev(value string) {
	s.Rev = value
}

func (s *SecRuleMetada) SetSeverity(value string) {
	s.Severity = value
}

func (s *SecRuleMetada) SetVer(value string) {
	s.Ver = value
}
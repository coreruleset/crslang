package types

import "strconv"

type UpdateActionDirective struct {
	Kind    Kind         `yaml:"kind"`
	Comment string       `yaml:"comment,omitempty"`
	Id      int          `yaml:"id"`
	Modify  ModifyAction `yaml:"modify"`
}

type ModifyAction struct {
	Metadata        *UpdateActionMetadata `yaml:"metadata,omitempty"`
	Transformations `yaml:",inline"`
	Actions         *SeclangActions `yaml:"actions,omitempty"`
}

type UpdateActionMetadata struct {
	Msg      string   `yaml:"message,omitempty"`
	Maturity string   `yaml:"maturity,omitempty"`
	Rev      string   `yaml:"revision,omitempty"`
	Severity string   `yaml:"severity,omitempty"`
	Tags     []string `yaml:"tags,omitempty"`
	Ver      string   `yaml:"version,omitempty"`
}

func NewUpdateActionDirective() *UpdateActionDirective {
	directive := new(UpdateActionDirective)
	directive.Kind = UpdateAction
	directive.Modify.Metadata = new(UpdateActionMetadata)
	directive.Modify.Actions = new(SeclangActions)
	return directive
}

func (d UpdateActionDirective) GetMetadata() Metadata {
	return d.Modify.Metadata
}

func (d UpdateActionDirective) GetActions() *SeclangActions {
	return d.Modify.Actions
}

func (d *UpdateActionDirective) AddTransformation(t string) error {
	return d.Modify.AddTransformation(t)
}

func (d UpdateActionDirective) ToSeclang() string {
	result := d.Comment + "SecRuleUpdateActionById " + strconv.Itoa(d.Id) + " \""
	actionString := ""
	actionString += d.Modify.Metadata.ToString()
	if actionString != "" {
		actionString += ","
	}
	actionString += d.Modify.Transformations.ToString()
	if actionString != "" {
		actionString += ","
	}
	actionString += d.Modify.Actions.ToString()
	result += actionString + "\"\n"
	return result
}

func (d UpdateActionDirective) AppendChainedDirective(directive ChainableDirective) {
	// Do nothing
}

func (m *UpdateActionMetadata) SetComment(value string) {
	// Do nothing
}

func (m *UpdateActionMetadata) SetId(value string) {
	// Do nothing
}

func (m *UpdateActionMetadata) SetPhase(value string) {
	// Do nothing
}

func (m *UpdateActionMetadata) SetMsg(value string) {
	m.Msg = value
}

func (m *UpdateActionMetadata) SetMaturity(value string) {
	m.Maturity = value
}

func (m *UpdateActionMetadata) SetRev(value string) {
	m.Rev = value
}

func (m *UpdateActionMetadata) SetSeverity(value string) {
	m.Severity = value
}

func (m *UpdateActionMetadata) AddTag(value string) {
	m.Tags = append(m.Tags, value)
}

func (m *UpdateActionMetadata) SetVer(value string) {
	m.Ver = value
}

func (s *UpdateActionMetadata) ToString() string {
	result := ""
	if s.Msg != "" {
		result += "msg:'" + s.Msg + "'"
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

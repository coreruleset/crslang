package types

type UpdateActionDirective struct {
	Kind            Kind `yaml:"kind"`
	CommentMetadata `yaml:",inline"`
	Id              int                   `yaml:"id,omitempty"`
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
	directive.Metadata = new(UpdateActionMetadata)
	directive.Actions = new(SeclangActions)
	return directive
}

func (d UpdateActionDirective) GetMetadata() Metadata {
	return d.Metadata
}

func (d UpdateActionDirective) GetActions() *SeclangActions {
	return d.Actions
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

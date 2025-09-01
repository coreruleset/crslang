package types

type SecAction struct {
	Metadata        *SecRuleMetadata `yaml:"metadata,omitempty"`
	Transformations `yaml:",inline"`
	Actions         *SeclangActions    `yaml:"actions"`
	ChainedRule     ChainableDirective `yaml:"chainedRule,omitempty"`
}

func NewSecAction() *SecAction {
	secAction := new(SecAction)
	secAction.Metadata = new(SecRuleMetadata)
	secAction.Actions = new(SeclangActions)
	return secAction
}

func (d SecAction) GetMetadata() Metadata {
	return d.Metadata
}

func (d SecAction) GetActions() *SeclangActions {
	return d.Actions
}

func (d SecAction) GetTransformations() Transformations {
	return d.Transformations
}

func (s *SecAction) AppendChainedDirective(chainedDirective ChainableDirective) {
	s.ChainedRule = chainedDirective
}

func (s SecAction) GetChainedDirective() ChainableDirective {
	return s.ChainedRule
}

func (s SecAction) NonDisruptiveActionsCount() int {
	return len(s.Actions.NonDisruptiveActions)
}

func (s SecAction) ToSeclang() string {
	result := ""
	result += s.Metadata.Comment + "SecAction \"phase:" + s.Metadata.Phase
	actions := s.Actions.ToString()
	transformations := s.Transformations.ToString()
	if actions != "" {
		result += "," + actions
	}
	if transformations != "" {
		result += ", " + transformations
	}
	result += "\"\n"
	return result
}

func (s SecAction) ToSeclangWithIdent(initialString string) string {
	return initialString + s.ToSeclang()
}

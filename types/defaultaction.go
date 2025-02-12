package types

type DefaultAction struct {
	Kind            string             `yaml:"kind,omitempty"`
	Metadata        *OnlyPhaseMetadata `yaml:"metadata"`
	Transformations `yaml:",inline"`
	Actions         *SeclangActions `yaml:"actions"`
}

func NewDefaultAction() *DefaultAction {
	defaultAction := new(DefaultAction)
	defaultAction.Kind = "defaultaction"
	defaultAction.Metadata = new(OnlyPhaseMetadata)
	defaultAction.Actions = new(SeclangActions)
	return defaultAction
}

func (d DefaultAction) GetMetadata() Metadata {
	return d.Metadata
}

func (d DefaultAction) GetActions() *SeclangActions {
	return d.Actions
}

func (s DefaultAction) ToSeclang() string {
	result := ""
	result += s.Metadata.Comment + "SecDefaultAction \"phase:" + s.Metadata.Phase
	actions := s.Actions.ToString()
	transformations := s.Transformations.ToString()
	if actions != "" {
		result += ", " + actions
	}
	if transformations != "" {
		result += ", " + transformations
	}
	result += "\"\n"
	return result
}

func (s DefaultAction) AppendChainedDirective(chainedDirective ChainableDirective) {
	return
}

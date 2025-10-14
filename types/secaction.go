package types

import "slices"

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

func (d SecAction) GetKind() Kind {
	return UnknownKind
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
	return s.ToSeclangWithIdent("")
}

func (s SecAction) ToSeclangWithIdent(initialString string) string {
	auxString := ",\\\n" + initialString + "    "
	endString := ""

	result := ""
	result += s.Metadata.Comment + initialString + "SecAction"
	sortedActions := sortActions(&s)
	for i, action := range sortedActions {
		if i == 0 {
			result += " \\\n" + initialString + "    \""
		} else {
			result += auxString
		}
		result += action
		if i == len(sortedActions)-1 {
			result += "\""
		} else {
			result += endString
		}
	}
	result += "\n"
	if slices.Contains(s.Actions.GetActionKeys(), "chain") {
		result += (s.ChainedRule).ToSeclangWithIdent(initialString + "    ")
	}
	return result
}

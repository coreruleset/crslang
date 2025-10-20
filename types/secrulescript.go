package types

import (
	"slices"
)

type SecRuleScript struct {
	Metadata        *SecRuleMetadata `yaml:"metadata,omitempty"`
	ScriptPath      string           `yaml:"scriptpath"`
	Transformations `yaml:",inline"`
	Actions         *SeclangActions    `yaml:"actions,omitempty"`
	ChainedRule     ChainableDirective `yaml:"chainedRule,omitempty"`
}

func NewSecRuleScript() *SecRuleScript {
	secRuleScript := new(SecRuleScript)
	secRuleScript.Metadata = new(SecRuleMetadata)
	secRuleScript.Actions = new(SeclangActions)
	return secRuleScript
}

func (d SecRuleScript) GetKind() Kind {
	return UnknownKind
}

func (d SecRuleScript) GetMetadata() Metadata {
	return d.Metadata
}

func (d SecRuleScript) GetActions() *SeclangActions {
	return d.Actions
}

func (d SecRuleScript) GetTransformations() Transformations {
	return d.Transformations
}

func (s SecRuleScript) ToSeclang() string {
	return s.ToSeclangWithIdent("")
}

func (s SecRuleScript) ToSeclangWithIdent(initialString string) string {
	auxString := ",\\\n" + initialString + "    "
	endString := ""
	result := ""
	result += s.Metadata.Comment + initialString + "SecRuleScript "
	result += s.ScriptPath + " "

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
		result += s.ChainedRule.ToSeclangWithIdent(initialString + "    ")
	}
	return result
}

func (s SecRuleScript) GetChainedDirective() ChainableDirective {
	return s.ChainedRule
}

func (s *SecRuleScript) AppendChainedDirective(chainedDirective ChainableDirective) {
	s.ChainedRule = chainedDirective
}

func (s SecRuleScript) NonDisruptiveActionsCount() int {
	return len(s.Actions.NonDisruptiveActions)
}

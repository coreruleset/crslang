package types

import (
	"slices"
	"strconv"
)

type SecRule struct {
	Metadata        *SecRuleMetadata `yaml:"metadata,omitempty"`
	Variables       `yaml:",inline"`
	Transformations `yaml:",inline"`
	Operator        `yaml:"operator"`
	Actions         *SeclangActions    `yaml:"actions,omitempty"`
	ChainedRule     ChainableDirective `yaml:"chainedRule,omitempty"`
}

func NewSecRule() *SecRule {
	secRule := new(SecRule)
	secRule.Metadata = new(SecRuleMetadata)
	secRule.Actions = new(SeclangActions)
	return secRule
}

func (d SecRule) GetMetadata() Metadata {
	return d.Metadata
}

func (d SecRule) GetActions() *SeclangActions {
	return d.Actions
}

func (s SecRule) ToSeclang() string {
	return s.ToSeclangWithIdent("")
}

func (s SecRule) ToSeclangWithIdent(initialString string) string {
	auxString := ",\\\n" + initialString + "    "
	endString := ""
	actions := s.Actions.GetActionKeys()
	auxSlice := []string{}
	chainedRule := false

	result := ""
	result += s.Metadata.Comment + initialString + "SecRule "
	result += s.Variables.ToString() + " "
	result += "\"" + s.Operator.ToString() + "\""
	if s.Metadata.Id != 0 {
		auxSlice = append(auxSlice, "id:"+strconv.Itoa(s.Metadata.Id))
	}
	if s.Metadata.Phase != "" {
		auxSlice = append(auxSlice, "phase:"+s.Metadata.Phase)
	}
	if s.Actions.DisruptiveAction.Action != "" {
		auxSlice = append(auxSlice, s.Actions.DisruptiveAction.ToString())
	}
	if slices.Contains(actions, "status") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("status").ToString())
	}
	if slices.Contains(actions, "capture") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("capture").ToString())
	}
	if len(s.Transformations.Transformations) > 0 {
		auxSlice = append(auxSlice, s.Transformations.ToString())
	}
	if slices.Contains(actions, "log") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("log").ToString())
	}
	if slices.Contains(actions, "nolog") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("nolog").ToString())
	}
	if slices.Contains(actions, "auditlog") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("auditlog").ToString())
	}
	if slices.Contains(actions, "noauditlog") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("noauditlog").ToString())
	}
	if s.Metadata.Msg != "" {
		auxSlice = append(auxSlice, "msg:'"+s.Metadata.Msg+"'")
	}
	if slices.Contains(actions, "logdata") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("logdata").ToString())
	}
	for _, tag := range s.Metadata.Tags {
		auxSlice = append(auxSlice, "tag:'"+tag+"'")
	}
	if slices.Contains(actions, "sanitiseArg") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("sanitiseArg").ToString())
	}
	if slices.Contains(actions, "sanitiseRequestHeader") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("sanitiseRequestHeader").ToString())
	}
	if slices.Contains(actions, "sanitiseMatched") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("sanitiseMatched").ToString())
	}
	if slices.Contains(actions, "sanitiseMatchedBytes") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("sanitiseMatchedBytes").ToString())
	}
	if slices.Contains(actions, "ctl") {
		ctlActions := s.Actions.GetActionsByKey("ctl")
		for _, action := range ctlActions {
			auxSlice = append(auxSlice, action.ToString())
		}
	}
	if s.Metadata.Ver != "" {
		auxSlice = append(auxSlice, "ver:'"+s.Metadata.Ver+"'")
	}
	if s.Metadata.Severity != "" {
		auxSlice = append(auxSlice, "severity:'"+s.Metadata.Severity+"'")
	}
	if slices.Contains(actions, "multiMatch") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("multiMatch").ToString())
	}
	if slices.Contains(actions, "initcol") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("initcol").ToString())
	}
	if slices.Contains(actions, "setenv") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("setenv").ToString())
	}
	if slices.Contains(actions, "setvar") {
		setvarActions := s.Actions.GetActionsByKey("setvar")
		for _, action := range setvarActions {
			auxSlice = append(auxSlice, action.ToString())
		}
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("setvar").ToString())
	}
	if slices.Contains(actions, "expirevar") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("expirevar").ToString())
	}
	if slices.Contains(actions, "chain") {
		chainedRule = true
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("chain").ToString())
	}
	if slices.Contains(actions, "skip") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("skip").ToString())
	}
	if slices.Contains(actions, "skipAfter") {
		auxSlice = append(auxSlice, s.Actions.GetActionByKey("skipAfter").ToString())
	}
	for i, action := range auxSlice {
		if i == 0 {
			result += " \\\n" + initialString + "    \""
		} else {
			result += auxString
		}
		result += action
		if i == len(auxSlice)-1 {
			result += "\""
		} else {
			result += endString
		}
	}
	result += "\n"
	if chainedRule {
		result += (s.ChainedRule).ToSeclangWithIdent(initialString + "    ")
	}
	return result
}

func (s SecRule) GetChainedDirective() ChainableDirective {
	return s.ChainedRule
}

func (s *SecRule) AppendChainedDirective(chainedDirective ChainableDirective) {
	s.ChainedRule = chainedDirective
}

func (s SecRule) NonDisruptiveActionsCount() int {
	return len(s.Actions.NonDisruptiveActions)
}

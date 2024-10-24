package types

import (
	"slices"
	"strconv"
)

type SeclangDirective interface {
	ToSeclang() string
}

type ConfigurationDirective struct {
	CommentMetadata `yaml:",inline"`
	DirectiveName   string `yaml:"directiveName"`
	Parameter       string `yaml:"parameter"`
}

func (c ConfigurationDirective) ToSeclang() string {
	return c.Comment + c.DirectiveName + " " + c.Parameter + "\n"
}

type SecDefaultAction struct {
	OnlyPhaseMetadata `yaml:"metadata"`
	Transformations   `yaml:",inline"`
	SeclangActions    `yaml:"actions"`
}

func (s SecDefaultAction) ToSeclang() string {
	result := ""
	result += s.Comment + "SecDefaultAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToString()
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

type SecAction struct {
	SecRuleMetadata `yaml:"metadata"`
	Transformations `yaml:",inline"`
	SeclangActions  `yaml:"actions"`
}

func (s SecAction) ToSeclang() string {
	result := ""
	result += s.Comment + "SecAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToString()
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

type SecRule struct {
	SecRuleMetadata `yaml:"metadata,omitempty"`
	Variables       `yaml:",inline"`
	Transformations `yaml:",inline"`
	Operator        `yaml:"operator"`
	SeclangActions  `yaml:"actions"`
	ChainedRule    *SecRule `yaml:"chainedRule,omitempty"`
}

func (s SecRule) ToSeclang() string {
	return s.ToSeclangWithParam("")
}

func (s SecRule) ToSeclangWithParam(initialString string) string {
	auxString := ",\\\n" + initialString + "\t"
	endString := ""
	actions := s.SeclangActions.GetActionKeys()
	auxSlice := []string{}
	chainedRule := false

	result := ""
	result += s.Comment + initialString + "SecRule "
	result += s.Variables.ToString() + " "
	result += "\"" + s.Operator.ToString() + "\""
	if s.SecRuleMetadata.Id != 0 {
		auxSlice = append(auxSlice, "id:" + strconv.Itoa(s.SecRuleMetadata.Id))
	}
	if s.SecRuleMetadata.Phase != "" {
		auxSlice = append(auxSlice, "phase:" + s.SecRuleMetadata.Phase)
	}
	if s.SeclangActions.DisruptiveAction != nil {
		auxSlice = append(auxSlice, s.SeclangActions.DisruptiveAction.ToString())
	}
	if slices.Contains(actions, "status") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("status").ToString())
	}
	if slices.Contains(actions, "capture") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("capture").ToString())
	}
	if len(s.Transformations.Transformations) > 0 {
		auxSlice = append(auxSlice, s.Transformations.ToString())
	}
	if slices.Contains(actions, "log") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("log").ToString())
	}
	if slices.Contains(actions, "nolog") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("nolog").ToString())
	}
	if slices.Contains(actions, "auditlog") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("auditlog").ToString())
	}
	if slices.Contains(actions, "noauditlog") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("noauditlog").ToString())
	}
	if s.SecRuleMetadata.Msg != "" {
		auxSlice = append(auxSlice, "msg:'" + s.SecRuleMetadata.Msg + "'")
	}
	if slices.Contains(actions, "logdata") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("logdata").ToString())
	}
	for _, tag := range s.SecRuleMetadata.Tags {
		auxSlice = append(auxSlice, "tag:'" + tag + "'")
	}
	// if s.SecRuleMetadata.Tag != "" {
	// 	auxSlice = append(auxSlice, "msg:'" + s.SecRuleMetadata.Msg + "'")
	// }
	if slices.Contains(actions, "sanitiseArg") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("sanitiseArg").ToString())
	}
	if slices.Contains(actions, "sanitiseRequestHeader") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("sanitiseRequestHeader").ToString())
	}
	if slices.Contains(actions, "sanitiseMatched") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("sanitiseMatched").ToString())
	}
	if slices.Contains(actions, "sanitiseMatchedBytes") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("sanitiseMatchedBytes").ToString())
	}
	if slices.Contains(actions, "ctl") {
		ctlActions := s.SeclangActions.GetActionsByKey("ctl")
		for _, action := range ctlActions {
			auxSlice = append(auxSlice, action.ToString())
		}
	}
	if s.SecRuleMetadata.Ver != "" {
		auxSlice = append(auxSlice, "ver:'" + s.SecRuleMetadata.Ver + "'")
	}
	if s.SecRuleMetadata.Severity != "" {
		auxSlice = append(auxSlice, "severity:'" + s.SecRuleMetadata.Severity + "'")
	}
	if slices.Contains(actions, "multiMatch") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("multiMatch").ToString())
	}
	if slices.Contains(actions, "initcol") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("initcol").ToString())
	}
	if slices.Contains(actions, "setenv") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("setenv").ToString())
	}
	if slices.Contains(actions, "setvar") {
		setvarActions := s.SeclangActions.GetActionsByKey("setvar")
		for _, action := range setvarActions {
			auxSlice = append(auxSlice, action.ToString())
		}
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("setvar").ToString())
	}
	if slices.Contains(actions, "expirevar") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("expirevar").ToString())
	}
	if slices.Contains(actions, "chain") {
		chainedRule = true
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("chain").ToString())
	}
	if slices.Contains(actions, "skip") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("skip").ToString())
	}
	if slices.Contains(actions, "skipAfter") {
		auxSlice = append(auxSlice, s.SeclangActions.GetActionByKey("skipAfter").ToString())
	}
	for i, action := range auxSlice {
		if i == 0 {
			result += " \\\n" + initialString + "\t\""
		} else {
			result += auxString
		}
		result += action
		if i == len(auxSlice) - 1 {
			result += "\""
		} else {
			result += endString
		}
	}
	result += "\n"
	if chainedRule {
		result += s.ChainedRule.ToSeclangWithParam(initialString + "\t")
	}
	return result
}


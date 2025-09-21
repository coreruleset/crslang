package listener

import (
	"fmt"

	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

// Helper functions to convert string to action types
func stringToDisruptiveAction(s string) (types.DisruptiveAction, error) {
	switch s {
	case "allow":
		return types.Allow, nil
	case "block":
		return types.Block, nil
	case "deny":
		return types.Deny, nil
	case "drop":
		return types.Drop, nil
	case "pass":
		return types.Pass, nil
	case "pause":
		return types.Pause, nil
	case "proxy":
		return types.Proxy, nil
	case "redirect":
		return types.Redirect, nil
	default:
		return types.Unknown, nil
	}
}

func stringToNonDisruptiveAction(s string) (types.NonDisruptiveAction, error) {
	switch s {
	case "append":
		return types.Append, nil
	case "auditlog":
		return types.AuditLog, nil
	case "capture":
		return types.Capture, nil
	case "ctl":
		return types.Ctl, nil
	case "deprecatevar":
		return types.DeprecateVar, nil
	case "exec":
		return types.Exec, nil
	case "expirevar":
		return types.ExpireVar, nil
	case "initcol":
		return types.InitCol, nil
	case "log":
		return types.Log, nil
	case "logdata":
		return types.LogData, nil
	case "multiMatch":
		return types.MultiMatch, nil
	case "noauditlog":
		return types.NoAuditLog, nil
	case "nolog":
		return types.NoLog, nil
	case "prepend":
		return types.Prepend, nil
	case "sanitiseArg":
		return types.SanitiseArg, nil
	case "sanitiseMatched":
		return types.SanitiseMatched, nil
	case "sanitiseMatchedBytes":
		return types.SanitiseMatchedBytes, nil
	case "sanitiseRequestHeader":
		return types.SanitiseRequestHeader, nil
	case "sanitiseResponseHeader":
		return types.SanitiseResponseHeader, nil
	case "setuid":
		return types.SetUid, nil
	case "setrsc":
		return types.SetRsc, nil
	case "setsid":
		return types.SetSid, nil
	case "setenv":
		return types.SetEnv, nil
	case "setvar":
		return types.SetVar, nil
	default:
		return types.NonDisruptiveUnknown, nil
	}
}

func stringToFlowAction(s string) (types.FlowAction, error) {
	switch s {
	case "chain":
		return types.Chain, nil
	case "skip":
		return types.Skip, nil
	case "skipAfter":
		return types.SkipAfter, nil
	default:
		return types.FlowUnknown, nil
	}
}

func stringToDataAction(s string) (types.DataAction, error) {
	switch s {
	case "data":
		return types.Data, nil
	case "status":
		return types.Status, nil
	case "xmlns":
		return types.XLMNS, nil
	default:
		return types.DataUnknown, nil
	}
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parser.Disruptive_action_onlyContext) {
	action, err := stringToDisruptiveAction(ctx.GetText())
	if err != nil {
		panic(fmt.Sprintf("failed to parse disruptive action: %v", err))
	}
	err = l.currentDirective.GetActions().SetDisruptiveActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to set disruptive action: %v", err))
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parser.Non_disruptive_action_onlyContext) {
	action, err := stringToNonDisruptiveAction(ctx.GetText())
	if err != nil {
		panic(fmt.Sprintf("failed to parse non-disruptive action: %v", err))
	}
	err = l.currentDirective.GetActions().AddNonDisruptiveActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to add non-disruptive action: %v", err))
	}
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parser.Flow_action_onlyContext) {
	action, err := stringToFlowAction(ctx.GetText())
	if err != nil {
		panic(fmt.Sprintf("failed to parse flow action: %v", err))
	}
	err = l.currentDirective.GetActions().AddFlowActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to add flow action: %v", err))
	}
	l.previousDirective = l.currentDirective
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parser.Disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		action, err := stringToDisruptiveAction(ctx.GetText())
		if err != nil {
			panic(fmt.Sprintf("failed to parse disruptive action: %v", err))
		}
		err = l.currentDirective.GetActions().SetDisruptiveActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to set disruptive action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parser.Non_disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		action, err := stringToNonDisruptiveAction(ctx.GetText())
		if err != nil {
			panic(fmt.Sprintf("failed to parse non-disruptive action: %v", err))
		}
		err = l.currentDirective.GetActions().AddNonDisruptiveActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to add non-disruptive action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parser.Flow_action_with_paramsContext) {
	l.setParam = func(value string) {
		action, err := stringToFlowAction(ctx.GetText())
		if err != nil {
			panic(fmt.Sprintf("failed to parse flow action: %v", err))
		}
		err = l.currentDirective.GetActions().AddFlowActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to add flow action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parser.Data_action_with_paramsContext) {
	l.setParam = func(value string) {
		action, err := stringToDataAction(ctx.GetText())
		if err != nil {
			panic(fmt.Sprintf("failed to parse data action: %v", err))
		}
		err = l.currentDirective.GetActions().AddDataActionWithParams(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to add data action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterAction_value_types(ctx *parser.Action_value_typesContext) {
	if l.setParam != nil {
		l.setParam(ctx.GetText())
		l.setParam = doNothingFuncString
	}
}

package listener

import (
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

// Helper functions to convert string to action types
func stringToDisruptiveAction(s string) types.DisruptiveAction {
	return types.DisruptiveAction(s)
}

func stringToNonDisruptiveAction(s string) types.NonDisruptiveAction {
	return types.NonDisruptiveAction(s)
}

func stringToFlowAction(s string) types.FlowAction {
	return types.FlowAction(s)
}

func stringToDataAction(s string) types.DataAction {
	return types.DataAction(s)
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parser.Disruptive_action_onlyContext) {
	l.currentDirective.GetActions().SetDisruptiveActionOnly(stringToDisruptiveAction(ctx.GetText()))
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parser.Non_disruptive_action_onlyContext) {
	l.currentDirective.GetActions().AddNonDisruptiveActionOnly(stringToNonDisruptiveAction(ctx.GetText()))
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parser.Flow_action_onlyContext) {
	l.currentDirective.GetActions().AddFlowActionOnly(stringToFlowAction(ctx.GetText()))
	l.previousDirective = l.currentDirective
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parser.Disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().SetDisruptiveActionWithParam(stringToDisruptiveAction(ctx.GetText()), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parser.Non_disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().AddNonDisruptiveActionWithParam(stringToNonDisruptiveAction(ctx.GetText()), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parser.Flow_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().AddFlowActionWithParam(stringToFlowAction(ctx.GetText()), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parser.Data_action_with_paramsContext) {
	l.setParam = func(value string) {
		err := l.currentDirective.GetActions().AddDataActionWithParams(stringToDataAction(ctx.GetText()), value)
		if err != nil {
			panic(err)
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterAction_value_types(ctx *parser.Action_value_typesContext) {
	if l.setParam != nil {
		l.setParam(ctx.GetText())
		l.setParam = doNothingFuncString
	}
}

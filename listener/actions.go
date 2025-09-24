package listener

import (
	"fmt"

	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parser.Disruptive_action_onlyContext) {
	action := types.StringToDisruptiveAction(ctx.GetText())
	err := l.currentDirective.GetActions().SetDisruptiveActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to set disruptive action: %v", err))
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parser.Non_disruptive_action_onlyContext) {
	action := types.StringToNonDisruptiveAction(ctx.GetText())
	err := l.currentDirective.GetActions().AddNonDisruptiveActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to add non-disruptive action: %v", err))
	}
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parser.Flow_action_onlyContext) {
	action := types.StringToFlowAction(ctx.GetText())
	err := l.currentDirective.GetActions().AddFlowActionOnly(action)
	if err != nil {
		panic(fmt.Sprintf("failed to add flow action: %v", err))
	}
	l.previousDirective = l.currentDirective
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parser.Disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		action := types.StringToDisruptiveAction(ctx.GetText())
		err := l.currentDirective.GetActions().SetDisruptiveActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to set disruptive action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parser.Non_disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		action := types.StringToNonDisruptiveAction(ctx.GetText())
		err := l.currentDirective.GetActions().AddNonDisruptiveActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to add non-disruptive action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parser.Flow_action_with_paramsContext) {
	l.setParam = func(value string) {
		action := types.StringToFlowAction(ctx.GetText())
		err := l.currentDirective.GetActions().AddFlowActionWithParam(action, value)
		if err != nil {
			panic(fmt.Sprintf("failed to add flow action with param: %v", err))
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parser.Data_action_with_paramsContext) {
	l.setParam = func(value string) {
		action := types.StringToDataAction(ctx.GetText())
		err := l.currentDirective.GetActions().AddDataActionWithParams(action, value)
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

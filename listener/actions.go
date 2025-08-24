package listener

import "github.com/coreruleset/seclang_parser/parser"

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parser.Disruptive_action_onlyContext) {
	l.currentDirective.GetActions().SetDisruptiveActionOnly(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parser.Non_disruptive_action_onlyContext) {
	l.currentDirective.GetActions().AddNonDisruptiveActionOnly(ctx.GetText())
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parser.Flow_action_onlyContext) {
	l.currentDirective.GetActions().AddFlowActionOnly(ctx.GetText())
	l.previousDirective = l.currentDirective
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parser.Disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().SetDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parser.Non_disruptive_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().AddNonDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parser.Flow_action_with_paramsContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetActions().AddFlowActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parser.Data_action_with_paramsContext) {
	l.setParam = func(value string) {
		err := l.currentDirective.GetActions().AddDataActionWithParams(ctx.GetText(), value)
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

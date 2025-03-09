package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterRules_directive(ctx *parsing.Rules_directiveContext) {
	l.currentDirective = types.NewSecRule()
	l.currentDirective.(*types.SecRule).SetOperatorName("rx")
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecRule))
		l.previousDirective = nil
	} else {
		l.appendDirective = func() {
			l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.SecRule))
		}
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComment
}

func (l *ExtendedSeclangParserListener) EnterVariables(ctx *parsing.VariablesContext) {
	l.targetDirective = l.currentDirective.(*types.SecRule)
}

func (l *ExtendedSeclangParserListener) EnterOperator_name(ctx *parsing.Operator_nameContext) {
	err := l.currentDirective.(*types.SecRule).SetOperatorName(ctx.GetText())
	if err != nil {
		panic(err)
	}
}

func (l *ExtendedSeclangParserListener) EnterOperator_value(ctx *parsing.Operator_valueContext) {
	l.currentDirective.(*types.SecRule).SetOperatorValue(ctx.GetText())
}

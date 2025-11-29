package listener

import (
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

func (l *ExtendedSeclangParserListener) EnterEngine_config_rule_directive(ctx *parser.Engine_config_rule_directiveContext) {
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

func (l *ExtendedSeclangParserListener) EnterVariables(ctx *parser.VariablesContext) {
	l.targetDirective = l.currentDirective.(*types.SecRule)
}

func (l *ExtendedSeclangParserListener) EnterOperator_name(ctx *parser.Operator_nameContext) {
	err := l.currentDirective.(*types.SecRule).SetOperatorName(ctx.GetText())
	if err != nil {
		panic(err)
	}
}

func (l *ExtendedSeclangParserListener) EnterOperator_value(ctx *parser.Operator_valueContext) {
	l.currentDirective.(*types.SecRule).SetOperatorValue(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterOperator_not(ctx *parser.Operator_notContext) {
	l.currentDirective.(*types.SecRule).SetOperatorNot(true)
}

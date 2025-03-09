package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterRule_script_directive(ctx *parsing.Rule_script_directiveContext) {
	l.currentDirective = types.NewSecRuleScript()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecRuleScript))
		l.previousDirective = nil
	} else {
		l.currentFunctionToAppendDirective = func() {
			l.Configuration.Directives = append(l.Configuration.Directives, l.currentDirective.(*types.SecRuleScript))
		}
	}
	l.currentFunctionToAppendComment = l.currentDirective.GetMetadata().SetComment
}

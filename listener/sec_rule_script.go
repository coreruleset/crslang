package listener

import (
	"github.com/coreruleset/seclang_parser/parser"
	"github.com/coreruleset/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterRule_script_directive(ctx *parser.Rule_script_directiveContext) {
	l.currentDirective = types.NewSecRuleScript()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecRuleScript))
		l.previousDirective = nil
	} else {
		l.appendDirective = func() {
			l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.SecRuleScript))
		}
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComment
}

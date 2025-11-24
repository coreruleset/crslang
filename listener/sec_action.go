package listener

import (
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

// SecDefaultAction
func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_default_action(ctx *parser.Config_dir_sec_default_actionContext) {
	l.currentDirective = types.NewDefaultAction()
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.currentDirective.(*types.DefaultAction))
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComments
}

// SecAction
func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_action(ctx *parser.Config_dir_sec_actionContext) {
	l.currentDirective = types.NewSecAction()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecAction))
		l.previousDirective = nil
	} else {
		l.appendDirective = func() {
			l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.SecAction))
		}
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComments
}

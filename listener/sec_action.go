package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

// SecDefaultAction
func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_default_action(ctx *parsing.Config_dir_sec_default_actionContext) {
	l.currentDirective = types.NewDefaultAction()
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.currentDirective.(*types.DefaultAction))
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComment
}

// SecAction
func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_action(ctx *parsing.Config_dir_sec_actionContext) {
	l.currentDirective = types.NewSecAction()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecAction))
		l.previousDirective = nil
	} else {
		l.appendDirective = func() {
			l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.SecAction))
		}
	}
	l.appendComment = l.currentDirective.GetMetadata().SetComment
}

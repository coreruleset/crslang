package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterValues(ctx *parsing.ValuesContext) {
	l.setParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_literal(ctx *parsing.String_literalContext) {
	if l.setParam != nil {
		l.setParam(ctx.GetText())
		l.setParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterFile_path(ctx *parsing.File_pathContext) {
	l.currentDirective.(*types.SecRuleScript).ScriptPath = ctx.GetText()
}

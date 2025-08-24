package listener

import (
	"github.com/coreruleset/seclang_parser/parser"
	"github.com/coreruleset/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterValues(ctx *parser.ValuesContext) {
	l.setParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_literal(ctx *parser.String_literalContext) {
	if l.setParam != nil {
		l.setParam(ctx.GetText())
		l.setParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterFile_path(ctx *parser.File_pathContext) {
	l.currentDirective.(*types.SecRuleScript).ScriptPath = ctx.GetText()
}

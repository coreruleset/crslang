package listener

import "github.com/coreruleset/seclang_parser/parser"

func (l *ExtendedSeclangParserListener) EnterACTION_ID(ctx *parser.ACTION_IDContext) {
	l.setParam = l.currentDirective.GetMetadata().SetId
}

func (l *ExtendedSeclangParserListener) EnterACTION_PHASE(ctx *parser.ACTION_PHASEContext) {
	l.setParam = l.currentDirective.GetMetadata().SetPhase
}

func (l *ExtendedSeclangParserListener) EnterACTION_MSG(ctx *parser.ACTION_MSGContext) {
	l.setParam = l.currentDirective.GetMetadata().SetMsg
}

func (l *ExtendedSeclangParserListener) EnterACTION_MATURITY(ctx *parser.ACTION_MATURITYContext) {
	l.setParam = l.currentDirective.GetMetadata().SetMaturity
}

func (l *ExtendedSeclangParserListener) EnterACTION_REV(ctx *parser.ACTION_REVContext) {
	l.setParam = l.currentDirective.GetMetadata().SetRev
}

func (l *ExtendedSeclangParserListener) EnterACTION_SEVERITY(ctx *parser.ACTION_SEVERITYContext) {
	l.setParam = l.currentDirective.GetMetadata().SetSeverity
}

func (l *ExtendedSeclangParserListener) EnterACTION_TAG(ctx *parser.ACTION_TAGContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetMetadata().AddTag(value)
	}
}

func (l *ExtendedSeclangParserListener) EnterACTION_VER(ctx *parser.ACTION_VERContext) {
	l.setParam = l.currentDirective.GetMetadata().SetVer
}

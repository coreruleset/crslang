package listener

import "gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"

func (l *ExtendedSeclangParserListener) EnterACTION_ID(ctx *parsing.ACTION_IDContext) {
	l.setParam = l.currentDirective.GetMetadata().SetId
}

func (l *ExtendedSeclangParserListener) EnterACTION_PHASE(ctx *parsing.ACTION_PHASEContext) {
	l.setParam = l.currentDirective.GetMetadata().SetPhase
}

func (l *ExtendedSeclangParserListener) EnterACTION_MSG(ctx *parsing.ACTION_MSGContext) {
	l.setParam = l.currentDirective.GetMetadata().SetMsg
}

func (l *ExtendedSeclangParserListener) EnterACTION_MATURITY(ctx *parsing.ACTION_MATURITYContext) {
	l.setParam = l.currentDirective.GetMetadata().SetMaturity
}

func (l *ExtendedSeclangParserListener) EnterACTION_REV(ctx *parsing.ACTION_REVContext) {
	l.setParam = l.currentDirective.GetMetadata().SetRev
}

func (l *ExtendedSeclangParserListener) EnterACTION_SEVERITY(ctx *parsing.ACTION_SEVERITYContext) {
	l.setParam = l.currentDirective.GetMetadata().SetSeverity
}

func (l *ExtendedSeclangParserListener) EnterACTION_TAG(ctx *parsing.ACTION_TAGContext) {
	l.setParam = func(value string) {
		l.currentDirective.GetMetadata().AddTag(value)
	}
}

func (l *ExtendedSeclangParserListener) EnterACTION_VER(ctx *parsing.ACTION_VERContext) {
	l.setParam = l.currentDirective.GetMetadata().SetVer
}

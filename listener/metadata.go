package listener

import "gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"

func (l *ExtendedSeclangParserListener) EnterACTION_ID(ctx *parsing.ACTION_IDContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetId
}

func (l *ExtendedSeclangParserListener) EnterACTION_PHASE(ctx *parsing.ACTION_PHASEContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetPhase
}

func (l *ExtendedSeclangParserListener) EnterACTION_MSG(ctx *parsing.ACTION_MSGContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetMsg
}

func (l *ExtendedSeclangParserListener) EnterACTION_MATURITY(ctx *parsing.ACTION_MATURITYContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetMaturity
}

func (l *ExtendedSeclangParserListener) EnterACTION_REV(ctx *parsing.ACTION_REVContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetRev
}

func (l *ExtendedSeclangParserListener) EnterACTION_SEVERITY(ctx *parsing.ACTION_SEVERITYContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetSeverity
}

func (l *ExtendedSeclangParserListener) EnterACTION_TAG(ctx *parsing.ACTION_TAGContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.GetMetadata().AddTag(value)
	}
}

func (l *ExtendedSeclangParserListener) EnterACTION_VER(ctx *parsing.ACTION_VERContext) {
	l.currentFunctionToSetParam = l.currentDirective.GetMetadata().SetVer
}

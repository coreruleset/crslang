package listener

import (
	"github.com/coreruleset/crslang/parsing"
	"github.com/coreruleset/crslang/types"
)

type BaseDirective interface {
	GetMetadata() types.Metadata
	GetActions() *types.SeclangActions
	AddTransformation(transformation string) error
	AppendChainedDirective(directive types.ChainableDirective)
}

type TargetDirective interface {
	AddVariable(variable string, excluded bool) error
	AddCollection(collection, value string, excluded bool, asCount bool) error
}

type BaseChainableDirective interface {
	BaseDirective
	types.ChainableDirective
}

type ExtendedSeclangParserListener struct {
	*parsing.BaseSecLangParserListener
	comment                string
	appendComment          func(value string)
	setParam               func(value string)
	addVariable            func() error
	appendDirective        func()
	configurationDirective *types.ConfigurationDirective
	targetDirective        TargetDirective
	currentDirective       BaseDirective
	previousDirective      BaseDirective
	removeDirective        types.RemoveRuleDirective
	idRange                types.IdRange
	updateTargetDirective  *types.UpdateTargetDirective
	varName                string
	varValue               string
	varExcluded            bool
	varCount               bool
	parameter              string
	DirectiveList          *types.DirectiveList
	ConfigurationList      types.ConfigurationList
}

func doNothingFunc() {}

func doNothingFuncString(value string) {}

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parsing.ConfigurationContext) {
	l.DirectiveList = new(types.DirectiveList)
	l.setParam = doNothingFuncString
	l.appendDirective = doNothingFunc
	l.appendComment = func(value string) {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, types.CommentMetadata{Comment: value})
	}
	l.previousDirective = nil
}

func (l *ExtendedSeclangParserListener) ExitConfiguration(ctx *parsing.ConfigurationContext) {
	if l.DirectiveList != nil && (len(l.DirectiveList.Directives) > 0 || l.DirectiveList.Marker.Name != "") {
		l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.DirectiveList)
	}
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
	if l.comment != "" {
		l.appendComment(l.comment)
		l.comment = ""
	}
	l.appendComment = func(value string) {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, types.CommentMetadata{Comment: value})
	}
	l.appendDirective()
	l.appendDirective = doNothingFunc
}

func (l *ExtendedSeclangParserListener) EnterComment(ctx *parsing.CommentContext) {
	l.comment = ctx.GetText()
}

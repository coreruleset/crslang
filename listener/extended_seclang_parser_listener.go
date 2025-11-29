package listener

import (
	"strings"

	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
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
	*parser.BaseSecLangParserListener
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

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parser.ConfigurationContext) {
	l.DirectiveList = new(types.DirectiveList)
	l.setParam = doNothingFuncString
	l.appendDirective = doNothingFunc
	l.appendComment = func(value string) {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, types.CommentMetadata{Comment: value})
	}
	l.previousDirective = nil
}

func (l *ExtendedSeclangParserListener) ExitConfiguration(ctx *parser.ConfigurationContext) {
	if l.DirectiveList != nil && (len(l.DirectiveList.Directives) > 0 || l.DirectiveList.Marker.Name != "") {
		l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.DirectiveList)
	}
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parser.StmtContext) {
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

func (l *ExtendedSeclangParserListener) EnterComment(ctx *parser.CommentContext) {
	// ctx.COMMENT() can be nil if there is only a HASH without comment text
	if ctx.COMMENT() != nil {
		// Remove leading space after the hash if any
		l.comment += strings.TrimPrefix(ctx.COMMENT().GetText(), " ") + "\n"
		// Ignore initial empty comment lines because YAML libraries do not work well with them
	} else if l.comment != "" {
		l.comment += "\n"
	}
}

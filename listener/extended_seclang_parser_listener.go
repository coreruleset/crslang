package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

type AuxDirective interface {
	GetMetadata() types.Metadata
	GetActions() *types.SeclangActions
	AddTransformation(transformation string) error
	AppendChainedDirective(directive types.ChainableDirective)
}

type TargetDirective interface {
	AddVariable(variable string) error
	AddCollection(collection, value string) error
}

type AuxChainableDirective interface {
	AuxDirective
	types.ChainableDirective
}

type ExtendedSeclangParserListener struct {
	*parsing.BaseSecLangParserListener
	currentComment                   string
	currentFunctionToAppendComment   func(value string)
	currentFunctionToSetParam        func(value string)
	currentFunctionToAddVariable     func() error
	currentFunctionToAppendDirective func()
	currentConfigurationDirective    *types.ConfigurationDirective
	targetDirective                  TargetDirective
	currentDirective                 AuxDirective
	previousDirective                AuxDirective
	removeDirective                  types.RemoveRuleDirective
	idRange                          types.IdRange
	updateTargetDirective            *types.UpdateTargetDirective
	varName                          string
	varValue                         string
	currentParameter                 string
	Configuration                    *types.DirectiveList
	ConfigurationList                types.ConfigurationList
}

func doNothingFunc() {}

func doNothingFuncString(value string) {}

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parsing.ConfigurationContext) {
	l.Configuration = new(types.DirectiveList)
	l.currentFunctionToSetParam = doNothingFuncString
	l.currentFunctionToAppendDirective = doNothingFunc
	l.currentFunctionToAppendComment = func(value string) {
		l.Configuration.Directives = append(l.Configuration.Directives, types.CommentMetadata{Comment: value})
	}
	l.previousDirective = nil
}

func (l *ExtendedSeclangParserListener) ExitConfiguration(ctx *parsing.ConfigurationContext) {
	l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.Configuration)
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
	if l.currentComment != "" {
		l.currentFunctionToAppendComment(l.currentComment)
		l.currentComment = ""
	}
	l.currentFunctionToAppendComment = func(value string) {
		l.Configuration.Directives = append(l.Configuration.Directives, types.CommentMetadata{Comment: value})
	}
	l.currentFunctionToAppendDirective()
	l.currentFunctionToAppendDirective = doNothingFunc
}

func (l *ExtendedSeclangParserListener) EnterComment(ctx *parsing.CommentContext) {
	l.currentComment = ctx.GetText()
}

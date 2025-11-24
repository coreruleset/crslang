package listener

import (
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_with_param(ctx *parser.Engine_config_directive_with_paramContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComments
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.configurationDirective)
	}
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_sec_cache_transformations(ctx *parser.Engine_config_sec_cache_transformationsContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComments
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = func(value2 string) {
			l.configurationDirective.Parameter += " " + value2
			l.setParam = doNothingFuncString
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterOption_list(ctx *parser.Option_listContext) {
	l.setParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_engine_config_directive(ctx *parser.String_engine_config_directiveContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComments
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.configurationDirective)
	}
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
}

// SecMarker
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parser.Sec_marker_directiveContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	err := l.configurationDirective.SetName(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.appendComment = l.configurationDirective.GetMetadata().SetComments
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
	l.appendDirective = func() {
		l.DirectiveList.Marker = *l.configurationDirective
		l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.DirectiveList)
		l.DirectiveList = new(types.DirectiveList)
	}
}

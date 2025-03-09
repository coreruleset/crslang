package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_with_param(ctx *parsing.Engine_config_directive_with_paramContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComment
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.configurationDirective)
	}
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_sec_cache_transformations(ctx *parsing.Engine_config_sec_cache_transformationsContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComment
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = func(value2 string) {
			l.configurationDirective.Parameter += " " + value2
			l.setParam = doNothingFuncString
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterOption_list(ctx *parsing.Option_listContext) {
	l.setParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_engine_config_directive(ctx *parsing.String_engine_config_directiveContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	l.configurationDirective.SetName(ctx.GetText())
	l.appendComment = l.configurationDirective.GetMetadata().SetComment
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, *l.configurationDirective)
	}
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
}

// SecMarker
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parsing.Sec_marker_directiveContext) {
	l.configurationDirective = types.NewConfigurationDirective()
	err := l.configurationDirective.SetName(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.appendComment = l.configurationDirective.GetMetadata().SetComment
	l.setParam = func(value string) {
		l.configurationDirective.Parameter = value
		l.setParam = doNothingFuncString
	}
	l.appendDirective = func() {
		l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.DirectiveList)
		l.DirectiveList = new(types.DirectiveList)
		l.DirectiveList.Marker = *l.configurationDirective
	}
}

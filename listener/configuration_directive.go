package listener

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_with_param(ctx *parsing.Engine_config_directive_with_paramContext) {
	// fmt.Println("String engine config directive: ", ctx.GetText())
	l.currentConfigurationDirective = types.NewConfigurationDirective()
	l.currentConfigurationDirective.SetName(ctx.GetText())
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.GetMetadata().SetComment
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentConfigurationDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_sec_cache_transformations(ctx *parsing.Engine_config_sec_cache_transformationsContext) {
	l.currentConfigurationDirective = types.NewConfigurationDirective()
	l.currentConfigurationDirective.SetName(ctx.GetText())
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.GetMetadata().SetComment
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = func(value2 string) {
			l.currentConfigurationDirective.Parameter += " " + value2
			l.currentFunctionToSetParam = doNothingFuncString
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterOption_list(ctx *parsing.Option_listContext) {
	l.currentFunctionToSetParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_engine_config_directive(ctx *parsing.String_engine_config_directiveContext) {
	l.currentConfigurationDirective = types.NewConfigurationDirective()
	l.currentConfigurationDirective.SetName(ctx.GetText())
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.GetMetadata().SetComment
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentConfigurationDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
}

// SecMarker
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parsing.Sec_marker_directiveContext) {
	l.currentConfigurationDirective = types.NewConfigurationDirective()
	err := l.currentConfigurationDirective.SetName(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.GetMetadata().SetComment
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
	l.currentFunctionToAppendDirective = func() {
		l.ConfigurationList.DirectiveList = append(l.ConfigurationList.DirectiveList, *l.Configuration)
		l.Configuration = new(types.DirectiveList)
		l.Configuration.Marker = *l.currentConfigurationDirective
	}
}

package main

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

type AuxDirective interface {
	SetId(value string)
	SetPhase(value string)
	SetMsg(value string)
	SetMaturity(value string)
	SetRev(value string)
	SetSeverity(value string)
	AddTag(value string)
	SetVer(value string)
	SetDisruptiveActionWithParam(action, value string)
	SetDisruptiveActionOnly(action string)
	AddNonDisruptiveActionWithParam(action, param string)
	AddNonDisruptiveActionOnly(action string)
	AddFlowActionWithParam(action, param string)
	AddFlowActionOnly(action string)
	AddDataActionWithParams(action, param string)
	AddTransformation(transformation string)
	SetComment(value string)
}

type ExtendedSeclangParserListener struct {
	*parsing.BaseSecLangParserListener
	currentComment 		  	  string
	currentFunctionToAppendComment func(value string)
	currentFunctionToSetParam func(value string)
	currentFunctionToAppendDirective func()
	currentConfigurationDirective *types.ConfigurationDirective
	currentDirective          AuxDirective
	currentParameter 		  string
	chainedNextRule 		  *types.SecRule
	Configuration 		*types.Configuration
	ConfigurationList 		types.ConfigurationList
}

func doNothingFunc(){}

func doNothingFuncString(value string) {}

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parsing.ConfigurationContext) {
	l.Configuration = new(types.Configuration)
	l.currentFunctionToSetParam = doNothingFuncString
	l.currentFunctionToAppendDirective = doNothingFunc
	l.currentFunctionToAppendComment = func (value string) {
		l.Configuration.Directives = append(l.Configuration.Directives, types.CommentMetadata{Comment: value})
	}
	l.chainedNextRule = nil
}

func (l *ExtendedSeclangParserListener) ExitConfiguration(ctx *parsing.ConfigurationContext) {
	l.ConfigurationList.Configurations = append(l.ConfigurationList.Configurations, *l.Configuration)
}

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_default_action(ctx *parsing.Config_dir_sec_default_actionContext) {
	l.currentDirective = new(types.SecDefaultAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentDirective.(*types.SecDefaultAction))
	}
	l.currentFunctionToAppendComment = l.currentDirective.SetComment
}

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_action(ctx *parsing.Config_dir_sec_actionContext) {
	l.currentDirective = new(types.SecAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentDirective.(*types.SecAction))
	}
	l.currentFunctionToAppendComment = l.currentDirective.SetComment
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_with_param(ctx *parsing.Engine_config_directive_with_paramContext) {
	// fmt.Println("String engine config directive: ", ctx.GetText())
	l.currentConfigurationDirective = new(types.ConfigurationDirective)
	l.currentConfigurationDirective.DirectiveName = ctx.GetText()
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.SetComment
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentConfigurationDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterValues(ctx *parsing.ValuesContext) {
	// fmt.Println("Config value types: ", ctx.GetText())
	l.currentFunctionToSetParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_sec_cache_transformations(ctx *parsing.Engine_config_sec_cache_transformationsContext) {
	l.currentConfigurationDirective = new(types.ConfigurationDirective)
	l.currentConfigurationDirective.DirectiveName = ctx.GetText()
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.SetComment
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

func (l *ExtendedSeclangParserListener) EnterRules_directive(ctx *parsing.Rules_directiveContext) {
	if l.chainedNextRule != nil {
		l.currentDirective = l.chainedNextRule
		l.chainedNextRule = nil
	} else {
		l.currentDirective = new(types.SecRule)
		l.currentDirective.(*types.SecRule).SetOperatorName("rx")
		l.currentFunctionToAppendDirective = func() {
			l.Configuration.Directives = append(l.Configuration.Directives, *l.currentDirective.(*types.SecRule))
		}
	}
	l.currentFunctionToAppendComment = l.currentDirective.SetComment
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
	if l.currentComment != "" {
		l.currentFunctionToAppendComment(l.currentComment)
		l.currentComment = ""
	}
	l.currentFunctionToAppendComment = func (value string) {
		l.Configuration.Directives = append(l.Configuration.Directives, types.CommentMetadata{Comment: value})
	}
	// fmt.Printf("Current directive: %v\n", l.currentDirective)
	l.currentFunctionToAppendDirective()
	l.currentFunctionToAppendDirective = doNothingFunc
}

func (l *ExtendedSeclangParserListener) EnterACTION_ID(ctx *parsing.ACTION_IDContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetId
}

func (l *ExtendedSeclangParserListener) EnterACTION_PHASE(ctx *parsing.ACTION_PHASEContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetPhase
}

func (l *ExtendedSeclangParserListener) EnterACTION_MSG(ctx *parsing.ACTION_MSGContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetMsg
}

func (l *ExtendedSeclangParserListener) EnterACTION_MATURITY(ctx *parsing.ACTION_MATURITYContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetMaturity
}

func (l *ExtendedSeclangParserListener) EnterACTION_REV(ctx *parsing.ACTION_REVContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetRev
}

func (l *ExtendedSeclangParserListener) EnterACTION_SEVERITY(ctx *parsing.ACTION_SEVERITYContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetSeverity
}

func (l *ExtendedSeclangParserListener) EnterACTION_TAG(ctx *parsing.ACTION_TAGContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.AddTag(value)
	}
}

func (l *ExtendedSeclangParserListener) EnterACTION_VER(ctx *parsing.ACTION_VERContext) {
	l.currentFunctionToSetParam = l.currentDirective.SetVer
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parsing.Disruptive_action_onlyContext) {
	l.currentDirective.SetDisruptiveActionOnly(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parsing.Non_disruptive_action_onlyContext) {
	l.currentDirective.AddNonDisruptiveActionOnly(ctx.GetText())
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parsing.Flow_action_onlyContext) {
	l.currentDirective.AddFlowActionOnly(ctx.GetText())
	l.currentDirective.(*types.SecRule).ChainedRule = new(types.SecRule)
	l.chainedNextRule = l.currentDirective.(*types.SecRule).ChainedRule
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parsing.Disruptive_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.SetDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parsing.Non_disruptive_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.AddNonDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parsing.Flow_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.AddFlowActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parsing.Data_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.AddDataActionWithParams(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterTransformation_action_value(ctx *parsing.Transformation_action_valueContext) {
	l.currentDirective.AddTransformation(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterAction_value_types(ctx *parsing.Action_value_typesContext) {
	if l.currentFunctionToSetParam != nil {
		l.currentFunctionToSetParam(ctx.GetText())
		l.currentFunctionToSetParam = doNothingFuncString
	} 
/* 	else {
		fmt.Println("No function to set param yet")
	} */
}

func (l *ExtendedSeclangParserListener) EnterString_literal(ctx *parsing.String_literalContext) {
	if l.currentFunctionToSetParam != nil {
		l.currentFunctionToSetParam(ctx.GetText())
		l.currentFunctionToSetParam = doNothingFuncString
	} 
/* 	else {
		fmt.Println("No function to set param yet")
	} */
}

func (l *ExtendedSeclangParserListener) EnterVar_stmt(ctx *parsing.Var_stmtContext) {
	l.currentDirective.(*types.SecRule).AddVariable(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterOperator_name(ctx *parsing.Operator_nameContext) {
	l.currentDirective.(*types.SecRule).SetOperatorName(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterOperator_value(ctx *parsing.Operator_valueContext) {

	l.currentDirective.(*types.SecRule).SetOperatorValue(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_engine_config_directive(ctx *parsing.String_engine_config_directiveContext) {
	l.currentConfigurationDirective = new(types.ConfigurationDirective)
	l.currentConfigurationDirective.DirectiveName = ctx.GetText()
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.SetComment
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentConfigurationDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
}

// This is the event function for the secmarker directive
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parsing.Sec_marker_directiveContext) {
	// fmt.Println("Sec marker directive: ", ctx.GetText())
	// l.currentParameter = ctx.GetText()
	l.currentConfigurationDirective = new(types.ConfigurationDirective)
	l.currentConfigurationDirective.DirectiveName = ctx.GetText()
	l.currentFunctionToAppendComment = l.currentConfigurationDirective.SetComment
	l.currentFunctionToSetParam = func(value string) {
		l.currentConfigurationDirective.Parameter = value
		l.currentFunctionToSetParam = doNothingFuncString
	}
	l.currentFunctionToAppendDirective = func() {
		l.ConfigurationList.Configurations = append(l.ConfigurationList.Configurations, *l.Configuration)
		l.Configuration = new(types.Configuration)
		l.Configuration.Marker = *l.currentConfigurationDirective
	}
}

func (l *ExtendedSeclangParserListener) EnterComment(ctx *parsing.CommentContext) {
	l.currentComment = ctx.GetText()
}
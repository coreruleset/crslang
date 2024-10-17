package main

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

type SeclangDirective interface {
	SetId(value string)
	SetPhase(value string)
	SetMsg(value string)
	SetMaturity(value string)
	SetRev(value string)
	SetSeverity(value string)
	SetVer(value string)
	SetDisruptiveActionWithParam(action, value string)
	SetDisruptiveActionOnly(action string)
	AddNonDisruptiveActionWithParam(action, param string)
	AddNonDisruptiveActionOnly(action string)
	AddFlowActionWithParam(action, param string)
	AddFlowActionOnly(action string)
	AddDataActionWithParams(action, param string)
	AddTransformation(transformation string)
	AddVariable(variable string)
	SetOperatorName(name string)
	SetOperatorValue(value string)
}

type ExtendedSeclangParserListener struct {
	*parsing.BaseSecLangParserListener
	currentFunctionToSetParam func(value string)
	currentFunctionToAppendDirective func()
	currentDirective          SeclangDirective
	currentParameter 		  string
	Configuration             types.Configuration
}

func doNothingFunc(){}

func doNothingFuncString(value string) {}

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parsing.ConfigurationContext) {
	l.Configuration = types.Configuration{
		ConfigDirectives: make(map[string]string),
	}
	l.currentFunctionToSetParam = doNothingFuncString
	l.currentFunctionToAppendDirective = doNothingFunc
}

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_default_action(ctx *parsing.Config_dir_sec_default_actionContext) {
	l.currentDirective = new(types.SecDefaultAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.DefaultActions = append(l.Configuration.DefaultActions, *l.currentDirective.(*types.SecDefaultAction))
	}
}

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_action(ctx *parsing.Config_dir_sec_actionContext) {
	l.currentDirective = new(types.SecAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.SecActions = append(l.Configuration.SecActions, *l.currentDirective.(*types.SecAction))
	}
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_only_param(ctx *parsing.Engine_config_directive_only_paramContext) {
	// fmt.Println("String engine config directive: ", ctx.GetText())
	l.currentParameter = ctx.GetText()
	l.currentFunctionToSetParam = func(value string) {
		l.Configuration.ConfigDirectives[l.currentParameter] = value
		l.currentParameter = ""
		l.currentFunctionToSetParam = doNothingFuncString
	}
}

func (l *ExtendedSeclangParserListener) EnterValues(ctx *parsing.ValuesContext) {
	// fmt.Println("Config value types: ", ctx.GetText())
	l.currentFunctionToSetParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterEngine_config_sec_cache_transformations(ctx *parsing.Engine_config_sec_cache_transformationsContext) {
	l.currentParameter = ctx.GetText()
	l.currentFunctionToSetParam = func(value string) {
		l.Configuration.ConfigDirectives[l.currentParameter] = value
		l.currentFunctionToSetParam = func(value2 string) {
			l.Configuration.ConfigDirectives[l.currentParameter] += " " + value2
			l.currentParameter = ""
			l.currentFunctionToSetParam = doNothingFuncString
		}
	}
}

func (l *ExtendedSeclangParserListener) EnterOption_list(ctx *parsing.Option_listContext) {
	l.currentFunctionToSetParam(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterRules_directive(ctx *parsing.Rules_directiveContext) {
	l.currentDirective = new(types.SecRule)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.SecRules = append(l.Configuration.SecRules, *l.currentDirective.(*types.SecRule))
	}
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
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
	l.currentDirective.AddVariable(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterOperator_name(ctx *parsing.Operator_nameContext) {
	l.currentDirective.SetOperatorName(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterOperator_value(ctx *parsing.Operator_valueContext) {

	l.currentDirective.SetOperatorValue(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterString_engine_config_directive(ctx *parsing.String_engine_config_directiveContext) {
	// fmt.Println("String engine config directive: ", ctx.GetText())
	l.currentParameter = ctx.GetText()
}

// This is the event function for the secmarker directive
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parsing.Sec_marker_directiveContext) {
	// fmt.Println("Sec marker directive: ", ctx.GetText())
	l.currentParameter = ctx.GetText()
}
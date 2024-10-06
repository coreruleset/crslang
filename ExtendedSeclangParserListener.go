package main

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

type ExtendedSeclangParserListener struct {
	*parsing.BaseSecLangParserListener
	currentFunctionToSetParam func(value string)
	currentFunctionToAppendDirective func()
	currentDirective          types.SeclangDirective
	currentParameter 		 string
	Configuration             types.Configuration
}

func (l *ExtendedSeclangParserListener) EnterConfiguration(ctx *parsing.ConfigurationContext) {
	l.Configuration = types.Configuration{
		ConfigDirectives: make(map[string]string),
	}
	l.currentFunctionToSetParam = nil
	l.currentFunctionToAppendDirective = func() {}
}

func (l *ExtendedSeclangParserListener) EnterCONFIG_DIR_SEC_DEFAULT_ACTION(ctx *parsing.CONFIG_DIR_SEC_DEFAULT_ACTIONContext) {
	l.currentDirective = new(types.SecDefaultAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.DefaultActions = append(l.Configuration.DefaultActions, *l.currentDirective.(*types.SecDefaultAction))
	}
}

func (l *ExtendedSeclangParserListener) EnterCONFIG_DIR_SEC_ACTION(ctx *parsing.CONFIG_DIR_SEC_ACTIONContext) {
	l.currentDirective = new(types.SecAction)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.SecActions = append(l.Configuration.SecActions, *l.currentDirective.(*types.SecAction))
	}
}

func (l *ExtendedSeclangParserListener) EnterRules_directive(ctx *parsing.Rules_directiveContext) {
	l.currentDirective = new(types.SecRule)
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.SecRules = append(l.Configuration.SecRules, *l.currentDirective.(*types.SecRule))
	}
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
	l.currentFunctionToAppendDirective()
	l.currentFunctionToAppendDirective = func () {}
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

func (l *ExtendedSeclangParserListener) EnterAction_value(ctx *parsing.Action_valueContext) {
	if l.currentFunctionToSetParam != nil {
		l.currentFunctionToSetParam(ctx.GetText())
		l.currentFunctionToSetParam = nil
	} 
/* 	else {
		fmt.Println("No function to set param yet")
	} */
}

func (l *ExtendedSeclangParserListener) EnterVariable_name(ctx *parsing.Variable_nameContext) {
	l.currentParameter = ctx.GetText()
}

// How do we want to store the variable?
func (l *ExtendedSeclangParserListener) EnterCollection_element_or_regexp(ctx *parsing.Collection_element_or_regexpContext) {
	l.currentParameter += ":" + ctx.GetText()
}

func (l *ExtendedSeclangParserListener) ExitVar_stmt(ctx *parsing.Var_stmtContext) {
	l.currentDirective.AddVariable(l.currentParameter)
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

func (l *ExtendedSeclangParserListener) EnterEngine_config_directive_only_param(ctx *parsing.Engine_config_directive_only_paramContext) {
	// fmt.Println("String engine config directive: ", ctx.GetText())
	l.currentParameter = ctx.GetText()
}

func (l *ExtendedSeclangParserListener) EnterConfig_value_types(ctx *parsing.Config_value_typesContext) {
	// fmt.Println("Config value types: ", ctx.GetText())
	l.Configuration.ConfigDirectives[l.currentParameter] = ctx.GetText()
}
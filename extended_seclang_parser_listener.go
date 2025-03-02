package main

import (
	"strconv"

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
	// chainedNextRule 		  *AuxChainableDirective
	Configuration     *types.DirectiveList
	ConfigurationList types.ConfigurationList
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

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_default_action(ctx *parsing.Config_dir_sec_default_actionContext) {
	l.currentDirective = types.NewDefaultAction()
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, *l.currentDirective.(*types.DefaultAction))
	}
	l.currentFunctionToAppendComment = l.currentDirective.GetMetadata().SetComment
}

func (l *ExtendedSeclangParserListener) EnterConfig_dir_sec_action(ctx *parsing.Config_dir_sec_actionContext) {
	l.currentDirective = types.NewSecAction()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecAction))
		l.previousDirective = nil
	} else {
		l.currentFunctionToAppendDirective = func() {
			l.Configuration.Directives = append(l.Configuration.Directives, l.currentDirective.(*types.SecAction))
		}
	}
	l.currentFunctionToAppendComment = l.currentDirective.GetMetadata().SetComment
}

func (l *ExtendedSeclangParserListener) EnterRule_script_directive(ctx *parsing.Rule_script_directiveContext) {
	l.currentDirective = types.NewSecRuleScript()
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecRuleScript))
		l.previousDirective = nil
	} else {
		l.currentFunctionToAppendDirective = func() {
			l.Configuration.Directives = append(l.Configuration.Directives, l.currentDirective.(*types.SecRuleScript))
		}
	}
	l.currentFunctionToAppendComment = l.currentDirective.GetMetadata().SetComment
}

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

func (l *ExtendedSeclangParserListener) EnterValues(ctx *parsing.ValuesContext) {
	// fmt.Println("Config value types: ", ctx.GetText())
	l.currentFunctionToSetParam(ctx.GetText())
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

func (l *ExtendedSeclangParserListener) EnterRules_directive(ctx *parsing.Rules_directiveContext) {
	l.currentDirective = types.NewSecRule()
	l.currentDirective.(*types.SecRule).SetOperatorName("rx")
	if l.previousDirective != nil {
		l.previousDirective.AppendChainedDirective(l.currentDirective.(*types.SecRule))
		l.previousDirective = nil
	} else {
		l.currentFunctionToAppendDirective = func() {
			l.Configuration.Directives = append(l.Configuration.Directives, l.currentDirective.(*types.SecRule))
		}
	}
	l.currentFunctionToAppendComment = l.currentDirective.GetMetadata().SetComment
}

func (l *ExtendedSeclangParserListener) ExitStmt(ctx *parsing.StmtContext) {
	if l.currentComment != "" {
		l.currentFunctionToAppendComment(l.currentComment)
		l.currentComment = ""
	}
	l.currentFunctionToAppendComment = func(value string) {
		l.Configuration.Directives = append(l.Configuration.Directives, types.CommentMetadata{Comment: value})
	}
	// fmt.Printf("Current directive: %v\n", l.currentDirective)
	l.currentFunctionToAppendDirective()
	l.currentFunctionToAppendDirective = doNothingFunc
}

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

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_only(ctx *parsing.Disruptive_action_onlyContext) {
	l.currentDirective.GetActions().SetDisruptiveActionOnly(ctx.GetText())
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_only(ctx *parsing.Non_disruptive_action_onlyContext) {
	l.currentDirective.GetActions().AddNonDisruptiveActionOnly(ctx.GetText())
}

// Event for chain action, the only flow action with no parameters is Chain
func (l *ExtendedSeclangParserListener) EnterFlow_action_only(ctx *parsing.Flow_action_onlyContext) {
	l.currentDirective.GetActions().AddFlowActionOnly(ctx.GetText())
	l.previousDirective = l.currentDirective
	// l.currentDirective.(*types.SecRule).ChainedRule = new(types.ChainableDirective)
}

func (l *ExtendedSeclangParserListener) EnterDisruptive_action_with_params(ctx *parsing.Disruptive_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.GetActions().SetDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterNon_disruptive_action_with_params(ctx *parsing.Non_disruptive_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.GetActions().AddNonDisruptiveActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterFlow_action_with_params(ctx *parsing.Flow_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.GetActions().AddFlowActionWithParam(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterData_action_with_params(ctx *parsing.Data_action_with_paramsContext) {
	l.currentFunctionToSetParam = func(value string) {
		l.currentDirective.GetActions().AddDataActionWithParams(ctx.GetText(), value)
	}
}

func (l *ExtendedSeclangParserListener) EnterTransformation_action_value(ctx *parsing.Transformation_action_valueContext) {
	err := l.currentDirective.AddTransformation(ctx.GetText())
	if err != nil {
		panic(err)
	}
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

func (l *ExtendedSeclangParserListener) EnterFile_path(ctx *parsing.File_pathContext) {
	l.currentDirective.(*types.SecRuleScript).ScriptPath = ctx.GetText()
}

func (l *ExtendedSeclangParserListener) EnterVar_stmt(ctx *parsing.Var_stmtContext) {
	l.varName = ""
	l.varValue = ""
}

func (l *ExtendedSeclangParserListener) EnterVariable_enum(ctx *parsing.Variable_enumContext) {
	l.varName = ctx.GetText()
	l.currentFunctionToAddVariable = func() error {
		err := l.targetDirective.AddVariable(l.varName)
		return err
	}
}

func (l *ExtendedSeclangParserListener) EnterCollection_enum(ctx *parsing.Collection_enumContext) {
	l.varName = ctx.GetText()
	l.currentFunctionToAddVariable = func() error {
		err := l.targetDirective.AddCollection(l.varName, "")
		return err
	}
}

func (l *ExtendedSeclangParserListener) EnterCollection_value(ctx *parsing.Collection_valueContext) {
	l.varValue = ctx.GetText()
	l.currentFunctionToAddVariable = func() error {
		err := l.targetDirective.AddCollection(l.varName, l.varValue)
		return err
	}
}

func (l *ExtendedSeclangParserListener) ExitVar_stmt(ctx *parsing.Var_stmtContext) {
	err := l.currentFunctionToAddVariable()
	if err != nil {
		panic(err)
	}
	l.varName = ""
	l.varValue = ""
}

func (l *ExtendedSeclangParserListener) EnterVariables(ctx *parsing.VariablesContext) {
	l.targetDirective = l.currentDirective.(*types.SecRule)
}

func (l *ExtendedSeclangParserListener) EnterOperator_name(ctx *parsing.Operator_nameContext) {
	err := l.currentDirective.(*types.SecRule).SetOperatorName(ctx.GetText())
	if err != nil {
		panic(err)
	}
}

func (l *ExtendedSeclangParserListener) EnterOperator_value(ctx *parsing.Operator_valueContext) {
	l.currentDirective.(*types.SecRule).SetOperatorValue(ctx.GetText())
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

// This is the event function for the secmarker directive
func (l *ExtendedSeclangParserListener) EnterSec_marker_directive(ctx *parsing.Sec_marker_directiveContext) {
	// fmt.Println("Sec marker directive: ", ctx.GetText())
	// l.currentParameter = ctx.GetText()
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

func (l *ExtendedSeclangParserListener) EnterComment(ctx *parsing.CommentContext) {
	l.currentComment = ctx.GetText()
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_msg(ctx *parsing.Remove_rule_by_msgContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.currentFunctionToAppendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.currentFunctionToSetParam = func(value string) {
		l.removeDirective.Msgs = append(l.removeDirective.Msgs, value)
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_tag(ctx *parsing.Remove_rule_by_tagContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.currentFunctionToAppendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.currentFunctionToSetParam = func(value string) {
		l.removeDirective.Tags = append(l.removeDirective.Tags, value)
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id(ctx *parsing.Remove_rule_by_idContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.currentFunctionToAppendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id_int(ctx *parsing.Remove_rule_by_id_intContext) {
	id, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.removeDirective.Ids = append(l.removeDirective.Ids, id)
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id_int_range(ctx *parsing.Remove_rule_by_id_int_rangeContext) {
	l.idRange = types.IdRange{}
}

func (l *ExtendedSeclangParserListener) EnterRange_start(ctx *parsing.Range_startContext) {
	start, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.idRange.Start = start
}

func (l *ExtendedSeclangParserListener) EnterRange_end(ctx *parsing.Range_endContext) {
	end, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.idRange.End = end
}

func (l *ExtendedSeclangParserListener) ExitRemove_rule_by_id_int_range(ctx *parsing.Remove_rule_by_id_int_rangeContext) {
	l.removeDirective.IdRanges = append(l.removeDirective.IdRanges, l.idRange)
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_rules(ctx *parsing.Update_target_rulesContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.currentFunctionToAppendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	switch ctx.GetText() {
	case "SecRuleUpdateTargetById":
		l.currentFunctionToSetParam = func(value string) {
			id, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			l.updateTargetDirective.Ids = append(l.updateTargetDirective.Ids, id)
		}
	case "SecRuleUpdateTargetByTag":
		l.currentFunctionToSetParam = func(value string) {
			l.updateTargetDirective.Tags = append(l.updateTargetDirective.Tags, value)
		}
	case "SecRuleUpdateTargetByMsg":
		l.currentFunctionToSetParam = func(value string) {
			l.updateTargetDirective.Msgs = append(l.updateTargetDirective.Msgs, value)
		}
	}

	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.updateTargetDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_rules_values(ctx *parsing.Update_target_rules_valuesContext) {
	l.currentFunctionToSetParam(ctx.GetText())
	l.currentFunctionToSetParam = doNothingFuncString
}

func (l *ExtendedSeclangParserListener) EnterUpdate_variables(ctx *parsing.Update_variablesContext) {
	l.targetDirective = l.updateTargetDirective
}

package listener

import (
	"strconv"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_id(ctx *parsing.Update_target_by_idContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.currentFunctionToAppendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.updateTargetDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		id, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		l.updateTargetDirective.Ids = append(l.updateTargetDirective.Ids, id)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_tag(ctx *parsing.Update_target_by_tagContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.currentFunctionToAppendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.updateTargetDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.updateTargetDirective.Tags = append(l.updateTargetDirective.Tags, value)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_msg(ctx *parsing.Update_target_by_msgContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.currentFunctionToAppendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.currentFunctionToAppendDirective = func() {
		l.Configuration.Directives = append(l.Configuration.Directives, l.updateTargetDirective)
	}
	l.currentFunctionToSetParam = func(value string) {
		l.updateTargetDirective.Msgs = append(l.updateTargetDirective.Msgs, value)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_rules_values(ctx *parsing.Update_target_rules_valuesContext) {
	l.currentFunctionToSetParam(ctx.GetText())
	l.currentFunctionToSetParam = doNothingFuncString
}

func (l *ExtendedSeclangParserListener) EnterUpdate_variables(ctx *parsing.Update_variablesContext) {
	l.targetDirective = l.updateTargetDirective
}

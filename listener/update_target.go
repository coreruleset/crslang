package listener

import (
	"strconv"

	"github.com/coreruleset/seclang_parser/parser"
	"github.com/coreruleset/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_id(ctx *parser.Update_target_by_idContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.appendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.updateTargetDirective)
	}
	l.setParam = func(value string) {
		id, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		l.updateTargetDirective.Ids = append(l.updateTargetDirective.Ids, id)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_tag(ctx *parser.Update_target_by_tagContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.appendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.updateTargetDirective)
	}
	l.setParam = func(value string) {
		l.updateTargetDirective.Tags = append(l.updateTargetDirective.Tags, value)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_by_msg(ctx *parser.Update_target_by_msgContext) {
	l.updateTargetDirective = types.NewUpdateTargetDirective()
	l.appendComment = func(comment string) {
		l.updateTargetDirective.Metadata.Comment = comment
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.updateTargetDirective)
	}
	l.setParam = func(value string) {
		l.updateTargetDirective.Msgs = append(l.updateTargetDirective.Msgs, value)
	}
}

func (l *ExtendedSeclangParserListener) EnterUpdate_target_rules_values(ctx *parser.Update_target_rules_valuesContext) {
	l.setParam(ctx.GetText())
	l.setParam = doNothingFuncString
}

func (l *ExtendedSeclangParserListener) EnterUpdate_variables(ctx *parser.Update_variablesContext) {
	l.targetDirective = l.updateTargetDirective
}

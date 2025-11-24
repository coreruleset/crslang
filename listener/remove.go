package listener

import (
	"strconv"

	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_msg(ctx *parser.Remove_rule_by_msgContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.appendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.setParam = func(value string) {
		l.removeDirective.Msgs = append(l.removeDirective.Msgs, value)
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_tag(ctx *parser.Remove_rule_by_tagContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.appendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.setParam = func(value string) {
		l.removeDirective.Tags = append(l.removeDirective.Tags, value)
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id(ctx *parser.Remove_rule_by_idContext) {
	l.removeDirective = types.RemoveRuleDirective{
		Kind: types.Remove,
	}
	l.appendComment = func(comment string) {
		l.removeDirective.Metadata.Comment = comment
	}
	l.appendDirective = func() {
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.removeDirective)
	}
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id_int(ctx *parser.Remove_rule_by_id_intContext) {
	id, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.removeDirective.Ids = append(l.removeDirective.Ids, id)
}

func (l *ExtendedSeclangParserListener) EnterRemove_rule_by_id_int_range(ctx *parser.Remove_rule_by_id_int_rangeContext) {
	l.idRange = types.IdRange{}
}

func (l *ExtendedSeclangParserListener) EnterRange_start(ctx *parser.Range_startContext) {
	start, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.idRange.Start = start
}

func (l *ExtendedSeclangParserListener) EnterRange_end(ctx *parser.Range_endContext) {
	end, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.idRange.End = end
}

func (l *ExtendedSeclangParserListener) ExitRemove_rule_by_id_int_range(ctx *parser.Remove_rule_by_id_int_rangeContext) {
	l.removeDirective.IdRanges = append(l.removeDirective.IdRanges, l.idRange)
}

func (l *ExtendedSeclangParserListener) EnterString_remove_rules_values(ctx *parser.String_remove_rules_valuesContext) {
	l.setParam(ctx.GetText())
}

package listener

import (
	"fmt"
	"strconv"

	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
)

func (l *ExtendedSeclangParserListener) EnterUpdate_action_rule(ctx *parser.Update_action_ruleContext) {
	l.currentDirective = types.NewUpdateActionDirective()
	l.appendDirective = func() {
		fmt.Printf("Appending directive: %v\n", l.currentDirective.(*types.UpdateActionDirective).ToSeclang())
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.UpdateActionDirective))
	}
	l.appendComment = func(comments []string) {
		l.currentDirective.(*types.UpdateActionDirective).Comments = comments
	}
}

func (l *ExtendedSeclangParserListener) EnterId(ctx *parser.IdContext) {
	id, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.currentDirective.(*types.UpdateActionDirective).Id = id
}

package listener

import (
	"fmt"
	"strconv"

	"github.com/coreruleset/crslang/parsing"
	"github.com/coreruleset/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterUpdate_action_rule(ctx *parsing.Update_action_ruleContext) {
	l.currentDirective = types.NewUpdateActionDirective()
	l.appendDirective = func() {
		fmt.Printf("Appending directive: %v\n", l.currentDirective.(*types.UpdateActionDirective).ToSeclang())
		l.DirectiveList.Directives = append(l.DirectiveList.Directives, l.currentDirective.(*types.UpdateActionDirective))
	}
	l.appendComment = func(comment string) {
		l.currentDirective.(*types.UpdateActionDirective).Comment = comment
	}
}

func (l *ExtendedSeclangParserListener) EnterId(ctx *parsing.IdContext) {
	id, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}
	l.currentDirective.(*types.UpdateActionDirective).Id = id
}

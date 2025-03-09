package listener

import (
	"fmt"
	"strconv"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func (l *ExtendedSeclangParserListener) EnterUpdate_action_rule(ctx *parsing.Update_action_ruleContext) {
	l.currentDirective = types.NewUpdateActionDirective()
	l.currentFunctionToAppendDirective = func() {
		fmt.Printf("Appending directive: %v\n", l.currentDirective.(*types.UpdateActionDirective).ToSeclang())
		l.Configuration.Directives = append(l.Configuration.Directives, l.currentDirective.(*types.UpdateActionDirective))
	}
	l.currentFunctionToAppendComment = func(comment string) {
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

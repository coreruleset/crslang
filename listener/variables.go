package listener

import (
	"github.com/coreruleset/crslang/parsing"
)

func (l *ExtendedSeclangParserListener) EnterVar_stmt(ctx *parsing.Var_stmtContext) {
	l.varName = ""
	l.varValue = ""
}

func (l *ExtendedSeclangParserListener) EnterVar_not(ctx *parsing.Var_notContext) {
	l.varExcluded = true
}

func (l *ExtendedSeclangParserListener) EnterVar_count(ctx *parsing.Var_countContext) {
	l.varCount = true
}

func (l *ExtendedSeclangParserListener) EnterVariable_enum(ctx *parsing.Variable_enumContext) {
	l.varName = ctx.GetText()
	l.addVariable = func() error {
		err := l.targetDirective.AddVariable(l.varName, l.varExcluded)
		return err
	}
}

func (l *ExtendedSeclangParserListener) EnterCollection_enum(ctx *parsing.Collection_enumContext) {
	l.varName = ctx.GetText()
	l.addVariable = func() error {
		err := l.targetDirective.AddCollection(l.varName, "", l.varExcluded, l.varCount)
		return err
	}
}

func (l *ExtendedSeclangParserListener) EnterCollection_value(ctx *parsing.Collection_valueContext) {
	l.varValue = ctx.GetText()
	l.addVariable = func() error {
		err := l.targetDirective.AddCollection(l.varName, l.varValue, l.varExcluded, l.varCount)
		return err
	}
}

func (l *ExtendedSeclangParserListener) ExitVar_stmt(ctx *parsing.Var_stmtContext) {
	err := l.addVariable()
	if err != nil {
		panic(err)
	}
	l.varName = ""
	l.varValue = ""
	l.varExcluded = false
	l.varCount = false
}

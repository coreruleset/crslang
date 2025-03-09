package listener

import "gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"

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

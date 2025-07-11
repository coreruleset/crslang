package listener

import "github.com/coreruleset/crslang/parsing"

func (l *ExtendedSeclangParserListener) EnterTransformation_action_value(ctx *parsing.Transformation_action_valueContext) {
	err := l.currentDirective.AddTransformation(ctx.GetText())
	if err != nil {
		panic(err)
	}
}

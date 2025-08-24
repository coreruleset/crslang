package listener

import "github.com/coreruleset/seclang_parser/parser"

func (l *ExtendedSeclangParserListener) EnterTransformation_action_value(ctx *parser.Transformation_action_valueContext) {
	err := l.currentDirective.AddTransformation(ctx.GetText())
	if err != nil {
		panic(err)
	}
}

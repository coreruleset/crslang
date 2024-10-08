// Copyright 2023 Felipe Zipitria
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
)

type TreeShapeListener struct {
	*parsing.BaseSecLangParserListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (t *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

var testfiles = []string{
// 	"seclang_parser/testdata/test1.conf",
// 	"seclang_parser/testdata/test2.conf",
// 	"seclang_parser/testdata/test3.conf",
// 	"seclang_parser/testdata/test4.conf",
// 	"seclang_parser/testdata/test5.conf",
	// "seclang_parser/testdata/test6.conf",
// 	"seclang_parser/testdata/test7.conf",
// 	"seclang_parser/testdata/crs-setup.conf",
// 	"seclang_parser/testdata/REQUEST-901-INITIALIZATION.conf",
	// "seclang_parser/testdata/REQUEST-905-COMMON-EXCEPTIONS.conf",
	// "seclang_parser/testdata/REQUEST-911-METHOD-ENFORCEMENT.conf",
	// "seclang_parser/testdata/REQUEST-913-SCANNER-DETECTION.conf",
	"seclang_parser/testdata/REQUEST-920-PROTOCOL-ENFORCEMENT.conf",
}

func TestSampleFiles(t *testing.T) {
	for _, file := range testfiles {
		input, err := antlr.NewFileStream(file)
		if err != nil {
			t.Errorf("Error reading file %s", file)
			continue
		}
		lexer := parsing.NewSecLangLexer(input)
		for {
			token := lexer.NextToken()
			if token.GetTokenType() == antlr.TokenEOF {
				break
			}
			fmt.Printf("%s (%q)\n",
				lexer.SymbolicNames[token.GetTokenType()], token.GetText())
		}
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parsing.NewSecLangParser(stream)
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
		p.BuildParseTrees = true
		tree := p.Configuration()
		antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
	}
}

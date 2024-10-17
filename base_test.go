// Copyright 2023 Felipe Zipitria
// SPDX-License-Identifier: Apache-2.0,

package main

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/stretchr/testify/require"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
)


type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func NewCustomErrorListenerV2() *CustomErrorListener {
	return &CustomErrorListener{antlr.NewDefaultErrorListener(), make([]error, 0)}
}

func (c *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	var err error
	if offendingSymbol == nil {
		err = fmt.Errorf("Recognition error at line %d, column %d: %s", line, column, msg)
	} else {
		err = fmt.Errorf("Syntax error at line %d, column %d: %v", line, column, offendingSymbol)
	}
	c.Errors = append(c.Errors, err)
}

type TreeShapeListener struct {
	*parsing.BaseSecLangParserListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (t *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	// if you need to debug, enable this one below
	//fmt.Println(ctx.GetText())
}

var crsTestFiles = []string{
	"seclang_parser/testdata/crs/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf",
	"seclang_parser/testdata/crs/REQUEST-901-INITIALIZATION.conf",
	"seclang_parser/testdata/crs/REQUEST-905-COMMON-EXCEPTIONS.conf",
	"seclang_parser/testdata/crs/REQUEST-911-METHOD-ENFORCEMENT.conf",
	"seclang_parser/testdata/crs/REQUEST-913-SCANNER-DETECTION.conf",
	"seclang_parser/testdata/crs/REQUEST-920-PROTOCOL-ENFORCEMENT.conf",
	"seclang_parser/testdata/crs/REQUEST-921-PROTOCOL-ATTACK.conf",
	"seclang_parser/testdata/crs/REQUEST-922-MULTIPART-ATTACK.conf",
	"seclang_parser/testdata/crs/REQUEST-930-APPLICATION-ATTACK-LFI.conf",
	"seclang_parser/testdata/crs/REQUEST-931-APPLICATION-ATTACK-RFI.conf",
	"seclang_parser/testdata/crs/REQUEST-932-APPLICATION-ATTACK-RCE.conf",
	"seclang_parser/testdata/crs/REQUEST-933-APPLICATION-ATTACK-PHP.conf",
	"seclang_parser/testdata/crs/REQUEST-934-APPLICATION-ATTACK-GENERIC.conf",
	"seclang_parser/testdata/crs/REQUEST-941-APPLICATION-ATTACK-XSS.conf",
	"seclang_parser/testdata/crs/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
	"seclang_parser/testdata/crs/REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION.conf",
	"seclang_parser/testdata/crs/REQUEST-944-APPLICATION-ATTACK-JAVA.conf",
	"seclang_parser/testdata/crs/REQUEST-949-BLOCKING-EVALUATION.conf",
	"seclang_parser/testdata/crs/RESPONSE-950-DATA-LEAKAGES.conf",
	"seclang_parser/testdata/crs/RESPONSE-951-DATA-LEAKAGES-SQL.conf",
	"seclang_parser/testdata/crs/RESPONSE-952-DATA-LEAKAGES-JAVA.conf",
	"seclang_parser/testdata/crs/RESPONSE-953-DATA-LEAKAGES-PHP.conf",
	"seclang_parser/testdata/crs/RESPONSE-954-DATA-LEAKAGES-IIS.conf",
	// "seclang_parser/testdata/crs/RESPONSE-955-WEB-SHELLS.conf",
	"seclang_parser/testdata/crs/RESPONSE-959-BLOCKING-EVALUATION.conf",
	"seclang_parser/testdata/crs/RESPONSE-980-CORRELATION.conf",
}

var genericTests = map[string]struct {
	errorCount int
	comment    string
}{
	"seclang_parser/testdata/REQUEST-901-INITIALIZATION.conf": {
		0,
		"Test file for REQUEST-901-INITIALIZATION.conf",
	},
	"seclang_parser/testdata/crs-setup.conf": {
		0,
		"Test file for crs-setup.conf",
	},
	"seclang_parser/testdata/test1.conf": {
		0,
		"Test SecDefaultAction",
	},
	"seclang_parser/testdata/test2.conf": {
		0,
		"Test SecAction and SecCollectionTimeout",
	},
	"seclang_parser/testdata/test3.conf": {
		0,
		"test comment and secaction",
	},
	"seclang_parser/testdata/test4.conf": {
		0,
		"test redefining SecCollectionTimeout",
	},
	"seclang_parser/testdata/test5.conf": {
		0,
		"Test comments only file",
	},
	"seclang_parser/testdata/test_01_comment.conf": {
		0,
		"Test comments only file",
	},
	"seclang_parser/testdata/test_02_seccompsignature.conf": {
		0,
		"test SecComponentSignature",
	},
	"seclang_parser/testdata/test_03_secruleengine.conf": {
		0,
		"test SecRuleEngine",
	},
	"seclang_parser/testdata/test_04_directives.conf": {
		0,
		"test directives",
	},
	"seclang_parser/testdata/test_05_secaction.conf": {
		0,
		"test SecAction",
	},
	"seclang_parser/testdata/test_06_secaction2.conf": {
		0,
		"test SecAction with and without continuation",
	},
	"seclang_parser/testdata/test_07_secaction3.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_08_secaction4.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_09_secaction_ctl_01.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_10_secaction_ctl_02.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_11_secaction_ctl_03.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_12_secaction_ctl_04.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_13_secaction_ctl_05.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_14_secaction_ctl_06.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_15_secaction_01.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_16_secrule_01.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_17_secrule_02.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_18_secrule_03.conf": {
		1,
		"test should fail with non-existent operator",
	},
	"seclang_parser/testdata/test_19_secrule_04.conf": {
		0,
		"test SecAction with ctl",
	},
	"seclang_parser/testdata/test_20_secrule_05.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_21_secrule_06.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_22_secrule_07.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_23_secrule_08.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_24_secrule_09.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_25_secrule_10.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_26_secrule_11.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_27_secrule_12.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_28_secrule_13.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_29_secrule_14.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_30_secrule_15.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_31_secaction_ctl_07.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_32_secrule_16.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_33_secrule_16.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_directive_unknown.conf": {
		4,
		"test should fail with unknown directive",
	},
	"seclang_parser/testdata/test_34_xml.conf": {
		0,
		"",
	},
	"seclang_parser/testdata/test_35_all_directives.conf": {
		0,
		"",
	},
}

func TestSecLang(t *testing.T) {
	for file, data := range genericTests {
		t.Logf("Testing file %s", file)
		input, err := antlr.NewFileStream(file)
		if err != nil {
			t.Errorf("Error reading file %s", file)
			continue
		}
		lexer := parsing.NewSecLangLexer(input)

		lexerErrors := NewCustomErrorListenerV2()
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(lexerErrors)

		
		parserErrors := NewCustomErrorListenerV2()
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parsing.NewSecLangParser(stream)
		p.RemoveErrorListeners()
		p.AddErrorListener(parserErrors)

		p.BuildParseTrees = true
		tree := p.Configuration()

		antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)

		// for {
		// 	token := lexer.NextToken()
		// 	if token.GetTokenType() == antlr.TokenEOF {
		// 		break
		// 	}
		// 	fmt.Printf("%s (%q)\n",
		// 		lexer.SymbolicNames[token.GetTokenType()], token.GetText())
		// }
		if data.errorCount == 0 && len(lexerErrors.Errors) > 0 {
			t.Logf("Lexer %d errors found\n", len(lexerErrors.Errors))
			t.Logf("First error: %v\n", lexerErrors.Errors[0])
			// for _, e := range lexerErrors.Errors {
			// 	t.Logf("%v\n", e.Error())
			// }
		}
		if data.errorCount == 0 && len(parserErrors.Errors) > 0 {
			t.Logf("Parser %d errors found\n", len(parserErrors.Errors))
			t.Logf("First error: %v\n", parserErrors.Errors[0])
			// for _, e := range parserErrors.Errors {
			// 	t.Logf("%v\n", e.Error())
			// }
		}
		require.Equalf(t, data.errorCount, (len(lexerErrors.Errors) + len(parserErrors.Errors)), "Error count mismatch for file %s -> we want: %s\n", file, data.comment)
	}
}

func TestCRSLang(t *testing.T) {
	for _, file := range crsTestFiles {
		t.Logf("Testing file %s", file)
		input, err := antlr.NewFileStream(file)
		if err != nil {
			t.Fatalf("Error reading file %s", file)
		}
		lexer := parsing.NewSecLangLexer(input)

		lexerErrors := NewCustomErrorListenerV2()
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(lexerErrors)

		parserErrors := NewCustomErrorListenerV2()
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parsing.NewSecLangParser(stream)
		p.RemoveErrorListeners()
		p.AddErrorListener(parserErrors)
		p.BuildParseTrees = true
		tree := p.Configuration()

		antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)

		// for {
		// 	token := lexer.NextToken()
		// 	if token.GetTokenType() == antlr.TokenEOF {
		// 		break
		// 	}
		// 	t.Logf("%s (%q)",
		// 		lexer.SymbolicNames[token.GetTokenType()], token.GetText())
		// }
		if len(lexerErrors.Errors) > 0 {
			t.Logf("Lexer %d errors found\n", len(lexerErrors.Errors))
			t.Logf("First error: %v\n", lexerErrors.Errors[0])
			// for _, e := range lexerErrors.Errors {
			// 	t.Logf("%v\n", e.Error())
			// }
		}
		if len(parserErrors.Errors) > 0 {
			t.Logf("Parser %d errors found\n", len(parserErrors.Errors))
			t.Logf("First error: %v\n", parserErrors.Errors[0])
			// for _, e := range parserErrors.Errors {
			// 	t.Logf("%v\n", e.Error())
			// }
		}
		require.Equalf(t, 0, (len(lexerErrors.Errors) + len(parserErrors.Errors)), "Error count mismatch for file %s -> we want no errors\n", file)
	}
}
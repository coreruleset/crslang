package main

import (
	"io"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

var files = []string{
	"seclang_parser/testdata/test1.conf",
	"seclang_parser/testdata/test2.conf",
	"seclang_parser/testdata/test3.conf",
	"seclang_parser/testdata/test4.conf",
	"seclang_parser/testdata/test5.conf",
	"seclang_parser/testdata/test6.conf",
	"seclang_parser/testdata/test7.conf",
	"seclang_parser/testdata/crs-setup.conf",
	"seclang_parser/testdata/REQUEST-901-INITIALIZATION.conf",
	"seclang_parser/testdata/REQUEST-905-COMMON-EXCEPTIONS.conf",
	"seclang_parser/testdata/REQUEST-911-METHOD-ENFORCEMENT.conf",
	"seclang_parser/testdata/REQUEST-913-SCANNER-DETECTION.conf",
	"seclang_parser/testdata/REQUEST-920-PROTOCOL-ENFORCEMENT.conf",
}

func main() {
	resultConfigs := []types.Configuration{}
	for _, file := range files {
		input, err := antlr.NewFileStream(file)
		if err != nil {
			panic("Error reading file" + file)
		}
		lexer := parsing.NewSecLangLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parsing.NewSecLangParser(stream)
		start := p.Configuration()
		var listener ExtendedSeclangParserListener
		antlr.ParseTreeWalkerDefault.Walk(&listener, start)
		resultConfigs = append(resultConfigs, listener.Configuration)
	}
	yamlFile, err := yaml.Marshal(resultConfigs)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Printing yaml", string(yamlFile))

	f, err := os.Create("seclang.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.WriteString(f, string(yamlFile))
	if err != nil {
		panic(err)
	}

}

package main

import (
	"io"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gopkg.in/yaml.v3"
)

func main() {
	file := "seclang_parser/testdata/test7.conf"
	input, err := antlr.NewFileStream("seclang_parser/testdata/test7.conf")
	if err != nil {
		panic("Error reading file" + file)
	}
	lexer := parsing.NewSecLangLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parsing.NewSecLangParser(stream)
	start := p.Configuration()
	var listener ExtendedSeclangParserListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, start)


	yamlFile, err := yaml.Marshal(listener.Configuration)
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

package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/listener"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

var files = "seclang_parser/testdata/crs"

func main() {
	resultConfigs := []types.DirectiveList{}
	filepath.Walk(files, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".conf" {
			input, err := antlr.NewFileStream(path)
			if err != nil {
				panic("Error reading file" + path)
			}
			lexer := parsing.NewSecLangLexer(input)
			stream := antlr.NewCommonTokenStream(lexer, 0)
			p := parsing.NewSecLangParser(stream)
			start := p.Configuration()
			var seclangListener listener.ExtendedSeclangParserListener
			antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
			for i := range seclangListener.ConfigurationList.DirectiveList {
				seclangListener.ConfigurationList.DirectiveList[i].Id = strings.TrimSuffix(filepath.Base(info.Name()), filepath.Ext(info.Name()))
				if i > 0 {
					seclangListener.ConfigurationList.DirectiveList[i].Id += "_" + strconv.Itoa(i+1)
				}
			}
			resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
		}
		return nil
	})
	configList := types.ConfigurationList{DirectiveList: resultConfigs}

	err := printCRSLang(configList, "crslang.yaml")
	if err != nil {
		panic(err)
	}

	/* 	loadedConfigList := types.LoadDirectivesWithConditionsFromFile("crslang.yaml")
	   	yamlFile, err := yaml.Marshal(loadedConfigList.DirectiveList)
	   	if err != nil {
	   		panic(err)
	   	}

	   	writeToFile(yamlFile, "crslang_loaded.yaml") */
}

// printSeclang writes seclang format directives to a file
func printSeclang(configList types.ConfigurationList, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	seclangDirectives := types.ToSeclang(configList)

	_, err = io.WriteString(f, seclangDirectives)
	if err != nil {
		return err
	}

	return nil
}

// printCRSLang writes crslang format directives (directives with conditions) to a file
func printCRSLang(configList types.ConfigurationList, filename string) error {
	configListWithConditions := types.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.DirectiveList)
	if err != nil {
		return err
	}

	err = writeToFile(yamlFile, filename)

	return err
}

func writeToFile(payload []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(payload))
	if err != nil {
		return err
	}

	return nil
}

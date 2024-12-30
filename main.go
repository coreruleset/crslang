package main

import (
	"io"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/exporters"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

var files = []string{
	"seclang_parser/testdata/crs-setup.conf",
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
	"seclang_parser/testdata/crs/RESPONSE-959-BLOCKING-EVALUATION.conf",
	"seclang_parser/testdata/crs/RESPONSE-980-CORRELATION.conf",
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
		resultConfigs = append(resultConfigs, listener.ConfigurationList.Configurations...)
	}
	configList := types.ConfigurationList{Configurations: resultConfigs}

	err := printCRSLang(configList, "crslang.yaml")
	if err != nil {
		panic(err)
	}

	/* 	loadedConfigList := exporters.LoadDirectivesWithConditionsFromFile("crslang.yaml")
	   	yamlFile, err := yaml.Marshal(loadedConfigList.Configurations)
	   	if err != nil {
	   		panic(err)
	   	}

	   	writeToFile(yamlFile, "crslang_loaded.yaml") */
}

// printDirectivesWithLabels writes alias format directives to a file
func printDirectivesWithLabels(configList types.ConfigurationList, filename string) error {
	wrappedConfigList := exporters.ToDirectivesWithLabels(configList)

	yamlFile, err := yaml.Marshal(wrappedConfigList.Configurations)
	if err != nil {
		return err
	}

	err = writeToFile(yamlFile, filename)

	return err
}

// printSeclang writes seclang format directives to a file
func printSeclang(configList types.ConfigurationList, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	seclangDirectives := exporters.ToSeclang(configList)

	_, err = io.WriteString(f, seclangDirectives)
	if err != nil {
		return err
	}

	return nil
}

// printCRSLang writes crslang format directives (directives with conditions) to a file
func printCRSLang(configList types.ConfigurationList, filename string) error {
	configListWithConditions := exporters.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.Configurations)
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

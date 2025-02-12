package main

import (
	"os"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/exporters"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/parsing"
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

var testFiles = []string{
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

func TestLoadCRS(t *testing.T) {
	resultConfigs := []types.DirectiveList{}
	for _, file := range testFiles {
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
		resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)
	}
	configList := types.ConfigurationList{DirectiveList: resultConfigs}

	configListWithConditions := exporters.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.Configurations)
	if err != nil {
		t.Errorf("Error marshalling yaml: %v", err)
	}

	err = writeToFile(yamlFile, "tmp_crslang.yaml")

	defer os.Remove("tmp_crslang.yaml")

	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	loadedConfigList := exporters.LoadDirectivesWithConditionsFromFile("tmp_crslang.yaml")
	yamlLoadedFile, err := yaml.Marshal(loadedConfigList.Configurations)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	if string(yamlFile) != string(yamlLoadedFile) {
		t.Errorf("Error: loaded file is different from original. Expected string length: %v, got: %v", len(string(yamlFile)), len(string(yamlLoadedFile)))
	}
}

func TestFromCRSLangToSeclang(t *testing.T) {
	resultConfigs := []types.DirectiveList{}
	for _, file := range testFiles {
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
		resultConfigs = append(resultConfigs, listener.ConfigurationList.DirectiveList...)
	}
	configList := types.ConfigurationList{DirectiveList: resultConfigs}

	seclangDirectives := exporters.ToSeclang(configList)

	configListWithConditions := exporters.ToDirectiveWithConditions(configList)

	configListFromConditions := exporters.FromCRSLangToUnformattedDirectives(*configListWithConditions)

	seclangDirectivesFromConditions := exporters.ToSeclang(*configListFromConditions)

	if seclangDirectives != seclangDirectivesFromConditions {
		t.Errorf("Error in CRSLang to Seclang directives convertion. Expected length: %v, got: %v", len(seclangDirectives), len(seclangDirectivesFromConditions))
	}

}

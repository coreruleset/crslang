package main

import (
	"os"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/coreruleset/crslang/listener"
	"github.com/coreruleset/seclang_parser/parser"
	"github.com/coreruleset/crslang/types"
	"gopkg.in/yaml.v3"
)

var testFiles = []string{
	"testdata/crs-setup.conf",
	"testdata/crs/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf",
	"testdata/crs/REQUEST-901-INITIALIZATION.conf",
	"testdata/crs/REQUEST-905-COMMON-EXCEPTIONS.conf",
	"testdata/crs/REQUEST-911-METHOD-ENFORCEMENT.conf",
	"testdata/crs/REQUEST-913-SCANNER-DETECTION.conf",
	"testdata/crs/REQUEST-920-PROTOCOL-ENFORCEMENT.conf",
	"testdata/crs/REQUEST-921-PROTOCOL-ATTACK.conf",
	"testdata/crs/REQUEST-922-MULTIPART-ATTACK.conf",
	"testdata/crs/REQUEST-930-APPLICATION-ATTACK-LFI.conf",
	"testdata/crs/REQUEST-931-APPLICATION-ATTACK-RFI.conf",
	"testdata/crs/REQUEST-932-APPLICATION-ATTACK-RCE.conf",
	"testdata/crs/REQUEST-933-APPLICATION-ATTACK-PHP.conf",
	"testdata/crs/REQUEST-934-APPLICATION-ATTACK-GENERIC.conf",
	"testdata/crs/REQUEST-941-APPLICATION-ATTACK-XSS.conf",
	"testdata/crs/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
	"testdata/crs/REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION.conf",
	"testdata/crs/REQUEST-944-APPLICATION-ATTACK-JAVA.conf",
	"testdata/crs/REQUEST-949-BLOCKING-EVALUATION.conf",
	"testdata/crs/RESPONSE-950-DATA-LEAKAGES.conf",
	"testdata/crs/RESPONSE-951-DATA-LEAKAGES-SQL.conf",
	"testdata/crs/RESPONSE-952-DATA-LEAKAGES-JAVA.conf",
	"testdata/crs/RESPONSE-953-DATA-LEAKAGES-PHP.conf",
	"testdata/crs/RESPONSE-954-DATA-LEAKAGES-IIS.conf",
	"testdata/crs/RESPONSE-959-BLOCKING-EVALUATION.conf",
	"testdata/crs/RESPONSE-980-CORRELATION.conf",
}

func TestLoadCRS(t *testing.T) {
	resultConfigs := []types.DirectiveList{}
	for _, file := range testFiles {
		input, err := antlr.NewFileStream(file)
		if err != nil {
			panic("Error reading file" + file)
		}
		lexer := parser.NewSecLangLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSecLangParser(stream)
		start := p.Configuration()
		var seclangListener listener.ExtendedSeclangParserListener
		antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
		resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
	}
	configList := types.ConfigurationList{DirectiveList: resultConfigs}

	configListWithConditions := types.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.DirectiveList)
	if err != nil {
		t.Errorf("Error marshalling yaml: %v", err)
	}

	err = writeToFile(yamlFile, "tmp_crslang.yaml")

	defer os.Remove("tmp_crslang.yaml")

	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	loadedConfigList := types.LoadDirectivesWithConditionsFromFile("tmp_crslang.yaml")
	yamlLoadedFile, err := yaml.Marshal(loadedConfigList.DirectiveList)
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
		lexer := parser.NewSecLangLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSecLangParser(stream)
		start := p.Configuration()
		var seclangListener listener.ExtendedSeclangParserListener
		antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
		resultConfigs = append(resultConfigs, seclangListener.ConfigurationList.DirectiveList...)
	}
	configList := types.ConfigurationList{DirectiveList: resultConfigs}

	seclangDirectives := types.ToSeclang(configList)

	configListWithConditions := types.ToDirectiveWithConditions(configList)

	configListFromConditions := types.FromCRSLangToUnformattedDirectives(*configListWithConditions)

	seclangDirectivesFromConditions := types.ToSeclang(*configListFromConditions)

	if seclangDirectives != seclangDirectivesFromConditions {
		t.Errorf("Error in CRSLang to Seclang directives convertion. Expected length: %v, got: %v", len(seclangDirectives), len(seclangDirectivesFromConditions))
	}

}

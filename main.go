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
	// "seclang_parser/testdata/test1.conf",
	// "seclang_parser/testdata/test2.conf",
	// "seclang_parser/testdata/test3.conf",
	// "seclang_parser/testdata/test4.conf",
	// "seclang_parser/testdata/test5.conf",
	// "seclang_parser/testdata/test6.conf",
	// "seclang_parser/testdata/test7.conf",
	// "seclang_parser/testdata/crs-setup.conf",
	// "seclang_parser/testdata/crs/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf",
	"seclang_parser/testdata/crs/REQUEST-901-INITIALIZATION.conf",
	// "seclang_parser/testdata/crs/REQUEST-905-COMMON-EXCEPTIONS.conf",
	// "seclang_parser/testdata/crs/REQUEST-911-METHOD-ENFORCEMENT.conf",
	// "seclang_parser/testdata/crs/REQUEST-913-SCANNER-DETECTION.conf",
	// "seclang_parser/testdata/crs/REQUEST-920-PROTOCOL-ENFORCEMENT.conf",
	// "seclang_parser/testdata/crs/REQUEST-921-PROTOCOL-ATTACK.conf",
	// "seclang_parser/testdata/crs/REQUEST-922-MULTIPART-ATTACK.conf",
	// "seclang_parser/testdata/crs/REQUEST-930-APPLICATION-ATTACK-LFI.conf",
	// "seclang_parser/testdata/crs/REQUEST-931-APPLICATION-ATTACK-RFI.conf",
	// "seclang_parser/testdata/crs/REQUEST-932-APPLICATION-ATTACK-RCE.conf",
	// "seclang_parser/testdata/crs/REQUEST-933-APPLICATION-ATTACK-PHP.conf",
	// "seclang_parser/testdata/crs/REQUEST-934-APPLICATION-ATTACK-GENERIC.conf",
	// "seclang_parser/testdata/crs/REQUEST-941-APPLICATION-ATTACK-XSS.conf",
	// "seclang_parser/testdata/crs/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
	// "seclang_parser/testdata/crs/REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION.conf",
	// "seclang_parser/testdata/crs/REQUEST-944-APPLICATION-ATTACK-JAVA.conf",
	// "seclang_parser/testdata/crs/REQUEST-949-BLOCKING-EVALUATION.conf",
	// "seclang_parser/testdata/crs/RESPONSE-950-DATA-LEAKAGES.conf",
	// "seclang_parser/testdata/crs/RESPONSE-951-DATA-LEAKAGES-SQL.conf",
	// "seclang_parser/testdata/crs/RESPONSE-952-DATA-LEAKAGES-JAVA.conf",
	// "seclang_parser/testdata/crs/RESPONSE-953-DATA-LEAKAGES-PHP.conf",
	// "seclang_parser/testdata/crs/RESPONSE-954-DATA-LEAKAGES-IIS.conf",
	// // "seclang_parser/testdata/crs/RESPONSE-955-WEB-SHELLS.conf",
	// "seclang_parser/testdata/crs/RESPONSE-959-BLOCKING-EVALUATION.conf",
	// "seclang_parser/testdata/crs/RESPONSE-980-CORRELATION.conf",
	// "seclang_parser/testdata/test_01_comment.conf",
	// "seclang_parser/testdata/test_02_seccompsignature.conf",
	// "seclang_parser/testdata/test_03_secruleengine.conf",
	// "seclang_parser/testdata/test_04_directives.conf",
	// "seclang_parser/testdata/test_05_secaction.conf",
	// "seclang_parser/testdata/test_06_secaction2.conf",
	// "seclang_parser/testdata/test_07_secaction3.conf",
	// "seclang_parser/testdata/test_08_secaction4.conf",
	// "seclang_parser/testdata/test_09_secaction_ctl_01.conf",
	// "seclang_parser/testdata/test_10_secaction_ctl_02.conf",
	// "seclang_parser/testdata/test_11_secaction_ctl_03.conf",
	// "seclang_parser/testdata/test_12_secaction_ctl_04.conf",
	// "seclang_parser/testdata/test_13_secaction_ctl_05.conf",
	// "seclang_parser/testdata/test_14_secaction_ctl_06.conf",
	// "seclang_parser/testdata/test_15_secaction_01.conf",
	// "seclang_parser/testdata/test_16_secrule_01.conf",
	// "seclang_parser/testdata/test_17_secrule_02.conf",
	// "seclang_parser/testdata/test_19_secrule_04.conf",
	// "seclang_parser/testdata/test_20_secrule_05.conf",
	// "seclang_parser/testdata/test_21_secrule_06.conf",
	// "seclang_parser/testdata/test_22_secrule_07.conf",
	// "seclang_parser/testdata/test_23_secrule_08.conf",
	// "seclang_parser/testdata/test_24_secrule_09.conf",
	// "seclang_parser/testdata/test_25_secrule_10.conf",
	// "seclang_parser/testdata/test_26_secrule_11.conf",
	// "seclang_parser/testdata/test_27_secrule_12.conf",
	// "seclang_parser/testdata/test_28_secrule_13.conf",
	// "seclang_parser/testdata/test_29_secrule_14.conf",
	// "seclang_parser/testdata/test_30_secrule_15.conf",
	// "seclang_parser/testdata/test_31_secaction_ctl_07.conf",
	// "seclang_parser/testdata/test_32_secrule_16.conf",
	// "seclang_parser/testdata/test_33_secrule_16.conf",
	// "seclang_parser/testdata/test_34_xml.conf",
	// "seclang_parser/testdata/test_35_all_directives.conf",
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

	// ReadYAML()

	f, err = os.Create("seclang.conf")
	if err != nil {
		panic(err)
	}

	for _, config := range resultConfigs {
		for _, directive := range config.Directives {
				_, err = io.WriteString(f, directive.ToSeclang() + "\n")
			if err != nil {
				panic(err)
			}
		}
	}
}

// func ReadYAML() {
// 	yamlFile, err := os.ReadFile("seclang.yaml")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var configs []types.Configuration
// 	err = yaml.Unmarshal(yamlFile, &configs)
// 	if err != nil {
// 		panic(err)
// 	}

// 	f, err := os.Create("seclang.conf")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, config := range configs {
// 		for _, directive := range config.Directives {
// 				_, err = io.WriteString(f, directive.ToSeclang())
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// }

func ToSeclang(configs []types.Configuration) string {
	result := ""
	for _, config := range configs {
		for _, directive := range config.Directives {
			result += directive.ToSeclang()
		}
	}
	return result
}
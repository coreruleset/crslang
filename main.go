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

	err := PrintCRSLang(configList, "crslang.yaml")
	if err != nil {
		panic(err)
	}
}

// printDirectivesWithLabels writes alias format directives to a file
func printDirectivesWithLabels(configList types.ConfigurationList, filename string) error {
	wrappedConfigList := exporters.ToDirectivesWithLabels(configList)

	yamlFile, err := yaml.Marshal(wrappedConfigList.Configurations)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(yamlFile))
	if err != nil {
		return err
	}

	return nil
}

// yamlLoader is a auxiliary struct to load and iterate over the yaml file
type yamlLoader struct {
	Marker     exporters.ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
	Directives []yaml.Node                             `yaml:"directives,omitempty"`
}

// directiveLoader is a auxiliary struct to load directives
type directiveLoader struct {
	types.SecRuleMetadata `yaml:"metadata,omitempty"`
	types.Variables       `yaml:",inline"`
	types.Transformations `yaml:",inline"`
	types.Operator        `yaml:"operator"`
	types.SeclangActions  `yaml:"actions"`
	ScriptPath            string    `yaml:"scriptpath"`
	ChainedRule           yaml.Node `yaml:"chainedRule"`
}

// loadDirectivesWithLabels loads alias format directives from a yaml file
func loadDirectivesWithLabels(filename string) types.ConfigurationList{
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var configs []yamlLoader
	err = yaml.Unmarshal(yamlFile, &configs)
	var resultConfigs []types.Configuration
	for _, config := range configs {
		var directives []types.SeclangDirective
		for _, yamlDirective := range config.Directives {
			directive := directivesWithLabelsAux(yamlDirective)
			if directive == nil {
				panic("Unknown directive type")
			} else {
				directives = append(directives, directive)
			}
		}
		resultConfigs = append(resultConfigs, types.Configuration{Marker: config.Marker.ConfigurationDirective, Directives: directives})
	}
	return types.ConfigurationList{Configurations: resultConfigs}
}

// directivesWithLabelsAux is a recursive function to load directives
func directivesWithLabelsAux(yamlDirective yaml.Node) types.SeclangDirective {
	if yamlDirective.Kind != yaml.MappingNode {
		panic("Unknown format type")
	}
	switch yamlDirective.Content[0].Value {
	case "comment":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := types.CommentMetadata{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "configurationdirective":
		rawDirective, err := yaml.Marshal(yamlDirective.Content[1])
		if err != nil {
			panic(err)
		}
		directive := types.ConfigurationDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "secdefaultaction":
		rawDirective, err := yaml.Marshal(yamlDirective.Content[1])
		if err != nil {
			panic(err)
		}
		directive := types.SecDefaultAction{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "secaction":
		rawDirective, err := yaml.Marshal(yamlDirective.Content[1])
		if err != nil {
			panic(err)
		}
		loaderDirective := directiveLoader{}
		err = yaml.Unmarshal(rawDirective, &loaderDirective)
		if err != nil {
			panic(err)
		}
		directive := types.SecAction{
			SecRuleMetadata: loaderDirective.SecRuleMetadata,
			Transformations: loaderDirective.Transformations,
			SeclangActions:  loaderDirective.SeclangActions,
		}
		var chainedRule types.SeclangDirective
		if len(loaderDirective.ChainedRule.Content) > 0 {
			chainedRule = directivesWithLabelsAux(loaderDirective.ChainedRule)
			directive.ChainedRule = castChainableDirective(chainedRule)
		}
		return &directive
	case "secrule":
		rawDirective, err := yaml.Marshal(yamlDirective.Content[1])
		if err != nil {
			panic(err)
		}
		loaderDirective := directiveLoader{}
		err = yaml.Unmarshal(rawDirective, &loaderDirective)
		if err != nil {
			panic(err)
		}
		directive := types.SecRule{
			SecRuleMetadata: loaderDirective.SecRuleMetadata,
			Variables:       loaderDirective.Variables,
			Transformations: loaderDirective.Transformations,
			Operator:        loaderDirective.Operator,
			SeclangActions:  loaderDirective.SeclangActions,
		}
		var chainedRule types.SeclangDirective
		if len(loaderDirective.ChainedRule.Content) > 0 {
			chainedRule = directivesWithLabelsAux(loaderDirective.ChainedRule)
			directive.ChainedRule = castChainableDirective(chainedRule)
		}
		return &directive
	case "secrulescript":
		rawDirective, err := yaml.Marshal(yamlDirective.Content[1])
		if err != nil {
			panic(err)
		}
		loaderDirective := directiveLoader{}
		err = yaml.Unmarshal(rawDirective, &loaderDirective)
		if err != nil {
			panic(err)
		}
		directive := types.SecRuleScript{
			SecRuleMetadata: loaderDirective.SecRuleMetadata,
			Transformations: loaderDirective.Transformations,
			SeclangActions:  loaderDirective.SeclangActions,
			ScriptPath:      loaderDirective.ScriptPath,
		}
		var chainedRule types.SeclangDirective
		if len(loaderDirective.ChainedRule.Content) > 0 {
			chainedRule = directivesWithLabelsAux(loaderDirective.ChainedRule)
			directive.ChainedRule = castChainableDirective(chainedRule)
		}
		return &directive
	}
	return nil
}

// castChainableDirective casts a seclang directive to a chainable directive
func castChainableDirective(directive types.SeclangDirective) types.ChainableDirective {
	switch directive.(type) {
	case *types.SecRule:
		return directive.(*types.SecRule)
	case *types.SecAction:
		return directive.(*types.SecAction)
	case *types.SecRuleScript:
		return directive.(*types.SecRuleScript)
	}
	return nil
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

// PrintCRSLang writes crslang format directives (directives with conditions) to a file
func PrintCRSLang(configList types.ConfigurationList, filename string) error {
	configListWithConditions := exporters.ToDirectiveWithConditions(configList)

	yamlFile, err := yaml.Marshal(configListWithConditions.Configurations)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(yamlFile))
	if err != nil {
		return err
	}

	return nil
}

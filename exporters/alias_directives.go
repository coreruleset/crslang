package exporters

import (
	"os"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

type ConfigurationDirectiveWrapper struct {
	types.ConfigurationDirective
}

type SecDefaultActionWrapper struct {
	types.SecDefaultAction
}

type SecActionWrapper struct {
	*types.SecAction
}

type SecRuleWrapper struct {
	*types.SecRule
}

type SecRuleScriptWrapper struct {
	*types.SecRuleScript
}

type ConfigurationListWrapper struct {
	Configurations []ConfigurationWrapper `yaml:"configurations,omitempty"`
}

type ConfigurationWrapper struct {
	Marker     ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
	Directives []types.SeclangDirective      `yaml:"directives,omitempty"`
}

func ToDirectivesWithLabels(configList types.ConfigurationList) *ConfigurationListWrapper {
	result := new(ConfigurationListWrapper)
	for _, config := range configList.Configurations {
		configWrapper := new(ConfigurationWrapper)
		configWrapper.Marker = ConfigurationDirectiveWrapper{config.Marker}
		for _, directive := range config.Directives {
			directiveWrapper := WrapDirective(directive)
			configWrapper.Directives = append(configWrapper.Directives, directiveWrapper)
		}
		result.Configurations = append(result.Configurations, *configWrapper)
	}
	return result
}

func WrapDirective(directive types.SeclangDirective) types.SeclangDirective {
	var directiveWrapper types.SeclangDirective
	switch directive.(type) {
	case types.CommentMetadata:
		directiveWrapper = directive.(types.CommentMetadata)
	case types.SecDefaultAction:
		directiveWrapper = SecDefaultActionWrapper{directive.(types.SecDefaultAction)}
	case *types.SecAction:
		// recursibly wrap chained rule
		if directive.(*types.SecAction).ChainedRule != nil {
			wrapedChainedRule := WrapDirective(directive.(*types.SecAction).ChainedRule)
			directive.(*types.SecAction).ChainedRule = castWrappedDirective(wrapedChainedRule)
		}
		directiveWrapper = SecActionWrapper{directive.(*types.SecAction)}
	case *types.SecRule:
		// recursibly wrap chained rule
		if directive.(*types.SecRule).ChainedRule != nil {
			wrapedChainedRule := WrapDirective(directive.(*types.SecRule).ChainedRule)
			directive.(*types.SecRule).ChainedRule = castWrappedDirective(wrapedChainedRule)
		}
		directiveWrapper = SecRuleWrapper{directive.(*types.SecRule)}
	case *types.SecRuleScript:
		// recursibly wrap chained rule
		if directive.(*types.SecRuleScript).ChainedRule != nil {
			wrapedChainedRule := WrapDirective(directive.(*types.SecRuleScript).ChainedRule)
			directive.(*types.SecRuleScript).ChainedRule = castWrappedDirective(wrapedChainedRule)
		}
		directiveWrapper = SecRuleScriptWrapper{directive.(*types.SecRuleScript)}
	case types.ConfigurationDirective:
		directiveWrapper = ConfigurationDirectiveWrapper{directive.(types.ConfigurationDirective)}
	}
	return directiveWrapper
}

func castWrappedDirective(directive types.SeclangDirective) types.ChainableDirective {
	switch directive.(type) {
	case SecRuleWrapper:
		return directive.(SecRuleWrapper)
	case SecActionWrapper:
		return directive.(SecActionWrapper)
	case SecRuleScriptWrapper:
		return directive.(SecRuleScriptWrapper)
	}
	return nil
}

// yamlLoader is a auxiliary struct to load and iterate over the yaml file
type yamlLoader struct {
	Marker     ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
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
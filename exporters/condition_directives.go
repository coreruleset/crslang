package exporters

import (
	"os"

	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
	"gopkg.in/yaml.v3"
)

type Condition interface {
	ConditionToSeclang() string
}

type SecRuleCondition struct {
	types.Transformations `yaml:",inline,omitempty"`
	types.Variables       `yaml:",inline,omitempty"`
	types.Operator        `yaml:",omitempty"`
}

func (s SecRuleCondition) ConditionToSeclang() string {
	return "New sec rule condition"
}

type SecActionCondition struct {
	AlwaysMatch bool `yaml:"alwaysMatch,omitempty"`
}

func (s SecActionCondition) ConditionToSeclang() string {
	return "New sec action condition"
}

type ScriptCondition struct {
	Script string `yaml:"script,omitempty"`
}

func (s ScriptCondition) ConditionToSeclang() string {
	return "New script condition"
}

type RuleWithCondition struct {
	types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions            []Condition `yaml:"conditions,omitempty"`
	types.SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule           *RuleWithCondition `yaml:"chainedRule,omitempty"`
}

type RuleWithConditionWrapper struct {
	RuleWithCondition `yaml:"rule"`
}

func (s RuleWithCondition) ToSeclang() string {
	return "New sec rule with conditions"
}

func ToDirectiveWithConditions(configList types.ConfigurationList) *ConfigurationListWrapper {
	result := new(ConfigurationListWrapper)
	for _, config := range configList.Configurations {
		configWrapper := new(ConfigurationWrapper)
		configWrapper.Marker = ConfigurationDirectiveWrapper{config.Marker}
		for _, directive := range config.Directives {
			var directiveWrapper types.SeclangDirective
			switch directive.(type) {
			case types.CommentMetadata:
				directiveWrapper = directive.(types.CommentMetadata)
			case types.SecDefaultAction:
				directiveWrapper = SecDefaultActionWrapper{directive.(types.SecDefaultAction)}
			case *types.SecAction:
				directiveWrapper = RuleWithConditionWrapper{
					RuleToCondition(directive.(*types.SecAction)),
				}
			case *types.SecRule:
				directiveWrapper = RuleWithConditionWrapper{
					RuleToCondition(directive.(*types.SecRule)),
				}
			case *types.SecRuleScript:
				directiveWrapper = RuleWithConditionWrapper{
					RuleToCondition(directive.(*types.SecRuleScript)),
				}
			case types.ConfigurationDirective:
				directiveWrapper = ConfigurationDirectiveWrapper{directive.(types.ConfigurationDirective)}
			}
			configWrapper.Directives = append(configWrapper.Directives, directiveWrapper)
		}
		result.Configurations = append(result.Configurations, *configWrapper)
	}
	return result
}

func RuleToCondition(directive types.ChainableDirective) RuleWithCondition {
	var ruleWithCondition RuleWithCondition
	switch directive.(type) {
	case *types.SecRule:
		rule := directive.(*types.SecRule)
		ruleWithCondition = RuleWithCondition{
			rule.SecRuleMetadata,
			[]Condition{
				SecRuleCondition{
					rule.Transformations,
					rule.Variables,
					rule.Operator,
				},
			},
			rule.SeclangActions,
			nil,
		}
	case *types.SecAction:
		action := directive.(*types.SecAction)
		ruleWithCondition = RuleWithCondition{
			action.SecRuleMetadata,
			[]Condition{
				SecActionCondition{
					AlwaysMatch: true,
				},
			},
			action.SeclangActions,
			nil,
		}
	case *types.SecRuleScript:
		script := directive.(*types.SecRuleScript)
		ruleWithCondition = RuleWithCondition{
			script.SecRuleMetadata,
			[]Condition{
				ScriptCondition{
					Script: script.ScriptPath,
				},
			},
			script.SeclangActions,
			nil,
		}
	}
	if directive.GetChainedDirective() != nil {
		chainedConditionRule := RuleToCondition(directive.GetChainedDirective())
		if directive.NonDisruptiveActionsCount() > 0 {
			ruleWithCondition.ChainedRule = &chainedConditionRule
		} else {
			ruleWithCondition.Conditions = append(ruleWithCondition.Conditions, chainedConditionRule.Conditions...)
			ruleWithCondition.NonDisruptiveActions = chainedConditionRule.NonDisruptiveActions
			if chainedConditionRule.ChainedRule != nil {
				ruleWithCondition.ChainedRule = chainedConditionRule.ChainedRule
			}
		}
	}
	return ruleWithCondition
}

// yamlLoaderConditionRules is a auxiliary struct to load and iterate over the yaml file
type yamlLoaderConditionRules struct {
	Marker     ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
	Directives []yaml.Node                   `yaml:"directives,omitempty"`
}

// conditionDirectiveLoader is a auxiliary struct to load condition directives
type conditionDirectiveLoader struct {
	types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions            yaml.Node `yaml:"conditions,omitempty"`
	types.SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule           yaml.Node `yaml:"chainedRule"`
}

// LoadDirectivesWithConditions loads condition format directives from a yaml file
func LoadDirectivesWithConditions(filename string) types.ConfigurationList {
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
			directive := loadConditionDirective(yamlDirective)
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

// loadConditionDirective loads the different kind of directives
func loadConditionDirective(yamlDirective yaml.Node) types.SeclangDirective {
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
		return ConfigurationDirectiveWrapper{directive} 
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
		return SecDefaultActionWrapper{directive}
	case "rule":
		return RuleWithConditionWrapper{loadRuleWithConditions(yamlDirective, false)}
	}
	return nil
}

// loadRuleWithConditions loads a rule with conditions in a recursive way
func loadRuleWithConditions(yamlDirective yaml.Node, isChained bool) RuleWithCondition {
	rawDirective := []byte{}
	var err error

	if !isChained {
		rawDirective, err = yaml.Marshal(yamlDirective.Content[1])
	} else {
		rawDirective, err = yaml.Marshal(yamlDirective)
	}
	if err != nil {
		panic(err)
	}

	loaderDirective := conditionDirectiveLoader{}
	err = yaml.Unmarshal(rawDirective, &loaderDirective)
	if err != nil {
		panic(err)
	}
	directive := RuleWithCondition{
		SecRuleMetadata: loaderDirective.SecRuleMetadata,
		SeclangActions:  loaderDirective.SeclangActions,
	}
	if loaderDirective.Conditions.Kind == yaml.SequenceNode {
		for _, condition := range loaderDirective.Conditions.Content {
			loadedCondition := castConditions(condition)
			directive.Conditions = append(directive.Conditions, loadedCondition)
		}
	}
	var loadedChainedRule RuleWithCondition
	if len(loaderDirective.ChainedRule.Content) > 0 {
		loadedChainedRule = loadRuleWithConditions(loaderDirective.ChainedRule, true)
		directive.ChainedRule = &loadedChainedRule
	}
	return directive
}

// castConditions casts a directive condition to the correct type
func castConditions(condition *yaml.Node) Condition {
	switch condition.Content[0].Value {
	case "alwaysMatch":
		return SecActionCondition{AlwaysMatch: true}
	case "script":
		return ScriptCondition{Script: condition.Content[1].Value}
	case "variables", "transformations", "operator":
		rawDirective, err := yaml.Marshal(condition)
		if err != nil {
			panic(err)
		}
		ruleCondition := SecRuleCondition{}
		err = yaml.Unmarshal(rawDirective, &ruleCondition)
		return ruleCondition
	}
	return nil
}

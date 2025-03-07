package types

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Condition interface {
	ConditionToSeclang() string
}

type SecRuleCondition struct {
	Variables       []Variable   `yaml:"variables,omitempty"`
	Collections     []Collection `yaml:"collections,omitempty"`
	Operator        Operator     `yaml:"operator"`
	Transformations `yaml:",inline,omitempty"`
}

func (s SecRuleCondition) ConditionToSeclang() string {
	return "New sec rule condition"
}

type SecActionCondition struct {
	AlwaysMatch     bool `yaml:"alwaysMatch,omitempty"`
	Transformations `yaml:",inline,omitempty"`
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
	Kind        Kind               `yaml:"kind"`
	Metadata    SecRuleMetadata    `yaml:"metadata,omitempty"`
	Conditions  []Condition        `yaml:"conditions,omitempty"`
	Actions     SeclangActions     `yaml:"actions,omitempty"`
	ChainedRule *RuleWithCondition `yaml:"chainedRule,omitempty"`
}

func (s RuleWithCondition) ToSeclang() string {
	return "New sec rule with conditions"
}

func ToDirectiveWithConditions(configList ConfigurationList) *ConfigurationList {
	result := new(ConfigurationList)
	for _, config := range configList.DirectiveList {
		configWrapper := new(DirectiveList)
		configWrapper.Marker = config.Marker
		for _, directive := range config.Directives {
			var directiveWrapper SeclangDirective
			switch directive.(type) {
			case CommentMetadata:
				directiveWrapper = CommentDirective{
					Kind:     CommentKind,
					Metadata: directive.(CommentMetadata),
				}
			case *SecAction:
				directiveWrapper = RuleToCondition(directive.(*SecAction))
			case *SecRule:
				directiveWrapper = RuleToCondition(directive.(*SecRule))
			case *SecRuleScript:
				directiveWrapper = RuleToCondition(directive.(*SecRuleScript))
			default:
				directiveWrapper = directive
			}
			configWrapper.Directives = append(configWrapper.Directives, directiveWrapper)
		}
		result.DirectiveList = append(result.DirectiveList, *configWrapper)
	}
	return result
}

func RuleToCondition(directive ChainableDirective) RuleWithCondition {
	var ruleWithCondition RuleWithCondition
	switch directive.(type) {
	case *SecRule:
		rule := directive.(*SecRule)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*rule.Metadata,
			[]Condition{
				SecRuleCondition{
					rule.Variables,
					rule.Collections,
					rule.Operator,
					rule.Transformations,
				},
			},
			*rule.Actions,
			nil,
		}
	case *SecAction:
		action := directive.(*SecAction)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*action.Metadata,
			[]Condition{
				SecActionCondition{
					AlwaysMatch:     true,
					Transformations: action.Transformations,
				},
			},
			*action.Actions,
			nil,
		}
	case *SecRuleScript:
		script := directive.(*SecRuleScript)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*script.Metadata,
			[]Condition{
				ScriptCondition{
					Script: script.ScriptPath,
				},
			},
			*script.Actions,
			nil,
		}
	}
	if directive.GetChainedDirective() != nil {
		chainedConditionRule := RuleToCondition(directive.GetChainedDirective())
		if directive.NonDisruptiveActionsCount() > 0 {
			ruleWithCondition.ChainedRule = &chainedConditionRule
		} else {
			ruleWithCondition.Conditions = append(ruleWithCondition.Conditions, chainedConditionRule.Conditions...)
			ruleWithCondition.Actions.NonDisruptiveActions = chainedConditionRule.Actions.NonDisruptiveActions
			if chainedConditionRule.ChainedRule != nil {
				ruleWithCondition.ChainedRule = chainedConditionRule.ChainedRule
			}
		}
	}
	return ruleWithCondition
}

// yamlLoaderConditionRules is a auxiliary struct to load and iterate over the yaml file
type yamlLoaderConditionRules struct {
	Marker     ConfigurationDirective `yaml:"marker,omitempty"`
	Directives []yaml.Node            `yaml:"directives,omitempty"`
}

// conditionDirectiveLoader is a auxiliary struct to load condition directives
type conditionDirectiveLoader struct {
	Kind        string          `yaml:"kind"`
	Metadata    SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions  yaml.Node       `yaml:"conditions,omitempty"`
	Actions     SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule yaml.Node       `yaml:"chainedRule"`
}

// LoadDirectivesWithConditionsFromFile loads condition format directives from a yaml file
func LoadDirectivesWithConditionsFromFile(filename string) ConfigurationList {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return LoadDirectivesWithConditions(yamlFile)
}

// LoadDirectivesWithConditions loads condition format directives from a yaml file
func LoadDirectivesWithConditions(yamlFile []byte) ConfigurationList {
	var configs []yamlLoaderConditionRules
	err := yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		panic(err)
	}
	var resultConfigs []DirectiveList
	for _, config := range configs {
		var directives []SeclangDirective
		for _, yamlDirective := range config.Directives {
			directive := loadConditionDirective(yamlDirective)
			if directive == nil {
				panic("Unknown directive type")
			} else {
				directives = append(directives, directive)
			}
		}
		resultConfigs = append(resultConfigs, DirectiveList{Marker: config.Marker, Directives: directives})
	}
	return ConfigurationList{DirectiveList: resultConfigs}
}

// loadConditionDirective loads the different kind of directives
func loadConditionDirective(yamlDirective yaml.Node) SeclangDirective {
	if yamlDirective.Kind != yaml.MappingNode {
		panic("Unknown format type")
	}
	switch yamlDirective.Content[1].Value {
	case "comment":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := CommentDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "configuration":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := ConfigurationDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "defaultaction":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := DefaultAction{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "rule":
		return loadRuleWithConditions(yamlDirective)
	case "remove":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := RemoveRuleDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	}
	return nil
}

// loadRuleWithConditions loads a rule with conditions in a recursive way
func loadRuleWithConditions(yamlDirective yaml.Node) RuleWithCondition {
	rawDirective := []byte{}
	var err error

	rawDirective, err = yaml.Marshal(yamlDirective)

	if err != nil {
		panic(err)
	}

	loaderDirective := conditionDirectiveLoader{}
	err = yaml.Unmarshal(rawDirective, &loaderDirective)
	if err != nil {
		print(string(rawDirective))
		panic(err)
	}
	directive := RuleWithCondition{
		Kind:     RuleKind,
		Metadata: loaderDirective.Metadata,
		Actions:  loaderDirective.Actions,
	}
	if loaderDirective.Conditions.Kind == yaml.SequenceNode {
		for _, condition := range loaderDirective.Conditions.Content {
			loadedCondition := castConditions(condition)
			directive.Conditions = append(directive.Conditions, loadedCondition)
		}
	}
	var loadedChainedRule RuleWithCondition
	if len(loaderDirective.ChainedRule.Content) > 0 {
		loadedChainedRule = loadRuleWithConditions(loaderDirective.ChainedRule)
		directive.ChainedRule = &loadedChainedRule
	}
	return directive
}

// castConditions casts a directive condition to the correct type
func castConditions(condition *yaml.Node) Condition {
	switch condition.Content[0].Value {
	case "alwaysMatch":
		rawDirective, err := yaml.Marshal(condition)
		if err != nil {
			panic(err)
		}
		ruleCondition := SecActionCondition{}
		err = yaml.Unmarshal(rawDirective, &ruleCondition)
		if err != nil {
			panic(err)
		}
		return ruleCondition
	case "script":
		return ScriptCondition{Script: condition.Content[1].Value}
	case "variables", "collections":
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

func FromCRSLangToUnformattedDirectives(configListWrapped ConfigurationList) *ConfigurationList {
	result := new(ConfigurationList)
	for _, config := range configListWrapped.DirectiveList {
		configList := new(DirectiveList)
		configList.Marker = config.Marker
		for _, directiveWrapped := range config.Directives {
			var directive SeclangDirective
			switch directiveWrapped.(type) {
			case CommentDirective:
				directive = directiveWrapped.(CommentDirective).Metadata
			case DefaultAction:
				directive = directiveWrapped
			case RuleWithCondition:
				directive = FromConditionToUnmorfattedDirective(directiveWrapped.(RuleWithCondition))
			case ConfigurationDirective:
				directive = ConfigurationDirective{
					Metadata:  directiveWrapped.(ConfigurationDirective).Metadata,
					Name:      directiveWrapped.(ConfigurationDirective).Name,
					Parameter: directiveWrapped.(ConfigurationDirective).Parameter,
				}
			}
			configList.Directives = append(configList.Directives, directive)
		}
		result.DirectiveList = append(result.DirectiveList, *configList)
	}
	return result
}

func FromConditionToUnmorfattedDirective(conditionDirective RuleWithCondition) ChainableDirective {
	var rootDirective ChainableDirective
	var directiveIterator ChainableDirective
	var chainedDirective ChainableDirective
	var directiveAux ChainableDirective

	chainedDirective = nil

	if conditionDirective.ChainedRule != nil {
		chainedDirective = FromConditionToUnmorfattedDirective(*conditionDirective.ChainedRule)
	}

	for i, condition := range conditionDirective.Conditions {
		switch condition.(type) {
		case SecRuleCondition:
			secruleDirective := new(SecRule)
			secruleDirective.Variables = condition.(SecRuleCondition).Variables
			secruleDirective.Collections = condition.(SecRuleCondition).Collections
			secruleDirective.Transformations = condition.(SecRuleCondition).Transformations
			secruleDirective.Operator = condition.(SecRuleCondition).Operator
			if i == 0 {
				secruleDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secruleDirective.Actions = CopyActions(conditionDirective.Actions)
				secruleDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secruleDirective
			} else {
				secruleDirective.Metadata = new(SecRuleMetadata)
				secruleDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secruleDirective.Actions.FlowActions = []Action{{Action: "chain"}}
				}
			}
			directiveAux = secruleDirective
		case SecActionCondition:
			secactionDirective := new(SecAction)
			secactionDirective.Transformations = condition.(SecActionCondition).Transformations
			if i == 0 {
				secactionDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secactionDirective.Actions = CopyActions(conditionDirective.Actions)
				secactionDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secactionDirective
			} else {
				secactionDirective.Metadata = new(SecRuleMetadata)
				secactionDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secactionDirective.Actions.FlowActions = []Action{{Action: "chain"}}
				}
			}
			directiveAux = secactionDirective
		case ScriptCondition:
			secscriptDirective := new(SecRuleScript)
			secscriptDirective.ScriptPath = condition.(ScriptCondition).Script
			if i == 0 {
				secscriptDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secscriptDirective.Actions = CopyActions(conditionDirective.Actions)
				secscriptDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secscriptDirective
			} else {
				secscriptDirective.Metadata = new(SecRuleMetadata)
				secscriptDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secscriptDirective.Actions.FlowActions = []Action{{Action: "chain"}}
				}
			}
			directiveAux = secscriptDirective
		}
		if i == 0 {
			directiveIterator = rootDirective
		} else {
			directiveIterator.AppendChainedDirective(directiveAux)
			directiveIterator = directiveAux
		}

	}

	switch directiveIterator.(type) {
	case *SecRule:
		directiveIterator.(*SecRule).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *SecAction:
		directiveIterator.(*SecAction).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *SecRuleScript:
		directiveIterator.(*SecRuleScript).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	}

	if chainedDirective != nil {
		directiveIterator.AppendChainedDirective(chainedDirective)
	}

	return rootDirective
}

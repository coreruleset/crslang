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
	types.Variables       `yaml:",inline,omitempty"`
	Operator              types.Operator `yaml:",omitempty"`
	types.Transformations `yaml:",inline,omitempty"`
}

func (s SecRuleCondition) ConditionToSeclang() string {
	return "New sec rule condition"
}

type SecActionCondition struct {
	AlwaysMatch           bool `yaml:"alwaysMatch,omitempty"`
	types.Transformations `yaml:",inline,omitempty"`
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
	Kind        string                `yaml:"kind"`
	Metadata    types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions  []Condition           `yaml:"conditions,omitempty"`
	Actions     types.SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule *RuleWithCondition    `yaml:"chainedRule,omitempty"`
}

func (s RuleWithCondition) ToSeclang() string {
	return "New sec rule with conditions"
}

func ToDirectiveWithConditions(configList types.ConfigurationList) *ConfigurationListWrapper {
	result := new(ConfigurationListWrapper)
	for _, config := range configList.DirectiveList {
		configWrapper := new(ConfigurationWrapper)
		configWrapper.Marker = ConfigurationDirectiveWrapper{config.Marker}
		for _, directive := range config.Directives {
			var directiveWrapper types.SeclangDirective
			switch directive.(type) {
			case types.CommentMetadata:
				directiveWrapper = types.CommentDirective{
					Kind:     "comment",
					Metadata: directive.(types.CommentMetadata),
				}
			case types.DefaultAction:
				defaultaction := directive.(types.DefaultAction)
				defaultaction.Kind = "defaultaction"
				directiveWrapper = defaultaction
			case *types.SecAction:
				directiveWrapper = RuleToCondition(directive.(*types.SecAction))
			case *types.SecRule:
				directiveWrapper = RuleToCondition(directive.(*types.SecRule))
			case *types.SecRuleScript:
				directiveWrapper = RuleToCondition(directive.(*types.SecRuleScript))
			case types.ConfigurationDirective:
				directiveWrapper = types.ConfigurationDirective{
					Kind:      "configuration",
					Metadata:  directive.(types.ConfigurationDirective).Metadata,
					Name:      directive.(types.ConfigurationDirective).Name,
					Parameter: directive.(types.ConfigurationDirective).Parameter,
				}
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
			"rule",
			*rule.Metadata,
			[]Condition{
				SecRuleCondition{
					rule.Variables,
					rule.Operator,
					rule.Transformations,
				},
			},
			*rule.Actions,
			nil,
		}
	case *types.SecAction:
		action := directive.(*types.SecAction)
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
	case *types.SecRuleScript:
		script := directive.(*types.SecRuleScript)
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
	Marker     ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
	Directives []yaml.Node                   `yaml:"directives,omitempty"`
}

// conditionDirectiveLoader is a auxiliary struct to load condition directives
type conditionDirectiveLoader struct {
	Kind        string                `yaml:"kind"`
	Metadata    types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions  yaml.Node             `yaml:"conditions,omitempty"`
	Actions     types.SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule yaml.Node             `yaml:"chainedRule"`
}

// LoadDirectivesWithConditionsFromFile loads condition format directives from a yaml file
func LoadDirectivesWithConditionsFromFile(filename string) ConfigurationListWrapper {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return LoadDirectivesWithConditions(yamlFile)
}

// LoadDirectivesWithConditions loads condition format directives from a yaml file
func LoadDirectivesWithConditions(yamlFile []byte) ConfigurationListWrapper {
	var configs []yamlLoaderConditionRules
	err := yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		panic(err)
	}
	var resultConfigs []ConfigurationWrapper
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
		resultConfigs = append(resultConfigs, ConfigurationWrapper{Marker: config.Marker, Directives: directives})
	}
	return ConfigurationListWrapper{Configurations: resultConfigs}
}

// loadConditionDirective loads the different kind of directives
func loadConditionDirective(yamlDirective yaml.Node) types.SeclangDirective {
	if yamlDirective.Kind != yaml.MappingNode {
		panic("Unknown format type")
	}
	switch yamlDirective.Content[1].Value {
	case "comment":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := types.CommentDirective{}
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
		directive := types.ConfigurationDirective{}
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
		directive := types.DefaultAction{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "rule":
		return loadRuleWithConditions(yamlDirective)
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
		Kind:     "rule",
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
	case "variables":
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

func FromCRSLangToUnformattedDirectives(configListWrapped ConfigurationListWrapper) *types.ConfigurationList {
	result := new(types.ConfigurationList)
	for _, config := range configListWrapped.Configurations {
		configList := new(types.DirectiveList)
		configList.Marker = config.Marker.ConfigurationDirective
		for _, directiveWrapped := range config.Directives {
			var directive types.SeclangDirective
			switch directiveWrapped.(type) {
			case types.CommentDirective:
				directive = directiveWrapped.(types.CommentDirective).Metadata
			case types.DefaultAction:
				directive = directiveWrapped
			case RuleWithCondition:
				directive = FromConditionToUnmorfattedDirective(directiveWrapped.(RuleWithCondition))
			case types.ConfigurationDirective:
				directive = types.ConfigurationDirective{
					Metadata:  directiveWrapped.(types.ConfigurationDirective).Metadata,
					Name:      directiveWrapped.(types.ConfigurationDirective).Name,
					Parameter: directiveWrapped.(types.ConfigurationDirective).Parameter,
				}
			}
			configList.Directives = append(configList.Directives, directive)
		}
		result.DirectiveList = append(result.DirectiveList, *configList)
	}
	return result
}

func FromConditionToUnmorfattedDirective(conditionDirective RuleWithCondition) types.ChainableDirective {
	var rootDirective types.ChainableDirective
	var directiveIterator types.ChainableDirective
	var chainedDirective types.ChainableDirective
	var directiveAux types.ChainableDirective

	chainedDirective = nil

	if conditionDirective.ChainedRule != nil {
		chainedDirective = FromConditionToUnmorfattedDirective(*conditionDirective.ChainedRule)
	}

	for i, condition := range conditionDirective.Conditions {
		switch condition.(type) {
		case SecRuleCondition:
			secruleDirective := new(types.SecRule)
			secruleDirective.Variables = condition.(SecRuleCondition).Variables
			secruleDirective.Transformations = condition.(SecRuleCondition).Transformations
			secruleDirective.Operator = condition.(SecRuleCondition).Operator
			if i == 0 {
				secruleDirective.Metadata = types.CopySecRuleMetadata(conditionDirective.Metadata)
				secruleDirective.Actions = types.CopyActions(conditionDirective.Actions)
				secruleDirective.Actions.NonDisruptiveActions = []types.Action{}
				rootDirective = secruleDirective
			} else {
				secruleDirective.Metadata = new(types.SecRuleMetadata)
				secruleDirective.Actions = new(types.SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secruleDirective.Actions.FlowActions = []types.Action{{Action: "chain"}}
				}
			}
			directiveAux = secruleDirective
		case SecActionCondition:
			secactionDirective := new(types.SecAction)
			secactionDirective.Transformations = condition.(SecActionCondition).Transformations
			if i == 0 {
				secactionDirective.Metadata = types.CopySecRuleMetadata(conditionDirective.Metadata)
				secactionDirective.Actions = types.CopyActions(conditionDirective.Actions)
				secactionDirective.Actions.NonDisruptiveActions = []types.Action{}
				rootDirective = secactionDirective
			} else {
				secactionDirective.Metadata = new(types.SecRuleMetadata)
				secactionDirective.Actions = new(types.SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secactionDirective.Actions.FlowActions = []types.Action{{Action: "chain"}}
				}
			}
			directiveAux = secactionDirective
		case ScriptCondition:
			secscriptDirective := new(types.SecRuleScript)
			secscriptDirective.ScriptPath = condition.(ScriptCondition).Script
			if i == 0 {
				secscriptDirective.Metadata = types.CopySecRuleMetadata(conditionDirective.Metadata)
				secscriptDirective.Actions = types.CopyActions(conditionDirective.Actions)
				secscriptDirective.Actions.NonDisruptiveActions = []types.Action{}
				rootDirective = secscriptDirective
			} else {
				secscriptDirective.Metadata = new(types.SecRuleMetadata)
				secscriptDirective.Actions = new(types.SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					secscriptDirective.Actions.FlowActions = []types.Action{{Action: "chain"}}
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
	case *types.SecRule:
		directiveIterator.(*types.SecRule).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *types.SecAction:
		directiveIterator.(*types.SecAction).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *types.SecRuleScript:
		directiveIterator.(*types.SecRuleScript).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	}

	if chainedDirective != nil {
		directiveIterator.AppendChainedDirective(chainedDirective)
	}

	return rootDirective
}

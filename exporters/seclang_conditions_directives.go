package exporters

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
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

type RuleWithCondititon struct {
	types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions            []Condition `yaml:"conditions,omitempty"`
	types.SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule           *RuleWithCondititon `yaml:"chainedRule,omitempty"`
}

type RuleWithCondititonWrapper struct {
	RuleWithCondititon `yaml:"rule"`
}

func (s RuleWithCondititon) ToSeclang() string {
	return "New sec rule with conditions"
}

func ConcreteRepr2(configList types.ConfigurationList) *ConfigurationListWrapper {
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
			// case types.SecAction:
			// 	directiveWrapper = SecActionWrapper{directive.(types.SecAction)}
			// case types.SecRule:
			// 	directiveWrapper = SecRuleWithCondititonWrapper{
			// 		ToConditions((directive.(types.SecRule))),
			// 	}
			case *types.SecAction:
				directiveWrapper = RuleWithCondititonWrapper{
					ToConditions(directive.(*types.SecAction)),
				}
			case *types.SecRule:
				directiveWrapper = RuleWithCondititonWrapper{
					ToConditions(directive.(*types.SecRule)),
				}
			case *types.SecRuleScript:
				directiveWrapper = RuleWithCondititonWrapper{
					ToConditions(directive.(*types.SecRuleScript)),
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

func ToConditions(directive types.ChainableDirective) RuleWithCondititon {
	var ruleWithCondition RuleWithCondititon
	switch directive.(type) {
	case *types.SecRule:
		rule := directive.(*types.SecRule)
		ruleWithCondition = RuleWithCondititon{
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
		action:= directive.(*types.SecAction)
		ruleWithCondition = RuleWithCondititon{
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
		script:= directive.(*types.SecRuleScript)
		ruleWithCondition = RuleWithCondititon{
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
		chainedConditionRule := ToConditions(directive.GetChainedDirective())
		// fmt.Printf("Condition rule: %v\n", conditionRule)
		if directive.NonDisruptiveActionsCount() > 0 {
			ruleWithCondition.ChainedRule = &chainedConditionRule
			// fmt.Printf("Condition rule: %v\n", chainedConditionRule)
			// fmt.Printf("Chained rule: %v\n", *conditionRule.chainedRule)
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

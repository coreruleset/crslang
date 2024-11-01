package exporters

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

type ConfigurationDirectiveWrapper struct {
	types.ConfigurationDirective
}

type SecDefaultActionWrapper struct {
	types.SecDefaultAction
}

type SecActionWrapper struct {
	types.SecAction
}

type SecRuleWrapper struct {
	types.SecRule
}

type ConfigurationListWrapper struct {
	Configurations []ConfigurationWrapper `yaml:"configurations,omitempty"`
}

type ConfigurationWrapper struct {
	Marker     ConfigurationDirectiveWrapper `yaml:"marker,omitempty"`
	Directives []types.SeclangDirective     `yaml:"directives,omitempty"`
}

func ToDirectivesWithLabels(configList types.ConfigurationList) *ConfigurationListWrapper {
	result := new(ConfigurationListWrapper)
	for _, config := range(configList.Configurations){
		configWrapper := new(ConfigurationWrapper)
		configWrapper.Marker = ConfigurationDirectiveWrapper{config.Marker}
		for _, directive := range(config.Directives){
			var directiveWrapper types.SeclangDirective
			switch directive.(type){
			case types.CommentMetadata:
				directiveWrapper = directive.(types.CommentMetadata)
			case types.SecDefaultAction:
				directiveWrapper = SecDefaultActionWrapper{directive.(types.SecDefaultAction)}
			case types.SecAction:
				directiveWrapper = SecActionWrapper{directive.(types.SecAction)}
			case types.SecRule:
				directiveWrapper = SecRuleWrapper{directive.(types.SecRule)}
			case types.ConfigurationDirective:
				directiveWrapper = ConfigurationDirectiveWrapper{directive.(types.ConfigurationDirective)}
			}
			configWrapper.Directives = append(configWrapper.Directives, directiveWrapper)
		}
		result.Configurations = append(result.Configurations, *configWrapper)
	}
	return result
}

type Condition struct {
	types.Transformations `yaml:",inline,omitempty"`
	types.Variables `yaml:",inline"`
	types.Operator
}

type SecRuleWithCondititon struct {
	types.SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions []Condition `yaml:"conditions,omitempty"`
	types.SeclangActions `yaml:"actions,omitempty"`
	ChainedRule *SecRuleWithCondititon `yaml:"chainedRule,omitempty"`
}

type SecRuleWithCondititonWrapper struct {
	SecRuleWithCondititon `yaml:"rule"`
}

func (s SecRuleWithCondititon) ToSeclang() string {
	return "New sec rule with conditions"
}

func ConcreteRepr2(configList types.ConfigurationList) *ConfigurationListWrapper {
	result := new(ConfigurationListWrapper)
	for _, config := range(configList.Configurations){
		configWrapper := new(ConfigurationWrapper)
		configWrapper.Marker = ConfigurationDirectiveWrapper{config.Marker}
		for _, directive := range(config.Directives){
			var directiveWrapper types.SeclangDirective
			switch directive.(type){
			case types.CommentMetadata:
				directiveWrapper = directive.(types.CommentMetadata)
			case types.SecDefaultAction:
				directiveWrapper = SecDefaultActionWrapper{directive.(types.SecDefaultAction)}
			case types.SecAction:
				directiveWrapper = SecActionWrapper{directive.(types.SecAction)}
			case types.SecRule:
				directiveWrapper = SecRuleWithCondititonWrapper{
					ToConditions(directive.(types.SecRule)),
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

func ToConditions(rule types.SecRule) SecRuleWithCondititon {
	if rule.ChainedRule == nil {
		return SecRuleWithCondititon{
			rule.SecRuleMetadata,
			[]Condition{
				{
					rule.Transformations,
					rule.Variables,
					rule.Operator,
				},
			},
			rule.SeclangActions,
			nil,
		}
	} else {
		chainedConditionRule := ToConditions(*rule.ChainedRule)
		// fmt.Printf("Chained rule: %v\n", chainedConditionRule)
		conditionRule := SecRuleWithCondititon{
			rule.SecRuleMetadata,
			[]Condition{
				{
					rule.Transformations,
					rule.Variables,
					rule.Operator,
				},
			},
			rule.SeclangActions,
			nil,
		}
		// fmt.Printf("Condition rule: %v\n", conditionRule)
		if len(rule.NonDisruptiveActions) > 0 {
			conditionRule.ChainedRule = &chainedConditionRule
			// fmt.Printf("Condition rule: %v\n", chainedConditionRule)
			// fmt.Printf("Chained rule: %v\n", *conditionRule.chainedRule)
		} else {
			conditionRule.Conditions = append(conditionRule.Conditions, chainedConditionRule.Conditions...)
			conditionRule.NonDisruptiveActions = chainedConditionRule.NonDisruptiveActions
			if chainedConditionRule.ChainedRule != nil {
				conditionRule.ChainedRule = chainedConditionRule.ChainedRule
			}
		}
		return conditionRule
	}
}
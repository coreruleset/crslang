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
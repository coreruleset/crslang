package converters

import "github.com/coreruleset/crslang/types"

// ToCRSLang process previously loaded seclang directives to CRSLang schema directives
func ToCRSLang(configList types.ConfigurationList) *types.ConfigurationList {
	configListWithConditions := types.ToDirectiveWithConditions(configList)

	configListWithConditions.ExtractDefaultValues()
	return configListWithConditions
}

package exporters

import (
	"gitlab.fing.edu.uy/gsi/seclang/crslang/types"
)

func ToSeclang(configList types.ConfigurationList) string {
	result := ""
	for _, config := range configList.Configurations {
		for _, directive := range config.Directives {
			result += directive.ToSeclang() + "\n"
		}
	}
	return result
}
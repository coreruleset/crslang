package types

type ConfigurationList struct {
	DirectiveList []DirectiveList `yaml:"directivelist,omitempty"`
}

type DirectiveList struct {
	Id         string                 `yaml:"id"`
	Marker     ConfigurationDirective `yaml:"marker,omitempty"`
	Directives []SeclangDirective     `yaml:"directives,omitempty"`
}

func ToSeclang(configList ConfigurationList) string {
	result := ""
	for _, config := range configList.DirectiveList {
		for _, directive := range config.Directives {
			result += directive.ToSeclang() + "\n"
		}
	}
	return result
}

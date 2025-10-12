package types

import (
	"slices"
)

type DefaultConfigs struct {
	Version string   `yaml:"version,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

type ConfigurationList struct {
	Global        DefaultConfigs  `yaml:"global,omitempty"`
	DirectiveList []DirectiveList `yaml:"directivelist,omitempty"`
}

type DirectiveList struct {
	Id         string                 `yaml:"id"`
	Directives []SeclangDirective     `yaml:"directives,omitempty"`
	Marker     ConfigurationDirective `yaml:"marker,omitempty"`
}

func (d DirectiveList) ToSeclang() string {
	result := ""
	for _, directive := range d.Directives {
		result += directive.ToSeclang() + "\n"
	}
	if d.Marker.Name != "" {
		result += d.Marker.ToSeclang() + "\n"
	}
	return result
}

func ToSeclang(configList ConfigurationList) string {
	result := ""
	for _, config := range configList.DirectiveList {
		for _, directive := range config.Directives {
			result += directive.ToSeclang() + "\n"
		}
		if config.Marker.Name != "" {
			result += config.Marker.ToSeclang() + "\n"
		}
	}
	return result
}

func (c *ConfigurationList) ExtractDefaultValues() {
	directiveFound := false
	version := ""
	tags := []string{}

	for _, directiveList := range c.DirectiveList {
		for _, directive := range directiveList.Directives {
			if directive.GetKind() == RuleKind {
				if !directiveFound {
					directiveFound = true
					version = directive.(RuleWithCondition).Metadata.Ver
					tags = directive.(RuleWithCondition).Metadata.Tags
				} else {
					if version != directive.(RuleWithCondition).Metadata.Ver {
						version = ""
					}
					auxTags := []string{}
					for _, tag := range tags {
						if slices.Contains(directive.(RuleWithCondition).Metadata.Tags, tag) {
							auxTags = append(auxTags, tag)
						}
					}
					tags = auxTags
				}
			}
		}
	}
	c.Global.Version = version
	c.Global.Tags = tags
}

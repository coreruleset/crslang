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

// ExtractDefaultValues extracts default values for version and tags from the rules in the configuration list
func (c *ConfigurationList) ExtractDefaultValues() {
	directiveFound := false
	version := ""
	tags := []string{}
	rules := []*RuleWithCondition{}

	for i := range c.DirectiveList {
		for j := range c.DirectiveList[i].Directives {
			// Only consider Rule directives
			if c.DirectiveList[i].Directives[j].GetKind() == RuleKind {
				rule := c.DirectiveList[i].Directives[j].(*RuleWithCondition)
				rules = append(rules, rule)
				auxTags := []string{}
				if !directiveFound {
					directiveFound = true
					version = rule.Metadata.Ver
					auxTags = append(auxTags, rule.Metadata.Tags...)
				} else {
					if version != rule.Metadata.Ver {
						version = ""
					}
					for _, tag := range tags {
						if slices.Contains(rule.Metadata.Tags, tag) {
							auxTags = append(auxTags, tag)
						}
					}
				}
				tags = auxTags
				// If both version and tags are empty after found a rule it means there is no common value
				// so we can stop the search
				if version == "" && len(tags) == 0 {
					return
				}
			}
		}
	}

	// Clear version and tags in rules since they are now in the global section
	for _, rule := range rules {
		if version != "" {
			rule.Metadata.Ver = ""
		}
		rule.Metadata.Tags = slices.DeleteFunc(rule.Metadata.Tags, func(s string) bool {
			return slices.Contains(tags, s)
		})
	}

	c.Global.Version = version
	c.Global.Tags = tags
}

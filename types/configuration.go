package types

import (
	"slices"
)

type DefaultConfigs struct {
	Version string   `yaml:"version,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

type Ruleset struct {
	Global    DefaultConfigs `yaml:"global,omitempty"`
	GroupsIds []string       `yaml:"rule_groups,omitempty"`
	Groups    []Group        `yaml:"groups,omitempty"`
}

type Group struct {
	Id             string                   `yaml:"id"`
	Tags           []string                 `yaml:"tags,omitempty"`
	Comments       []string                 `yaml:"comments,omitempty"`
	Configurations []ConfigurationDirective `yaml:"configurations,omitempty"`
	Directives     []SeclangDirective       `yaml:"directives,omitempty"`
	Rules          []int                    `yaml:"rules,omitempty"`
	Marker         ConfigurationDirective   `yaml:"marker,omitempty"`
}

func (d Group) ToSeclang() string {
	result := ""
	for _, directive := range d.Directives {
		result += directive.ToSeclang() + "\n"
	}
	if d.Marker.Name != "" {
		result += d.Marker.ToSeclang() + "\n"
	}
	return result
}

func ToSeclang(configList Ruleset) string {
	result := ""
	for _, config := range configList.Groups {
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
func (c *Ruleset) ExtractDefaultValues() {
	directiveFound := false
	version := ""
	tags := []string{}
	rules := []*RuleWithCondition{}

	for i := range c.Groups {
		for j := range c.Groups[i].Directives {
			// Only consider Rule directives
			if c.Groups[i].Directives[j].GetKind() == RuleKind {
				rule := c.Groups[i].Directives[j].(*RuleWithCondition)
				rules = append(rules, rule)
				if !directiveFound {
					directiveFound = true
					version = rule.Metadata.Ver
					tags = rule.Metadata.Tags
				} else {
					if version != rule.Metadata.Ver {
						version = ""
					}
					auxTags := []string{}
					for _, tag := range tags {
						if slices.Contains(rule.Metadata.Tags, tag) {
							auxTags = append(auxTags, tag)
						}
					}
					tags = auxTags
				}
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
		rule.Metadata.Ver = ""
		rule.Metadata.Tags = slices.DeleteFunc(rule.Metadata.Tags, func(s string) bool {
			return slices.Contains(tags, s)
		})
	}

	c.Global.Version = version
	c.Global.Tags = tags
}

// TODO: merge this method with the one in Ruleset
// ExtractDefaultValues extracts default values for version and tags from the rules in the configuration list
func (g *Group) ExtractDefaultValues() {
	directiveFound := false
	tags := []string{}
	rules := []*RuleWithCondition{}

	for j := range g.Directives {
		// Only consider Rule directives
		if g.Directives[j].GetKind() == RuleKind {
			// Ignore paranoia level check rules
			lastDigits := g.Directives[j].(*RuleWithCondition).Metadata.Id % 1000
			if lastDigits < 20 {
				rule := g.Directives[j].(*RuleWithCondition)
				rules = append(rules, rule)
				auxTags := []string{}
				if !directiveFound {
					directiveFound = true
					auxTags = append(auxTags, rule.Metadata.Tags...)
				} else {
					for _, tag := range tags {
						if slices.Contains(rule.Metadata.Tags, tag) {
							auxTags = append(auxTags, tag)
						}
					}
				}
				tags = auxTags
				// If tags are empty after found a rule it means there is no common value
				// so we can stop the search
				if len(tags) == 0 {
					return
				}
			}
		}
	}

	// Clear tags in rules since they are now in the global section
	for _, rule := range rules {
		rule.Metadata.Tags = slices.DeleteFunc(rule.Metadata.Tags, func(s string) bool {
			return slices.Contains(tags, s)
		})
	}

	g.Tags = append([]string{}, tags...)

}

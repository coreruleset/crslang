package types

import (
	"fmt"
	"slices"
)

type PhaseDefaults struct {
	Comment         string          `yaml:"comment,omitempty"`
	Transformations Transformations `yaml:",inline"`
	Actions         *SeclangActions `yaml:"actions"`
}

type DefaultConfigs struct {
	Version string                   `yaml:"version,omitempty"`
	Tags    []string                 `yaml:"tags,omitempty"`
	Phases  map[string]PhaseDefaults `yaml:"phases,omitempty"`
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

// ExtractPhaseDefaults extract default actions from the directive list and add it to the global config
func (c *ConfigurationList) ExtractPhaseDefaults() error {
	defaultActions := map[string]PhaseDefaults{}

	for _, dirList := range c.DirectiveList {
		for _, directive := range dirList.Directives {
			if directive.GetKind() == DefaultActionKind {
				da, ok := directive.(DefaultAction)
				if !ok {
					return fmt.Errorf("Error: casting directive to DefaultAction")
				}
				pd := PhaseDefaults{
					Comment:         da.Metadata.Comment,
					Transformations: da.Transformations,
					Actions:         CopyActions(*da.Actions),
				}
				_, ok = defaultActions[da.Metadata.Phase]
				if ok {
					return fmt.Errorf("Error: duplicate default actions for phase %s", da.Metadata.Phase)
				}
				defaultActions[da.Metadata.Phase] = pd
			}
		}
	}

	for i := range c.DirectiveList {
		c.DirectiveList[i].Directives = slices.DeleteFunc(c.DirectiveList[i].Directives, func(d SeclangDirective) bool {
			return d.GetKind() == DefaultActionKind
		})
	}

	c.Global.Phases = defaultActions
	return nil
}

// ExtractPhaseDefaults extract default actions from the directive list and add it to the global config
func (c *ConfigurationList) PhaseDefaultsToSeclang() error {
	if len(c.Global.Phases) > 0 {
		seclangDirs := []SeclangDirective{}
		for p, v := range c.Global.Phases {
			da := NewDefaultAction()
			da.Metadata.SetComment(v.Comment)
			da.Metadata.SetPhase(p)
			da.Actions = CopyActions(*v.Actions)
			seclangDirs = append(seclangDirs, *da)
		}

		if len(c.DirectiveList) != 0 {
			c.DirectiveList[0].Directives = append(seclangDirs, c.DirectiveList[0].Directives...)
		} else {
			dirList := DirectiveList{}
			dirList.Id = "crs-setup"
			dirList.Directives = seclangDirs
			c.DirectiveList = append(c.DirectiveList, dirList)
		}
	}

	return nil
}

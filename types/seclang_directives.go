package types

type SeclangDirective interface {
	ToSeclang() string
}

type ConfigurationDirective struct {
	CommentMetadata `yaml:",inline"`
	DirectiveName   string `yaml:"directiveName"`
	Parameter       string `yaml:"parameter"`
}

func (c ConfigurationDirective) ToSeclang() string {
	return c.Comment + c.DirectiveName + " " + c.Parameter + "\n"
}

type SecDefaultAction struct {
	OnlyPhaseMetadata `yaml:"metadata"`
	Transformations   `yaml:",inline"`
	SeclangActions    `yaml:"actions"`
}

func (s SecDefaultAction) ToSeclang() string {
	result := ""
	result += s.Comment + "SecDefaultAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToString()
	transformations := s.Transformations.ToString()
	if actions != "" {
		result += ", " + actions
	}
	if transformations != "" {
		result += ", " + transformations
	}
	result += "\"\n"
	return result
}

type SecAction struct {
	SecRuleMetadata `yaml:"metadata"`
	Transformations `yaml:",inline"`
	SeclangActions  `yaml:"actions"`
}

func (s SecAction) ToSeclang() string {
	result := ""
	result += s.Comment + "SecAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToString()
	transformations := s.Transformations.ToString()
	if actions != "" {
		result += ", " + actions
	}
	if transformations != "" {
		result += ", " + transformations
	}
	result += "\"\n"
	return result
}

type SecRule struct {
	SecRuleMetadata `yaml:"metadata"`
	Variables       `yaml:",inline"`
	Transformations `yaml:",inline"`
	Operator        `yaml:"operator"`
	SeclangActions  `yaml:"actions"`
}

func (s SecRule) ToSeclang() string {
	result := ""
	result += s.Comment + "SecRule "
	result += s.Variables.ToString() + " "
	result += "\"" + s.Operator.ToString() + "\""
	result += " \"" + s.SecRuleMetadata.ToSeclang() + ", " + s.SeclangActions.ToString()
	if s.Transformations.ToString() != "" {
		result += ", " + s.Transformations.ToString()
	}
	result += "\"\n"
	return result
}
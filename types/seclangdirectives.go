package types

type SeclangDirective interface {
	ToSeclang() string
}

type ConfigDirective struct {
	CommentMetadata `yaml:",inline"`
	DirectiveName   string `yaml:"directiveName"`
	Parameter       string `yaml:"parameter"`
}

type ConfigurationDirective struct {
	ConfigDirective
}

func (c ConfigurationDirective) ToSeclang() string {
	return c.Comment + c.DirectiveName + " " + c.Parameter + "\n"
}

type SecDefaultAction struct {
	OnlyPhaseMetadata `yaml:"metadata"`
	EmptyVariables    `yaml:"-"`
	Transformations   `yaml:",inline"`
	EmptyOperator     `yaml:"-"`
	SeclangActions    `yaml:"actions"`
}

type SecDefaultActionDirective struct {
	SecDefaultAction
}

func (s SecDefaultActionDirective) ToSeclang() string {
	result := ""
	result += s.Comment + "SecDefaultAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToSeclang()
	transformations := s.Transformations.ToSeclang()
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
	EmptyVariables  `yaml:"-"`
	Transformations `yaml:",inline"`
	EmptyOperator   `yaml:"-"`
	SeclangActions  `yaml:"actions"`
}

type SecActionDirective struct {
	SecAction
}

func (s SecActionDirective) ToSeclang() string {
	result := ""
	result += s.Comment + "SecAction \"phase:" + s.Phase
	actions := s.SeclangActions.ToSeclang()
	transformations := s.Transformations.ToSeclang()
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
	StringOperator  `yaml:"operator"`
	SeclangActions  `yaml:"actions"`
}

type SecRuleDirective struct {
	SecRule
}

func (s SecRuleDirective) ToSeclang() string {
	result := ""
	result += s.Comment + "SecRule "
	result += s.Variables.ToSeclang() + " "
	result += "\"" + s.StringOperator.ToSeclang() + "\""
	result += " \"" + s.SecRuleMetadata.ToSeclang() + ", " + s.SeclangActions.ToSeclang()
	if s.Transformations.ToSeclang() != "" {
		result += ", " + s.Transformations.ToSeclang()
	}
	result += "\"\n"
	return result
}
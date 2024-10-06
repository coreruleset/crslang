package types

type SecAction struct {
	SecRuleMetada   `yaml:"metadata"`
	EmptyVariables  `yaml:"-"`
	Transformations `yaml:",inline"`
	EmptyOperator   `yaml:"-"`
	SeclangActions  `yaml:"actions"`
}

type SecRule struct {
	SecRuleMetada   `yaml:"metadata"`
	Variables       `yaml:",inline"`
	Transformations `yaml:",inline"`
	StringOperator  `yaml:"operator"`
	SeclangActions  `yaml:"actions"`
}
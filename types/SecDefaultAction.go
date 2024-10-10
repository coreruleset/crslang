package types

type SecDefaultAction struct {
	OnlyPhaseMetadata `yaml:"metadata"`
	EmptyVariables    `yaml:"-"`
	Transformations   `yaml:",inline"`
	EmptyOperator     `yaml:"-"`
	SeclangActions    `yaml:"actions"`
}
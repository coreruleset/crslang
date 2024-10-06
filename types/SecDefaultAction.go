package types

type SecDefaultActionMetadata struct {
	EmptyMetadata `yaml:"-"`
	Phase         string `yaml:"phase"`
}

func (s *SecDefaultActionMetadata) SetPhase(value string) {
	s.Phase = value
}

type SecDefaultAction struct {
	SecDefaultActionMetadata `yaml:"metadata"`
	EmptyVariables           `yaml:"-"`
	Transformations          `yaml:",inline"`
	EmptyOperator            `yaml:"-"`
	SeclangActions           `yaml:"actions"`
}
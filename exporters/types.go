package exporters

import "gitlab.fing.edu.uy/gsi/seclang/crslang/types"

type DefaultAction struct {
	Kind                  string                  `yaml:"kind"`
	Metadata              types.OnlyPhaseMetadata `yaml:"metadata"`
	types.Transformations `yaml:",inline"`
	Actions               types.SeclangActions `yaml:"actions"`
}

func (d DefaultAction) ToSeclang() string {
	return "Interface implementation"
}

type ConfigurationDirective struct {
	Kind      string                `yaml:"kind"`
	Metadata  types.CommentMetadata `yaml:",inline"`
	Name      string                `yaml:"name"`
	Parameter string                `yaml:"parameter"`
}

func (d ConfigurationDirective) ToSeclang() string {
	return "Interface implementation"
}

type CommentDirective struct {
	Kind     string                `yaml:"kind"`
	Metadata types.CommentMetadata `yaml:"metadata"`
}

func (d CommentDirective) ToSeclang() string {
	return "Interface implementation"
}

package types

type Variables struct {
	Variables []string `yaml:"variables,omitempty"`
}

func (v *Variables) AddVariable(variable string) {
	v.Variables = append(v.Variables, variable)
}

type EmptyVariables struct {
}

func (e *EmptyVariables) AddVariable(variable string) {
}
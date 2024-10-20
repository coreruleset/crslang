package types

type Variables struct {
	Variables []string `yaml:"variables,omitempty"`
}

func (v *Variables) ToSeclang() string {
	result := ""
	for i, variable := range v.Variables {
		if i != len(v.Variables)-1 {
			result += variable + "|"
		} else {
			result += variable
		}
	}
	return result
}

func (v *Variables) AddVariable(variable string) {
	v.Variables = append(v.Variables, variable)
}

type EmptyVariables struct {
}

func (e *EmptyVariables) AddVariable(variable string) {
}
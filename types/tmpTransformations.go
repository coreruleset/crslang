package types

type Transformations struct {
	Transformations []string `yaml:"transformations,omitempty"`
}

func (t *Transformations) AddTransformation(transformation string) {
	t.Transformations = append(t.Transformations, transformation)
}

func (t *Transformations) ToSeclang() string {
	results := []string{}
	for _, transformation := range t.Transformations {
		results = append(results, transformation)
	}
	result := ""
	for i, value := range results {
		if i == 0 {
			result += "t:" + value
		} else {
			result += ", t:" + value
		}
	}
	return result
}
package types

type Transformations struct {
	Transformations []string `yaml:"transformations,omitempty"`
}

func (t *Transformations) AddTransformation(transformation string) {
	t.Transformations = append(t.Transformations, transformation)
}
package types

type Operator struct {
	Name  string
	Value string `yaml:",omitempty"`
}

func (o *Operator) SetOperatorName(name string) {
	o.Name = name
}

func (o *Operator) SetOperatorValue(value string) {
	o.Value = value
}

func (o *Operator) ToString() string {
	return "@" + o.Name + " " + o.Value
}
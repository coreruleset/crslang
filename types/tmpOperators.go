package types

type StringOperator struct {
	Name  string
	Value string `yaml:",omitempty"`
}

func (o *StringOperator) SetOperatorName(name string) {
	o.Name = name
}

func (o *StringOperator) SetOperatorValue(value string) {
	o.Value = value
}

func (o *StringOperator) ToSeclang() string {
	return "@" + o.Name + " " + o.Value
}

type EmptyOperator struct {
}

func (e *EmptyOperator) SetOperatorName(name string) {
}

func (e *EmptyOperator) SetOperatorValue(value string) {
}
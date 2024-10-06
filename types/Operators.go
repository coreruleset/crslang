package types

import "strconv"

type StringOperator struct {
	Name  string
	Value string `yaml:",omitempty"`
}

type NumericOperator struct {
	Name  string
	Value int `yaml:",omitempty"`
}

func (o *StringOperator) SetOperatorName(name string) {
	o.Name = name
}

func (o *StringOperator) SetOperatorValue(value string) {
	o.Value = value
}

type EmptyOperator struct {
}

func (e *EmptyOperator) SetOperatorName(name string) {
}

func (e *EmptyOperator) SetOperatorValue(value string) {
}

func (o *NumericOperator) SetOperatorName(name string) {
	o.Name = name
}


func (o *NumericOperator) SetOperatorValue(value string) {
	num, err := strconv.Atoi(value)
	if err != nil {
		panic(err) // TODO: handle error in less disruptive way
	}
	o.Value = num
}
package types

import (
	"strconv"
)

type UpdateTargetDirective struct {
	Kind        Kind             `yaml:"kind"`
	Metadata    *CommentMetadata `yaml:",inline"`
	Ids         []int            `yaml:"ids,omitempty"`
	Tags        []string         `yaml:"tags,omitempty"`
	Msgs        []string         `yaml:"msgs,omitempty"`
	Variables   []Variable       `yaml:"variables,omitempty"`
	Collections []Collection     `yaml:"collections,omitempty"`
}

func NewUpdateTargetDirective() *UpdateTargetDirective {
	directive := new(UpdateTargetDirective)
	directive.Metadata = new(CommentMetadata)
	directive.Kind = UpdateTarget
	return directive
}

func (d *UpdateTargetDirective) AddVariable(name string, excluded bool) error {
	variable, err := GetVariable(name)
	if err != nil {
		return err
	}
	d.Variables = append(d.Variables, Variable{Name: variable, Excluded: false})
	return nil
}

func (d *UpdateTargetDirective) AddCollection(name, value string, excluded bool) error {
	col, err := GetCollection(name)
	if err != nil {
		return err
	}
	d.Collections = append(d.Collections, Collection{Name: col, Argument: value, Excluded: excluded})
	return nil
}

func (d *UpdateTargetDirective) ToSeclang() string {
	varResult := ""
	vars := VariablesToString(d.Variables, ",")
	colls := CollectionsToString(d.Collections, ",")
	if vars != "" && colls != "" {
		varResult += vars + "|" + colls
	} else if vars != "" {
		varResult += vars
	} else if colls != "" {
		varResult += colls
	}
	results := ""
	if len(d.Ids) > 0 {
		for _, id := range d.Ids {
			results += "SecRuleUpdateTargetById " + strconv.Itoa(id) + varResult + "\n"
		}
	}
	if len(d.Tags) > 0 {
		for _, tag := range d.Tags {
			results += "SecRuleUpdateTargetByTag \"" + tag + "\"" + varResult + "\n"
		}
	}
	if len(d.Msgs) > 0 {
		for _, msg := range d.Msgs {
			results += "SecRuleUpdateTargetByMsg \"" + msg + "\"" + varResult + "\n"
		}
	}
	return results
}

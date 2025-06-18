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
	if excluded {
		vars := []Variable{}
		for _, v := range d.Variables {
			if v.Name != variable {
				vars = append(vars, v)
			}
		}
		d.Variables = vars
	} else {
		d.Variables = append(d.Variables, Variable{Name: variable, Excluded: false})
	}
	return nil
}

func (d *UpdateTargetDirective) AddCollection(name, value string, excluded, asCount bool) error {
	col, err := GetCollection(name)
	if err != nil {
		return err
	}
	if excluded && !asCount {
		results := []Collection{}
		for _, collection := range d.Collections {
			if collection.Name != col {
				results = append(results, collection)
			} else if value != "" && !collection.Count {
				for i, arg := range collection.Arguments {
					if arg == value {
						collection.Arguments = append(collection.Arguments[:i], collection.Arguments[i+1:]...)
					}
				}
				collection.Excluded = append(collection.Excluded, value)
				results = append(results, collection)
			}
		}
		d.Collections = results
	} else if value != "" && !asCount {
		i := len(d.Collections) - 1
		for i >= 0 && !(!d.Collections[i].Count && d.Collections[i].Name == col && len(d.Collections[i].Arguments) > 0 && len(d.Collections[i].Excluded) == 0) {
			i--
		}
		if i >= 0 {
			d.Collections[i].Arguments = append(d.Collections[i].Arguments, value)
		} else {
			d.Collections = append(d.Collections, Collection{Name: col, Arguments: []string{value}, Excluded: []string{}, Count: asCount})
		}
	} else if value != "" {
		d.Collections = append(d.Collections, Collection{Name: col, Arguments: []string{value}, Excluded: []string{}, Count: asCount})
	} else {
		d.Collections = append(d.Collections, Collection{Name: col, Arguments: []string{}, Excluded: []string{}, Count: asCount})
	}

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

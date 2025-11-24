package types

import (
	"fmt"
	"slices"
)

type SecRule struct {
	Metadata        *SecRuleMetadata `yaml:"metadata,omitempty"`
	Variables       []Variable       `yaml:"variables"`
	Collections     []Collection     `yaml:"collections,omitempty"`
	Transformations `yaml:",inline"`
	Operator        `yaml:"operator"`
	Actions         *SeclangActions    `yaml:"actions,omitempty"`
	ChainedRule     ChainableDirective `yaml:"chainedRule,omitempty"`
}

func NewSecRule() *SecRule {
	secRule := new(SecRule)
	secRule.Metadata = new(SecRuleMetadata)
	secRule.Actions = new(SeclangActions)
	return secRule
}

func (d SecRule) GetKind() Kind {
	return UnknownKind
}

func (d SecRule) GetMetadata() Metadata {
	return d.Metadata
}

func (d SecRule) GetActions() *SeclangActions {
	return d.Actions
}

func (d SecRule) GetTransformations() Transformations {
	return d.Transformations
}

func (s *SecRule) AddVariable(name string, excluded bool) error {
	variable := stringToVariableName(name)
	if variable == UNKNOWN_VAR {
		return fmt.Errorf("Invalid variable name: %s", name)
	}
	if excluded {
		vars := []Variable{}
		for _, v := range s.Variables {
			if v.Name != variable {
				vars = append(vars, v)
			}
		}
		s.Variables = vars
	} else {
		s.Variables = append(s.Variables, Variable{Name: variable, Excluded: false})
	}
	return nil
}

func (s *SecRule) AddCollection(name, value string, excluded, asCount bool) error {
	col := stringToCollectionName(name)
	if col == UNKNOWN_COLLECTION {
		return fmt.Errorf("Invalid collection name: %s", name)
	}
	if excluded && !asCount {
		results := []Collection{}
		for _, collection := range s.Collections {
			if collection.Name != col {
				results = append(results, collection)
			} else if value != "" && !collection.Count && len(collection.Arguments) == 0 {
				collection.Excluded = append(collection.Excluded, value)
				results = append(results, collection)
			} else if value != "" && !collection.Count {
				for i, arg := range collection.Arguments {
					if arg == value {
						collection.Arguments = append(collection.Arguments[:i], collection.Arguments[i+1:]...)
					}
				}
				if len(collection.Arguments) > 0 {
					collection.Excluded = append(collection.Excluded, value)
					results = append(results, collection)
				}
			}
		}
		s.Collections = results
	} else if value != "" && !asCount {
		i := len(s.Collections) - 1
		for i >= 0 && !(!s.Collections[i].Count && s.Collections[i].Name == col && len(s.Collections[i].Arguments) > 0 && len(s.Collections[i].Excluded) == 0) {
			i--
		}
		if i >= 0 {
			s.Collections[i].Arguments = append(s.Collections[i].Arguments, value)
		} else {
			s.Collections = append(s.Collections, Collection{Name: col, Arguments: []string{value}, Excluded: []string{}, Count: asCount})
		}
	} else if value != "" {
		s.Collections = append(s.Collections, Collection{Name: col, Arguments: []string{value}, Excluded: []string{}, Count: asCount})
	} else {
		s.Collections = append(s.Collections, Collection{Name: col, Arguments: []string{}, Excluded: []string{}, Count: asCount})
	}

	return nil
}

func (s SecRule) ToSeclang() string {
	return s.ToSeclangWithIdent("")
}

func (s SecRule) ToSeclangWithIdent(initialString string) string {
	auxString := ",\\\n" + initialString + "    "
	endString := ""

	result := ""
	result += commentsToSeclang(s.Metadata.Comments) + "SecRule "
	vars := VariablesToString(s.Variables, "|")
	colls := CollectionsToString(s.Collections, "|")
	if vars != "" && colls != "" {
		result += vars + "|" + colls
	} else if vars != "" {
		result += vars
	} else if colls != "" {
		result += colls
	}
	result += " \"" + s.Operator.ToString() + "\""
	sortedActions := sortActions(&s)
	for i, action := range sortedActions {
		if i == 0 {
			result += " \\\n" + initialString + "    \""
		} else {
			result += auxString
		}
		result += action
		if i == len(sortedActions)-1 {
			result += "\""
		} else {
			result += endString
		}
	}
	result += "\n"
	if slices.Contains(s.Actions.GetActionKeys(), "chain") {
		result += (s.ChainedRule).ToSeclangWithIdent(initialString + "    ")
	}
	return result
}

func (s SecRule) GetChainedDirective() ChainableDirective {
	return s.ChainedRule
}

func (s *SecRule) AppendChainedDirective(chainedDirective ChainableDirective) {
	s.ChainedRule = chainedDirective
}

func (s SecRule) NonDisruptiveActionsCount() int {
	return len(s.Actions.NonDisruptiveActions)
}

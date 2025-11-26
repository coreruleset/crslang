package types

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Condition represents a condition defined in a rule. It can represent SecActions, SecRules, or Script conditions.
type Condition struct {
	// SecRule conditions are represented by Variables, Collections, Operator, and Transformations.
	Variables       []Variable   `yaml:"variables,omitempty"`
	Collections     []Collection `yaml:"collections,omitempty"`
	Operator        Operator     `yaml:"operator,omitempty"`
	Transformations `yaml:",inline,omitempty"`
	// SecAction conditions are represented by an AlwaysMatch flag and it can also have Transformations.
	AlwaysMatch bool `yaml:"always-match,omitempty"`
	// Script conditions are represented by a ScriptPath.
	Script string `yaml:"script,omitempty"`
}

type RuleWithCondition struct {
	Kind        Kind               `yaml:"kind"`
	Metadata    SecRuleMetadata    `yaml:"metadata,omitempty"`
	Conditions  []Condition        `yaml:"conditions,omitempty"`
	Actions     SeclangActions     `yaml:"actions,omitempty"`
	ChainedRule *RuleWithCondition `yaml:"chainedRule,omitempty"`
}

func (s *RuleWithCondition) ToSeclang() string {
	return "New sec rule with conditions"
}

func (s *RuleWithCondition) GetKind() Kind {
	return s.Kind
}

func ToDirectiveWithConditions(configList ConfigurationList) *ConfigurationList {
	result := new(ConfigurationList)
	for _, config := range configList.DirectiveList {
		configWrapper := new(DirectiveList)
		configWrapper.Id = config.Id
		configWrapper.Marker = config.Marker
		for _, directive := range config.Directives {
			var directiveWrapper SeclangDirective
			switch directive.(type) {
			case CommentMetadata:
				directiveWrapper = CommentDirective{
					Kind:     CommentKind,
					Metadata: directive.(CommentMetadata),
				}
			case *SecAction:
				directiveWrapper = RuleToCondition(directive.(*SecAction))
			case *SecRule:
				directiveWrapper = RuleToCondition(directive.(*SecRule))
			case *SecRuleScript:
				directiveWrapper = RuleToCondition(directive.(*SecRuleScript))
			default:
				directiveWrapper = directive
			}
			configWrapper.Directives = append(configWrapper.Directives, directiveWrapper)
		}
		result.DirectiveList = append(result.DirectiveList, *configWrapper)
	}
	return result
}

func RuleToCondition(directive ChainableDirective) *RuleWithCondition {
	var ruleWithCondition RuleWithCondition
	switch directive.(type) {
	case *SecRule:
		rule := directive.(*SecRule)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*rule.Metadata,
			[]Condition{
				{
					Variables:       rule.Variables,
					Collections:     rule.Collections,
					Operator:        rule.Operator,
					Transformations: rule.Transformations,
				},
			},
			*rule.Actions,
			nil,
		}
	case *SecAction:
		action := directive.(*SecAction)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*action.Metadata,
			[]Condition{
				{
					AlwaysMatch:     true,
					Transformations: action.Transformations,
				},
			},
			*action.Actions,
			nil,
		}
	case *SecRuleScript:
		script := directive.(*SecRuleScript)
		ruleWithCondition = RuleWithCondition{
			"rule",
			*script.Metadata,
			[]Condition{
				{
					Script: script.ScriptPath,
				},
			},
			*script.Actions,
			nil,
		}
	}
	if directive.GetChainedDirective() != nil {
		chainedConditionRule := RuleToCondition(directive.GetChainedDirective())
		if directive.NonDisruptiveActionsCount() > 0 {
			ruleWithCondition.ChainedRule = chainedConditionRule
		} else {
			ruleWithCondition.Conditions = append(ruleWithCondition.Conditions, chainedConditionRule.Conditions...)
			ruleWithCondition.Actions.NonDisruptiveActions = chainedConditionRule.Actions.NonDisruptiveActions
			if chainedConditionRule.ChainedRule != nil {
				ruleWithCondition.ChainedRule = chainedConditionRule.ChainedRule
			}
		}
	}
	return &ruleWithCondition
}

// configurationYamlLoader is a auxiliary struct to load the whole yaml file
type configurationYamlLoader struct {
	Global        DefaultConfigs             `yaml:"global,omitempty"`
	DirectiveList []yamlLoaderConditionRules `yaml:"directivelist,omitempty"`
}

// yamlLoaderConditionRules is a auxiliary struct to load and iterate over the yaml file
type yamlLoaderConditionRules struct {
	Id         string                 `yaml:"id"`
	Directives []yaml.Node            `yaml:"directives,omitempty"`
	Marker     ConfigurationDirective `yaml:"marker,omitempty"`
}

// conditionDirectiveLoader is a auxiliary struct to load condition directives
type conditionDirectiveLoader struct {
	Kind        string          `yaml:"kind"`
	Metadata    SecRuleMetadata `yaml:"metadata,omitempty"`
	Conditions  yaml.Node       `yaml:"conditions,omitempty"`
	Actions     SeclangActions  `yaml:"actions,omitempty"`
	ChainedRule yaml.Node       `yaml:"chainedRule"`
}

// UnmarshalYAML unmarshals a YAML node into a SeclangActions struct
// it converts the actions to their respective types
func (s *SeclangActions) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	for k, v := range raw {
		switch k {
		case "disruptive":
			switch val := v.(type) {
			case string:
				s.SetDisruptiveActionOnly(StringToDisruptiveAction(val))
			case map[string]interface{}:
				if len(v.(map[string]interface{})) != 1 {
					return fmt.Errorf("Error: invalid format for disruptive action")
				}
				for a, p := range val {
					s.SetDisruptiveActionWithParam(StringToDisruptiveAction(a), p.(string))
				}
			}
		case "non-disruptive", "flow", "data":
			for _, action := range v.([]interface{}) {
				switch act := action.(type) {
				case string:
					switch k {
					case "non-disruptive":
						s.AddNonDisruptiveActionOnly(StringToNonDisruptiveAction(act))
					case "flow":
						s.AddFlowActionOnly(StringToFlowAction(act))
					}
				case map[string]interface{}:
					if len(act) != 1 {
						return fmt.Errorf("Error: invalid format for non-disruptive action")
					}
					for a, p := range act {
						if a == "setvar" {
							switch p := p.(type) {
							case map[string]interface{}:
								colName, ok := p["collection"]
								if _, parseOk := colName.(string); !ok || !parseOk {
									colName = "tx"
								}
								op, ok := p["operation"]
								if _, parseOk := op.(string); !ok || !parseOk {
									op = "="
								}
								assigns, ok := p["assignments"]
								if _, parseOk := assigns.([]interface{}); !ok || !parseOk {
									return fmt.Errorf("Error: setvar actions must have assignments and assignments must be a list")
								}
								parsedAssigns := []VarAssignment{}
								for _, v := range assigns.([]interface{}) {
									castedAssign, ok := v.(map[string]interface{})
									if !ok || len(castedAssign) != 1 {
										return fmt.Errorf("Error: invalid variable assignment format: %T", v)
									}
									for varName, varValue := range castedAssign {
										if sVarValue, ok := varValue.(string); ok {
											parsedAssigns = append(parsedAssigns, VarAssignment{Variable: varName, Value: sVarValue})
										} else {
											return fmt.Errorf("Error: assignment must be a string")
										}
									}
								}
								cName := stringToCollectionName(strings.ToUpper(colName.(string)))
								if cName == UNKNOWN_COLLECTION {
									return fmt.Errorf("Collection name %s is not valid", cName)
								}
								vOp := stringToVarOperation(op.(string))
								if vOp == UnknownOp {
									return fmt.Errorf("invalid setvar action: invalid operation '%s'", op)
								}
								newAct, err := NewSetvarAction(cName, vOp, parsedAssigns)
								if err != nil {
									return err
								}
								s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAct)
							case []interface{}:
								parsedAssigns := []VarAssignment{}
								for _, v := range p {
									castedAssign, ok := v.(map[string]interface{})
									if !ok || len(castedAssign) != 1 {
										return fmt.Errorf("Error: invalid variable assignment format: %T", v)
									}
									for varName, varValue := range castedAssign {
										if sVarValue, ok := varValue.(string); ok {
											parsedAssigns = append(parsedAssigns, VarAssignment{Variable: varName, Value: sVarValue})
										} else {
											return fmt.Errorf("Error: assignment must be a string")
										}
									}
								}
								newAct, err := NewSetvarAction(TX, Assign, parsedAssigns)
								if err != nil {
									return err
								}
								s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAct)
							default:
								return fmt.Errorf("Error: invalid format for setvar action: %T", p)
							}
						} else {
							switch k {
							case "non-disruptive":
								s.AddNonDisruptiveActionWithParam(StringToNonDisruptiveAction(a), p.(string))
							case "flow":
								s.AddFlowActionWithParam(StringToFlowAction(a), p.(string))
							case "data":
								s.AddDataActionWithParams(StringToDataAction(a), p.(string))
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// LoadDirectivesWithConditionsFromFile loads condition format directives from a yaml file
func LoadDirectivesWithConditionsFromFile(filename string) ConfigurationList {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return LoadDirectivesWithConditions(yamlFile)
}

// LoadDirectivesWithConditions loads condition format directives from a yaml file
func LoadDirectivesWithConditions(yamlFile []byte) ConfigurationList {
	var config configurationYamlLoader
	err := yaml.Unmarshal(yamlFile, &config)
	configs := config.DirectiveList
	if err != nil {
		panic(err)
	}
	var resultConfigs []DirectiveList
	for _, config := range configs {
		var directives []SeclangDirective
		for _, yamlDirective := range config.Directives {
			directive := loadConditionDirective(yamlDirective)
			if directive == nil {
				panic("Unknown directive type")
			} else {
				directives = append(directives, directive)
			}
		}
		resultConfigs = append(resultConfigs, DirectiveList{Id: config.Id, Directives: directives, Marker: config.Marker})
	}
	return ConfigurationList{Global: config.Global, DirectiveList: resultConfigs}
}

// loadConditionDirective loads the different kind of directives
func loadConditionDirective(yamlDirective yaml.Node) SeclangDirective {
	if yamlDirective.Kind != yaml.MappingNode {
		panic("Unknown format type")
	}
	switch yamlDirective.Content[1].Value {
	case "comment":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := CommentDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "configuration":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := ConfigurationDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "defaultaction":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := DefaultAction{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "rule":
		return loadRuleWithConditions(yamlDirective)
	case "remove":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := RemoveRuleDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return directive
	case "update_target":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := UpdateTargetDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return &directive
	case "update_action":
		rawDirective, err := yaml.Marshal(yamlDirective)
		if err != nil {
			panic(err)
		}
		directive := UpdateActionDirective{}
		err = yaml.Unmarshal(rawDirective, &directive)
		if err != nil {
			panic(err)
		}
		return &directive
	}
	return nil
}

// loadRuleWithConditions loads a rule with conditions in a recursive way
func loadRuleWithConditions(yamlDirective yaml.Node) *RuleWithCondition {
	rawDirective := []byte{}
	var err error

	rawDirective, err = yaml.Marshal(yamlDirective)

	if err != nil {
		panic(err)
	}

	loaderDirective := conditionDirectiveLoader{}
	err = yaml.Unmarshal(rawDirective, &loaderDirective)
	if err != nil {
		print(string(rawDirective))
		panic(err)
	}
	directive := &RuleWithCondition{
		Kind:     RuleKind,
		Metadata: loaderDirective.Metadata,
		Actions:  loaderDirective.Actions,
	}
	if loaderDirective.Conditions.Kind == yaml.SequenceNode {
		for _, condition := range loaderDirective.Conditions.Content {
			loadedCondition, err := castConditions(condition)
			if err != nil {
				panic(err)
			}
			directive.Conditions = append(directive.Conditions, loadedCondition)
		}
	}
	var loadedChainedRule *RuleWithCondition
	if len(loaderDirective.ChainedRule.Content) > 0 {
		loadedChainedRule = loadRuleWithConditions(loaderDirective.ChainedRule)
		directive.ChainedRule = loadedChainedRule
	}
	return directive
}

// castConditions casts a directive condition to the correct type
func castConditions(condition *yaml.Node) (Condition, error) {
	rawDirective, err := yaml.Marshal(condition)
	if err != nil {
		return Condition{}, err
	}
	ruleCondition := Condition{}
	err = yaml.Unmarshal(rawDirective, &ruleCondition)
	if err != nil {
		return Condition{}, err
	}
	// switch condition.Content[0].Value {
	// case "alwaysMatch":
	// 	rawDirective, err := yaml.Marshal(condition)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ruleCondition := Condition{}
	// 	err = yaml.Unmarshal(rawDirective, &ruleCondition)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return ruleCondition
	// case "script":
	// 	return Condition{Script: condition.Content[1].Value}
	// case "variables", "collections":
	// 	rawDirective, err := yaml.Marshal(condition)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ruleCondition := Condition{}
	// 	err = yaml.Unmarshal(rawDirective, &ruleCondition)
	// 	return ruleCondition
	// }
	return ruleCondition, nil
}

func FromCRSLangToUnformattedDirectives(configListWrapped ConfigurationList) *ConfigurationList {
	result := new(ConfigurationList)
	for _, config := range configListWrapped.DirectiveList {
		configList := new(DirectiveList)
		configList.Id = config.Id
		configList.Marker = config.Marker
		for _, directiveWrapped := range config.Directives {
			var directive SeclangDirective
			switch directiveWrapped.(type) {
			case CommentDirective:
				directive = directiveWrapped.(CommentDirective).Metadata
			case DefaultAction:
				directive = directiveWrapped
			case *RuleWithCondition:
				chainableDir := FromConditionToUnmorfattedDirective(*directiveWrapped.(*RuleWithCondition))
				if configListWrapped.Global.Version != "" {
					chainableDir.GetMetadata().SetVer(configListWrapped.Global.Version)
				}
				for _, tag := range configListWrapped.Global.Tags {
					chainableDir.GetMetadata().AddTag(tag)
				}
				directive = chainableDir
			case ConfigurationDirective:
				directive = ConfigurationDirective{
					Metadata:  directiveWrapped.(ConfigurationDirective).Metadata,
					Name:      directiveWrapped.(ConfigurationDirective).Name,
					Parameter: directiveWrapped.(ConfigurationDirective).Parameter,
				}
			}
			configList.Directives = append(configList.Directives, directive)
		}
		result.DirectiveList = append(result.DirectiveList, *configList)
	}
	return result
}

func FromConditionToUnmorfattedDirective(conditionDirective RuleWithCondition) ChainableDirective {
	var rootDirective ChainableDirective
	var directiveIterator ChainableDirective
	var chainedDirective ChainableDirective
	var directiveAux ChainableDirective

	chainedDirective = nil

	if conditionDirective.ChainedRule != nil {
		chainedDirective = FromConditionToUnmorfattedDirective(*conditionDirective.ChainedRule)
	}

	for i, condition := range conditionDirective.Conditions {
		if condition.Operator.Name != UnknownOperator {
			secruleDirective := new(SecRule)
			secruleDirective.Variables = condition.Variables
			secruleDirective.Collections = condition.Collections
			secruleDirective.Transformations = condition.Transformations
			secruleDirective.Operator = condition.Operator
			if i == 0 {
				secruleDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secruleDirective.Actions = CopyActions(conditionDirective.Actions)
				secruleDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secruleDirective
			} else {
				secruleDirective.Metadata = new(SecRuleMetadata)
				secruleDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					chainAction, err := NewActionOnly(Chain)
					if err != nil {
						panic(fmt.Sprintf("failed to create chain action: %v", err))
					}
					secruleDirective.Actions.FlowActions = []Action{chainAction}
				}
			}
			directiveAux = secruleDirective
		} else if condition.AlwaysMatch {
			secactionDirective := new(SecAction)
			secactionDirective.Transformations = condition.Transformations
			if i == 0 {
				secactionDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secactionDirective.Actions = CopyActions(conditionDirective.Actions)
				secactionDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secactionDirective
			} else {
				secactionDirective.Metadata = new(SecRuleMetadata)
				secactionDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					chainAction, err := NewActionOnly(Chain)
					if err != nil {
						panic(fmt.Sprintf("failed to create chain action: %v", err))
					}
					secactionDirective.Actions.FlowActions = []Action{chainAction}
				}
			}
			directiveAux = secactionDirective
		} else if condition.Script != "" {
			secscriptDirective := new(SecRuleScript)
			secscriptDirective.ScriptPath = condition.Script
			if i == 0 {
				secscriptDirective.Metadata = CopySecRuleMetadata(conditionDirective.Metadata)
				secscriptDirective.Actions = CopyActions(conditionDirective.Actions)
				secscriptDirective.Actions.NonDisruptiveActions = []Action{}
				rootDirective = secscriptDirective
			} else {
				secscriptDirective.Metadata = new(SecRuleMetadata)
				secscriptDirective.Actions = new(SeclangActions)
				if i < len(conditionDirective.Conditions)-1 || chainedDirective != nil {
					chainAction, err := NewActionOnly(Chain)
					if err != nil {
						panic(fmt.Sprintf("failed to create chain action: %v", err))
					}
					secscriptDirective.Actions.FlowActions = []Action{chainAction}
				}
			}
			directiveAux = secscriptDirective
		}
		if i == 0 {
			directiveIterator = rootDirective
		} else {
			directiveIterator.AppendChainedDirective(directiveAux)
			directiveIterator = directiveAux
		}

	}

	switch directiveIterator.(type) {
	case *SecRule:
		directiveIterator.(*SecRule).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *SecAction:
		directiveIterator.(*SecAction).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	case *SecRuleScript:
		directiveIterator.(*SecRuleScript).Actions.NonDisruptiveActions = conditionDirective.Actions.NonDisruptiveActions
	}

	if chainedDirective != nil {
		directiveIterator.AppendChainedDirective(chainedDirective)
	}

	return rootDirective
}

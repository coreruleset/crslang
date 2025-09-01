package types

import (
	"fmt"
)

type SeclangActions struct {
	DisruptiveAction     Action   `yaml:"disruptive,omitempty"`
	NonDisruptiveActions []Action `yaml:"non-disruptive,omitempty"`
	FlowActions          []Action `yaml:"flow,omitempty"`
	DataActions          []Action `yaml:"data,omitempty"`
}

func (s *SeclangActions) ToString() string {
	results := []string{}
	if len(s.DisruptiveAction) > 0 {
		results = append(results, s.DisruptiveAction.ToString())
	}
	for _, action := range s.NonDisruptiveActions {
		results = append(results, action.ToString())
	}
	for _, action := range s.FlowActions {
		results = append(results, action.ToString())
	}
	for _, action := range s.DataActions {
		results = append(results, action.ToString())
	}
	result := ""
	for i, value := range results {
		if i == 0 {
			result += value
		} else {
			result += ", " + value
		}
	}
	return result
}

func (s *SeclangActions) String() string {
	return fmt.Sprintf("Disruptive: %v, NonDisruptive: %v, Flow: %v, Data: %v", s.DisruptiveAction, s.NonDisruptiveActions, s.FlowActions, s.DataActions)
}

// Action represents a single action with its parameters
// It's a map where the key is the action name and the value is the parameter
type Action map[string]string

// ToString converts the action to its string representation
func (a Action) ToString() string {
	if len(a) == 0 {
		return ""
	}

	// Get the first (and should be only) key-value pair
	for action, param := range a {
		if param == "" {
			return action
		} else {
			return action + ":'" + param + "'"
		}
	}
	return ""
}

// GetKey returns the action name (first key in the map)
func (a Action) GetKey() string {
	for key := range a {
		return key
	}
	return ""
}

// GetParam returns the parameter value for the action
func (a Action) GetParam() string {
	for _, value := range a {
		return value
	}
	return ""
}

// ActionType is a constraint for all action types
type ActionType interface {
	DisruptiveAction | FlowAction | DataAction | NonDisruptiveAction
}

// NewAction creates a new Action with the given action type and parameter
// It uses generics to accept DisruptiveAction, FlowAction, DataAction, or NonDisruptiveAction
func NewAction[T ActionType](action T, param string) Action {
	actionStr := string(action)
	if param == "" {
		return Action{actionStr: ""}
	}
	return Action{actionStr: param}
}

type DisruptiveAction string

const (
	Allow    DisruptiveAction = "allow"
	Block    DisruptiveAction = "block"
	Deny     DisruptiveAction = "deny"
	Drop     DisruptiveAction = "drop"
	Pass     DisruptiveAction = "pass"
	Pause    DisruptiveAction = "pause"
	Proxy    DisruptiveAction = "proxy"
	Redirect DisruptiveAction = "redirect"
)

type FlowAction string

const (
	Chain     FlowAction = "chain"
	Skip      FlowAction = "skip"
	SkipAfter FlowAction = "skipAfter"
)

type DataAction string

const (
	Data   DataAction = "data"
	Status DataAction = "status"
	XLMNS  DataAction = "xmlns"
)

type NonDisruptiveAction string

const (
	Append                 NonDisruptiveAction = "append"
	AuditLog               NonDisruptiveAction = "auditlog"
	Capture                NonDisruptiveAction = "capture"
	Ctl                    NonDisruptiveAction = "ctl"
	DeprecateVar           NonDisruptiveAction = "deprecatevar"
	Exec                   NonDisruptiveAction = "exec"
	ExpireVar              NonDisruptiveAction = "expirevar"
	InitCol                NonDisruptiveAction = "initcol"
	Log                    NonDisruptiveAction = "log"
	LogData                NonDisruptiveAction = "logdata"
	MultiMatch             NonDisruptiveAction = "multiMatch"
	NoAuditLog             NonDisruptiveAction = "noauditlog"
	NoLog                  NonDisruptiveAction = "nolog"
	Prepend                NonDisruptiveAction = "prepend"
	SanitiseArg            NonDisruptiveAction = "sanitiseArg"
	SanitiseMatched        NonDisruptiveAction = "sanitiseMatched"
	SanitiseMatchedBytes   NonDisruptiveAction = "sanitiseMatchedBytes"
	SanitiseRequestHeader  NonDisruptiveAction = "sanitiseRequestHeader"
	SanitiseResponseHeader NonDisruptiveAction = "sanitiseResponseHeader"
	SetUid                 NonDisruptiveAction = "setuid"
	SetRsc                 NonDisruptiveAction = "setrsc"
	SetSid                 NonDisruptiveAction = "setsid"
	SetEnv                 NonDisruptiveAction = "setenv"
	SetVar                 NonDisruptiveAction = "setvar"
)

func (s *SeclangActions) SetDisruptiveActionWithParam(action DisruptiveAction, value string) error {
	s.DisruptiveAction = NewAction(action, value)
	return nil
}

func (s *SeclangActions) SetDisruptiveActionOnly(action DisruptiveAction) error {
	s.DisruptiveAction = NewAction(action, "")
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action NonDisruptiveAction, param string) error {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, NewAction(action, param))
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action NonDisruptiveAction) error {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, NewAction(action, ""))
	return nil
}

func (s *SeclangActions) AddFlowActionWithParam(action FlowAction, param string) error {
	s.FlowActions = append(s.FlowActions, NewAction(action, param))
	return nil
}

func (s *SeclangActions) AddFlowActionOnly(action FlowAction) error {
	s.FlowActions = append(s.FlowActions, NewAction(action, ""))
	return nil
}

func (s *SeclangActions) AddDataActionWithParams(action DataAction, param string) error {
	s.DataActions = append(s.DataActions, NewAction(action, param))
	return nil
}

func CopyActions(a SeclangActions) *SeclangActions {
	newActions := new(SeclangActions)
	newActions.DisruptiveAction = a.DisruptiveAction
	newActions.NonDisruptiveActions = make([]Action, len(a.NonDisruptiveActions))
	copy(newActions.NonDisruptiveActions, a.NonDisruptiveActions)
	newActions.FlowActions = make([]Action, len(a.FlowActions))
	copy(newActions.FlowActions, a.FlowActions)
	newActions.DataActions = make([]Action, len(a.DataActions))
	copy(newActions.DataActions, a.DataActions)
	return newActions
}

func (s *SeclangActions) GetActionKeys() []string {
	keys := []string{}
	// if s.DisruptiveAction != nil {
	// 	keys = append(keys, s.DisruptiveAction.ToString())
	// }
	for _, action := range s.NonDisruptiveActions {
		keys = append(keys, action.GetKey())
	}
	for _, action := range s.FlowActions {
		keys = append(keys, action.GetKey())
	}
	for _, action := range s.DataActions {
		keys = append(keys, action.GetKey())
	}
	return keys
}

func (s *SeclangActions) GetActionByKey(key string) Action {
	// if s.DisruptiveAction != nil {
	// 	if s.DisruptiveAction.ToString() == key {
	// 		return s.DisruptiveAction
	// 	}
	// }
	for _, action := range s.NonDisruptiveActions {
		if action.GetKey() == key {
			return action
		}
	}
	for _, action := range s.FlowActions {
		if action.GetKey() == key {
			return action
		}
	}
	for _, action := range s.DataActions {
		if action.GetKey() == key {
			return action
		}
	}
	return Action{}
}

func (s *SeclangActions) GetActionsByKey(key string) []Action {
	actions := []Action{}
	// if s.DisruptiveAction != nil {
	// 	if s.DisruptiveAction.ToString() == key {
	// 		actions = append(actions, s.DisruptiveAction)
	// 	}
	// }
	for _, action := range s.NonDisruptiveActions {
		if action.GetKey() == key {
			actions = append(actions, action)
		}
	}
	for _, action := range s.FlowActions {
		if action.GetKey() == key {
			actions = append(actions, action)
		}
	}
	for _, action := range s.DataActions {
		if action.GetKey() == key {
			actions = append(actions, action)
		}
	}
	return actions
}

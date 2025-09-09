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
	if s.DisruptiveAction != nil {
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

// Action interface represents a generic action
type Action interface {
	GetKey() string
	ToString() string
}

// ActionOnly represents an action without parameters
type ActionOnly string

// GetKey returns the action name
func (a ActionOnly) GetKey() string {
	return string(a)
}

// ToString returns the string representation of the action
func (a ActionOnly) ToString() string {
	return string(a)
}

// NewActionOnly creates a new NewActionOnly with the given action type
// It uses generics to accept DisruptiveAction, FlowAction, DataAction, or NonDisruptiveAction
func NewActionOnly[T ActionType](action T) Action {
	actionStr := string(action)
	return ActionOnly(actionStr)
}

// ActionWithParam represents a single action with its parameters
// It's a map where the key is the action name and the value is the parameter
type ActionWithParam map[string]string

// ToString converts the action to its string representation
func (a ActionWithParam) ToString() string {
	if len(a) == 0 {
		return ""
	}

	// Get the first (and should be only) key-value pair
	for action, param := range a {
		if param == "" {
			return action
		} else if action == "ctl" {
			return action + ":" + param + ""
		} else {
			return action + ":'" + param + "'"
		}
	}
	return ""
}

// GetKey returns the action name (first key in the map)
func (a ActionWithParam) GetKey() string {
	for key := range a {
		return key
	}
	return ""
}

// GetParam returns the parameter value for the action
func (a ActionWithParam) GetParam() string {
	for _, value := range a {
		return value
	}
	return ""
}

// ActionType is a constraint for all action types
type ActionType interface {
	DisruptiveAction | FlowAction | DataAction | NonDisruptiveAction
}

// NewActionWithParam creates a new NewActionWithParam with the given action type and parameter
// It uses generics to accept DisruptiveAction, FlowAction, DataAction, or NonDisruptiveAction
func NewActionWithParam[T ActionType](action T, param string) Action {
	actionStr := string(action)
	if param == "" {
		return ActionWithParam{actionStr: ""}
	}
	return ActionWithParam{actionStr: param}
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
	s.DisruptiveAction = NewActionWithParam(action, value)
	return nil
}

func (s *SeclangActions) SetDisruptiveActionOnly(action DisruptiveAction) error {
	s.DisruptiveAction = NewActionOnly(action)
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action NonDisruptiveAction, param string) error {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, NewActionWithParam(action, param))
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action NonDisruptiveAction) error {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, NewActionOnly(action))
	return nil
}

func (s *SeclangActions) AddFlowActionWithParam(action FlowAction, param string) error {
	s.FlowActions = append(s.FlowActions, NewActionWithParam(action, param))
	return nil
}

func (s *SeclangActions) AddFlowActionOnly(action FlowAction) error {
	s.FlowActions = append(s.FlowActions, NewActionOnly(action))
	return nil
}

func (s *SeclangActions) AddDataActionWithParams(action DataAction, param string) error {
	s.DataActions = append(s.DataActions, NewActionWithParam(action, param))
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
	if s.DisruptiveAction != nil {
		keys = append(keys, s.DisruptiveAction.GetKey())
	}
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
	if s.DisruptiveAction != nil && s.DisruptiveAction.GetKey() == key {
		return s.DisruptiveAction
	}
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
	return ActionWithParam{}
}

func (s *SeclangActions) GetActionsByKey(key string) []ActionWithParam {
	actions := []ActionWithParam{}
	// if s.DisruptiveAction != nil {
	// 	if s.DisruptiveAction.ToString() == key {
	// 		actions = append(actions, s.DisruptiveAction)
	// 	}
	// }
	for _, action := range s.NonDisruptiveActions {
		if action.GetKey() == key {
			aP, ok := action.(ActionWithParam)
			if ok {
				actions = append(actions, aP)
			}
		}
	}
	for _, action := range s.FlowActions {
		if action.GetKey() == key {
			aP, ok := action.(ActionWithParam)
			if ok {
				actions = append(actions, aP)
			}
		}
	}
	for _, action := range s.DataActions {
		if action.GetKey() == key {
			aP, ok := action.(ActionWithParam)
			if ok {
				actions = append(actions, aP)
			}
		}
	}
	return actions
}

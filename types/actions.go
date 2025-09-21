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
		} else if action == "ctl" {
			return action + ":" + param + ""
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
	String() string
}

// NewAction creates a new Action with the given action type and parameter
// It uses generics to accept DisruptiveAction, FlowAction, DataAction, or NonDisruptiveAction
func NewAction[T ActionType](action T, param string) (Action, error) {
	// Use the String() method to get the string representation
	actionStr := action.String()

	// Check if the action is an Unknown value
	if actionStr == "unknown" {
		return Action{}, fmt.Errorf("invalid action: unknown action type")
	}

	if param == "" {
		return Action{actionStr: ""}, nil
	}
	return Action{actionStr: param}, nil
}

type DisruptiveAction int

const (
	Allow DisruptiveAction = iota
	Block
	Deny
	Drop
	Pass
	Pause
	Proxy
	Redirect
	Unknown
)

func (d DisruptiveAction) String() string {
	switch d {
	case Allow:
		return "allow"
	case Block:
		return "block"
	case Deny:
		return "deny"
	case Drop:
		return "drop"
	case Pass:
		return "pass"
	case Pause:
		return "pause"
	case Proxy:
		return "proxy"
	case Redirect:
		return "redirect"
	case Unknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type FlowAction int

const (
	Chain FlowAction = iota
	Skip
	SkipAfter
	FlowUnknown
)

func (f FlowAction) String() string {
	switch f {
	case Chain:
		return "chain"
	case Skip:
		return "skip"
	case SkipAfter:
		return "skipAfter"
	case FlowUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type DataAction int

const (
	Data DataAction = iota
	Status
	XLMNS
	DataUnknown
)

func (d DataAction) String() string {
	switch d {
	case Data:
		return "data"
	case Status:
		return "status"
	case XLMNS:
		return "xmlns"
	case DataUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type NonDisruptiveAction int

const (
	Append NonDisruptiveAction = iota
	AuditLog
	Capture
	Ctl
	DeprecateVar
	Exec
	ExpireVar
	InitCol
	Log
	LogData
	MultiMatch
	NoAuditLog
	NoLog
	Prepend
	SanitiseArg
	SanitiseMatched
	SanitiseMatchedBytes
	SanitiseRequestHeader
	SanitiseResponseHeader
	SetUid
	SetRsc
	SetSid
	SetEnv
	SetVar
	NonDisruptiveUnknown
)

func (n NonDisruptiveAction) String() string {
	switch n {
	case Append:
		return "append"
	case AuditLog:
		return "auditlog"
	case Capture:
		return "capture"
	case Ctl:
		return "ctl"
	case DeprecateVar:
		return "deprecatevar"
	case Exec:
		return "exec"
	case ExpireVar:
		return "expirevar"
	case InitCol:
		return "initcol"
	case Log:
		return "log"
	case LogData:
		return "logdata"
	case MultiMatch:
		return "multiMatch"
	case NoAuditLog:
		return "noauditlog"
	case NoLog:
		return "nolog"
	case Prepend:
		return "prepend"
	case SanitiseArg:
		return "sanitiseArg"
	case SanitiseMatched:
		return "sanitiseMatched"
	case SanitiseMatchedBytes:
		return "sanitiseMatchedBytes"
	case SanitiseRequestHeader:
		return "sanitiseRequestHeader"
	case SanitiseResponseHeader:
		return "sanitiseResponseHeader"
	case SetUid:
		return "setuid"
	case SetRsc:
		return "setrsc"
	case SetSid:
		return "setsid"
	case SetEnv:
		return "setenv"
	case SetVar:
		return "setvar"
	case NonDisruptiveUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (s *SeclangActions) SetDisruptiveActionWithParam(action DisruptiveAction, value string) error {
	newAction, err := NewAction(action, value)
	if err != nil {
		return err
	}
	s.DisruptiveAction = newAction
	return nil
}

func (s *SeclangActions) SetDisruptiveActionOnly(action DisruptiveAction) error {
	newAction, err := NewAction(action, "")
	if err != nil {
		return err
	}
	s.DisruptiveAction = newAction
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action NonDisruptiveAction, param string) error {
	newAction, err := NewAction(action, param)
	if err != nil {
		return err
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action NonDisruptiveAction) error {
	newAction, err := NewAction(action, "")
	if err != nil {
		return err
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
	return nil
}

func (s *SeclangActions) AddFlowActionWithParam(action FlowAction, param string) error {
	newAction, err := NewAction(action, param)
	if err != nil {
		return err
	}
	s.FlowActions = append(s.FlowActions, newAction)
	return nil
}

func (s *SeclangActions) AddFlowActionOnly(action FlowAction) error {
	newAction, err := NewAction(action, "")
	if err != nil {
		return err
	}
	s.FlowActions = append(s.FlowActions, newAction)
	return nil
}

func (s *SeclangActions) AddDataActionWithParams(action DataAction, param string) error {
	newAction, err := NewAction(action, param)
	if err != nil {
		return err
	}
	s.DataActions = append(s.DataActions, newAction)
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
	keys = append(keys, s.DisruptiveAction.GetKey())
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
	if s.DisruptiveAction.GetKey() == key {
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

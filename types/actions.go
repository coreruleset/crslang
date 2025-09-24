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
func NewActionOnly[T ActionType](action T) (Action, error) {
	actionStr := action.String()

	// Check if the action is an Unknown value
	if actionStr == "unknown" {
		return ActionOnly(""), fmt.Errorf("invalid action: unknown action type")
	}
	return ActionOnly(actionStr), nil
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
	String() string
}

// NewActionWithParam creates a new NewActionWithParam with the given action type and parameter
// It uses generics to accept DisruptiveAction, FlowAction, DataAction, or NonDisruptiveAction
func NewActionWithParam[T ActionType](action T, param string) (ActionWithParam, error) {
	// Use the String() method to get the string representation
	actionStr := action.String()

	// Check if the action is an Unknown value
	if actionStr == "unknown" {
		return ActionWithParam{}, fmt.Errorf("invalid action: unknown action type")
	}

	return ActionWithParam{actionStr: param}, nil
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

// Helper functions to convert string to action types
func StringToDisruptiveAction(s string) DisruptiveAction {
	switch s {
	case "allow":
		return Allow
	case "block":
		return Block
	case "deny":
		return Deny
	case "drop":
		return Drop
	case "pass":
		return Pass
	case "pause":
		return Pause
	case "proxy":
		return Proxy
	case "redirect":
		return Redirect
	default:
		return Unknown
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

func StringToFlowAction(s string) FlowAction {
	switch s {
	case "chain":
		return Chain
	case "skip":
		return Skip
	case "skipAfter":
		return SkipAfter
	default:
		return FlowUnknown
	}
}

type DataAction int

const (
	DataUnknown DataAction = iota
	Status
	XLMNS
)

func (d DataAction) String() string {
	switch d {
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

func StringToDataAction(s string) DataAction {
	switch s {
	case "status":
		return Status
	case "xmlns":
		return XLMNS
	default:
		return DataUnknown
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

func StringToNonDisruptiveAction(s string) NonDisruptiveAction {
	switch s {
	case "append":
		return Append
	case "auditlog":
		return AuditLog
	case "capture":
		return Capture
	case "ctl":
		return Ctl
	case "deprecatevar":
		return DeprecateVar
	case "exec":
		return Exec
	case "expirevar":
		return ExpireVar
	case "initcol":
		return InitCol
	case "log":
		return Log
	case "logdata":
		return LogData
	case "multiMatch":
		return MultiMatch
	case "noauditlog":
		return NoAuditLog
	case "nolog":
		return NoLog
	case "prepend":
		return Prepend
	case "sanitiseArg":
		return SanitiseArg
	case "sanitiseMatched":
		return SanitiseMatched
	case "sanitiseMatchedBytes":
		return SanitiseMatchedBytes
	case "sanitiseRequestHeader":
		return SanitiseRequestHeader
	case "sanitiseResponseHeader":
		return SanitiseResponseHeader
	case "setuid":
		return SetUid
	case "setrsc":
		return SetRsc
	case "setsid":
		return SetSid
	case "setenv":
		return SetEnv
	case "setvar":
		return SetVar
	default:
		return NonDisruptiveUnknown
	}
}

func (s *SeclangActions) SetDisruptiveActionWithParam(action DisruptiveAction, value string) error {
	newAction, err := NewActionWithParam(action, value)
	if err != nil {
		return err
	}
	s.DisruptiveAction = newAction
	return nil
}

func (s *SeclangActions) SetDisruptiveActionOnly(action DisruptiveAction) error {
	newAction, err := NewActionOnly(action)
	if err != nil {
		return err
	}
	s.DisruptiveAction = newAction
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action NonDisruptiveAction, param string) error {
	newAction, err := NewActionWithParam(action, param)
	if err != nil {
		return err
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action NonDisruptiveAction) error {
	newAction, err := NewActionOnly(action)
	if err != nil {
		return err
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
	return nil
}

func (s *SeclangActions) AddFlowActionWithParam(action FlowAction, param string) error {
	newAction, err := NewActionWithParam(action, param)
	if err != nil {
		return err
	}
	s.FlowActions = append(s.FlowActions, newAction)
	return nil
}

func (s *SeclangActions) AddFlowActionOnly(action FlowAction) error {
	newAction, err := NewActionOnly(action)
	if err != nil {
		return err
	}
	s.FlowActions = append(s.FlowActions, newAction)
	return nil
}

func (s *SeclangActions) AddDataActionWithParams(action DataAction, param string) error {
	newAction, err := NewActionWithParam(action, param)
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

package types

import (
	"fmt"
	"strings"
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

type VarAssignment struct {
	Variable string `yaml:"variable"`
	Value    string `yaml:"value"`
}

type SetvarAction struct {
	Collection  CollectionName  `yaml:"collection,omitempty"`
	Operation   string          `yaml:"operation,omitempty"`
	Assignments []VarAssignment `yaml:"assignments,omitempty"`
}

// GetKey returns the action name (it is always "setvar")
func (a SetvarAction) GetKey() string {
	return SetVar.String()
}

// ToString allows to implement the Action interface
func (a SetvarAction) ToString() string {
	if len(a.Assignments) == 0 {
		return ""
	}

	var result []string
	// Reconstruct the setvar actions
	for _, asg := range a.Assignments {
		result = append(result, SetVar.String()+":"+a.Collection.String()+"."+asg.Variable+a.Operation+asg.Value)
	}
	return strings.Join(result, ", ")
}

func (a *SetvarAction) AppendAssignment(variable, value string) error {
	a.Assignments = append(a.Assignments, VarAssignment{Variable: variable, Value: value})
	return nil
}

func (a SetvarAction) GetAllParams() []string {
	if len(a.Assignments) == 0 {
		return []string{}
	}

	var result []string
	// Get all the variables
	for _, asg := range a.Assignments {
		res := SetVar.String() + ":" + a.Collection.String() + "." + asg.Variable + a.Operation + asg.Value
		result = append(result, res)
	}
	return result
}

func (s VarAssignment) MarshalYAML() (interface{}, error) {
	if s.Variable == "" || s.Value == "" {
		return nil, fmt.Errorf("invalid variable assignment: missing variable name or value")
	}
	return map[string]string{s.Variable: s.Value}, nil
}

func (s SetvarAction) MarshalYAML() (interface{}, error) {
	if s.Collection == UNKNOWN_COLLECTION || s.Operation == "" || len(s.Assignments) == 0 {
		return nil, fmt.Errorf("invalid setvar action: missing collection name, operation, or assignments")
	}
	if s.Collection == TX && s.Operation == "=" {
		// Default case
		res := map[string][]VarAssignment{}
		res["setvar"] = s.Assignments
		return res, nil
	} else {
		// Non-default case, collection is different to `TX` or operation is different to `=`.
		// Fields are re-mapped to a mirrored struct in order to preserve the order in the YAML
		res := map[string]struct {
			Collection  CollectionName
			Operation   string
			Assignments []VarAssignment
		}{}
		res["setvar"] = struct {
			Collection  CollectionName
			Operation   string
			Assignments []VarAssignment
		}{
			Collection:  s.Collection,
			Operation:   s.Operation,
			Assignments: s.Assignments,
		}
		return res, nil
	}
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

// NewSetvarAction creates a new SetvarAction with the given collection name, operation, and variable assignments
func NewSetvarAction(collection CollectionName, operation string, vars []VarAssignment) (SetvarAction, error) {
	if collection != GLOBAL && collection != IP && collection != RESOURCE && collection != SESSION && collection != TX && collection != USER {
		return SetvarAction{}, fmt.Errorf("invalid setvar action: invalid collection name '%s'", collection)
	}
	return SetvarAction{Collection: collection, Operation: operation, Assignments: vars}, nil
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

// AddSetvarAction adds a setvar action to the NonDisruptiveActions list
func (s *SeclangActions) AddSetvarAction(collection, variable, operation, value string) error {
	colName := stringToCollectionName(strings.ToUpper(collection))
	// Check if there is already a setvar action in the last position
	if len(s.NonDisruptiveActions) > 0 {
		lastAction := s.NonDisruptiveActions[len(s.NonDisruptiveActions)-1]
		if lastAction.GetKey() != SetVar.String() || lastAction.(SetvarAction).Collection != colName || lastAction.(SetvarAction).Operation != operation {
			// If the last action is not setvar, we need to create a new one
			newAction, err := NewSetvarAction(colName, operation, []VarAssignment{{Variable: variable, Value: value}})
			if err != nil {
				return err
			}
			s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
		} else {
			// If the last action is setvar, we need to append the param to it
			sv, ok := lastAction.(SetvarAction)
			if !ok {
				return fmt.Errorf("invalid action type: expected SetvarAction, got %T", lastAction)
			}
			err := sv.AppendAssignment(variable, value)
			if err != nil {
				return err
			}
			s.NonDisruptiveActions[len(s.NonDisruptiveActions)-1] = sv
		}
	} else {
		// If there are no actions yet, we need to create a new setvar action
		newAction, err := NewSetvarAction(colName, operation, []VarAssignment{{Variable: variable, Value: value}})
		if err != nil {
			return err
		}
		s.NonDisruptiveActions = append(s.NonDisruptiveActions, newAction)
	}
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

func (s *SeclangActions) GetActionsByKey(key string) []Action {
	actions := []Action{}
	if s.DisruptiveAction != nil && s.DisruptiveAction.GetKey() == key {
		actions = append(actions, s.DisruptiveAction)
	}
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

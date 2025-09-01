package types

import "fmt"

type SeclangActions struct {
	DisruptiveAction     Action   `yaml:"disruptiveAction,omitempty"`
	NonDisruptiveActions []Action `yaml:"non-disruptiveActions,omitempty"`
	FlowActions          []Action `yaml:"flowActions,omitempty"`
	DataActions          []Action `yaml:"dataActions,omitempty"`
}

func (s *SeclangActions) ToString() string {
	results := []string{}
	if s.DisruptiveAction.Action != "" {
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

type Action struct {
	Action string `yaml:"action"`
	Param  string `yaml:"param,omitempty"`
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

var (
	disruptiveActions = map[string]DisruptiveAction{
		"allow":    Allow,
		"block":    Block,
		"deny":     Deny,
		"drop":     Drop,
		"pass":     Pass,
		"pause":    Pause,
		"proxy":    Proxy,
		"redirect": Redirect,
	}

	flowActions = map[string]FlowAction{
		"chain":     Chain,
		"skip":      Skip,
		"skipAfter": SkipAfter,
	}

	dataActions = map[string]DataAction{
		"data":   Data,
		"status": Status,
		"xmlns":  XLMNS,
	}

	nonDisruptiveActions = map[string]NonDisruptiveAction{
		"append":                 Append,
		"auditlog":               AuditLog,
		"capture":                Capture,
		"ctl":                    Ctl,
		"deprecatevar":           DeprecateVar,
		"exec":                   Exec,
		"expirevar":              ExpireVar,
		"initcol":                InitCol,
		"log":                    Log,
		"logdata":                LogData,
		"multiMatch":             MultiMatch,
		"noauditlog":             NoAuditLog,
		"nolog":                  NoLog,
		"prepend":                Prepend,
		"sanitiseArg":            SanitiseArg,
		"sanitiseMatched":        SanitiseMatched,
		"sanitiseMatchedBytes":   SanitiseMatchedBytes,
		"sanitiseRequestHeader":  SanitiseRequestHeader,
		"sanitiseResponseHeader": SanitiseResponseHeader,
		"setuid":                 SetUid,
		"setrsc":                 SetRsc,
		"setsid":                 SetSid,
		"setenv":                 SetEnv,
		"setvar":                 SetVar,
	}
)

func (a Action) ToString() string {
	if a.Param == "" {
		return a.Action
	} else {
		return a.Action + ":" + a.Param
	}
}

func (a Action) GetKey() string {
	return a.Action
}

func (s *SeclangActions) SetDisruptiveActionWithParam(action, value string) error {
	_, ok := disruptiveActions[action]
	if !ok {
		return fmt.Errorf("Disruptive action %s not found", action)
	}
	s.DisruptiveAction = Action{Action: action, Param: value}
	return nil
}

func (s *SeclangActions) SetDisruptiveActionOnly(action string) error {
	_, ok := disruptiveActions[action]
	if !ok {
		return fmt.Errorf("Disruptive action %s not found", action)
	}
	s.DisruptiveAction = Action{Action: action}
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action, param string) error {
	_, ok := nonDisruptiveActions[action]
	if !ok {
		return fmt.Errorf("Non-disruptive action %s not found", action)
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, Action{Action: action, Param: param})
	return nil
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action string) error {
	_, ok := nonDisruptiveActions[action]
	if !ok {
		return fmt.Errorf("Non-disruptive action %s not found", action)
	}
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, Action{Action: action})
	return nil
}

func (s *SeclangActions) AddFlowActionWithParam(action, param string) error {
	_, ok := flowActions[action]
	if !ok {
		return fmt.Errorf("Flow action %s not found", action)
	}
	s.FlowActions = append(s.FlowActions, Action{Action: action, Param: param})
	return nil
}

func (s *SeclangActions) AddFlowActionOnly(action string) error {
	_, ok := flowActions[action]
	if !ok {
		return fmt.Errorf("Flow action %s not found", action)
	}
	s.FlowActions = append(s.FlowActions, Action{Action: action})
	return nil
}

func (s *SeclangActions) AddDataActionWithParams(action, param string) error {
	_, ok := dataActions[action]
	if !ok {
		return fmt.Errorf("Data action %s not found", action)
	}
	s.DataActions = append(s.DataActions, Action{Action: action, Param: param})
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
	if s.DisruptiveAction.Action == key {
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

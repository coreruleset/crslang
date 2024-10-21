package types

import "fmt"

type SeclangActions struct {
	DisruptiveAction     SeclangAction   `yaml:"disruptiveAction,omitempty"`
	NonDisruptiveActions []SeclangAction `yaml:"non-disruptiveActions,omitempty"`
	FlowActions          []SeclangAction `yaml:"flowActions,omitempty"`
	DataActions          []SeclangAction `yaml:"dataActions,omitempty"`
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

type SeclangAction interface {
	SetAction(action, param string)
	ToString() string
	GetKey() string
}

type ActionOnly struct {
	Action string
}

func (a *ActionOnly) SetAction(action, param string) {
	a.Action = action
}

func (a *ActionOnly) ToString() string {
	return a.Action
}

func (a *ActionOnly) GetKey() string {
	return a.Action
}

type ActionWithParam struct {
	Action string
	Param  string
}

func (a *ActionWithParam) SetAction(action, param string) {
	a.Action = action
	a.Param = param
}

func (a *ActionWithParam) ToString() string {
	return a.Action + ":" + a.Param
}

func (a *ActionWithParam) GetKey() string {
	return a.Action
}

func (s *SeclangActions) SetDisruptiveActionWithParam(action, value string) {
	s.DisruptiveAction = &ActionWithParam{Action: action, Param: value}
}

func (s *SeclangActions) SetDisruptiveActionOnly(action string) {
	s.DisruptiveAction = &ActionOnly{Action: action}
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action, param string) {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, &ActionWithParam{Action: action, Param: param})
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action string) {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, &ActionOnly{Action: action})
}

func (s *SeclangActions) AddFlowActionWithParam(action, param string) {
	s.FlowActions = append(s.FlowActions, &ActionWithParam{Action: action, Param: param})
}

func (s *SeclangActions) AddFlowActionOnly(action string) {
	s.FlowActions = append(s.FlowActions, &ActionOnly{Action: action})
}

func (s *SeclangActions) AddDataActionWithParams(action, param string) {
	s.DataActions = append(s.DataActions, &ActionWithParam{Action: action, Param: param})
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

func (s *SeclangActions) GetActionByKey(key string) SeclangAction {
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
	return nil
}

func (s *SeclangActions) GetActionsByKey(key string) []SeclangAction {
	actions := []SeclangAction{}
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
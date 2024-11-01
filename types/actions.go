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

func (a Action) SetAction(action, param string) {
	a.Action = action
	a.Param = param
}

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

func (s *SeclangActions) SetDisruptiveActionWithParam(action, value string) {
	s.DisruptiveAction = Action{Action: action, Param: value}
}

func (s *SeclangActions) SetDisruptiveActionOnly(action string) {
	s.DisruptiveAction = Action{Action: action}
}

func (s *SeclangActions) AddNonDisruptiveActionWithParam(action, param string) {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, Action{Action: action, Param: param})
}

func (s *SeclangActions) AddNonDisruptiveActionOnly(action string) {
	s.NonDisruptiveActions = append(s.NonDisruptiveActions, Action{Action: action})
}

func (s *SeclangActions) AddFlowActionWithParam(action, param string) {
	s.FlowActions = append(s.FlowActions, Action{Action: action, Param: param})
}

func (s *SeclangActions) AddFlowActionOnly(action string) {
	s.FlowActions = append(s.FlowActions, Action{Action: action})
}

func (s *SeclangActions) AddDataActionWithParams(action, param string) {
	s.DataActions = append(s.DataActions, Action{Action: action, Param: param})
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
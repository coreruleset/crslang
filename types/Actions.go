package types

import "fmt"

type SeclangActions struct {
	DisruptiveAction     SeclangAction   `yaml:"disruptiveAction,omitempty"`
	NonDisruptiveActions []SeclangAction `yaml:"non-disruptiveActions,omitempty"`
	FlowActions          []SeclangAction `yaml:"flowActions,omitempty"`
	DataActions          []SeclangAction `yaml:"dataActions,omitempty"`
}

func (s *SeclangActions) String() string {
	return fmt.Sprintf("Disruptive: %v, NonDisruptive: %v, Flow: %v, Data: %v", s.DisruptiveAction, s.NonDisruptiveActions, s.FlowActions, s.DataActions)
}

type SeclangAction interface {
	SetAction(action, param string)
}

type ActionOnly struct {
	Action string
}

func (a *ActionOnly) SetAction(action, param string) {
	a.Action = action
}

type ActionWithParam struct {
	Action string
	Param  string
}

func (a *ActionWithParam) SetAction(action, param string) {
	a.Action = action
	a.Param = param
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
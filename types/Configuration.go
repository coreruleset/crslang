package types

type SeclangDirective interface {
	SetId(value string)
	SetPhase(value string)
	SetMsg(value string)
	SetMaturity(value string)
	SetRev(value string)
	SetSeverity(value string)
	SetVer(value string)
	SetDisruptiveActionWithParam(action, value string)
	SetDisruptiveActionOnly(action string)
	AddNonDisruptiveActionWithParam(action, param string)
	AddNonDisruptiveActionOnly(action string)
	AddFlowActionWithParam(action, param string)
	AddFlowActionOnly(action string)
	AddDataActionWithParams(action, param string)
	AddTransformation(transformation string)
	AddVariable(variable string)
	SetOperatorName(name string)
	SetOperatorValue(value string)
}

type Configuration struct {
	ConfigDirectives map[string]string  `yaml:"configDirectives,omitempty"`
	DefaultActions   []SecDefaultAction `yaml:"defaultActions,omitempty"`
	SecActions       []SecAction        `yaml:"secActions,omitempty"`
	SecRules         []SecRule          `yaml:"secRules,omitempty"`
}
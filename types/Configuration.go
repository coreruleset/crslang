package types

type Configuration struct {
	ConfigDirectives map[string]string  `yaml:"configDirectives,omitempty"`
	DefaultActions   []SecDefaultAction `yaml:"defaultActions,omitempty"`
	SecActions       []SecAction        `yaml:"secActions,omitempty"`
	SecRules         []SecRule          `yaml:"secRules,omitempty"`
}
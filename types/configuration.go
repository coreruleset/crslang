package types

// type Configuration struct {
// 	Comments         []CommentMetadata  `yaml:"comments,omitempty"`
// 	ConfigDirectives map[string]string  `yaml:"configDirectives,omitempty"`
// 	DefaultActions   []SecDefaultAction `yaml:"defaultActions,omitempty"`
// 	SecActions       []SecAction        `yaml:"secActions,omitempty"`
// 	SecRules         []SecRule          `yaml:"secRules,omitempty"`
// }

type ConfigurationList struct {
	Configurations []Configuration `yaml:"configurations,omitempty"`
}

type Configuration struct {
	Marker     ConfigurationDirective `yaml:"marker,omitempty"`
	Directives []SeclangDirective     `yaml:"directives,omitempty"`
}
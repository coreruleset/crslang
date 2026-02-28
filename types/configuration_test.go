package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	extractDefValues = []struct {
		name     string
		input    ConfigurationList
		expected ConfigurationList
	}{
		{
			name: "Extract version and tags from only one rule",
			input: ConfigurationList{
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
						},
					},
				},
			},
			expected: ConfigurationList{
				Global: DefaultConfigs{
					Version: "OWASP_CRS/4.0.0",
					Tags:    []string{"application-multi2", "attack-protocol3"},
				},
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Tags: []string{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Extract version and tags from multiple rules without common values",
			input: ConfigurationList{
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Ver:  "OWASP_CRS/4.0.1",
									Tags: []string{"test-tag"},
								},
							},
						},
					},
				},
			},
			expected: ConfigurationList{
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Ver:  "OWASP_CRS/4.0.1",
									Tags: []string{"test-tag"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Extract version and tags from multiple rules in different groups",
			input: ConfigurationList{
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
						},
					},
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
						},
					},
				},
			},
			expected: ConfigurationList{
				Global: DefaultConfigs{
					Version: "OWASP_CRS/4.0.0",
					Tags:    []string{"application-multi2", "attack-protocol3"},
				},
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Tags: []string{},
								},
							},
						},
					},
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Tags: []string{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Extract tags from multiple rules in different groups",
			input: ConfigurationList{
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
						},
					},
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Ver:  "OWASP_CRS/4.0.1",
									Tags: []string{"application-multi2", "attack-protocol3"},
								},
							},
						},
					},
				},
			},
			expected: ConfigurationList{
				Global: DefaultConfigs{
					Tags: []string{"application-multi2", "attack-protocol3"},
				},
				DirectiveList: []DirectiveList{
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   1,
									Ver:  "OWASP_CRS/4.0.0",
									Tags: []string{},
								},
							},
						},
					},
					{
						Directives: []SeclangDirective{
							&RuleWithCondition{
								Kind: RuleKind,
								Metadata: SecRuleMetadata{
									Id:   2,
									Ver:  "OWASP_CRS/4.0.1",
									Tags: []string{},
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestExtractDefaultValues(t *testing.T) {
	for _, tt := range extractDefValues {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.ExtractDefaultValues()
			require.Equalf(t, tt.expected, tt.input, tt.name)
		})
	}
}

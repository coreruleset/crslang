package translator

import (
	"os"
	"testing"

	"github.com/coreruleset/crslang/types"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

var testFilepath = "../testdata/crs/"

func TestLoadCRS(t *testing.T) {

	configList, err := LoadSeclang(testFilepath)
	if err != nil {
		t.Errorf("Error loading seclang directives: %v", err)
	}

	configList = *ToCRSLang(configList)

	yamlFile, err := yaml.Marshal(configList)
	if err != nil {
		t.Errorf("Error marshalling yaml: %v", err)
	}

	err = writeToFile(yamlFile, "tmp_crslang.yaml")

	defer os.Remove("tmp_crslang.yaml")

	loadedConfigList := types.LoadDirectivesWithConditionsFromFile("tmp_crslang.yaml")
	yamlLoadedFile, err := yaml.Marshal(loadedConfigList)
	if err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	if string(yamlFile) != string(yamlLoadedFile) {
		t.Errorf("Error: loaded file is different from original. Expected string length: %v, got: %v", len(string(yamlFile)), len(string(yamlLoadedFile)))
	}
}

func TestLoadSeclangFromString(t *testing.T) {
	const seclangInput = `SecRule TX:DETECTION_PARANOIA_LEVEL "@lt 1" "id:911011,phase:1,pass,nolog,tag:'OWASP_CRS',ver:'OWASP_CRS/4.0.1-dev',skipAfter:END-REQUEST-911-METHOD-ENFORCEMENT"`

	configList, err := LoadSeclangFromString(seclangInput, "test")
	if err != nil {
		t.Fatalf("Error loading seclang from string: %v", err)
	}

	if len(configList.Groups) == 0 {
		t.Fatal("Expected at least one directive list, got none")
	}

	configList = *ToCRSLang(configList)

	yamlFile, err := yaml.Marshal(configList)
	if err != nil {
		t.Fatalf("Error marshalling yaml: %v", err)
	}

	if len(yamlFile) == 0 {
		t.Error("Expected non-empty YAML output")
	}
}

func TestFromCRSLangToSeclang(t *testing.T) {
	configList, err := LoadSeclang(testFilepath)
	if err != nil {
		t.Errorf("Error loading seclang directives: %v", err)
	}

	seclangDirectives := types.ToSeclang(configList)

	crslangConfigList := ToCRSLang(configList)
	unformattedDirectives := types.FromCRSLangToUnformattedDirectives(*crslangConfigList)
	seclangDirectivesFromConditions := types.ToSeclang(*unformattedDirectives)

	if len(seclangDirectives) != len(seclangDirectivesFromConditions) {
		t.Errorf("Error in CRSLang to Seclang directives conversion. Expected length: %v, got: %v", len(seclangDirectives), len(seclangDirectivesFromConditions))
	}

}

func TestWriteAndLoadRuleSeparately(t *testing.T) {
	testCases := []struct {
		name     string
		input    types.Ruleset
		expected types.Ruleset
	}{
		{
			name: "Simple ruleset with comment, config and rule",
			input: types.Ruleset{
				Global: types.DefaultConfigs{
					Version: "4.0.0",
					Tags:    []string{"OWASP_CRS"},
				},
				Groups: []types.Group{
					{
						Id:   "test-group-1",
						Tags: []string{"tag1", "tag2"},
						Directives: []types.SeclangDirective{
							types.CommentDirective{
								Kind: types.CommentKind,
								Metadata: types.CommentMetadata{
									Comment: "Test comment",
								},
							},
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Name:      types.SecRuleEngine,
								Parameter: "On",
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "1",
									},
									Id:  12345,
									Msg: "Test rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.REQUEST_URI,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "/test",
										},
									},
								},
							},
						},
					},
				},
			},
			expected: types.Ruleset{
				Global: types.DefaultConfigs{
					Version: "4.0.0",
					Tags:    []string{"OWASP_CRS"},
				},
				Groups: []types.Group{
					{
						Id:   "test-group-1",
						Tags: []string{"tag1", "tag2"},
						Directives: []types.SeclangDirective{
							types.CommentDirective{
								Kind: types.CommentKind,
								Metadata: types.CommentMetadata{
									Comment: "Test comment",
								},
							},
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Name:      types.SecRuleEngine,
								Parameter: "On",
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "1",
									},
									Id:  12345,
									Msg: "Test rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.REQUEST_URI,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "/test",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple groups with different directive types",
			input: types.Ruleset{
				Global: types.DefaultConfigs{
					Version: "4.0.0",
					Tags:    []string{"OWASP_CRS"},
				},
				Groups: []types.Group{
					{
						Id:   "group-a",
						Tags: []string{"security", "testing"},
						Directives: []types.SeclangDirective{
							types.CommentDirective{
								Kind: types.CommentKind,
								Metadata: types.CommentMetadata{
									Comment: "This is a test group",
								},
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "1",
									},
									Id:  100,
									Msg: "First rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.REQUEST_URI,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "/admin",
										},
									},
								},
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "2",
									},
									Id:  200,
									Msg: "Second rule",
								},
								Conditions: []types.Condition{
									{
										Collections: []types.Collection{
											{
												Name: types.REQUEST_HEADERS,
											},
										},
										Operator: types.Operator{
											Name:  types.Contains,
											Value: "test",
										},
									},
								},
							},
						},
					},
					{
						Id: "group-b",
						Directives: []types.SeclangDirective{
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Name:      types.SecRuleEngine,
								Parameter: "DetectionOnly",
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "3",
									},
									Id:  300,
									Msg: "Third rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.RESPONSE_BODY,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "sensitive",
										},
									},
								},
							},
						},
					},
				},
			},
			expected: types.Ruleset{
				Global: types.DefaultConfigs{
					Version: "4.0.0",
					Tags:    []string{"OWASP_CRS"},
				},
				Groups: []types.Group{
					{
						Id:   "group-a",
						Tags: []string{"security", "testing"},
						Directives: []types.SeclangDirective{
							types.CommentDirective{
								Kind: types.CommentKind,
								Metadata: types.CommentMetadata{
									Comment: "This is a test group",
								},
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "1",
									},
									Id:  100,
									Msg: "First rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.REQUEST_URI,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "/admin",
										},
									},
								},
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "2",
									},
									Id:  200,
									Msg: "Second rule",
								},
								Conditions: []types.Condition{
									{
										Collections: []types.Collection{
											{
												Name: types.REQUEST_HEADERS,
											},
										},
										Operator: types.Operator{
											Name:  types.Contains,
											Value: "test",
										},
									},
								},
							},
						},
					},
					{
						Id: "group-b",
						Directives: []types.SeclangDirective{
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Name:      types.SecRuleEngine,
								Parameter: "DetectionOnly",
							},
							&types.RuleWithCondition{
								Kind: types.RuleKind,
								Metadata: types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										Phase: "3",
									},
									Id:  300,
									Msg: "Third rule",
								},
								Conditions: []types.Condition{
									{
										Variables: []types.Variable{
											{
												Name: types.RESPONSE_BODY,
											},
										},
										Operator: types.Operator{
											Name:  types.Rx,
											Value: "sensitive",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary directory for output
			tmpDir := t.TempDir()

			// Write ruleset separately
			err := WriteRuleSeparately(tc.input, tmpDir)
			if err != nil {
				t.Fatalf("WriteRuleSeparately failed: %v", err)
			}

			// Load from directory
			loadedRuleset, err := LoadRulesFromDirectory(tmpDir)
			if err != nil {
				t.Fatalf("LoadRulesFromDirectory failed: %v", err)
			}

			require.Equal(t, tc.expected, loadedRuleset, tc.name)
		})
	}
}

package main

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/coreruleset/crslang/listener"
	"github.com/coreruleset/crslang/types"
	"github.com/coreruleset/seclang_parser/parser"
	"github.com/stretchr/testify/require"
)

// Helper functions to create actions for tests, panicking on error
func mustNewActionOnly[T types.ActionType](action T) types.Action {
	newAction, err := types.NewActionOnly(action)
	if err != nil {
		panic(err)
	}
	return newAction
}

func mustNewActionWithParam[T types.ActionType](action T, param string) types.Action {
	newAction, err := types.NewActionWithParam(action, param)
	if err != nil {
		panic(err)
	}
	return newAction
}

func mustNewActionMultipleParam[T types.ActionType](action T, params []string) types.Action {
	newAction, err := types.NewActionMultipleParam(action, params)
	if err != nil {
		panic(err)
	}
	return newAction
}

type testCase struct {
	name     string
	payload  string
	expected types.ConfigurationList
}

var (
	listenerTestCases = []testCase{
		{
			name: "LoadComment",
			payload: `#
# -- [[ Introduction ]] --------------------------------------------------------
#
# The OWASP ModSecurity Core Rule Set (CRS) is a set of generic attack
# detection rules that provide a base level of protection for any web
# application. They are written for the open source, cross-platform
# ModSecurity Web Application Firewall.
#
# See also:
# https://coreruleset.org/
# https://github.com/coreruleset/coreruleset
# https://owasp.org/www-project-modsecurity-core-rule-set/
#
`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							types.CommentMetadata{
								Comment: `#
# -- [[ Introduction ]] --------------------------------------------------------
#
# The OWASP ModSecurity Core Rule Set (CRS) is a set of generic attack
# detection rules that provide a base level of protection for any web
# application. They are written for the open source, cross-platform
# ModSecurity Web Application Firewall.
#
# See also:
# https://coreruleset.org/
# https://github.com/coreruleset/coreruleset
# https://owasp.org/www-project-modsecurity-core-rule-set/
#
`,
							},
						},
					},
				},
			},
		},
		{
			name: "LoadConfigurationDirective",
			payload: `
SecRuleEngine On
`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Metadata:  &types.CommentMetadata{},
								Name:      "SecRuleEngine",
								Parameter: "On",
							},
						},
					},
				},
			},
		},
		{
			name: "LoadSecAction",
			payload: `
# Initialize anomaly scoring variables.
# All _score variables start at 0, and are incremented by the various rules
# upon detection of a possible attack.

SecAction \
    "id:901200,\
    phase:1,\
    pass,\
    t:none,\
    nolog,\
    tag:'OWASP_CRS',\
    ver:'OWASP_CRS/4.0.1-dev',\
    setvar:'tx.blocking_inbound_anomaly_score=0',\
    setvar:'tx.detection_inbound_anomaly_score=0',\
    setvar:'tx.inbound_anomaly_score_pl1=0',\
    setvar:'tx.inbound_anomaly_score_pl2=0',\
    setvar:'tx.inbound_anomaly_score_pl3=0',\
    setvar:'tx.inbound_anomaly_score_pl4=0',\
    setvar:'tx.sql_injection_score=0',\
    setvar:'tx.xss_score=0',\
    setvar:'tx.rfi_score=0',\
    setvar:'tx.lfi_score=0',\
    setvar:'tx.rce_score=0',\
    setvar:'tx.php_injection_score=0',\
    setvar:'tx.http_violation_score=0',\
    setvar:'tx.session_fixation_score=0',\
    setvar:'tx.blocking_outbound_anomaly_score=0',\
    setvar:'tx.detection_outbound_anomaly_score=0',\
    setvar:'tx.outbound_anomaly_score_pl1=0',\
    setvar:'tx.outbound_anomaly_score_pl2=0',\
    setvar:'tx.outbound_anomaly_score_pl3=0',\
    setvar:'tx.outbound_anomaly_score_pl4=0',\
    setvar:'tx.anomaly_score=0'"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecAction{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{
											Comment: `# Initialize anomaly scoring variables.
# All _score variables start at 0, and are incremented by the various rules
# upon detection of a possible attack.

`,
										},
										Phase: "1",
									},
									Id:   901200,
									Tags: []string{"OWASP_CRS"},
									Ver:  "OWASP_CRS/4.0.1-dev",
								},
								Transformations: types.Transformations{
									Transformations: []types.Transformation{
										types.None,
									},
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
									NonDisruptiveActions: []types.Action{
										mustNewActionOnly(types.NoLog),
										mustNewActionMultipleParam(types.SetVar, []string{
											"tx.blocking_inbound_anomaly_score=0",
											"tx.detection_inbound_anomaly_score=0",
											"tx.inbound_anomaly_score_pl1=0",
											"tx.inbound_anomaly_score_pl2=0",
											"tx.inbound_anomaly_score_pl3=0",
											"tx.inbound_anomaly_score_pl4=0",
											"tx.sql_injection_score=0",
											"tx.xss_score=0",
											"tx.rfi_score=0",
											"tx.lfi_score=0",
											"tx.rce_score=0",
											"tx.php_injection_score=0",
											"tx.http_violation_score=0",
											"tx.session_fixation_score=0",
											"tx.blocking_outbound_anomaly_score=0",
											"tx.detection_outbound_anomaly_score=0",
											"tx.outbound_anomaly_score_pl1=0",
											"tx.outbound_anomaly_score_pl2=0",
											"tx.outbound_anomaly_score_pl3=0",
											"tx.outbound_anomaly_score_pl4=0",
											"tx.anomaly_score=0"},
										),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LoadSecRule",
			payload: `
#
# Validate request line against the format specified in the HTTP RFC
#
# -=[ Rule Logic ]=-
#
# Uses rule negation against the regex for positive security.   The regex specifies the proper
# construction of URI request lines such as:
#
#   "http" "://" authority path-abempty [ "?" query ]
#
# It also outlines proper construction for CONNECT, OPTIONS and GET requests.
#
# Regular expression generated from regex-assembly/920100.ra.
# To update the regular expression run the following shell script
# (consult https://coreruleset.org/docs/development/regex_assembly/ for details):
#   crs-toolchain regex update 920100
#
# -=[ References ]=-
# https://www.rfc-editor.org/rfc/rfc9110.html#section-4.2.1
# http://capec.mitre.org/data/definitions/272.html
#
SecRule REQUEST_LINE "@rx (?i)^(?:get /[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?|(?:connect (?:(?:[0-9]{1,3}\.){3}[0-9]{1,3}\.?(?::[0-9]+)?|[\--9A-Z_a-z]+:[0-9]+)|options \*|[a-z]{3,10}[\s\v]+(?:[0-9A-Z_a-z]{3,7}?://[\--9A-Z_a-z]*(?::[0-9]+)?)?/[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?)[\s\v]+[\.-9A-Z_a-z]+)$" \
    "id:920100,\
    phase:1,\
    block,\
    t:none,\
    msg:'Invalid HTTP Request Line',\
    logdata:'%{request_line}',\
    tag:'application-multi',\
    tag:'language-multi',\
    tag:'platform-multi',\
    tag:'attack-protocol',\
    tag:'paranoia-level/1',\
    tag:'OWASP_CRS',\
    tag:'capec/1000/210/272',\
    ver:'OWASP_CRS/4.0.1-dev',\
    severity:'WARNING',\
    setvar:'tx.inbound_anomaly_score_pl1=+%{tx.warning_anomaly_score}'"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{
											Comment: `#
# Validate request line against the format specified in the HTTP RFC
#
# -=[ Rule Logic ]=-
#
# Uses rule negation against the regex for positive security.   The regex specifies the proper
# construction of URI request lines such as:
#
#   "http" "://" authority path-abempty [ "?" query ]
#
# It also outlines proper construction for CONNECT, OPTIONS and GET requests.
#
# Regular expression generated from regex-assembly/920100.ra.
# To update the regular expression run the following shell script
# (consult https://coreruleset.org/docs/development/regex_assembly/ for details):
#   crs-toolchain regex update 920100
#
# -=[ References ]=-
# https://www.rfc-editor.org/rfc/rfc9110.html#section-4.2.1
# http://capec.mitre.org/data/definitions/272.html
#
`,
										},
										Phase: "1",
									},
									Id:       920100,
									Tags:     []string{"application-multi", "language-multi", "platform-multi", "attack-protocol", "paranoia-level/1", "OWASP_CRS", "capec/1000/210/272"},
									Ver:      "OWASP_CRS/4.0.1-dev",
									Msg:      "Invalid HTTP Request Line",
									Severity: "WARNING",
								},
								Variables: []types.Variable{
									{Name: types.REQUEST_LINE},
								},
								Operator: types.Operator{
									Name:  types.Rx,
									Value: `(?i)^(?:get /[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?|(?:connect (?:(?:[0-9]{1,3}\.){3}[0-9]{1,3}\.?(?::[0-9]+)?|[\--9A-Z_a-z]+:[0-9]+)|options \*|[a-z]{3,10}[\s\v]+(?:[0-9A-Z_a-z]{3,7}?://[\--9A-Z_a-z]*(?::[0-9]+)?)?/[^#\?]*(?:\?[^\s\v#]*)?(?:#[^\s\v]*)?)[\s\v]+[\.-9A-Z_a-z]+)$`,
								},
								Transformations: types.Transformations{
									Transformations: []types.Transformation{
										types.None,
									},
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Block),
									NonDisruptiveActions: []types.Action{
										mustNewActionWithParam(types.LogData, "%{request_line}"),
										mustNewActionMultipleParam(types.SetVar, []string{"tx.inbound_anomaly_score_pl1=+%{tx.warning_anomaly_score}"}),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LoadSecRuleWithCollection",
			payload: `
#
# -=[ Exclusion rule for 942440 ]=-
#
# Prevent FPs against Facebook click identifier
#
SecRule ARGS_GET:fbclid "@rx [a-zA-Z0-9_-]{61,61}" \
    "id:942441,\
    phase:2,\
    pass,\
    t:none,t:urlDecodeUni,\
    nolog,\
    tag:'OWASP_CRS',\
    ctl:ruleRemoveTargetById=942440;ARGS:fbclid,\
    ver:'OWASP_CRS/4.0.1-dev'"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{
											Comment: `#
# -=[ Exclusion rule for 942440 ]=-
#
# Prevent FPs against Facebook click identifier
#
`,
										},
										Phase: "2",
									},
									Id:   942441,
									Tags: []string{"OWASP_CRS"},
									Ver:  "OWASP_CRS/4.0.1-dev",
								},
								Collections: []types.Collection{
									{
										Name:      types.ARGS_GET,
										Arguments: []string{"fbclid"},
										Excluded:  []string{},
									},
								},
								Operator: types.Operator{
									Name:  types.Rx,
									Value: "[a-zA-Z0-9_-]{61,61}",
								},
								Transformations: types.Transformations{
									Transformations: []types.Transformation{
										types.None,
										types.UrlDecodeUni,
									},
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
									NonDisruptiveActions: []types.Action{
										mustNewActionOnly(types.NoLog),
										mustNewActionWithParam(types.Ctl, "ruleRemoveTargetById=942440;ARGS:fbclid"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Load SecRule with differents exclusions",
			payload: `
SecRule ARGS_GET:fbclid|!ARGS_GET|ARGS_GET:fbclid|!ARGS_GET:fbclid|ARGS_GET:test "@rx test" \
    "id:942441,\
    phase:2,\
    pass"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{},
										Phase:           "2",
									},
									Id: 942441,
								},
								Collections: []types.Collection{
									{
										Name:      types.ARGS_GET,
										Arguments: []string{"test"},
										Excluded:  []string{},
									},
								},
								Operator: types.Operator{
									Name:  types.Rx,
									Value: "test",
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Load SecRule with collections and excluded values",
			payload: `
SecRule ARGS:/^id_/|!ARGS:id_1 "@rx test" \
    "id:942441,\
    phase:2,\
    pass"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{},
										Phase:           "2",
									},
									Id: 942441,
								},
								Collections: []types.Collection{
									{
										Name:      types.ARGS,
										Arguments: []string{"/^id_/"},
										Excluded:  []string{"id_1"},
									},
								},
								Operator: types.Operator{
									Name:  types.Rx,
									Value: "test",
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LoadChain",
			payload: `
# This file is used as an exception mechanism to remove common false positives
# that may be encountered.
#
# Exception for Apache SSL pinger
#
SecRule REQUEST_LINE "@streq GET /" \
    "id:905100,\
    phase:1,\
    pass,\
    t:none,\
    nolog,\
    tag:'application-multi',\
    tag:'language-multi',\
    tag:'platform-apache',\
    tag:'attack-generic',\
    tag:'OWASP_CRS',\
    ver:'OWASP_CRS/4.6.0-dev',\
    chain"
    SecRule REMOTE_ADDR "@ipMatch 127.0.0.1,::1" \
        "t:none,\
        ctl:ruleRemoveByTag=OWASP_CRS,\
        ctl:auditEngine=Off"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{
											Comment: `# This file is used as an exception mechanism to remove common false positives
# that may be encountered.
#
# Exception for Apache SSL pinger
#
`,
										},
										Phase: "1",
									},
									Id:   905100,
									Tags: []string{"application-multi", "language-multi", "platform-apache", "attack-generic", "OWASP_CRS"},
									Ver:  "OWASP_CRS/4.6.0-dev",
								},
								Variables: []types.Variable{
									{Name: types.REQUEST_LINE},
								},
								Operator: types.Operator{
									Name:  types.StrEq,
									Value: "GET /",
								},
								Transformations: types.Transformations{
									Transformations: []types.Transformation{
										types.None,
									},
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
									NonDisruptiveActions: []types.Action{
										mustNewActionOnly(types.NoLog),
									},
									FlowActions: []types.Action{
										mustNewActionOnly(types.Chain),
									},
								},
								ChainedRule: &types.SecRule{
									Metadata: &types.SecRuleMetadata{},
									Variables: []types.Variable{
										{Name: types.REMOTE_ADDR},
									},
									Operator: types.Operator{
										Name:  types.IpMatch,
										Value: "127.0.0.1,::1",
									},
									Transformations: types.Transformations{
										Transformations: []types.Transformation{
											types.None,
										},
									},
									Actions: &types.SeclangActions{
										NonDisruptiveActions: []types.Action{
											mustNewActionWithParam(types.Ctl, "ruleRemoveByTag=OWASP_CRS"),
											mustNewActionWithParam(types.Ctl, "auditEngine=Off"),
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
			name: "LoadRemoveRules",
			payload: `
SecRuleRemoveByID 1 2 9000-9010

SecRuleRemoveByTag "attack-sqli"

SecRuleRemoveByMsg FAIL 
`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							types.RemoveRuleDirective{
								Kind: types.Remove,
								Ids:  []int{1, 2},
								IdRanges: []types.IdRange{
									{
										Start: 9000,
										End:   9010,
									},
								},
							},
							types.RemoveRuleDirective{
								Kind: types.Remove,
								Tags: []string{"attack-sqli"},
							},
							types.RemoveRuleDirective{
								Kind: types.Remove,
								Msgs: []string{"FAIL"},
							},
						},
					},
				},
			},
		},
		{
			name: "Load SecRule with negated operator",
			payload: `
# Force body variable
SecRule REQBODY_PROCESSOR "!@rx (?:URLENCODED|MULTIPART|XML|JSON)" \
    "id:901340,\
    phase:1,\
    pass,\
    nolog,\
    noauditlog,\
    msg:'Enabling body inspection',\
    ctl:forceRequestBodyVariable=On,\
    ver:'OWASP_CRS/4.0.0-rc1'"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							&types.SecRule{
								Metadata: &types.SecRuleMetadata{
									OnlyPhaseMetadata: types.OnlyPhaseMetadata{
										CommentMetadata: types.CommentMetadata{
											Comment: `# Force body variable
`,
										},
										Phase: "1",
									},
									Id:  901340,
									Msg: "Enabling body inspection",
									Ver: "OWASP_CRS/4.0.0-rc1",
								},
								Variables: []types.Variable{
									{Name: types.REQBODY_PROCESSOR},
								},
								Operator: types.Operator{
									Negate: true,
									Name:   types.Rx,
									Value:  "(?:URLENCODED|MULTIPART|XML|JSON)",
								},
								Actions: &types.SeclangActions{
									DisruptiveAction: mustNewActionOnly(types.Pass),
									NonDisruptiveActions: []types.Action{
										mustNewActionOnly(types.NoLog),
										mustNewActionOnly(types.NoAuditLog),
										mustNewActionWithParam(types.Ctl, "forceRequestBodyVariable=On"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Load directive lists with markers",
			payload: `
SecCollectionTimeout 600

SecMarker "END-TEST"

SecComponentSignature "OWASP_CRS/4.0.1-dev"`,
			expected: types.ConfigurationList{
				DirectiveList: []types.DirectiveList{
					{
						Directives: []types.SeclangDirective{
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Metadata:  &types.CommentMetadata{},
								Name:      "SecCollectionTimeout",
								Parameter: "600",
							},
						},
						Marker: types.ConfigurationDirective{
							Kind:      types.ConfigurationKind,
							Metadata:  &types.CommentMetadata{},
							Name:      "SecMarker",
							Parameter: "END-TEST",
						},
					},
					{
						Directives: []types.SeclangDirective{
							types.ConfigurationDirective{
								Kind:      types.ConfigurationKind,
								Metadata:  &types.CommentMetadata{},
								Name:      "SecComponentSignature",
								Parameter: "OWASP_CRS/4.0.1-dev",
							},
						},
					},
				},
			},
		},
	}
)

func TestLoadSecLang(t *testing.T) {
	for _, test := range listenerTestCases {
		got := types.ConfigurationList{}
		input := antlr.NewInputStream(test.payload)
		lexer := parser.NewSecLangLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		p := parser.NewSecLangParser(stream)
		start := p.Configuration()
		var seclangListener listener.ExtendedSeclangParserListener
		antlr.ParseTreeWalkerDefault.Walk(&seclangListener, start)
		got = seclangListener.ConfigurationList

		require.Equalf(t, test.expected, got, test.name)
	}
}

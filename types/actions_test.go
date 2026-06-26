package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

// Helper functions to create actions for tests, panicking on error
func mustNewActionOnly[T ActionType](action T) Action {
	newAction, err := NewActionOnly(action)
	if err != nil {
		panic(err)
	}
	return newAction
}

func mustNewActionWithParam[T ActionType](action T, param string) Action {
	newAction, err := NewActionWithParam(action, param)
	if err != nil {
		panic(err)
	}
	return newAction
}

func mustNewSetvarAction(collection CollectionName, operation VarOperation, vars []VarAssignment) Action {
	newAction, err := NewSetvarAction(collection, operation, vars)
	if err != nil {
		panic(err)
	}
	return newAction
}

var (
	unmarshalTests = []struct {
		input    string
		expected SeclangActions
	}{
		{
			input: `
disruptive: block
non-disruptive:
    - capture
    - logdata: 'Matched Data: %{TX.0} found within %{MATCHED_VAR_NAME}: %{MATCHED_VAR}'
    - setvar: 
        - rfi_parameter_%{MATCHED_VAR_NAME}: .%{tx.1}
flow:
    - chain`,
			expected: SeclangActions{
				DisruptiveAction: mustNewActionOnly(Block),
				NonDisruptiveActions: []Action{
					mustNewActionOnly(Capture),
					mustNewActionWithParam(LogData, "Matched Data: %{TX.0} found within %{MATCHED_VAR_NAME}: %{MATCHED_VAR}"),
					mustNewSetvarAction(TX, Assign, []VarAssignment{{Variable: "rfi_parameter_%{MATCHED_VAR_NAME}", Value: ".%{tx.1}"}}),
				},
				FlowActions: []Action{
					mustNewActionOnly(Chain),
				},
			},
		},
		{
			input: `
disruptive: pass
non-disruptive:
    - nolog
    - setvar:
        - blocking_inbound_anomaly_score: "0"
        - detection_inbound_anomaly_score: "0"
        - inbound_anomaly_score_pl1: "0"
        - inbound_anomaly_score_pl2: "0"
        - inbound_anomaly_score_pl3: "0"
        - inbound_anomaly_score_pl4: "0"
        - sql_injection_score: "0"
        - xss_score: "0"
        - rfi_score: "0"
        - lfi_score: "0"
        - rce_score: "0"
        - php_injection_score: "0"
        - http_violation_score: "0"
        - session_fixation_score: "0"
        - blocking_outbound_anomaly_score: "0"
        - detection_outbound_anomaly_score: "0"
        - outbound_anomaly_score_pl1: "0"
        - outbound_anomaly_score_pl2: "0"
        - outbound_anomaly_score_pl3: "0"
        - outbound_anomaly_score_pl4: "0"
        - anomaly_score: "0"`,
			expected: SeclangActions{
				DisruptiveAction: mustNewActionOnly(Pass),
				NonDisruptiveActions: []Action{
					mustNewActionOnly(NoLog),
					mustNewSetvarAction(TX, Assign, []VarAssignment{
						{Variable: "blocking_inbound_anomaly_score", Value: "0"},
						{Variable: "detection_inbound_anomaly_score", Value: "0"},
						{Variable: "inbound_anomaly_score_pl1", Value: "0"},
						{Variable: "inbound_anomaly_score_pl2", Value: "0"},
						{Variable: "inbound_anomaly_score_pl3", Value: "0"},
						{Variable: "inbound_anomaly_score_pl4", Value: "0"},
						{Variable: "sql_injection_score", Value: "0"},
						{Variable: "xss_score", Value: "0"},
						{Variable: "rfi_score", Value: "0"},
						{Variable: "lfi_score", Value: "0"},
						{Variable: "rce_score", Value: "0"},
						{Variable: "php_injection_score", Value: "0"},
						{Variable: "http_violation_score", Value: "0"},
						{Variable: "session_fixation_score", Value: "0"},
						{Variable: "blocking_outbound_anomaly_score", Value: "0"},
						{Variable: "detection_outbound_anomaly_score", Value: "0"},
						{Variable: "outbound_anomaly_score_pl1", Value: "0"},
						{Variable: "outbound_anomaly_score_pl2", Value: "0"},
						{Variable: "outbound_anomaly_score_pl3", Value: "0"},
						{Variable: "outbound_anomaly_score_pl4", Value: "0"},
						{Variable: "anomaly_score", Value: "0"},
					}),
				},
			},
		},
		{
			input: `
disruptive: pass
non-disruptive:
    - nolog
    - setvar:
        collection: TX
        operation: =+
        assignments:
            - paramcounter_%{MATCHED_VAR_NAME}: "1"`,
			expected: SeclangActions{
				DisruptiveAction: mustNewActionOnly(Pass),
				NonDisruptiveActions: []Action{
					mustNewActionOnly(NoLog),
					mustNewSetvarAction(TX, Increment, []VarAssignment{{Variable: "paramcounter_%{MATCHED_VAR_NAME}", Value: "1"}}),
				},
			},
		},
		{
			input: `
disruptive: deny
non-disruptive:
    - log
data:
    - status: "500"`,
			expected: SeclangActions{
				DisruptiveAction: mustNewActionOnly(Deny),
				NonDisruptiveActions: []Action{
					mustNewActionOnly(Log),
				},
				DataActions: []Action{
					mustNewActionWithParam(Status, "500"),
				},
			},
		},
	}
)

func TestUnmarshalActions(t *testing.T) {
	for _, tt := range unmarshalTests {
		t.Run("Unmarshal actions from YAML", func(t *testing.T) {
			var result SeclangActions
			err := yaml.Unmarshal([]byte(tt.input), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewActionWithUnknownActions(t *testing.T) {
	t.Run("DisruptiveAction Unknown should return error", func(t *testing.T) {
		action, err := NewActionOnly(Unknown)
		assert.Error(t, err)
		assert.Equal(t, "invalid action: unknown action type", err.Error())
		assert.Empty(t, action)
	})

	t.Run("FlowAction Unknown should return error", func(t *testing.T) {
		action, err := NewActionOnly(FlowUnknown)
		assert.Error(t, err)
		assert.Equal(t, "invalid action: unknown action type", err.Error())
		assert.Empty(t, action)
	})

	t.Run("DataAction Unknown should return error", func(t *testing.T) {
		action, err := NewActionOnly(DataUnknown)
		assert.Error(t, err)
		assert.Equal(t, "invalid action: unknown action type", err.Error())
		assert.Empty(t, action)
	})

	t.Run("NonDisruptiveAction Unknown should return error", func(t *testing.T) {
		action, err := NewActionOnly(NonDisruptiveUnknown)
		assert.Error(t, err)
		assert.Equal(t, "invalid action: unknown action type", err.Error())
		assert.Empty(t, action)
	})
}

func TestNewActionWithValidActions(t *testing.T) {
	t.Run("DisruptiveAction Pass should work", func(t *testing.T) {
		action, err := NewActionOnly(Pass)
		assert.NoError(t, err)
		assert.Equal(t, "pass", action.GetKey())
	})

	t.Run("FlowAction Chain should work", func(t *testing.T) {
		action, err := NewActionOnly(Chain)
		assert.NoError(t, err)
		assert.Equal(t, "chain", action.GetKey())
	})

	t.Run("DataAction Status should work", func(t *testing.T) {
		action, err := NewActionWithParam(Status, "400")
		assert.NoError(t, err)
		assert.Equal(t, "status", action.GetKey())
		assert.Equal(t, "400", action.GetParam())
	})

	t.Run("NonDisruptiveAction Log should work", func(t *testing.T) {
		action, err := NewActionOnly(Log)
		assert.NoError(t, err)
		assert.Equal(t, "log", action.GetKey())
	})
}

func TestActionStringMethods(t *testing.T) {
	t.Run("DisruptiveAction Unknown string", func(t *testing.T) {
		assert.Equal(t, "unknown", Unknown.String())
	})

	t.Run("FlowAction Unknown string", func(t *testing.T) {
		assert.Equal(t, "unknown", FlowUnknown.String())
	})

	t.Run("DataAction Unknown string", func(t *testing.T) {
		assert.Equal(t, "unknown", DataUnknown.String())
	})

	t.Run("NonDisruptiveAction Unknown string", func(t *testing.T) {
		assert.Equal(t, "unknown", NonDisruptiveUnknown.String())
	})
}

func TestActionToStringMethods(t *testing.T) {
	t.Run("ActionOnly ToString for all action keys", func(t *testing.T) {
		tests := []struct {
			name     string
			action   ActionOnly
			expected string
		}{
			{name: "disruptive allow", action: ActionOnly(Allow.String()), expected: Allow.String()},
			{name: "disruptive block", action: ActionOnly(Block.String()), expected: Block.String()},
			{name: "disruptive deny", action: ActionOnly(Deny.String()), expected: Deny.String()},
			{name: "disruptive drop", action: ActionOnly(Drop.String()), expected: Drop.String()},
			{name: "disruptive pass", action: ActionOnly(Pass.String()), expected: Pass.String()},
			{name: "disruptive pause", action: ActionOnly(Pause.String()), expected: Pause.String()},
			{name: "flow chain", action: ActionOnly(Chain.String()), expected: Chain.String()},
			{name: "non disruptive auditlog", action: ActionOnly(AuditLog.String()), expected: AuditLog.String()},
			{name: "non disruptive capture", action: ActionOnly(Capture.String()), expected: Capture.String()},
			{name: "non disruptive log", action: ActionOnly(Log.String()), expected: Log.String()},
			{name: "non disruptive multiMatch", action: ActionOnly(MultiMatch.String()), expected: MultiMatch.String()},
			{name: "non disruptive noauditlog", action: ActionOnly(NoAuditLog.String()), expected: NoAuditLog.String()},
			{name: "non disruptive nolog", action: ActionOnly(NoLog.String()), expected: NoLog.String()},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.action.ToString())
			})
		}
	})

	t.Run("ActionWithParam ToString for all action keys", func(t *testing.T) {
		tests := []struct {
			name     string
			action   ActionWithParam
			expected string
		}{
			{name: "proxy", action: ActionWithParam{Proxy.String(): "value"}, expected: "proxy:'value'"},
			{name: "redirect", action: ActionWithParam{Redirect.String(): "value"}, expected: "redirect:'value'"},
			{name: "skip", action: ActionWithParam{Skip.String(): "value"}, expected: "skip:'value'"},
			{name: "skipAfter", action: ActionWithParam{SkipAfter.String(): "value"}, expected: "skipAfter:'value'"},
			{name: "status", action: ActionWithParam{Status.String(): "500"}, expected: "status:'500'"},
			{name: "xmlns", action: ActionWithParam{XLMNS.String(): "ns"}, expected: "xmlns:'ns'"},
			{name: "append", action: ActionWithParam{Append.String(): "value"}, expected: "append:'value'"},
			{name: "ctl", action: ActionWithParam{Ctl.String(): "ruleEngine=Off"}, expected: "ctl:ruleEngine=Off"},
			{name: "deprecatevar", action: ActionWithParam{DeprecateVar.String(): "value"}, expected: "deprecatevar:'value'"},
			{name: "exec", action: ActionWithParam{Exec.String(): "value"}, expected: "exec:'value'"},
			{name: "expirevar", action: ActionWithParam{ExpireVar.String(): "value"}, expected: "expirevar:'value'"},
			{name: "initcol", action: ActionWithParam{InitCol.String(): "value"}, expected: "initcol:'value'"},
			{name: "logdata", action: ActionWithParam{LogData.String(): "value"}, expected: "logdata:'value'"},
			{name: "prepend", action: ActionWithParam{Prepend.String(): "value"}, expected: "prepend:'value'"},
			{name: "sanitiseArg", action: ActionWithParam{SanitiseArg.String(): "value"}, expected: "sanitiseArg:'value'"},
			{name: "sanitiseMatched", action: ActionWithParam{SanitiseMatched.String(): "value"}, expected: "sanitiseMatched:'value'"},
			{name: "sanitiseMatchedBytes", action: ActionWithParam{SanitiseMatchedBytes.String(): "value"}, expected: "sanitiseMatchedBytes:'value'"},
			{name: "sanitiseRequestHeader", action: ActionWithParam{SanitiseRequestHeader.String(): "value"}, expected: "sanitiseRequestHeader:'value'"},
			{name: "sanitiseResponseHeader", action: ActionWithParam{SanitiseResponseHeader.String(): "value"}, expected: "sanitiseResponseHeader:'value'"},
			{name: "setuid", action: ActionWithParam{SetUid.String(): "value"}, expected: "setuid:'value'"},
			{name: "setrsc", action: ActionWithParam{SetRsc.String(): "value"}, expected: "setrsc:'value'"},
			{name: "setsid", action: ActionWithParam{SetSid.String(): "value"}, expected: "setsid:'value'"},
			{name: "setenv", action: ActionWithParam{SetEnv.String(): "value"}, expected: "setenv:'value'"},
			{name: "setvar", action: ActionWithParam{SetVar.String(): "tx.flag=on"}, expected: "setvar:'tx.flag=on'"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.action.ToString())
			})
		}
	})

	t.Run("SetvarAction ToString cases", func(t *testing.T) {
		tests := []struct {
			name     string
			action   SetvarAction
			expected string
		}{
			{
				name: "with string assignments",
				action: SetvarAction{
					Collection: TX,
					Operation:  Assign,
					Assignments: []VarAssignment{
						{Variable: "test", Value: "critical"},
						{Variable: "test2", Value: "payload with spaces"},
					},
				},
				expected: "setvar:'TX.test=critical',setvar:'TX.test2=payload with spaces'",
			},
			{
				name: "numeric assignments",
				action: SetvarAction{
					Collection: TX,
					Operation:  Assign,
					Assignments: []VarAssignment{
						{Variable: "counter", Value: "1"},
						{Variable: "score", Value: "5"},
					},
				},
				expected: "setvar:'TX.counter=1',setvar:'TX.score=5'",
			},
			{
				name: "numeric assignments with increment operation",
				action: SetvarAction{
					Collection: TX,
					Operation:  Increment,
					Assignments: []VarAssignment{
						{Variable: "counter", Value: "1"},
						{Variable: "score", Value: "5"},
					},
				},
				expected: "setvar:'TX.counter=+1',setvar:'TX.score=+5'",
			},
			{
				name: "numeric assignments with decrement operation",
				action: SetvarAction{
					Collection: TX,
					Operation:  Decrement,
					Assignments: []VarAssignment{
						{Variable: "counter", Value: "1"},
						{Variable: "score", Value: "5"},
					},
				},
				expected: "setvar:'TX.counter=-1',setvar:'TX.score=-5'",
			},
			{
				name:     "without assignments",
				action:   SetvarAction{Collection: TX, Operation: Assign},
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.action.ToString())
			})
		}
	})
}

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

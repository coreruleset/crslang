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
    - setvar: tx.rfi_parameter_%{MATCHED_VAR_NAME}=.%{tx.1}
flow:
    - chain`,
			expected: SeclangActions{
				DisruptiveAction: mustNewActionOnly(Block),
				NonDisruptiveActions: []Action{
					mustNewActionOnly(Capture),
					mustNewActionWithParam(LogData, "Matched Data: %{TX.0} found within %{MATCHED_VAR_NAME}: %{MATCHED_VAR}"),
					mustNewActionWithParam(SetVar, "tx.rfi_parameter_%{MATCHED_VAR_NAME}=.%{tx.1}"),
				},
				FlowActions: []Action{
					mustNewActionOnly(Chain),
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

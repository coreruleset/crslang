package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

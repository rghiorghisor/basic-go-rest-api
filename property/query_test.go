package property

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestNewFields(t *testing.T) {
	fields := NewFields([]string{" a ", "B "})

	assert.Equal(t, true, fields.IsEnabled())
	assert.Equal(t, true, fields.Contains("a"))
	assert.Equal(t, true, fields.Contains("A"))
	assert.Equal(t, true, fields.Contains("b"))
	assert.Equal(t, true, fields.Contains("B"))
	assert.Equal(t, false, fields.Contains("3"))

	fields.Disable()
	assert.Equal(t, false, fields.IsEnabled())
}

func TestNewFieldsEmpty(t *testing.T) {
	fields := NewFields([]string{})

	assert.Equal(t, false, fields.IsEnabled())
	fields.Disable()
	assert.Equal(t, false, fields.IsEnabled())
}

func TestHasSet(t *testing.T) {
	query := &Query{}

	assert.Equal(t, false, query.HasSet())
	query.Set = "set-test"
	assert.Equal(t, true, query.HasSet())
	assert.Equal(t, "set-test", query.GetSet())
}

package container

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestProvide(t *testing.T) {
	container := New()

	container.Provide(NewStructure1)
	err := container.Invoke(func(structure Structure) {
		assert.Equal(t, "1", structure.id)
	})
	assert.Equal(t, err, nil)
}

type Structure struct {
	id string
}

func NewStructure1() Structure {
	return Structure{"1"}
}

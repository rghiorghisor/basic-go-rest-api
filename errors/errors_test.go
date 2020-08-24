package errors

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rghiorghisor/basic-go-rest-api/model"
)

func TestEntityNotFound(t *testing.T) {
	err := NewEntityNotFound(reflect.TypeOf(model.Property{}), "123")
	actual := err.(*Error)
	assert.Equal(t, 404, actual.Code)
	assert.Equal(t, "Cannot find model.Property entity (id='123')", actual.Message)
	assert.Equal(t, "[code=404][Cannot find model.Property entity (id='123')]", actual.Error())
}

func TestConflict(t *testing.T) {
	err := NewConflict(reflect.TypeOf(model.Property{}), "name", "TestName")
	actual := err.(*Error)
	assert.Equal(t, 409, actual.Code)
	assert.Equal(t, "Found model.Property with same unique property (name='TestName')", actual.Message)
	assert.Equal(t, "[code=409][Found model.Property with same unique property (name='TestName')]", actual.Error())
}

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

func TestInvalidEntityEmpty(t *testing.T) {
	err := NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name")
	actual := err.(*Error)
	assert.Equal(t, 400, actual.Code)
	assert.Equal(t, "Invalid properties for model.Property entity. Property 'name' cannot be empty", actual.Message)
	assert.Equal(t, "[code=400][Invalid properties for model.Property entity. Property 'name' cannot be empty]", actual.Error())
}

func TestInvalidEntityCustom(t *testing.T) {
	err := NewInvalidEntityCustom(reflect.TypeOf(model.Property{}), "test message")
	actual := err.(*Error)
	assert.Equal(t, 400, actual.Code)
	assert.Equal(t, "Invalid properties for model.Property entity. test message", actual.Message)
	assert.Equal(t, "[code=400][Invalid properties for model.Property entity. test message]", actual.Error())
}

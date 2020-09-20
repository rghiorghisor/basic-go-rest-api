package errors

import (
	"fmt"
	"reflect"
)

// Error signals that something went wrong during the business actions.
//
// Such errors may signal both client and server errors.
// For the time being, here there are used the same codes as the HTTP response codes.
type Error struct {
	Code    int
	Message string
}

const (
	notFound          = 100
	conflict          = 101
	invalidEmptyField = 102
	invalidCustom     = 103
)

var errorTemplates = map[int]errorTemplate{
	notFound:          errorTemplate{404, "Cannot find %s entity (id='%s')"},
	conflict:          errorTemplate{409, "Found %s with same unique property (%s='%s')"},
	invalidEmptyField: errorTemplate{400, "Invalid %s entity. Property '%s' cannot be empty"},
	invalidCustom:     errorTemplate{400, "Invalid %s entity. %s"},
}

func (e *Error) Error() string {
	return fmt.Sprintf("[code=%d][%s]", e.Code, e.Message)
}

// NewEntityNotFound retrieves a new Error, signaling that a certain entity is
// not available.
func NewEntityNotFound(entity interface{}, identifier string) error {

	return createError(notFound, entity, identifier)
}

// NewConflict retrieves a new Error, signaling that an entity cannot be processed
// because there was found another such entity.
func NewConflict(entity interface{}, propName string, propValue string) error {

	return createError(conflict, entity, propName, propValue)
}

// NewInvalidEntityEmpty retrieves a new Error, signaling that a certain entity is
// not valid due to one of its property being empty.
func NewInvalidEntityEmpty(entity interface{}, propName string) error {

	return createError(invalidEmptyField, entity, propName)
}

// NewInvalidEntityCustom retrieves a new Error, signaling that a certain entity is
// not valid. The details of the invalid status must be found in the passed message
func NewInvalidEntityCustom(entity interface{}, message string) error {

	return createError(invalidCustom, entity, message)
}

func createError(errorType int, entity interface{}, args ...interface{}) error {
	errorTemplate := errorTemplates[errorType]

	message := errorTemplate.message
	typ := reflect.TypeOf(entity)
	args = append([]interface{}{typ}, args...)
	message = fmt.Sprintf(message, args...)

	return &Error{errorTemplate.code, message}
}

type errorTemplate struct {
	code    int
	message string
}

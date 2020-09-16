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

func (e *Error) Error() string {
	return fmt.Sprintf("[code=%d][%s]", e.Code, e.Message)
}

// NewEntityNotFound retrieves a new Error, signaling that a certain entity is
// not available.
func NewEntityNotFound(t reflect.Type, identifier string) error {
	s := "Cannot find %s entity (id='%s')"
	s = fmt.Sprintf(s, t, identifier)

	return &Error{404, s}
}

// NewConflict retrieves a new Error, signaling that an entity cannot be processed
// because there was found another such entity.
func NewConflict(t reflect.Type, propName string, propValue string) error {
	s := "Found %s with same unique property (%s='%s')"
	s = fmt.Sprintf(s, t, propName, propValue)

	return &Error{409, s}
}

// NewInvalidEntityEmpty retrieves a new Error, signaling that a certain entity is
// not valid due to one of its property being empty.
func NewInvalidEntityEmpty(t reflect.Type, propName string) error {
	s := "Invalid properties for %s entity. Property '%s' cannot be empty"
	s = fmt.Sprintf(s, t, propName)

	return &Error{400, s}
}

// NewInvalidEntityCustom retrieves a new Error, signaling that a certain entity is
// not valid. The details of the invalid status must be found in the passed message
func NewInvalidEntityCustom(t reflect.Type, message string) error {
	s := "Invalid properties for %s entity. %s"
	s = fmt.Sprintf(s, t, message)

	return &Error{400, s}
}

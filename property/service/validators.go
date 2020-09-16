package service

import (
	"reflect"
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
)

type validator interface {
	check(prop *model.Property) error
}

type validators struct {
	values []validator
}

func newValidators() validators {
	return validators{
		values: []validator{nameValidator{}},
	}
}

func (v validators) check(prop *model.Property) error {
	for _, v := range v.values {
		err := v.check(prop)
		if err != nil {
			return err
		}
	}

	return nil
}

type nameValidator struct {
}

func (v nameValidator) check(prop *model.Property) error {
	if prop.Name == "" {
		return errors.NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name")
	}

	if strings.ContainsAny(prop.Name, " ") {
		return errors.NewInvalidEntityCustom(reflect.TypeOf(model.Property{}), "'name' cannot contain spaces.")
	}

	return nil
}

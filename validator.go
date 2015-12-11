// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"reflect"

	"gopkg.in/go-playground/validator.v8"
)

type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the obj is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(obj interface{}) error
}

var DefaultValidator StructValidator = (*defaultValidator)(validator.New(&validator.Config{TagName: "validate"}))

type defaultValidator validator.Validate

func (p *defaultValidator) ValidateStruct(obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	return (*validator.Validate)(p).Struct(obj)
}

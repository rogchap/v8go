// Copyright 2020 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"errors"
)

type ValueTemplate struct {
	*Value
}

// NewValueTemplate is nothing but sugar to allow better composition.
func NewValueTemplate(iso *Isolate, val interface{}) (*ValueTemplate, error) {
	v, err := NewValue(iso, val)
	if err != nil {
		return nil, err
	}
	if v.IsObject() || v.IsExternal() {
		return nil, errors.New("v8go: unsupported property: value type must be a primitive or use a template")
	}
	return &ValueTemplate{Value: v}, nil
}

// ToValue creates a new Value based on the template.
func (o *ValueTemplate) ToValue(ctx *ExecContext) (Valuer, error) {
	return o.Value, nil
}

// Set is here to comply with Templater interface, othwerwise useless.
func (t *ValueTemplate) Set(name string, val Templater, attributes ...PropertyAttribute) error {
	return errors.New("ValueTemplate does not allow to set any property")
}

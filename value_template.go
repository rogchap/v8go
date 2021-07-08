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

// GetObject creates a new Object based on the template.
func (o *ValueTemplate) GetValue(ctx *ExecContext) (Valuer, error) {
	return o.Value, nil
}

func (o *ValueTemplate) templater() {}

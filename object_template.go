package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"

// PropertyAttribute are the attribute flags for a property on an Object.
// Typical usage when setting an Object or TemplateObject property, and
// can also be validated when accessing a property.
type PropertyAttribute uint8

const (
	// None.
	None PropertyAttribute = 0
	// ReadOnly, ie. not writable.
	ReadOnly PropertyAttribute = 1 << iota
	// DontEnum, ie. not enumerable.
	DontEnum
	// DontDelete, ie. not configurable.
	DontDelete
)

// ObjectTemplate is used to create objects at runtime.
// Properties added to an ObjectTemplate are added to each object created from the ObjectTemplate.
type ObjectTemplate struct {
	*template
}

// NewObjectTemplate creates a new ObjectTemplate.
// The *ObjectTemplate can be used as a v8go.ContextOption to create a global object in a Context.
func NewObjectTemplate(iso *Isolate) (*ObjectTemplate, error) {
	tmpl, err := newTemplate(iso)
	if err != nil {
		return nil, err
	}
	return &ObjectTemplate{tmpl}, nil
}

func (o *ObjectTemplate) apply(opts *contextOptions) {
	opts.gTmpl = o
}

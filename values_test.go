package v8go

import (
	"testing"
)

func TestValues(t *testing.T) {
	iso, _ := NewIsolate()
	vals, err := NewValues(iso, 1, "hello", true)
	if err != nil {
		t.Error("failed to get values")
	}
	err = vals.Validate(
		ValueCondLen("must have 3 arguments", 3),
		ValueCondMinLen("must have 3 arguments", 3),

		ValueCondType("must be number", 0, IsNumber),
		ValueCondType("must be string", 1, IsString),
		ValueCondType("must be bool", 2, IsBoolean),

		ValueCondOptionalType("must be number", 0, IsNumber),
		ValueCondOptionalType("must be string", 1, IsString),
		ValueCondOptionalType("must be bool", 2, IsBoolean),
		// optional
		ValueCondOptionalType("must be bool", 3, IsBoolean),
	)
	if err != nil {
		t.Errorf("failed to vlidate %s", err.Error())
	}

}

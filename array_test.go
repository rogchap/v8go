package v8go_test

import (
	"errors"
	"fmt"
	"testing"

	v8 "rogchap.com/v8go"
)

func reverseUint8ArrayFunctionCallback(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	args := info.Args()

	if len(args) != 1 {
		e, _ := v8.NewValue(iso, "Function ReverseUint8Array expects 1 parameter")
		return iso.ThrowException(e)
	}
	if !args[0].IsUint8Array() {
		e, _ := v8.NewValue(iso, "Function ReverseUint8Array expects Uint8Array parameter")
		return iso.ThrowException(e)
	}
	array := args[0].Uint8Array()
	length := len(array)
	reversed := make([]uint8, length)
	for i := 0; i < length; i++ {
		reversed[i] = array[length-i-1]
	}
	val, err := v8.NewValue(iso, reversed)
	if err != nil {
		e, _ := v8.NewValue(iso, fmt.Sprintf("Could not get value for array: %v\n", err))
		return iso.ThrowException(e)
	}
	return val
}

func injectUint8ArrayTester(ctx *v8.Context) error {
	if ctx == nil {
		return errors.New("injectUint8ArrayTester: ctx is required")
	}

	iso := ctx.Isolate()

	con := v8.NewObjectTemplate(iso)

	reverseFn := v8.NewFunctionTemplate(iso, reverseUint8ArrayFunctionCallback)

	if err := con.Set("reverseUint8Array", reverseFn, v8.ReadOnly); err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	nativeObj, err := con.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	global := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	return nil
}

// Test that a script can call a go function to reverse a []uint8 array
func TestUint8Array(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)

	if err := injectUint8ArrayTester(ctx); err != nil {
		t.Error(err)
	}

	if val, err := ctx.RunScript("native.reverseUint8Array(new Uint8Array([0,1,2,3,4,5,6,7,8,9]))", ""); err != nil {
		t.Error(err)
	} else {
		if !val.IsUint8Array() {
			t.Errorf("Expected uint8 array return value")
		}
		arr := val.Uint8Array()
		if len(arr) != 10 {
			t.Errorf("Got wrong array length %d, expected 10", len(arr))
		}
		for i := 0; i < 10; i++ {
			if arr[i] != uint8(10-i-1) {
				t.Errorf("Incorrect byte at index %d (whole array: %v)", i, arr)
			}
		}
	}
}

// Test that a native go function can throw exceptions that make it back to the script runner
func TestUint8ArrayException(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)

	if err := injectUint8ArrayTester(ctx); err != nil {
		t.Error(err)
	}

	if _, err := ctx.RunScript("native.reverseUint8Array(\"notanarray\")", ""); err != nil {
		t.Logf("Got expected error from script: %v", err)
	} else {
		t.Errorf("Should have received an error from the script")
	}
}

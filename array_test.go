package v8go

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func reverseUint8ArrayFunctionCallback(info *FunctionCallbackInfo) *Value {
	iso, err := info.ExecContext().Isolate()
	if err != nil {
		log.Fatalf("Could not get isolate from context: %v\n", err)
	}
	args := info.Args()
	if len(args) != 1 {
		return iso.ThrowException("Function ReverseUint8Array expects 1 parameter")
	}
	if !args[0].IsUint8Array() {
		return iso.ThrowException("Function ReverseUint8Array expects Uint8Array parameter")
	}
	array := args[0].Uint8Array()
	length := len(array)
	reversed := make([]uint8, length)
	for i := 0; i < length; i++ {
		reversed[i] = array[length-i-1]
	}
	val, err := NewValue(iso, reversed)
	if err != nil {
		return iso.ThrowException(fmt.Sprintf("Could not get value for array: %v\n", err))
	}
	return val
}

func injectUint8ArrayTester(ctx *ExecContext) error {
	if ctx == nil {
		return errors.New("injectUint8ArrayTester: ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	con, err := NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	reverseFn, err := NewFunctionTemplate(iso, reverseUint8ArrayFunctionCallback)
	if err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	if err := con.Set("reverseUint8Array", reverseFn, ReadOnly); err != nil {
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

	iso, _ := NewIsolate()
	ctx, _ := NewExecContext(iso)

	if err := injectUint8ArrayTester(ctx); err != nil {
		t.Error(err)
	}

	if val, err := ctx.RunScript("native.reverseUint8Array(new Uint8Array([0,1,2,3,4,5,6,7,8,9]))", ""); err != nil {
		t.Error(err)
	} else {
		if !val.IsUint8Array() {
			t.Errorf("Expected uint8 array return value")
		}
		fmt.Printf("Reversed array: %v\n", val.Uint8Array())
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

	iso, _ := NewIsolate()
	ctx, _ := NewExecContext(iso)

	if err := injectUint8ArrayTester(ctx); err != nil {
		t.Error(err)
	}

	if _, err := ctx.RunScript("native.reverseUint8Array(\"notanarray\")", ""); err != nil {
		t.Logf("Got expected error from script: %v", err)
	} else {
		t.Errorf("Should have received an error from the script")
	}
}

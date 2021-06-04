package v8go

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

type NativeObject interface {
	GetReverseUint8ArrayFunctionCallback() FunctionCallback
}

type nativeObject struct {
}

func NewNativeObject() NativeObject {
	return &nativeObject{}
}

func (nto *nativeObject) GetReverseUint8ArrayFunctionCallback() FunctionCallback {
	return func(info *FunctionCallbackInfo) *Value {
		iso, err := info.Context().Isolate()
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
}

func injectNativeObject(ctx *Context) error {
	if ctx == nil {
		return errors.New("injectNativeObject: ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	c := NewNativeObject()

	con, err := NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	reverseFn, err := NewFunctionTemplate(iso, c.GetReverseUint8ArrayFunctionCallback())
	if err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	if err := con.Set("reverseUint8Array", reverseFn, ReadOnly); err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	nativeObj, err := con.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	global := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("injectNativeObject: %w", err)
	}

	return nil
}

// Test that a script can call a go function to reverse a []uint8 array
func TestNativeUint8Array(t *testing.T) {
	t.Parallel()

	iso, _ := NewIsolate()
	ctx, _ := NewContext(iso)

	if err := injectNativeObject(ctx); err != nil {
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
func TestNativeUint8ArrayException(t *testing.T) {
	t.Parallel()

	iso, _ := NewIsolate()
	ctx, _ := NewContext(iso)

	if err := injectNativeObject(ctx); err != nil {
		t.Error(err)
	}

	if _, err := ctx.RunScript("native.reverseUint8Array(\"notanarray\")", ""); err != nil {
		t.Logf("Got expected error from script: %v", err)
	} else {
		t.Errorf("Should have received an error from the script")
	}
}

func TestNativeUint8ArrayManyCalls(t *testing.T) {
	t.Parallel()
	iso, _ := NewIsolate()
	ctx, _ := NewContext(iso)
	if err := injectNativeObject(ctx); err != nil {
		t.Error(err)
	}
	stats := iso.GetHeapStatistics()
	fmt.Printf("MEMSTATS BEFORE: %+v\n", stats)

	if _, err := ctx.RunScript("for(i = 0; i < 100000; i++) native.reverseUint8Array(new Uint8Array([0,1,2,3,4,5,6,7,8,9]));", ""); err != nil {
		t.Error(err)
	}

	stats = iso.GetHeapStatistics()
	fmt.Printf("MEMSTATS AFTER: %+v\n", stats)
}

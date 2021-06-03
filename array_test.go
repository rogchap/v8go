package v8go

import (
	"errors"
	"fmt"
	"os"
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
		args := info.Args()
		if len(args) != 1 {
			os.Stderr.WriteString("Function ReverseUint8Array expects 1 parameter\n") //TODO
			return nil
		}
		iso, err := info.Context().Isolate()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Could not get isolate from context: %v\n", err)) //TODO
			return nil
		}
		if !args[0].IsUint8Array() {
			os.Stderr.WriteString("Function ReverseUint8Array expects Uint8Array parameter\n") //TODO
			return nil
		}
		inarray := args[0].Uint8Array() //TODO who frees this
		length := len(inarray)
		reversed := make([]uint8, length)
		for i := 0; i < length; i++ {
			reversed[i] = inarray[length-i-1]
		}
		val, err := NewValue(iso, reversed)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Could not get value for array: %v\n", err))
			return nil
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

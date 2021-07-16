package v8go

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func TestNewArray(t *testing.T) {
	ctx, _ := NewExecContext()
	arr := NewArray(ctx, 2)
	err := arr.SetIdx(0, true)
	if err != nil {
		t.Error(err.Error())
	}
	str, err := JSONStringify(ctx, arr)
	if err != nil {
		t.Error(err.Error())
	}
	if str != "[true,null]" {
		t.Error("invalid array output")
	}
}

func TestNewArrayFromStrings(t *testing.T) {
	iso, _ := NewIsolate()
	ctx, _ := NewExecContext(iso)

	v1, err := NewValue(iso, int32(1))
	if err != nil {
		t.Error("failed to create value")
	}

	arr := NewArrayFromValues(ctx, []Valuer{v1})
	str, err := JSONStringify(ctx, arr)
	if err != nil {
		t.Error(err.Error())
	}
	if str != `[1]` {
		t.Error("invalid array output")
	}

	v1, err = NewValue(iso, "1")
	if err != nil {
		t.Error("failed to create value")
	}

	arr = NewArrayFromValues(ctx, []Valuer{v1})
	str, err = JSONStringify(ctx, arr)
	if err != nil {
		t.Error(err.Error())
	}
	if str != `["1"]` {
		t.Error("invalid array output")
	}
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
		t.Logf("Reversed array: %v\n", val.Uint8Array())
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

	if _, err := ctx.RunScript("native.reverseUint8Array(\"notanarray\")", ""); err == nil {
		t.Errorf("Should have received an error from the script")
	}
}

// Test that a script can call a go function to reverse an ArrayBuffer.
// The function reverses the ArrayBuffer in-place, i.e. this is a call-by-reference.
func TestModifyArrayBuffer(t *testing.T) {
	t.Parallel()

	iso, _ := NewIsolate()
	ctx, _ := NewExecContext(iso)
	if err := injectArrayBufferTester(ctx, "reverseArrayBuffer", reverseArrayBufferFunctionCallback); err != nil {
		t.Error(err)
	}

	js := `
		let ab = new ArrayBuffer(10);
		let view = new Uint8Array(ab);
		for (let i = 0; i < 10; i++) view[i] = i;
		native.reverseArrayBuffer(ab);
		ab;
	`

	if val, err := ctx.RunScript(js, ""); err != nil {
		t.Error(err)
	} else {
		if !val.IsArrayBuffer() {
			t.Errorf("Expected ArrayBuffer return value")
		}
		ab := val.ArrayBuffer()
		if ab.Len() != 10 {
			t.Errorf("Got wrong ArrayBuffer length %d, expected 10", ab.Len())
		}
		bytes := ab.Bytes()
		t.Logf("Got reversed ArrayBuffer from script: %v\n", bytes)
		for i := int64(0); i < ab.Len(); i++ {
			if bytes[i] != uint8(10-i-1) {
				t.Errorf("Incorrect byte at index %d (whole array: %v)", i, ab)
			}
		}
	}
}

func TestCreateArrayBuffer(t *testing.T) {
	t.Parallel()

	iso, _ := NewIsolate()
	ctx, _ := NewExecContext(iso)
	if err := injectArrayBufferTester(ctx, "createArrayBuffer", createArrayBufferFunctionCallback); err != nil {
		t.Error(err)
	}

	js := `
		native.createArrayBuffer(16);
	`

	if val, err := ctx.RunScript(js, ""); err != nil {
		t.Error(err)
	} else {
		if !val.IsArrayBuffer() {
			t.Errorf("Expected ArrayBuffer return value")
		}
		ab := val.ArrayBuffer()
		if ab.Len() != 16 {
			t.Errorf("Got wrong ArrayBuffer length %d, expected 16", ab.Len())
		}
		bytes := ab.Bytes()
		fmt.Printf("Got ArrayBuffer from script: %v\n", bytes)
		for i := int64(0); i < ab.Len(); i++ {
			if bytes[i] != uint8(i) {
				t.Errorf("Incorrect byte at index %d (whole array: %v)", i, bytes)
			}
		}
	}
}

func reverseArrayBufferFunctionCallback(info *FunctionCallbackInfo) (Valuer, error) {
	iso, err := info.ExecContext().Isolate()
	if err != nil {
		log.Fatalf("Could not get isolate from context: %v\n", err)
	}
	args := info.Args()
	if len(args) != 1 {
		iso.ThrowException("Function ReverseArrayBuffer expects 1 parameter")
		return nil, nil
	}
	if !args[0].IsArrayBuffer() {
		iso.ThrowException("Function ReverseArrayBuffer expects ArrayBuffer parameter")
		return nil, nil
	}
	ab := args[0].ArrayBuffer() // "cast" to ArrayBuffer
	length := int(ab.Len())
	bytes := ab.Bytes() // get a copy of the bytes from the ArrayBuffer
	reversed := make([]uint8, length)
	for i := 0; i < length; i++ {
		reversed[i] = bytes[length-i-1]
	}
	ab.Write(reversed) // update the bytes in the ArrayBuffer (length must match!)
	return nil, nil
}

func createArrayBufferFunctionCallback(info *FunctionCallbackInfo) (Valuer, error) {
	iso, err := info.ExecContext().Isolate()
	if err != nil {
		log.Fatalf("Could not get isolate from context: %v\n", err)
	}
	args := info.Args()
	if len(args) != 1 {
		iso.ThrowException("Function CreateArrayBuffer expects 1 parameter")
		return nil, nil
	}
	if !args[0].IsInt32() {
		iso.ThrowException("Function CreateArrayBuffer expects Int32 parameter")
		return nil, nil
	}
	length := args[0].Int32()
	ab := NewArrayBuffer(info.ExecContext(), int(length)) // create ArrayBuffer object of given length
	bytes := make([]uint8, length)
	for i := uint8(0); i < uint8(length); i++ {
		bytes[i] = i
	}
	ab.Write(bytes)      // copy these bytes into it. Caller is responsible for avoiding overruns!
	return ab.Value, nil // return the ArrayBuffer to javascript
}

func injectArrayBufferTester(ctx *ExecContext, funcName string, funcCb FunctionCallback) error {
	if ctx == nil {
		return errors.New("injectArrayBufferTester: ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	con, err := NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	funcTempl, err := NewFunctionTemplate(iso, funcCb)
	if err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	if err := con.Set(funcName, funcTempl, ReadOnly); err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	nativeObj, err := con.GetObject(ctx)
	if err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	global := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("injectArrayBufferTester: %v", err)
	}

	return nil
}

func reverseUint8ArrayFunctionCallback(info *FunctionCallbackInfo) (Valuer, error) {
	args := info.Args()
	if len(args) != 1 {
		return nil, errors.New("Function ReverseUint8Array expects 1 parameter")
	}
	if !args[0].IsUint8Array() {
		return nil, errors.New("Function ReverseUint8Array expects Uint8Array parameter")
	}
	array := args[0].Bytes()
	length := len(array)
	reversed := make([]byte, length)
	for i := 0; i < length; i++ {
		reversed[i] = array[length-i-1]
	}
	buf := NewArrayBufferFromBytes(info.ExecContext(), reversed)
	val := NewTypedUint8ArrayFromBuffer(buf)
	return val, nil
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

	nativeObj, err := con.GetObject(ctx)
	if err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	global := ctx.Global()

	if err := global.Set("native", nativeObj); err != nil {
		return fmt.Errorf("injectUint8ArrayTester: %v", err)
	}

	return nil
}

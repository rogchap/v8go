// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"runtime"
	"testing"

	"rogchap.com/v8go"
)

func TestValueNewBaseCases(t *testing.T) {
	t.Parallel()
	if _, err := v8go.NewValue(nil, ""); err == nil {
		t.Error("expected error, but got <nil>")
	}
	iso, _ := v8go.NewIsolate()
	if _, err := v8go.NewValue(iso, nil); err == nil {
		t.Error("expected error, but got <nil>")
	}
	if _, err := v8go.NewValue(iso, struct{}{}); err == nil {
		t.Error("expected error, but got <nil>")
	}

}

func TestValueFormatting(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source          string
		defaultVerb     string
		defaultVerbFlag string
		stringVerb      string
		quoteVerb       string
	}{
		{"new Object()", "[object Object]", "#<Object>", "[object Object]", `"[object Object]"`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			if s := fmt.Sprintf("%v", val); s != tt.defaultVerb {
				t.Errorf("incorrect format for %%v: %s", s)
			}
			if s := fmt.Sprintf("%+v", val); s != tt.defaultVerbFlag {
				t.Errorf("incorrect format for %%+v: %s", s)
			}
			if s := fmt.Sprintf("%s", val); s != tt.stringVerb {
				t.Errorf("incorrect format for %%s: %s", s)
			}
			if s := fmt.Sprintf("%q", val); s != tt.quoteVerb {
				t.Errorf("incorrect format for %%q: %s", s)
			}
		})
	}
}

func TestValueString(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		name   string
		source string
		out    string
	}{
		{"Number", `13 * 2`, "26"},
		{"String", `"string"`, "string"},
		{"Object", `let obj = {}; obj`, "[object Object]"},
		{"Function", `let fn = function(){}; fn`, "function(){}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, _ := ctx.RunScript(tt.source, "test.js")
			str := result.String()
			if str != tt.out {
				t.Errorf("unexpected result: expected %q, got %q", tt.out, str)
			}
		})
	}
}

func TestValueDetailString(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		name   string
		source string
		out    string
	}{
		{"Number", `13 * 2`, "26"},
		{"String", `"a string"`, "a string"},
		{"Object", `let obj = {}; obj`, "#<Object>"},
		{"Function", `let fn = function(){}; fn`, "function(){}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, _ := ctx.RunScript(tt.source, "test.js")
			str := result.DetailString()
			if str != tt.out {
				t.Errorf("unexpected result: expected %q, got %q", tt.out, str)
			}
		})
	}
}

func TestValueBoolean(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source string
		out    bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"null", false},
		{"undefined", false},
		{"''", false},
		{"'foo'", true},
		{"() => {}", true},
		{"{}", false},
		{"{foo:'bar'}", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			if b := val.Boolean(); b != tt.out {
				t.Errorf("unexpected value: expected %v, got %v", tt.out, b)
			}
		})
	}
}

func TestValueArrayIndex(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source string
		idx    uint32
		ok     bool
	}{
		{"1", 1, true},
		{"0", 0, true},
		{"-1", 0, false},
		{"'1'", 1, true},
		{"'-1'", 0, false},
		{"'a'", 0, false},
		{"[1]", 1, true},
		{"['1']", 1, true},
		{"[1, 1]", 0, false},
		{"{}", 0, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			idx, ok := val.ArrayIndex()
			if ok != tt.ok {
				t.Errorf("unexpected ok: expected %v, got %v", tt.ok, ok)
			}
			if idx != tt.idx {
				t.Errorf("unexpected array index: expected %v, got %v", tt.idx, idx)
			}
		})
	}
}

func TestValueInt32(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source   string
		expected int32
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"'1'", 1},
		{"1.5", 1},
		{"-1.5", -1},
		{"'a'", 0},
		{"[1]", 1},
		{"[1,1]", 0},
		{"Infinity", 0},
		{"Number.MAX_SAFE_INTEGER", -1},
		{"Number.MIN_SAFE_INTEGER", 1},
		{"Number.NaN", 0},
		{"2_147_483_647", 1<<31 - 1},
		{"-2_147_483_648", -1 << 31},
		{"2_147_483_648", -1 << 31}, // overflow
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			if i32 := val.Int32(); i32 != tt.expected {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, i32)
			}
		})
	}
}

func TestValueInteger(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source   string
		expected int64
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"'1'", 1},
		{"1.5", 1},
		{"-1.5", -1},
		{"'a'", 0},
		{"[1]", 1},
		{"[1,1]", 0},
		{"Infinity", 1<<63 - 1},
		{"Number.MAX_SAFE_INTEGER", 1<<53 - 1},
		{"Number.MIN_SAFE_INTEGER", -(1<<53 - 1)},
		{"Number.NaN", 0},
		{"9_007_199_254_740_991", 1<<53 - 1},
		{"-9_007_199_254_740_991", -(1<<53 - 1)},
		{"9_223_372_036_854_775_810", 1<<63 - 1}, // does not overflow, pinned at 2^64 -1
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			if i64 := val.Integer(); i64 != tt.expected {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, i64)
			}
		})
	}
}

func TestValueNumber(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source   string
		expected float64
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"'1'", 1},
		{"1.5", 1.5},
		{"-1.5", -1.5},
		{"'a'", math.NaN()},
		{"[1]", 1},
		{"[1,1]", math.NaN()},
		{"Infinity", math.Inf(0)},
		{"Number.MAX_VALUE", 1.7976931348623157e+308},
		{"Number.MIN_VALUE", 5e-324},
		{"Number.MAX_SAFE_INTEGER", 1<<53 - 1},
		{"Number.NaN", math.NaN()},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			f64 := val.Number()
			if math.IsNaN(tt.expected) {
				if !math.IsNaN(f64) {
					t.Errorf("unexpected value: expected NaN, got %v", f64)
				}
				return
			}
			if f64 != tt.expected {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, f64)
			}
		})
	}
}

func TestValueUint32(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
	defer ctx.Close()

	tests := [...]struct {
		source   string
		expected uint32
	}{
		{"0", 0},
		{"1", 1},
		{"-1", 1<<32 - 1}, // overflow
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			val, _ := ctx.RunScript(tt.source, "test.js")
			if u32 := val.Uint32(); u32 != tt.expected {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, u32)
			}
		})
	}
}

func TestValueBigInt(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()

	x, _ := new(big.Int).SetString("36893488147419099136", 10) // larger than a single word size (64bit)

	tests := [...]struct {
		source   string
		expected *big.Int
	}{
		{"BigInt(0)", &big.Int{}},
		{"-1n", big.NewInt(-1)},
		{"new BigInt(1)", nil}, // bad syntax
		{"BigInt(Number.MAX_SAFE_INTEGER)", big.NewInt(1<<53 - 1)},
		{"BigInt(Number.MIN_SAFE_INTEGER)", new(big.Int).Neg(big.NewInt(1<<53 - 1))},
		{"BigInt(Number.MAX_SAFE_INTEGER) * 2n", big.NewInt(1<<54 - 2)},
		{"BigInt(Number.MAX_SAFE_INTEGER) * 4096n", x},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			t.Parallel()
			ctx, _ := v8go.NewContext(iso)
			defer ctx.Close()

			val, _ := ctx.RunScript(tt.source, "test.js")
			b := val.BigInt()
			if b == nil && tt.expected != nil {
				t.Errorf("uexpected <nil> value")
				return
			}
			if b != nil && tt.expected == nil {
				t.Errorf("expected <nil>, but got value: %v", b)
				return
			}
			if b != nil && b.Cmp(tt.expected) != 0 {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, b)
			}
		})
	}
}

func TestValueObject(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	defer ctx.Close()

	val, _ := ctx.RunScript("1", "")
	if _, err := val.AsObject(); err == nil {
		t.Error("Expected error but got <nil>")
	}
	if obj := val.Object(); obj.String() != "1" {
		t.Errorf("unexpected object value: %v", obj)
	}
}

func TestValuePromise(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	defer ctx.Close()

	val, _ := ctx.RunScript("1", "")
	if _, err := val.AsPromise(); err == nil {
		t.Error("Expected error but got <nil>")
	}
	if _, err := ctx.RunScript("new Promise(()=>{})", ""); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

}

func TestValueFunction(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	defer ctx.Close()

	val, _ := ctx.RunScript("1", "")
	if _, err := val.AsFunction(); err == nil {
		t.Error("Expected error but got <nil>")
	}
	val, err := ctx.RunScript("(a, b) => { return a + b; }", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if _, err := val.AsFunction(); err != nil {
		t.Errorf("Expected success but got: %v", err)
	}

}

func TestValueIsXXX(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()
	tests := [...]struct {
		source string
		assert func(*v8go.Value) bool
	}{
		{"", (*v8go.Value).IsUndefined},
		{"let v; v", (*v8go.Value).IsUndefined},
		{"null", (*v8go.Value).IsNull},
		{"let v; v", (*v8go.Value).IsNullOrUndefined},
		{"let v = null; v", (*v8go.Value).IsNullOrUndefined},
		{"true", (*v8go.Value).IsTrue},
		{"false", (*v8go.Value).IsFalse},
		{"'name'", (*v8go.Value).IsName},
		{"Symbol()", (*v8go.Value).IsName},
		{`"double quote"`, (*v8go.Value).IsString},
		{"'single quote'", (*v8go.Value).IsString},
		{"`string literal`", (*v8go.Value).IsString},
		{"Symbol()", (*v8go.Value).IsSymbol},
		{"Symbol('foo')", (*v8go.Value).IsSymbol},
		{"() => {}", (*v8go.Value).IsFunction},
		{"function v() {}; v", (*v8go.Value).IsFunction},
		{"const v = function() {}; v", (*v8go.Value).IsFunction},
		{"console.log", (*v8go.Value).IsFunction},
		{"Object", (*v8go.Value).IsFunction},
		{"class Foo {}; Foo", (*v8go.Value).IsFunction},
		{"class Foo { bar() {} }; (new Foo()).bar", (*v8go.Value).IsFunction},
		{"function* v(){}; v", (*v8go.Value).IsFunction},
		{"async function v(){}; v", (*v8go.Value).IsFunction},
		{"Object()", (*v8go.Value).IsObject},
		{"new Object", (*v8go.Value).IsObject},
		{"var v = {}; v", (*v8go.Value).IsObject},
		{"10n", (*v8go.Value).IsBigInt},
		{"BigInt(1)", (*v8go.Value).IsBigInt},
		{"true", (*v8go.Value).IsBoolean},
		{"false", (*v8go.Value).IsBoolean},
		{"Boolean()", (*v8go.Value).IsBoolean},
		{"(new Boolean).valueOf()", (*v8go.Value).IsBoolean},
		{"1", (*v8go.Value).IsNumber},
		{"1.1", (*v8go.Value).IsNumber},
		{"1_1", (*v8go.Value).IsNumber},
		{".1", (*v8go.Value).IsNumber},
		{"2e4", (*v8go.Value).IsNumber},
		{"0x2", (*v8go.Value).IsNumber},
		{"NaN", (*v8go.Value).IsNumber},
		{"Infinity", (*v8go.Value).IsNumber},
		{"Number(1)", (*v8go.Value).IsNumber},
		{"(new Number()).valueOf()", (*v8go.Value).IsNumber},
		{"1", (*v8go.Value).IsInt32},
		{"-1", (*v8go.Value).IsInt32},
		{"1", (*v8go.Value).IsUint32},
		{"new Date", (*v8go.Value).IsDate},
		{"function foo(){ return arguments }; foo()", (*v8go.Value).IsArgumentsObject},
		{"Object(1n)", (*v8go.Value).IsBigIntObject},
		{"Object(1)", (*v8go.Value).IsNumberObject},
		{"new Number", (*v8go.Value).IsNumberObject},
		{"new String", (*v8go.Value).IsStringObject},
		{"Object('')", (*v8go.Value).IsStringObject},
		{"Object(Symbol())", (*v8go.Value).IsSymbolObject},
		{"Error()", (*v8go.Value).IsNativeError},
		{"TypeError()", (*v8go.Value).IsNativeError},
		{"SyntaxError()", (*v8go.Value).IsNativeError},
		{"/./", (*v8go.Value).IsRegExp},
		{"RegExp()", (*v8go.Value).IsRegExp},
		{"async function v(){}; v", (*v8go.Value).IsAsyncFunction},
		{"let v = async () => {}; v", (*v8go.Value).IsAsyncFunction},
		{"function* v(){}; v", (*v8go.Value).IsGeneratorFunction},
		{"function* v(){}; v()", (*v8go.Value).IsGeneratorObject},
		{"new Promise(()=>{})", (*v8go.Value).IsPromise},
		{"new Map", (*v8go.Value).IsMap},
		{"new Set", (*v8go.Value).IsSet},
		{"(new Map).entries()", (*v8go.Value).IsMapIterator},
		{"(new Set).entries()", (*v8go.Value).IsSetIterator},
		{"new WeakMap", (*v8go.Value).IsWeakMap},
		{"new WeakSet", (*v8go.Value).IsWeakSet},
		{"new Array", (*v8go.Value).IsArray},
		{"Array()", (*v8go.Value).IsArray},
		{"[]", (*v8go.Value).IsArray},
		{"new ArrayBuffer", (*v8go.Value).IsArrayBuffer},
		{"new Int8Array", (*v8go.Value).IsArrayBufferView},
		{"new Int8Array", (*v8go.Value).IsTypedArray},
		{"new Uint32Array", (*v8go.Value).IsTypedArray},
		{"new Uint8Array", (*v8go.Value).IsUint8Array},
		{"new Uint8ClampedArray", (*v8go.Value).IsUint8ClampedArray},
		{"new Int8Array", (*v8go.Value).IsInt8Array},
		{"new Uint16Array", (*v8go.Value).IsUint16Array},
		{"new Int16Array", (*v8go.Value).IsInt16Array},
		{"new Uint32Array", (*v8go.Value).IsUint32Array},
		{"new Int32Array", (*v8go.Value).IsInt32Array},
		{"new Float32Array", (*v8go.Value).IsFloat32Array},
		{"new Float64Array", (*v8go.Value).IsFloat64Array},
		{"new BigInt64Array", (*v8go.Value).IsBigInt64Array},
		{"new BigUint64Array", (*v8go.Value).IsBigUint64Array},
		{"new DataView(new ArrayBuffer)", (*v8go.Value).IsDataView},
		{"new SharedArrayBuffer", (*v8go.Value).IsSharedArrayBuffer},
		{"new Proxy({},{})", (*v8go.Value).IsProxy},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.source, func(t *testing.T) {
			t.Parallel()
			ctx, _ := v8go.NewContext(iso)
			defer ctx.Close()

			val, err := ctx.RunScript(tt.source, "test.js")
			if err != nil {
				t.Fatalf("failed to run script: %v", err)
			}
			if !tt.assert(val) {
				t.Errorf("value is false for %s", runtime.FuncForPC(reflect.ValueOf(tt.assert).Pointer()).Name())
			}
		})
	}
}

func TestValueMarshalJSON(t *testing.T) {
	t.Parallel()
	iso, _ := v8go.NewIsolate()

	tests := [...]struct {
		name     string
		val      func() *v8go.Value
		expected []byte
	}{
		{
			"primitive",
			func() *v8go.Value {
				val, _ := v8go.NewValue(iso, int32(0))
				return val
			},
			[]byte("0"),
		},
		{
			"object",
			func() *v8go.Value {
				ctx, _ := v8go.NewContext(iso)
				val, _ := ctx.RunScript("let foo = {a:1, b:2}; foo", "test.js")
				return val
			},
			[]byte(`{"a":1,"b":2}`),
		},
		{
			"objectFunc",
			func() *v8go.Value {
				ctx, _ := v8go.NewContext(iso)
				val, _ := ctx.RunScript("let foo = {a:1, b:()=>{}}; foo", "test.js")
				return val
			},
			[]byte(`{"a":1}`),
		},
		{
			"nil",
			func() *v8go.Value { return nil },
			[]byte(""),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val := tt.val()
			json, _ := val.MarshalJSON()
			if !bytes.Equal(json, tt.expected) {
				t.Errorf("unexpected JSON value: %s", string(json))
			}

		})
	}
}

package v8go_test

import (
	"reflect"
	"runtime"
	"testing"

	"rogchap.com/v8go"
)

func TestValueString(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
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
			t.Parallel()
			result, _ := ctx.RunScript(tt.source, "test.js")
			str := result.String()
			if str != tt.out {
				t.Errorf("unexpected result: expected %q, got %q", tt.out, str)
			}
		})
	}
}

func TestValueBoolean(t *testing.T) {
	t.Parallel()
	ctx, _ := v8go.NewContext(nil)
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
			t.Parallel()
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
			t.Parallel()
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
	tests := [...]struct {
		source   string
		expected int32
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"'1'", 1},
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
			t.Parallel()
			val, _ := ctx.RunScript(tt.source, "test.js")
			if i32 := val.Int32(); i32 != tt.expected {
				t.Errorf("unexpected value: expected %v, got %v", tt.expected, i32)
			}
		})
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

package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// Value represents all Javascript values and objects
type Value struct {
	ptr C.ValuePtr
}

// ArrayIndex attempts to converts a string to an array index. Returns ok false if conversion fails.
func (v *Value) ArrayIndex() (idx uint32, ok bool) {
	arrayIdx := C.ValueToArrayIndex(v.ptr)
	defer C.free(unsafe.Pointer(arrayIdx))
	if arrayIdx == nil {
		return 0, false
	}
	return uint32(*arrayIdx), true
}

// BigInt perform the equivalent of `BigInt(value)` in JS.
func (v *Value) BigInt() struct{} { // *BigInt
	panic("not implemented")
}

// Boolean perform the equivalent of `Boolean(value)` in JS. This can never fail.
func (v *Value) Boolean() bool {
	return C.ValueToBoolean(v.ptr) != 0
}

// DetailString provide a string representation of this value usable for debugging.
func (v *Value) DetailString() string {
	panic("not implemented")
}

// Int32 perform the equivalent of `Number(value)` in JS and convert the result to a
// signed 32-bit integer by performing the steps in https://tc39.es/ecma262/#sec-toint32.
func (v *Value) Int32() int32 {
	return int32(C.ValueToInt32(v.ptr))
}

// Integer perform the equivalent of `Number(value)` in JS and convert the result to an integer.
// Negative values are rounded up, positive values are rounded down. NaN is converted to 0.
// Infinite values yield undefined results.
func (v *Value) Integer() int64 {
	return int64(C.ValueToInteger(v.ptr))
}

// Number perform the equivalent of `Number(value)` in JS.
func (v *Value) Number() float64 {
	panic("not implemented")
}

func (v *Value) Object() struct{} { // *Object
	panic("not implemented")
}

// String perform the equivalent of `String(value)` in JS. Primitive values
// are returned as-is, objects will return `[object Object]` and functions will
// print their definition.
func (v *Value) String() string {
	s := C.ValueToString(v.ptr)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// Uint32 perform the equivalent of `Number(value)` in JS and convert the result to an
// unsigned 32-bit integer by performing the steps in https://tc39.es/ecma262/#sec-touint32.
func Uint32() uint32 {
	panic("not implemented")
}

// IsUndefined returns true if this value is the undefined value. See ECMA-262 4.3.10.
func (v *Value) IsUndefined() bool {
	return C.ValueIsUndefined(v.ptr) != 0
}

// IsNull returns true if this value is the null value. See ECMA-262 4.3.11.
func (v *Value) IsNull() bool {
	return C.ValueIsNull(v.ptr) != 0
}

// IsNullOrUndefined returns true if this value is either the null or the undefined value.
// See ECMA-262 4.3.11. and 4.3.12
// This is equivalent to `value == null` in JS.
func (v *Value) IsNullOrUndefined() bool {
	return C.ValueIsNullOrUndefined(v.ptr) != 0
}

// IsTrue returns true if this value is true.
// This is not the same as `BooleanValue()`. The latter performs a conversion to boolean,
// i.e. the result of `Boolean(value)` in JS, whereas this checks `value === true`.
func (v *Value) IsTrue() bool {
	return C.ValueIsTrue(v.ptr) != 0
}

// IsFalse returns true if this value is false.
// This is not the same as `!BooleanValue()`. The latter performs a conversion to boolean,
// i.e. the result of `!Boolean(value)` in JS, whereas this checks `value === false`.
func (v *Value) IsFalse() bool {
	return C.ValueIsFalse(v.ptr) != 0
}

// IsName returns true if this value is a symbol or a string.
// This is equivalent to `typeof value === 'string' || typeof value === 'symbol'` in JS.
func (v *Value) IsName() bool {
	return C.ValueIsName(v.ptr) != 0
}

// IsString returns true if this value is an instance of the String type. See ECMA-262 8.4.
// This is equivalent to `typeof value === 'string'` in JS.
func (v *Value) IsString() bool {
	return C.ValueIsString(v.ptr) != 0
}

// IsSymbol returns true if this value is a symbol.
// This is equivalent to `typeof value === 'symbol'` in JS.
func (v *Value) IsSymbol() bool {
	return C.ValueIsSymbol(v.ptr) != 0
}

// IsFunction returns true if this value is a function.
// This is equivalent to `typeof value === 'function'` in JS.
func (v *Value) IsFunction() bool {
	return C.ValueIsFunction(v.ptr) != 0
}

// IsObject returns true if this value is an object.
func (v *Value) IsObject() bool {
	return C.ValueIsObject(v.ptr) != 0
}

// IsBigInt returns true if this value is a bigint.
// This is equivalent to `typeof value === 'bigint'` in JS.
func (v *Value) IsBigInt() bool {
	return C.ValueIsBigInt(v.ptr) != 0
}

// IsBoolean returns true if this value is boolean.
// This is equivalent to `typeof value === 'boolean'` in JS.
func (v *Value) IsBoolean() bool {
	return C.ValueIsBoolean(v.ptr) != 0
}

// IsNumber returns true if this value is a number.
// This is equivalent to `typeof value === 'number'` in JS.
func (v *Value) IsNumber() bool {
	return C.ValueIsNumber(v.ptr) != 0
}

// IsExternal returns true if this value is an `External` object.
func (v *Value) IsExternal() bool {
	// TODO(rogchap): requires test case
	return C.ValueIsExternal(v.ptr) != 0
}

// IsInt32 returns true if this value is a 32-bit signed integer.
func (v *Value) IsInt32() bool {
	return C.ValueIsInt32(v.ptr) != 0
}

// IsUint32 returns true if this value is a 32-bit unsigned integer.
func (v *Value) IsUint32() bool {
	return C.ValueIsUint32(v.ptr) != 0
}

// IsDate returns true if this value is a `Date`.
func (v *Value) IsDate() bool {
	return C.ValueIsDate(v.ptr) != 0
}

// IsArgumentsObject returns true if this value is an Arguments object.
func (v *Value) IsArgumentsObject() bool {
	return C.ValueIsArgumentsObject(v.ptr) != 0
}

// IsBigIntObject returns true if this value is a BigInt object.
func (v *Value) IsBigIntObject() bool {
	return C.ValueIsBigIntObject(v.ptr) != 0
}

// IsNumberObject returns true if this value is a `Number` object.
func (v *Value) IsNumberObject() bool {
	return C.ValueIsNumberObject(v.ptr) != 0
}

// IsStringObject returns true if this value is a `String` object.
func (v *Value) IsStringObject() bool {
	return C.ValueIsStringObject(v.ptr) != 0
}

// IsSymbolObject returns true if this value is a `Symbol` object.
func (v *Value) IsSymbolObject() bool {
	return C.ValueIsSymbolObject(v.ptr) != 0
}

// IsNativeError returns true if this value is a NativeError.
func (v *Value) IsNativeError() bool {
	return C.ValueIsNativeError(v.ptr) != 0
}

// IsRegExp returns true if this value is a `RegExp`.
func (v *Value) IsRegExp() bool {
	return C.ValueIsRegExp(v.ptr) != 0
}

// IsAsyncFunc returns true if this value is an async function.
func (v *Value) IsAsyncFunction() bool {
	return C.ValueIsAsyncFunction(v.ptr) != 0
}

// Is IsGeneratorFunc returns true if this value is a Generator function.
func (v *Value) IsGeneratorFunction() bool {
	return C.ValueIsGeneratorFunction(v.ptr) != 0
}

// IsGeneratorObject returns true if this value is a Generator object (iterator).
func (v *Value) IsGeneratorObject() bool {
	return C.ValueIsGeneratorObject(v.ptr) != 0
}

// IsPromise returns true if this value is a `Promise`.
func (v *Value) IsPromise() bool {
	return C.ValueIsPromise(v.ptr) != 0
}

// IsMap returns true if this value is a `Map`.
func (v *Value) IsMap() bool {
	return C.ValueIsMap(v.ptr) != 0
}

// IsSet returns true if this value is a `Set`.
func (v *Value) IsSet() bool {
	return C.ValueIsSet(v.ptr) != 0
}

// IsMapIterator returns true if this value is a `Map` Iterator.
func (v *Value) IsMapIterator() bool {
	return C.ValueIsMapIterator(v.ptr) != 0
}

// IsSetIterator returns true if this value is a `Set` Iterator.
func (v *Value) IsSetIterator() bool {
	return C.ValueIsSetIterator(v.ptr) != 0
}

// IsWeakMap returns true if this value is a `WeakMap`.
func (v *Value) IsWeakMap() bool {
	return C.ValueIsWeakMap(v.ptr) != 0
}

// IsWeakSet returns true if this value is a `WeakSet`.
func (v *Value) IsWeakSet() bool {
	return C.ValueIsWeakSet(v.ptr) != 0
}

// IsArray returns true if this value is an array.
// Note that it will return false for a `Proxy` of an array.
func (v *Value) IsArray() bool {
	return C.ValueIsArray(v.ptr) != 0
}

// IsArrayBuffer returns true if this value is an `ArrayBuffer`.
func (v *Value) IsArrayBuffer() bool {
	return C.ValueIsArrayBuffer(v.ptr) != 0
}

// IsArrayBufferView returns true if this value is an `ArrayBufferView`.
func (v *Value) IsArrayBufferView() bool {
	return C.ValueIsArrayBufferView(v.ptr) != 0
}

// IsTypedArray returns true if this value is one of TypedArrays.
func (v *Value) IsTypedArray() bool {
	return C.ValueIsTypedArray(v.ptr) != 0
}

// IsUint8Array returns true if this value is an `Uint8Array`.
func (v *Value) IsUint8Array() bool {
	return C.ValueIsUint8Array(v.ptr) != 0
}

// IsUint8ClampedArray returns true if this value is an `Uint8ClampedArray`.
func (v *Value) IsUint8ClampedArray() bool {
	return C.ValueIsUint8ClampedArray(v.ptr) != 0
}

// IsInt8Array returns true if this value is an `Int8Array`.
func (v *Value) IsInt8Array() bool {
	return C.ValueIsInt8Array(v.ptr) != 0
}

// IsUint16Array returns true if this value is an `Uint16Array`.
func (v *Value) IsUint16Array() bool {
	return C.ValueIsUint16Array(v.ptr) != 0
}

// IsInt16Array returns true if this value is an `Int16Array`.
func (v *Value) IsInt16Array() bool {
	return C.ValueIsInt16Array(v.ptr) != 0
}

// IsUint32Array returns true if this value is an `Uint32Array`.
func (v *Value) IsUint32Array() bool {
	return C.ValueIsUint32Array(v.ptr) != 0
}

// IsInt32Array returns true if this value is an `Int32Array`.
func (v *Value) IsInt32Array() bool {
	return C.ValueIsInt32Array(v.ptr) != 0
}

// IsFloat32Array returns true if this value is a `Float32Array`.
func (v *Value) IsFloat32Array() bool {
	return C.ValueIsFloat32Array(v.ptr) != 0
}

// IsFloat64Array returns true if this value is a `Float64Array`.
func (v *Value) IsFloat64Array() bool {
	return C.ValueIsFloat64Array(v.ptr) != 0
}

// IsBigInt64Array returns true if this value is a `BigInt64Array`.
func (v *Value) IsBigInt64Array() bool {
	return C.ValueIsBigInt64Array(v.ptr) != 0
}

// IsBigUint64Array returns true if this value is a BigUint64Array`.
func (v *Value) IsBigUint64Array() bool {
	return C.ValueIsBigUint64Array(v.ptr) != 0
}

// IsDataView returns true if this value is a `DataView`.
func (v *Value) IsDataView() bool {
	return C.ValueIsDataView(v.ptr) != 0
}

// IsSharedArrayBuffer returns true if this value is a `SharedArrayBuffer`.
func (v *Value) IsSharedArrayBuffer() bool {
	return C.ValueIsSharedArrayBuffer(v.ptr) != 0
}

// IsProxy returns true if this value is a JavaScript `Proxy`.
func (v *Value) IsProxy() bool {
	return C.ValueIsProxy(v.ptr) != 0
}

// IsWasmModuleObject returns true if this value is a `WasmModuleObject`.
func (v *Value) IsWasmModuleObject() bool {
	// TODO(rogchap): requires test case
	return C.ValueIsWasmModuleObject(v.ptr) != 0
}

// IsModuleNamespaceObject returns true if the value is a `Module` Namespace `Object`.
func (v *Value) IsModuleNamespaceObject() bool {
	// TODO(rogchap): requires test case
	return C.ValueIsModuleNamespaceObject(v.ptr) != 0
}

func (v *Value) finalizer() {
	C.ValueDispose(v.ptr)
	v.ptr = nil
	runtime.SetFinalizer(v, nil)
}

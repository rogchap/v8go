package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import (
	"io"
	"unsafe"
)

// Array represent javascript array
type Array struct {
	*Object
}

// NewArray is an similar to `new Array()` from javascript.
func NewArray(ctx *ExecContext, len int) *Array {
	val := &Value{C.NewArray(ctx.ptr, C.size_t(len)), ctx}
	return &Array{Object: &Object{Value: val}}
}

// NewArrayFromStrings creates new array from strings
func NewArrayFromValues(ctx *ExecContext, strs []Valuer) *Array {
	arr := NewArray(ctx, len(strs))
	for i, x := range strs {
		err := arr.SetIdx(uint32(i), x)
		if err != nil {
			panic(err.Error())
		}
	}
	return arr
}

// ArrayBuffer represents javascript array buffer. It is intended
// for work with raw data such as []byte slices.
type ArrayBuffer struct {
	*Value
	written int
}

var _ io.Writer = (*ArrayBuffer)(nil)

// NewArrayBuffer creates new ArrayBuffer.
func NewArrayBuffer(ctx *ExecContext, len int) *ArrayBuffer {
	val := &Value{
		ptr: C.NewArrayBuffer(ctx.ptr, C.size_t(len)),
		ctx: ctx,
	}
	return &ArrayBuffer{val, 0}
}

// NewArrayBufferFromBytes creates new ArrayBuffer prepopulated with data.
// Size of the buffer is determined based on slice size.
func NewArrayBufferFromBytes(ctx *ExecContext, v []byte) *ArrayBuffer {
	arr := NewArrayBuffer(ctx, len(v))
	_, _ = arr.Write(v)
	return arr
}

func (ab *ArrayBuffer) Len() int64 {
	return int64(C.ArrayBufferByteLength(ab.ptr))
}

func (ab *ArrayBuffer) Bytes() []byte {
	len := C.ArrayBufferByteLength(ab.ptr)
	cbytes := unsafe.Pointer(C.GetArrayBufferBytes(ab.ptr)) // points into BackingStore
	return C.GoBytes(cbytes, C.int(len))
}

// Write is io.Writer interface method. Allows to write on buffer and
// returns amount of written bytes.
func (ab *ArrayBuffer) Write(bytes []byte) (n int, err error) {
	cbytes := C.CBytes(bytes) //FIXME is there really no way to avoid this malloc+memcpy?
	defer C.free(cbytes)
	C.PutArrayBufferBytes(ab.ptr, 0, (*C.char)(cbytes), C.size_t(len(bytes)))
	ab.written += len(bytes)
	return ab.written, nil
}

// ArrayTypedUint8 is an javascript Uint8Array array.
type ArrayTypedUint8 struct {
	*Value
	buf *ArrayBuffer
}

// NewTypedUint8ArrayFromBuffer creates new typed array from buffer.
func NewTypedUint8ArrayFromBuffer(buf *ArrayBuffer) *ArrayTypedUint8 {
	val := &Value{
		ptr: C.NewTypedUint8ArrayFromBuffer(buf.ptr, C.size_t(buf.Len())),
		ctx: buf.ctx,
	}
	return &ArrayTypedUint8{val, buf}
}

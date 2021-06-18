package v8go

// #include <stdlib.h>
// #include "v8go.h"
import "C"
import "unsafe"

type ArrayBuffer struct {
	*Value
}

func NewArrayBuffer(ctx *Context, len int64) *ArrayBuffer {
	return &ArrayBuffer{&Value{C.NewArrayBuffer(ctx.iso.ptr, C.size_t(len)), ctx}}
}

func (ab *ArrayBuffer) ByteLength() int64 {
	return int64(C.ArrayBufferByteLength(ab.ptr))
}

func (ab *ArrayBuffer) GetBytes() []uint8 {
	len := C.ArrayBufferByteLength(ab.ptr)
	cbytes := unsafe.Pointer(C.GetArrayBufferBytes(ab.ptr)) // points into BackingStore
	return C.GoBytes(cbytes, C.int(len))
}

func (ab *ArrayBuffer) PutBytes(bytes []uint8) {
	cbytes := C.CBytes(bytes) //FIXME is there really no way to avoid this malloc+memcpy?
	defer C.free(cbytes)
	C.PutArrayBufferBytes(ab.ptr, 0, (*C.char)(cbytes), C.size_t(len(bytes)))
}

//go:build cgo

package rcl

/*
#include <stdlib.h>

#include <rcutils/types/uint8_array.h>
#include <rmw/serialized_message.h>
#include <rcl/types.h>
#include <rcl/allocator.h>
*/
import "C"
import (
	"unsafe"
)

type Allocator struct {
	v C.rcutils_allocator_t
}

func NewAllocator() Allocator {
	return Allocator{v: C.rcl_get_default_allocator()}
}

func (a *Allocator) New(size int) (*SerializedMessage, error) {
	v := (*C.rcl_serialized_message_t)(C.malloc(C.sizeof_rcl_serialized_message_t))
	*v = C.rmw_get_zero_initialized_serialized_message()

	// #define rmw_serialized_message_t -> rcutils_uint8_array_init.
	if rc := C.rcutils_uint8_array_init(v, C.size_t(size), &a.v); rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(v))
		return nil, RclError(rc)
	}
	v.buffer_length = C.size_t(size)

	return &SerializedMessage{v}, nil
}

type SerializedMessage struct {
	v *C.rcl_serialized_message_t
}

func (m *SerializedMessage) AsSlice() []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(m.v.buffer)), m.v.buffer_length)
}

func (m *SerializedMessage) ToSlice() []byte {
	buf := m.AsSlice()
	v := make([]byte, len(buf))
	copy(v, buf)
	return v
}

func (m *SerializedMessage) Close() error {
	if m.v == nil {
		return nil
	}

	// #define rmw_serialized_message_fini -> rcutils_uint8_array_fini.
	if rc := C.rcutils_uint8_array_fini(m.v); rc != C.RCL_RET_OK {
		return RclError(rc)
	}
	C.free(unsafe.Pointer(m.v))
	m.v = nil
	return nil
}

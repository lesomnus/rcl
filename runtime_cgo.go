//go:build cgo

package rcl

/*
#include <rcl/context.h>
#include <rcl/init_options.h>
#include <rcl/init.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Runtime struct {
	allocator Allocator
	context   *C.rcl_context_t
}

func (b RuntimeBuilder) Build(args ...string) (*Runtime, error) {
	b.evaluate()
	if args == nil {
		args = []string{"--ros-args"}
	} else if args[0] != "--ros-args" {
		args = append([]string{"--ros-args"}, args...)
	}

	opts := C.rcl_get_zero_initialized_init_options()
	if rc := C.rcl_init_options_init(&opts, b.Allocator.v); rc != C.RCL_RET_OK {
		return nil, RclError(rc)
	}
	defer C.rcl_init_options_fini(&opts)

	if b.DomainId != nil {
		if rc := C.rcl_init_options_set_domain_id(&opts, C.size_t(*b.DomainId)); rc != C.RCL_RET_OK {
			return nil, RclError(rc)
		}
	}

	context := (*C.rcl_context_t)(C.malloc(C.sizeof_rcl_context_t))
	*context = C.rcl_get_zero_initialized_context()

	args_c := []*C.char{}
	for _, a := range args {
		args_c = append(args_c, C.CString(a))
	}
	defer func() {
		for _, a := range args_c {
			C.free(unsafe.Pointer(a))
		}
	}()

	rc := C.rcl_init(C.int(len(args_c)), (**C.char)(unsafe.SliceData(args_c)), &opts, context)
	runtime.KeepAlive(args_c)
	if rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(context))
		return nil, RclError(rc)
	}

	return &Runtime{
		allocator: b.Allocator,
		context:   context,
	}, nil
}

func (r *Runtime) Close() error {
	if r.context == nil {
		return nil
	}

	C.rcl_shutdown(r.context)
	C.rcl_context_fini(r.context)
	C.free(unsafe.Pointer(r.context))
	r.context = nil
	return nil
}

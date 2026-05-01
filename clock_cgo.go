//go:build cgo

package rcl

/*
#include <rcl/time.h>
#include <rcl/timer.h>
*/
import "C"
import (
	"time"
	"unsafe"
)

type Clock struct {
	v *C.rcl_clock_t
}

func (c Clock) Now() (time.Duration, error) {
	var t C.rcl_time_point_value_t
	if rc := C.rcl_clock_get_now(c.v, &t); rc != C.RCL_RET_OK {
		return 0, RclError(rc)
	}
	return time.Duration(t), nil
}

func (c *Clock) Close() (err error) {
	if rc := C.rcl_clock_fini(c.v); rc != C.RCL_RET_OK {
		err = RclError(rc)
	}
	C.free(unsafe.Pointer(c.v))
	c.v = nil
	return err
}

func (b ClockBuilder) Build() (*Clock, error) {
	b.evaluate()

	v := (*C.rcl_clock_t)(C.malloc(C.sizeof_rcl_clock_t))
	C.rcl_clock_init(C.rcl_clock_type_t(b.Type), v, &b.Allocator.v)

	return &Clock{v}, nil
}

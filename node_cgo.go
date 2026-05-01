//go:build cgo

package rcl

/*
#include <rcl/node.h>
#include <rcl/node_options.h>
*/
import "C"
import (
	"unsafe"
)

type Node struct {
	r *Runtime
	v *C.rcl_node_t
}

func (r *Runtime) NewNode(name, ns string) (*Node, error) {
	v := (*C.rcl_node_t)(C.malloc(C.sizeof_rcl_node_t))
	*v = C.rcl_get_zero_initialized_node()

	name_c := C.CString(name)
	defer C.free(unsafe.Pointer(name_c))

	ns_c := C.CString(ns)
	defer C.free(unsafe.Pointer(ns_c))

	opts := C.rcl_node_get_default_options()
	if rc := C.rcl_node_init(v, name_c, ns_c, r.context, &opts); rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(v))
		return nil, RclError(rc)
	}

	return &Node{r, v}, nil
}

func (n *Node) Close() (err error) {
	if n.v == nil {
		return
	}

	rc := C.rcl_node_fini(n.v)
	if rc != C.RCL_RET_OK {
		err = RclError(rc)
	}

	C.free(unsafe.Pointer(n.v))
	n.v = nil
	return
}

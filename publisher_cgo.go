//go:build cgo

package rcl

/*
#include <stdlib.h>
#include <string.h>

#include <rmw/rmw.h>
#include <rcl/publisher.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type RawPublisher struct {
	n *Node
	v *C.rcl_publisher_t
}

func (n *Node) NewPublisher(topic string, typename Typename, opts ...Option) (*RawPublisher, error) {
	pkg_name, kind, name := typename.Split()
	if kind != "msg" {
		return nil, fmt.Errorf("expected a kind of \"msg\", got %q", kind)
	}

	pkg, err := DefaultTypeRegistry.LoadMessage(pkg_name)
	if err != nil {
		return nil, fmt.Errorf("load message type package: %w", err)
	}

	msg_ts, err := pkg.Get(name)
	if err != nil {
		return nil, fmt.Errorf("get message type support: %w", err)
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	opt := C.rcl_publisher_get_default_options()
	opt.allocator = n.r.allocator.v
	if o.qos != nil {
		o.qos.toCStruct(&opt.qos)
	}

	v := (*C.rcl_publisher_t)(C.malloc(C.sizeof_rcl_publisher_t))
	*v = C.rcl_get_zero_initialized_publisher()

	topic_c := C.CString(topic)
	defer C.free(unsafe.Pointer(topic_c))

	if rc := C.rcl_publisher_init(v, n.v, msg_ts.v, topic_c, &opt); rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(v))
		return nil, RclError(rc)
	}

	return &RawPublisher{n, v}, nil
}

func (p *RawPublisher) Send(data []byte) error {
	buf, err := p.n.r.allocator.From(data)
	if err != nil {
		return fmt.Errorf("alloc message: %w", err)
	}
	defer buf.Close()

	if rc := C.rcl_publish_serialized_message(p.v, buf.v, nil); rc != C.RCL_RET_OK {
		return RclError(rc)
	}

	return nil
}

func (p *RawPublisher) Close() (err error) {
	if p.n == nil {
		return nil
	}
	if rc := C.rcl_publisher_fini(p.v, p.n.v); rc != C.RCL_RET_OK {
		err = RclError(rc)
	}

	C.free(unsafe.Pointer(p.v))
	p.n = nil
	p.v = nil

	return err
}

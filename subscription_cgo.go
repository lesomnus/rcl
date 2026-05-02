//go:build cgo

package rcl

/*
#include <stdlib.h>
#include <string.h>

#include <rmw/rmw.h>
#include <rcl/subscription.h>
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type RawSubscription struct {
	n *Node
	v *C.rcl_subscription_t
}

func (n *Node) NewSubscription(topic string, typename Typename, opts ...Option) (*RawSubscription, error) {
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

	opt := C.rcl_subscription_get_default_options()
	opt.allocator = n.r.allocator.v
	if o.qos != nil {
		o.qos.toCStruct(&opt.qos)
	}

	v := (*C.rcl_subscription_t)(C.malloc(C.sizeof_rcl_subscription_t))
	*v = C.rcl_get_zero_initialized_subscription()

	topic_c := C.CString(topic)
	defer C.free(unsafe.Pointer(topic_c))

	if rc := C.rcl_subscription_init(v, n.v, msg_ts.v, topic_c, &opt); rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(v))
		return nil, RclError(rc)
	}

	return &RawSubscription{n, v}, nil
}

func (s RawSubscription) Recv() ([]byte, *MessageInfo, error) {
	buf, err := s.n.r.allocator.New(0)
	if err != nil {
		return nil, nil, fmt.Errorf("alloc message: %w", err)
	}
	defer buf.Close()

	info := C.rmw_get_zero_initialized_message_info()
	switch rc := C.rcl_take_serialized_message(s.v, buf.v, &info, nil); rc {
	case C.RCL_RET_OK:
		return buf.ToSlice(), &MessageInfo{
			SourceTimestamp:   time.Unix(0, int64(info.source_timestamp)),
			ReceivedTimestamp: time.Unix(0, int64(info.received_timestamp)),
			FromIntraProcess:  bool(info.from_intra_process),
		}, nil
	default:
		return nil, nil, RclError(rc)
	}
}

func (s *RawSubscription) Close() (err error) {
	if s.v == nil {
		return nil
	}

	if rc := C.rcl_subscription_fini(s.v, s.n.v); rc != C.RCL_RET_OK {
		err = RclError(rc)
	}

	C.free(unsafe.Pointer(s.v))
	s.v = nil

	return
}

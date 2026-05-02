//go:build cgo

package rcl

/*
#include <rcl/client.h>
#include <rcl/service.h>
#include <rcl/graph.h>
#include <rmw/rmw.h>
*/
import "C"
import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

// It converts Go struct into ROS message via CDR, so it may not be the most efficient way. But it's simple and works for now.
// TODO: To optimize it, we should implement a un/marshaller using introspection.
type RawClient struct {
	n *Node
	v *C.rcl_client_t
	t ServiceTypeSupport
}

func (n *Node) NewClient(service string, typename Typename, opts ...Option) (*RawClient, error) {
	pkg_name, kind, svc_name := typename.Split()
	if kind != "srv" {
		return nil, fmt.Errorf("expected a kind of \"srv\", got %q", kind)
	}

	pkg, err := DefaultTypeRegistry.LoadService(pkg_name)
	if err != nil {
		return nil, fmt.Errorf("load service type package: %w", err)
	}

	svc_ts, err := pkg.Get(svc_name)
	if err != nil {
		return nil, fmt.Errorf("get service type support: %w", err)
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	opt := C.rcl_client_get_default_options()
	opt.allocator = n.r.allocator.v
	if o.qos != nil {
		o.qos.toCStruct(&opt.qos)
	}

	v := (*C.rcl_client_t)(C.malloc(C.sizeof_rcl_client_t))
	*v = C.rcl_get_zero_initialized_client()

	service_c := C.CString(service)
	defer C.free(unsafe.Pointer(service_c))

	// ts.request_typesupport.

	if rc := C.rcl_client_init(v, n.v, svc_ts.v, service_c, &opt); rc != C.RCL_RET_OK {
		C.free(unsafe.Pointer(v))
		return nil, RclError(rc)
	}

	return &RawClient{n, v, svc_ts}, nil
}

func (c *RawClient) SendRequest(data []byte) (int64, error) {
	msg, err := c.t.Request().Create()
	if err != nil {
		return 0, fmt.Errorf("create request message: %w", err)
	}
	defer msg.Close()

	buf, err := c.n.r.allocator.From(data)
	if err != nil {
		return 0, fmt.Errorf("alloc buffer: %w", err)
	}
	defer buf.Close()

	if err := Deserialize(msg, buf); err != nil {
		return 0, fmt.Errorf("deserialize request: %w", err)
	}

	var seq C.int64_t
	if rc := C.rcl_send_request(c.v, msg.v, &seq); rc != C.RCL_RET_OK {
		return 0, RclError(rc)
	}

	return int64(seq), nil
}

func (c *RawClient) TakeResponse() ([]byte, ServiceInfo, error) {
	msg, err := c.t.Response().Create()
	if err != nil {
		return nil, ServiceInfo{}, fmt.Errorf("create response message: %w", err)
	}
	defer msg.Close()

	var header C.rmw_service_info_t

	switch rc := C.rcl_take_response_with_info(c.v, &header, msg.v); rc {
	case C.RCL_RET_OK:
		buf, err := c.n.r.allocator.New(0)
		if err != nil {
			return nil, ServiceInfo{}, fmt.Errorf("alloc serialized message: %w", err)
		}
		defer buf.Close()

		if err := Serialize(buf, msg); err != nil {
			return nil, ServiceInfo{}, fmt.Errorf("serialize response: %w", err)
		}

		return buf.ToSlice(), ServiceInfo{
			SourceTimestamp:   time.Unix(0, int64(header.source_timestamp)),
			ReceivedTimestamp: time.Unix(0, int64(header.received_timestamp)),
			RequestId: RequestId{
				WriterGuid:     *(*[16]int8)(unsafe.Pointer(&header.request_id.writer_guid)),
				SequenceNumber: int64(header.request_id.sequence_number),
			},
		}, nil
	default:
		return nil, ServiceInfo{}, RclError(rc)
	}
}

// Wait waits for the service to be available.
func (c *RawClient) Wait(ctx context.Context) error {
	// TODO: use rcl_wait_set_t to wait for the service to be available.
	// For now, polling.
	var ok C.bool
	for !ok {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			C.rcl_service_server_is_available(c.n.v, c.v, &ok)
		}
	}

	return nil
}

func (c *RawClient) Close() (err error) {
	if c.v == nil {
		return nil
	}

	if rc := C.rcl_client_fini(c.v, c.n.v); rc != C.RCL_RET_OK {
		return RclError(rc)
	}

	C.free(unsafe.Pointer(c.v))
	c.v = nil

	return err
}

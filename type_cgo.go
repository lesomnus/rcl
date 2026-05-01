//go:build cgo

package rcl

/*
#include <stdlib.h>
#include <rmw/rmw.h>

#include "type.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type typePackage struct {
	name   string
	lib_g  unsafe.Pointer
	lib_ts unsafe.Pointer
	lib_ti unsafe.Pointer
}

type MessageTypePackage typePackage
type ServiceTypePackage typePackage

func (r TypeLoader) load(name string) (typePackage, error) {
	v, ok := r[name]
	if ok {
		return v, nil
	}

	v, err := loadTypePackage(name)
	if err != nil {
		return typePackage{}, err
	}

	r[name] = v
	return v, nil
}

func (r TypeLoader) LoadMessage(name string) (MessageTypePackage, error) {
	v, err := r.load(name)
	return MessageTypePackage(v), err
}

func (r TypeLoader) LoadService(name string) (ServiceTypePackage, error) {
	v, err := r.load(name)
	return ServiceTypePackage(v), err
}

func loadTypePackage(name string) (typePackage, error) {
	g, err := dlopen(fmt.Sprintf("lib%s__rosidl_generator_c.so", name))
	if err != nil {
		return typePackage{}, err
	}

	ts, err := dlopen(fmt.Sprintf("lib%s__rosidl_typesupport_c.so", name))
	if err != nil {
		return typePackage{}, err
	}

	ti, err := dlopen(fmt.Sprintf("lib%s__rosidl_typesupport_introspection_c.so", name))
	if err != nil {
		return typePackage{}, err
	}

	return typePackage{
		name:   name,
		lib_g:  g,
		lib_ts: ts,
		lib_ti: ti,
	}, nil
}

type MessageTypeSupport struct {
	lib_g unsafe.Pointer
	v     *C.rosidl_message_type_support_t

	// E.g.:
	//
	//	std_msgs__msg__
	//	std_srvs__srv__
	prefix string

	// E.g.:
	//
	//	SetBool
	name string

	// E.g.:
	//
	//	UInt8__Sequence__
	//	Empty_Request__
	suffix string
}

func (p MessageTypePackage) Get(name string) (MessageTypeSupport, error) {
	sym := fmt.Sprintf("rosidl_typesupport_c__get_message_type_support_handle__%s__msg__%s", p.name, name)

	sym_c := C.CString(sym)
	defer C.free(unsafe.Pointer(sym_c))

	var err [1024]C.char
	ts := C.rclgo_get_message_type_support(p.lib_ts, sym_c, &err[0], C.size_t(len(err)))
	if ts == nil {
		return MessageTypeSupport{}, errors.New(C.GoString(&err[0]))
	}

	return MessageTypeSupport{
		lib_g: p.lib_g,
		v:     ts,

		prefix: p.name + "__msg__",
		name:   name,
		suffix: "__",
	}, nil
}

type ServiceTypeSupport struct {
	pkg  ServiceTypePackage
	name string
	v    *C.rosidl_service_type_support_t
}

func (p ServiceTypePackage) Get(name string) (ServiceTypeSupport, error) {
	sym := fmt.Sprintf("rosidl_typesupport_c__get_service_type_support_handle__%s__srv__%s", p.name, name)

	sym_c := C.CString(sym)
	defer C.free(unsafe.Pointer(sym_c))

	var err [1024]C.char
	ts := C.rclgo_get_service_type_support(p.lib_ts, sym_c, &err[0], C.size_t(len(err)))
	if ts == nil {
		return ServiceTypeSupport{}, errors.New(C.GoString(&err[0]))
	}

	return ServiceTypeSupport{
		pkg:  p,
		name: name,
		v:    ts,
	}, nil
}

// std_msgs__msg__UInt8__Sequence__create
// std_srvs__srv__Empty_Request__create

func (s ServiceTypeSupport) Request() MessageTypeSupport {
	return MessageTypeSupport{
		lib_g: s.pkg.lib_g,
		v:     s.v.request_typesupport,

		prefix: s.pkg.name + "__srv__",
		name:   s.name,
		suffix: "_Request__",
	}
}

func (s ServiceTypeSupport) Response() MessageTypeSupport {
	return MessageTypeSupport{
		lib_g: s.pkg.lib_g,
		v:     s.v.response_typesupport,

		prefix: s.pkg.name + "__srv__",
		name:   s.name,
		suffix: "_Response__",
	}
}

type Message struct {
	ts MessageTypeSupport
	v  unsafe.Pointer
}

func (m Message) Close() error {
	return m.ts.Destroy(m.v)
}

func (t MessageTypeSupport) Create() (Message, error) {
	name := t.prefix + t.name + t.suffix + "create"
	name_c := C.CString(name)
	defer C.free(unsafe.Pointer(name_c))

	var err [1024]C.char
	msg := C.rclgo_msg_create(t.lib_g, name_c, &err[0], C.size_t(len(err)))
	if msg == nil {
		return Message{}, errors.New(C.GoString(&err[0]))
	}

	return Message{ts: t, v: msg}, nil
}

func (t MessageTypeSupport) Destroy(msg unsafe.Pointer) error {
	name := t.prefix + t.name + t.suffix + "destroy"
	name_c := C.CString(name)
	defer C.free(unsafe.Pointer(name_c))

	var err [1024]C.char
	C.rclgo_msg_destroy(t.lib_g, name_c, msg, &err[0], C.size_t(len(err)))
	if msg == nil {
		return errors.New(C.GoString(&err[0]))
	}

	return nil
}

func Serialize(buff *SerializedMessage, msg Message) error {
	if rt := C.rmw_serialize(msg.v, msg.ts.v, buff.v); rt != C.RMW_RET_OK {
		return RclError(rt)
	}
	return nil
}

func Deserialize(msg Message, buff *SerializedMessage) error {
	if rt := C.rmw_deserialize(buff.v, msg.ts.v, msg.v); rt != C.RMW_RET_OK {
		return RclError(rt)
	}
	return nil
}

func dlopen(name string) (unsafe.Pointer, error) {
	name_c := C.CString(name)
	defer C.free(unsafe.Pointer(name_c))

	var err [1024]C.char
	l := C.rclgo_dlopen(name_c, &err[0], C.size_t(len(err)))
	if l == nil {
		return nil, errors.New(C.GoString(&err[0]))
	}
	return l, nil
}

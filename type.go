package rcl

import (
	"strings"
	"time"
)

type Typename string

func (t Typename) Split() (pkg, kind, name string) {
	parts := strings.SplitN(string(t), "/", 3)
	if len(parts) > 0 {
		pkg = parts[0]
	}
	if len(parts) > 1 {
		kind = parts[1]
	}
	if len(parts) > 2 {
		name = parts[2]
	}
	return
}

func (t Typename) Package() string {
	pkg, _, _ := t.Split()
	return pkg
}

func (t Typename) Kind() string {
	_, kind, _ := t.Split()
	return kind
}

func (t Typename) Name() string {
	_, _, name := t.Split()
	return name
}

type TypeRegistry interface {
	LoadMessage(name string) (MessageTypePackage, error)
	LoadService(name string) (ServiceTypePackage, error)
}

type TypeLoader map[string]typePackage

var DefaultTypeRegistry TypeRegistry = TypeLoader{}

func LoadMessage(name string) (MessageTypePackage, error) {
	return DefaultTypeRegistry.LoadMessage(name)
}

func LoadService(name string) (ServiceTypePackage, error) {
	return DefaultTypeRegistry.LoadService(name)
}

type MessageInfo struct {
	SourceTimestamp   time.Time
	ReceivedTimestamp time.Time
	FromIntraProcess  bool
}

type ServiceInfo struct {
	SourceTimestamp   time.Time
	ReceivedTimestamp time.Time
	RequestId         RequestId
}

type RequestId struct {
	WriterGuid     [16]int8
	SequenceNumber int64
}

//go:build !cgo

package rcl

type Allocator struct{}

func NewAllocator() Allocator {
	return Allocator{}
}

func (a Allocator) New(size int) (*SerializedMessage, error) {
	return &SerializedMessage{buf: make([]byte, size)}, nil
}

type SerializedMessage struct {
	buf []byte
}

func (m *SerializedMessage) AsSlice() []byte {
	return m.buf
}

func (m *SerializedMessage) ToSlice() []byte {
	buf := make([]byte, len(m.buf))
	copy(buf, m.buf)
	return buf
}

func (m *SerializedMessage) Close() error {
	m.buf = nil
	return nil
}

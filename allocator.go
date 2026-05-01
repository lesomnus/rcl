package rcl

var DefaultAllocator = NewAllocator()

func (a Allocator) isZero() bool {
	return a == (Allocator{})
}

func (a *Allocator) From(data []byte) (*SerializedMessage, error) {
	m, err := a.New(len(data))
	if err != nil {
		return nil, err
	}
	copy(m.AsSlice(), data)
	return m, nil
}

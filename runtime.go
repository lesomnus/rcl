package rcl

func NewRuntime(args ...string) (*Runtime, error) {
	return (&RuntimeBuilder{}).Build()
}

type RuntimeBuilder struct {
	Allocator Allocator
	DomainId  *int
}

func (b *RuntimeBuilder) evaluate() error {
	if b.Allocator.isZero() {
		b.Allocator = DefaultAllocator
	}
	return nil
}

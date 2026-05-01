//go:build !cgo

package rcl

type Runtime struct {
	allocator Allocator
}

func (b RuntimeBuilder) Build(args ...string) (*Runtime, error) {
	b.evaluate()
	return &Runtime{allocator: b.Allocator}, nil
}

func (r *Runtime) Close() error {
	return nil
}

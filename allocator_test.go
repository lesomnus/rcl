package rcl_test

import (
	"testing"

	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

func TestAllocator(t *testing.T) {
	asan.NoLeak(t)

	a := rcl.NewAllocator()

	b, err := a.New(42)
	x.NoError(t, err)

	err = b.Close()
	x.NoError(t, err)
}

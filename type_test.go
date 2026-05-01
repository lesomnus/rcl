package rcl_test

import (
	"testing"

	"github.com/lesomnus/cdr"
	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

func TestMessage(t *testing.T) {
	asan.NoLeak(t)

	msg_pkg, err := rcl.DefaultTypeRegistry.LoadMessage("geometry_msgs")
	x.NoError(t, err)

	md, err := msg_pkg.Get("Point")
	x.NoError(t, err)

	m, err := md.Create()
	x.NoError(t, err)
	defer m.Close()

	type Point struct {
		X float64
		Y float64
		Z float64
	}
	data, err := cdr.Marshal(Point{X: 1, Y: 2, Z: 3})
	x.NoError(t, err)

	buf, err := rcl.DefaultAllocator.From(data)
	x.NoError(t, err)
	defer buf.Close()

	err = rcl.Deserialize(m, buf)
	x.NoError(t, err)

	err = rcl.Serialize(buf, m)
	x.NoError(t, err)
	x.Eq(t, data, buf.AsSlice())
}

package rcl_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

func TestSubscription(t *testing.T) {
	const Topic = "/rclgo/test/publisher"
	const Type = "geometry_msgs/msg/Point"

	type Point struct {
		X float64
		Y float64
		Z float64
	}

	type Args struct {
		v *Point
		e error
	}

	t.Run("tick with no message", func(t *testing.T) {
		asan.NoLeak(t)

		r, err := rcl.NewRuntime()
		x.NoError(t, err)
		defer r.Close()

		node, err := r.NewNode("foo", "bar")
		x.NoError(t, err)
		defer node.Close()

		args := []Args{}
		sub, err := rcl.Subscribe(node, Topic, Type, func(v *Point, info *rcl.MessageInfo, err error) {
			args = append(args, Args{v, err})
		})
		x.NoError(t, err)
		defer sub.Close()

		sub.Tick(0)
		sub.Tick(0)
		sub.Tick(0)
		x.Len(t, args, 3)
		x.Eq(t, Args{nil, rcl.RclError(rcl.CodeSubscriptionTakeFailed)}, args[0])
		x.Eq(t, Args{nil, rcl.RclError(rcl.CodeSubscriptionTakeFailed)}, args[1])
		x.Eq(t, Args{nil, rcl.RclError(rcl.CodeSubscriptionTakeFailed)}, args[2])
	})
	t.Run("tick with a message", func(t *testing.T) {
		asan.NoLeak(t)

		r, err := rcl.NewRuntime()
		x.NoError(t, err)
		defer r.Close()

		node, err := r.NewNode("foo", "bar")
		x.NoError(t, err)
		defer node.Close()

		args := []Args{}
		sub, err := rcl.Subscribe(node, Topic, Type, func(v *Point, info *rcl.MessageInfo, err error) {
			args = append(args, Args{v, err})
		})
		x.NoError(t, err)
		defer sub.Close()

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		err = exec.CommandContext(ctx, "ros2", "topic", "pub", "--once", Topic, Type, "{x: 1.0, y: 2.0, z: 3.0}").Run()
		x.NoError(t, err)
		err = exec.CommandContext(ctx, "ros2", "topic", "pub", "--once", Topic, Type, "{x: 4.0, y: 5.0, z: 6.0}").Run()
		x.NoError(t, err)
		err = exec.CommandContext(ctx, "ros2", "topic", "pub", "--once", Topic, Type, "{x: 7.0, y: 8.0, z: 9.0}").Run()
		x.NoError(t, err)

		sub.Tick(0)
		sub.Tick(0)
		sub.Tick(0)
		x.Len(t, args, 3)
		x.Eq(t, Args{&Point{X: 1, Y: 2, Z: 3}, nil}, args[0])
		x.Eq(t, Args{&Point{X: 4, Y: 5, Z: 6}, nil}, args[1])
		x.Eq(t, Args{&Point{X: 7, Y: 8, Z: 9}, nil}, args[2])
	})
}

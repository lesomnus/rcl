package rcl_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

func TestPublisher(t *testing.T) {
	asan.NoLeak(t)

	const Topic = "/rclgo/test/publisher"
	const Type = "geometry_msgs/msg/Point"

	r, err := rcl.NewRuntime()
	x.NoError(t, err)
	defer r.Close()

	node, err := r.NewNode("foo", "bar")
	x.NoError(t, err)
	defer node.Close()

	type Point struct {
		X float64
		Y float64
		Z float64
	}
	pub, err := rcl.NewPublisher[Point](node, Topic, Type)
	x.NoError(t, err)
	defer pub.Close()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ros2", "topic", "echo", "--once", Topic, Type)
	stdout, err := cmd.StdoutPipe()
	x.NoError(t, err)

	err = cmd.Start()
	x.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	err = pub.Send(&Point{41, 42, 43})
	x.NoError(t, err)

	output, err := io.ReadAll(stdout)
	x.NoError(t, err)

	err = cmd.Wait()
	x.NoError(t, err)
	x.Eq(t, "x: 41.0\ny: 42.0\nz: 43.0\n---\n", string(output))
}

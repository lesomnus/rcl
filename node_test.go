package rcl_test

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

func NewNode(t *testing.T) *rcl.Node {
	r, err := rcl.NewRuntime()
	x.NoError(t, err)
	t.Cleanup(func() { r.Close() })

	node, err := r.NewNode("foo", "bar")
	x.NoError(t, err)
	t.Cleanup(func() { node.Close() })

	return node
}

func TestNode(t *testing.T) {
	asan.NoLeak(t)

	_ = NewNode(t)

	time.Sleep(100 * time.Millisecond)

	output, err := exec.CommandContext(t.Context(), "ros2", "node", "list").Output()
	x.NoError(t, err)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "/bar/foo" {
			return
		}
	}

	t.Fatalf("node not found")
}

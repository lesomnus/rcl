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

func TestNode(t *testing.T) {
	asan.NoLeak(t)

	r, err := rcl.NewRuntime()
	x.NoError(t, err)
	defer r.Close()

	node, err := r.NewNode("rclgo_name", "rclgo_ns")
	x.NoError(t, err)
	defer node.Close()

	time.Sleep(100 * time.Millisecond)

	output, err := exec.CommandContext(t.Context(), "ros2", "node", "list").Output()
	x.NoError(t, err)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "/rclgo_ns/rclgo_name" {
			return
		}
	}

	t.Fatalf("node not found")
}

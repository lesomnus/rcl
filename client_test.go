package rcl_test

import (
	"context"
	_ "embed"
	"os/exec"
	"testing"
	"time"

	"github.com/lesomnus/rcl"
	"github.com/lesomnus/rcl/internal/asan"
	"github.com/lesomnus/rcl/internal/x"
)

//go:embed scripts/service.py
var service_server_py string

func TestClient(t *testing.T) {
	const Topic = "/rclgo/test/add"
	const Type = "example_interfaces/srv/AddTwoInts"

	type Req struct {
		A int64
		B int64
	}
	type Res struct {
		Sum int64
	}
	type Args struct {
		res *Res
		err error
	}

	t.Run("call with no server", func(t *testing.T) {
		asan.NoLeak(t)

		node := NewNode(t)

		client, err := rcl.NewClient[Req, Res](node, Topic, Type)
		x.NoError(t, err)
		defer client.Close()

		args := []Args{}
		cancel, err := client.Call(&Req{A: 1, B: 2}, func(res *Res, info rcl.ServiceInfo, err error) {
			args = append(args, Args{res, err})
		})
		x.NoError(t, err)
		defer cancel()

		client.Tick(0)
		client.Tick(0)
		client.Tick(0)
		x.Len(t, args, 0)
	})
	t.Run("call with server", func(t *testing.T) {
		asan.NoLeak(t)

		server := exec.CommandContext(t.Context(), "python3", "-c", service_server_py, Topic)
		err := server.Start()
		x.NoError(t, err)
		defer server.Wait()
		defer server.Process.Kill()

		time.Sleep(2 * time.Second)

		node := NewNode(t)

		client, err := rcl.NewClient[Req, Res](node, Topic, Type)
		x.NoError(t, err)
		defer client.Close()

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		err = client.Wait(ctx)
		x.NoError(t, err)

		go func() {
			tick := time.Tick(100 * time.Millisecond)
			for {
				select {
				case <-ctx.Done():
					return
				case <-tick:
					client.Tick(0)
				}
			}
		}()

		res, _, err := rcl.Call(ctx, client, &Req{A: 1, B: 2})
		x.NoError(t, err)
		x.Eq(t, &Res{Sum: 3}, res)
	})
}

package rcl

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lesomnus/cdr"
)

type Client[Req, Res any] interface {
	Call(req *Req, f func(res *Res, info ServiceInfo, err error)) (cancel func(), err error)
	Wait(ctx context.Context) error
	TickCloser
}

type client[Req, Res any] struct {
	r *RawClient

	mu sync.Mutex
	cs map[int64]clientCallback[Res]

	ctx    context.Context
	cancel context.CancelFunc
}

type clientCallback[Res any] func(*Res, ServiceInfo, error)

func NewClient[Req, Res any](n *Node, service string, typename Typename, opts ...Option) (Client[Req, Res], error) {
	raw, err := n.NewClient(service, typename, opts...)
	if err != nil {
		return nil, fmt.Errorf("create raw client: %w", err)
	}

	c := &client[Req, Res]{
		r:  raw,
		cs: make(map[int64]clientCallback[Res]),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	return c, nil
}

func Call[Req, Res any](ctx context.Context, c Client[Req, Res], req *Req) (*Res, ServiceInfo, error) {
	done := make(chan struct{})

	var (
		cb_res  *Res
		cb_info ServiceInfo
		cb_err  error
	)
	cancel, err := c.Call(req, func(res_ *Res, info_ ServiceInfo, err_ error) {
		cb_res = res_
		cb_info = info_
		cb_err = err_
		close(done)
	})
	if err != nil {
		return nil, ServiceInfo{}, err
	}
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ServiceInfo{}, ctx.Err()
	case <-done:
		return cb_res, cb_info, cb_err
	}
}

func (c *client[Req, Res]) Call(req *Req, f func(res *Res, info ServiceInfo, err error)) (func(), error) {
	data, err := cdr.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var done atomic.Bool

	stop := context.AfterFunc(c.ctx, func() {
		if done.Swap(true) {
			return
		}
		f(nil, ServiceInfo{}, fmt.Errorf("client closed"))
	})
	cb := clientCallback[Res](func(res *Res, info ServiceInfo, err error) {
		if done.Swap(true) {
			return
		}
		stop()
		f(res, info, err)
	})
	if _, err := (func() (int64, error) {
		c.mu.Lock()
		defer c.mu.Unlock()

		seq, err := c.r.SendRequest(data)
		if err != nil {
			return 0, fmt.Errorf("send request: %w", err)
		}

		c.cs[int64(seq)] = cb
		return int64(seq), nil
	})(); err != nil {
		return nil, err
	}

	cancel := func() {
		if done.Swap(true) {
			return
		}
		stop()
		f(nil, ServiceInfo{}, fmt.Errorf("canceled"))
	}

	return cancel, nil
}

func (c *client[Req, Res]) Tick(int) {
	buf, info, err := c.r.TakeResponse()
	if err != nil {
		return
	}

	cb, ok := c.fetchCall(info.RequestId.SequenceNumber)
	if !ok {
		return
	}

	var v Res
	if err := cdr.Unmarshal(buf, &v); err != nil {
		cb(nil, ServiceInfo{}, fmt.Errorf("unmarshal response: %w", err))
	} else {
		cb(&v, info, nil)
	}
}

func (c *client[Req, Res]) fetchCall(seq int64) (clientCallback[Res], bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.cs)
	if n == 0 {
		return nil, false
	}

	s, ok := c.cs[int64(seq)]
	if !ok {
		// Seems that the call is canceled.
		return nil, false
	}
	delete(c.cs, int64(seq))

	return s, true
}

func (c *client[Req, Res]) Wait(ctx context.Context) error {
	return c.r.Wait(ctx)
}

func (c *client[Req, Res]) Close() error {
	c.cancel()
	return c.r.Close()
}

package rcl

import (
	"fmt"

	"github.com/lesomnus/cdr"
)

type Subscription interface {
	TickCloser
}

type subscription[T any] struct {
	r *RawSubscription
	f func(v *T, info *MessageInfo, err error)
}

func Subscribe[T any](n *Node, topic string, typename Typename, f func(v *T, info *MessageInfo, err error), opts ...Option) (Subscription, error) {
	raw, err := n.NewSubscription(topic, typename, opts...)
	if err != nil {
		return subscription[T]{}, fmt.Errorf("create raw subscription: %w", err)
	}

	return subscription[T]{raw, f}, nil
}

func (s subscription[T]) Tick(int) {
	data, info, err := s.r.Recv()
	if err != nil {
		s.f(nil, nil, err)
		return
	}

	var v T
	if err := cdr.Unmarshal(data, &v); err != nil {
		s.f(nil, nil, err)
		return
	}

	s.f(&v, info, nil)
}

func (s subscription[T]) Close() error {
	return s.r.Close()
}

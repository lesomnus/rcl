package rcl

import (
	"fmt"

	"github.com/lesomnus/cdr"
)

type Publisher[T any] interface {
	Send(value *T) error
	Close() error
}

type publisher[T any] struct {
	*RawPublisher
}

func NewPublisher[T any](n *Node, topic string, typename Typename, opts ...Option) (Publisher[T], error) {
	raw, err := n.NewPublisher(topic, typename, opts...)
	if err != nil {
		return publisher[T]{}, fmt.Errorf("create raw publisher: %w", err)
	}
	return publisher[T]{raw}, nil
}

func (p publisher[T]) Send(value *T) error {
	data, err := cdr.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	return p.RawPublisher.Send(data)
}

func (p publisher[T]) Close() error {
	return p.RawPublisher.Close()
}

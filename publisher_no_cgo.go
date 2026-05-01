//go:build !cgo

package rcl

type RawPublisher struct{}

func (n *Node) NewPublisher(topic string, typename Typename, opts ...Option) (*RawPublisher, error) {
	return &RawPublisher{}, nil
}

func (p RawPublisher) Send(data []byte) error {
	return nil
}

func (p RawPublisher) Close() error {
	return nil
}

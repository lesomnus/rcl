//go:build !cgo

package rcl

type RawSubscription struct{}

func (n *Node) NewSubscription(topic string, typename Typename, opts ...Option) (*RawSubscription, error) {
	return &RawSubscription{}, nil
}

func (s RawSubscription) Recv() ([]byte, *MessageInfo, error) {
	return nil, nil, RclError(CodeSubscriptionTakeFailed)
}

func (s RawSubscription) Close() error {
	return nil
}

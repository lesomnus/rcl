//go:build !cgo

package rcl

import "context"

type RawClient struct{}

func (n *Node) NewClient(service string, typename Typename, opts ...Option) (*RawClient, error) {
	return &RawClient{}, nil
}

func (c RawClient) SendRequest(data []byte) (int64, error) {
	return 0, nil
}

func (c RawClient) TakeResponse() ([]byte, ServiceInfo, error) {
	return nil, ServiceInfo{}, nil
}

func (c RawClient) Wait(ctx context.Context) error {
	return ctx.Err()
}

func (c RawClient) Close() error {
	return nil
}

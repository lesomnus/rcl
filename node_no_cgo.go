//go:build !cgo

package rcl

type Node struct {
	r *Runtime
}

func (r *Runtime) NewNode(name, ns string) (*Node, error) {
	return &Node{r}, nil
}

func (n *Node) Close() error {
	return nil
}

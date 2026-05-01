//go:build !cgo

package rcl

import "time"

type Clock struct {
	v any
}

func (c Clock) Now() (time.Duration, error) {
	t := time.Now()
	return time.Duration(t.UnixNano()), nil
}

func (c Clock) Close() error {
	return nil
}

func (b ClockBuilder) Build() (*Clock, error) {
	b.evaluate()
	return &Clock{}, nil
}

package rcl

type Ticker interface {
	// Tick process ROS message queue once.
	// Argument `x` is not used but only for interface compatibility.
	Tick(x int)
}

type TickCloser interface {
	Ticker
	Close() error
}

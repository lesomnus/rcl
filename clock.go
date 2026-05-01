package rcl

var DefaultClock, _ = ClockBuilder{Allocator: NewAllocator()}.Build()

type ClockType uint32

const (
	ClockTypeUninitialized ClockType = 0
	ClockTypeROSTime       ClockType = 1
	ClockTypeSystemTime    ClockType = 2
	ClockTypeSteadyTime    ClockType = 3
)

func NewClock() (*Clock, error) {
	return ClockBuilder{}.Build()
}

type ClockBuilder struct {
	Type      ClockType
	Allocator Allocator
}

func (b ClockBuilder) evaluate() {
	if b.Type == ClockTypeUninitialized {
		b.Type = ClockTypeROSTime
	}
	if b.Allocator.isZero() {
		b.Allocator = DefaultAllocator
	}
}

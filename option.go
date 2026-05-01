package rcl

type Options struct {
	qos *Qos
}

type Option func(*Options)

func WithQos(qos Qos) Option {
	return func(o *Options) {
		o.qos = &qos
	}
}

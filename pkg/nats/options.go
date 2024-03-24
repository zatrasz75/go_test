package nats

import "time"

// Option -.
type Option func(*Nats)

func OptionSet(allow bool, size int, wait, timeout time.Duration) Option {
	return func(n *Nats) {
		AllowReconnect(allow)(n)
		MaxSize(size)(n)
		WaitSize(wait)(n)
		ConnTimeout(timeout)(n)
	}
}

func AllowReconnect(allow bool) Option {
	return func(n *Nats) {
		n.allowReconnect = allow
	}
}

func MaxSize(size int) Option {
	return func(n *Nats) {
		n.maxReconnect = size
	}
}

func WaitSize(wait time.Duration) Option {
	return func(n *Nats) {
		n.reconnectWait = wait
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(n *Nats) {
		n.timeout = timeout
	}
}

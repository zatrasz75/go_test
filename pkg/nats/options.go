package nats

import "time"

// Option -.
type Option func(*Nats)

func OptionSet(size int, wait, timeout time.Duration) Option {
	return func(n *Nats) {
		MaxSize(size)(n)
		WaitSize(wait)(n)
		Timeout(timeout)(n)
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

func Timeout(timeout time.Duration) Option {
	return func(n *Nats) {
		n.timeout = timeout
	}
}

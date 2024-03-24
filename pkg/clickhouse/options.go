package clickhouse

import "time"

type Option func(*Clickhouse)

func OptionSet(attempts int, timeout time.Duration) Option {
	return func(c *Clickhouse) {
		WithConnAttempts(attempts)(c)
		WithConnTimeout(timeout)(c)
	}
}

// WithConnAttempts устанавливает количество попыток подключения.
func WithConnAttempts(attempts int) Option {
	return func(ch *Clickhouse) {
		ch.connAttempts = attempts
	}
}

// WithConnTimeout устанавливает таймаут подключения.
func WithConnTimeout(timeout time.Duration) Option {
	return func(ch *Clickhouse) {
		ch.connTimeout = timeout
	}
}

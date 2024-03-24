package redisdb

import (
	"time"
)

type Option func(*Redis)

func OptionSet(addr, password string, db, protocol, poolSize, attempts int, poolTimeout, timeout, exp time.Duration) Option {
	return func(c *Redis) {
		WithAddr(addr)(c)
		WithPassword(password)(c)
		WithDB(db)(c)
		WithProtocol(protocol)(c)
		WithPoolSize(poolSize)(c)
		WithPoolTimeout(poolTimeout)(c)

		ConnAttempts(attempts)(c)
		ConnTimeout(timeout)(c)
		Expiration(exp)(c)
	}
}

func WithAddr(addr string) Option {
	return func(r *Redis) {
		r.Rds.Options().Addr = addr
	}
}

func WithPassword(password string) Option {
	return func(r *Redis) {
		r.Rds.Options().Password = password
	}
}

func WithDB(db int) Option {
	return func(r *Redis) {
		r.Rds.Options().DB = db
	}
}

func WithProtocol(protocol int) Option {
	return func(r *Redis) {
		r.Rds.Options().Protocol = protocol
	}
}

func WithPoolSize(poolSize int) Option {
	return func(r *Redis) {
		r.Rds.Options().PoolSize = poolSize
	}
}

func WithPoolTimeout(poolTimeout time.Duration) Option {
	return func(r *Redis) {
		r.Rds.Options().PoolTimeout = poolTimeout
	}
}

// ConnAttempts Попытки соединения
func ConnAttempts(attempts int) Option {
	return func(c *Redis) {
		c.connAttempts = attempts
	}
}

// ConnTimeout Время ожидания соединения
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Redis) {
		c.connTimeout = timeout
	}
}

func Expiration(exp time.Duration) Option {
	return func(c *Redis) {
		c.Expiration = exp
	}
}

package nats

import (
	"github.com/nats-io/nats.go"
	"time"
)

type Nats struct {
	nc *nats.Conn

	allowReconnect bool
	maxReconnect   int
	reconnectWait  time.Duration
	timeout        time.Duration
}

func New(natsURL string, opts ...Option) (*Nats, error) {
	n := &Nats{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(n)
	}

	var err error
	n.nc, err = nats.Connect(natsURL, nats.ReconnectWait(n.reconnectWait), nats.MaxReconnects(n.maxReconnect), nats.Timeout(n.timeout))
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *Nats) Publish(subject string, data []byte) error {
	return n.nc.Publish(subject, data)
}

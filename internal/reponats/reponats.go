package reponats

import (
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/nats"
)

type Store struct {
	*nats.Nats
	l logger.LoggersInterface
}

func New(nc *nats.Nats, l logger.LoggersInterface) *Store {
	return &Store{nc, l}
}

func (s *Store) Publish(subject string, data []byte) error {
	return s.Nats.Publish(subject, data)
}

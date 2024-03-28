package click

import (
	"zatrasz75/go_test/pkg/clickhouse"
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/nats"
)

type Store struct {
	*clickhouse.Clickhouse
	n *nats.Nats
	l logger.LoggersInterface
}

func New(ch *clickhouse.Clickhouse, nc *nats.Nats, l logger.LoggersInterface) *Store {
	return &Store{ch, nc, l}
}

func (s *Store) InsertLogsClickhouse() error {
	msgChan, err := s.n.ReceiveLog()
	if err != nil {
		s.l.Error("ReceiveLog", err)
		return err
	}

	for msg := range msgChan {
		err = s.InsertData(msg)
		if err != nil {
			s.l.Error("InsertData", err)
			return err
		}
	}

	return nil
}

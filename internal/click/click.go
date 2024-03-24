package click

import (
	"zatrasz75/go_test/pkg/clickhouse"
	"zatrasz75/go_test/pkg/logger"
)

type Store struct {
	*clickhouse.Clickhouse
	l logger.LoggersInterface
}

func New(ch *clickhouse.Clickhouse, l logger.LoggersInterface) *Store {
	return &Store{ch, l}
}

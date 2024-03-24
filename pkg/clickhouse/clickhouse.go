package clickhouse

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"zatrasz75/go_test/pkg/logger"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

// Clickhouse Хранилище логов
type Clickhouse struct {
	connAttempts int
	connTimeout  time.Duration

	ch *sql.DB
}

func New(dsn string, l logger.LoggersInterface, opts ...Option) (*Clickhouse, error) {
	c := &Clickhouse{}

	// Пользовательские параметры.
	for _, opt := range opts {
		opt(c)
	}

	var err error
	for c.connAttempts > 0 {
		c.ch, err = sql.Open("clickhouse", dsn)
		if err == nil {
			if err = c.ch.Ping(); err == nil {
				break
			}
		}
		l.Info("Clickhouse пытается подключиться, попыток осталось: %d", c.connAttempts)

		time.Sleep(c.connTimeout)

		c.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("clickhouse - New - connAttempts == 0: %w", err)
	}

	return c, nil
}

// Close Закрыть
func (c *Clickhouse) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
}

// Migrate Миграция таблиц
func (c *Clickhouse) Migrate() error {
	migrationScript, err := ioutil.ReadFile("initScriptClickhouse/up.sql")
	if err != nil {
		log.Fatal(err)
	}
	migrationScriptStr := string(migrationScript)

	statements := strings.Split(migrationScriptStr, ";")

	for _, statement := range statements {
		if strings.TrimSpace(statement) != "" {
			_, err = c.ch.Exec(statement)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

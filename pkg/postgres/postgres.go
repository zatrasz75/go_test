package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"zatrasz75/go_test/pkg/logger"
)

// Postgres Хранилище данных
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(connStr string, l logger.LoggersInterface, opts ...Option) (*Postgres, error) {
	pg := &Postgres{}

	// Пользовательские параметры
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		l.Info("Postgres пытается подключиться, попыток осталось: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close Закрыть
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Migrate Миграция таблиц
func (p *Postgres) Migrate() error {
	migrationScript, err := ioutil.ReadFile("initScriptPostgres/up.sql")
	if err != nil {
		log.Fatal(err)
	}
	migrationScriptStr := string(migrationScript)

	statements := strings.Split(migrationScriptStr, ";")

	for _, statement := range statements {
		if strings.TrimSpace(statement) != "" {
			_, err = p.Pool.Exec(context.Background(), statement)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RollingUp Добавление записи о миграции
func (p *Postgres) RollingUp() error {
	var id int
	query := "INSERT INTO projects (name) VALUES ($1) RETURNING id"
	row := p.Pool.QueryRow(context.Background(), query, "запись")
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	up := "UPDATE projects SET name=$1 WHERE id=$2"
	_, err = p.Pool.Exec(context.Background(), up, "новая запись"+" "+"№ "+fmt.Sprint(id), id)
	if err != nil {
		return err
	}

	return nil
}

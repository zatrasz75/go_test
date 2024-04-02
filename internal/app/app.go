package app

import (
	redis2 "github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
	"zatrasz75/go_test/configs"
	"zatrasz75/go_test/internal/click"
	"zatrasz75/go_test/internal/controller"
	"zatrasz75/go_test/internal/redis"
	"zatrasz75/go_test/internal/repository"
	"zatrasz75/go_test/pkg/clickhouse"
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/nats"
	"zatrasz75/go_test/pkg/postgres"
	"zatrasz75/go_test/pkg/redisdb"
	"zatrasz75/go_test/pkg/server"
)

func Run(cfg *configs.Config, l logger.LoggersInterface) {
	pg, err := postgres.New(cfg.DataBase.ConnStr, l, postgres.OptionSet(cfg.DataBase.PoolMax, cfg.DataBase.ConnAttempts, cfg.DataBase.ConnTimeout))
	if err != nil {
		l.Fatal("ошибка запуска - postgres.New:", err)
	}
	defer pg.Close()

	err = pg.Migrate()
	if err != nil {
		l.Fatal("ошибка миграции", err)
	}
	err = pg.RollingUp()
	if err != nil {
		l.Fatal("ошибка добавления записи", err)
	}

	rds, err := redisdb.New(cfg.RedisDB.Addr, l, redisdb.OptionSet(cfg.RedisDB.Addr, cfg.RedisDB.Password, cfg.RedisDB.DB, cfg.RedisDB.Protocol, cfg.RedisDB.PoolSize, cfg.RedisDB.ConnAttempts, cfg.RedisDB.PoolTimeout, cfg.RedisDB.ConnTimeout, cfg.RedisDB.Expiration))
	if err != nil {
		l.Fatal("ошибка запуска - redis.New:", err)
	}
	defer func(Rds *redis2.Client) {
		err = Rds.Close()
		if err != nil {
			l.Warn("ошибка закрытия соединения redis:", err)
		}
	}(rds.Rds)

	nc, err := nats.New(cfg.Nats.NatsURL, l, nats.OptionSet(cfg.Nats.MaxReconnect, cfg.Nats.ReconnectWait, cfg.Nats.Timeout))
	if err != nil {
		l.Fatal("ошибка запуска - nats.New", err)
	}

	ch, err := clickhouse.New(cfg.Clickhouse.DSN, l, clickhouse.OptionSet(cfg.Clickhouse.ConnAttempts, cfg.Clickhouse.ConnTimeout))
	if err != nil {
		l.Fatal("ошибка запуска clickhouse.New", err)
	}
	defer ch.Close()

	err = ch.Migrate()
	if err != nil {
		l.Fatal("ошибка миграции", err)
	}

	repo := repository.New(pg, l)
	rd := redis.New(rds, l)
	cl := click.New(ch, nc, l)

	go func() {
		err = cl.InsertLogsClickhouse()
		if err != nil {
			return
		}
	}()

	router := controller.NewRouter(cfg, l, repo, rd, nc)

	srv := server.New(router, server.OptionSet(cfg.Server.AddrHost, cfg.Server.AddrPort, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, cfg.Server.IdleTimeout, cfg.Server.ShutdownTime))

	go func() {
		err = srv.Start()
		if err != nil {
			l.Error("Остановка сервера:", err)
		}
	}()

	l.Info("Запуск сервера на http://" + cfg.Server.AddrHost + ":" + cfg.Server.AddrPort)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("принят сигнал прерывания прерывание %s", s.String())
	case err = <-srv.Notify():
		l.Error("получена ошибка сигнала прерывания сервера", err)
	}

	err = srv.Shutdown()
	if err != nil {
		l.Error("не удалось завершить работу сервера", err)
	}
}

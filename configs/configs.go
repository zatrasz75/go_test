package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"time"
	"zatrasz75/go_test/pkg/logger"
)

type Config struct {
	Server struct {
		AddrPort     string        `yaml:"port" env:"APP_PORT" env-description:"Server port" env-default:"8585"`
		AddrHost     string        `yaml:"host" env:"APP_IP" env-description:"Server host" env-default:"0.0.0.0"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env:"READ_TIMEOUT" env-description:"Server ReadTimeout" env-default:"3s"`
		WriteTimeout time.Duration `yaml:"write-timeout" env:"WRITE_TIMEOUT" env-description:"Server WriteTimeout" env-default:"3s"`
		IdleTimeout  time.Duration `yaml:"idle-timeout" env:"IDLE_TIMEOUT" env-description:"Server IdleTimeout" env-default:"6s"`
		ShutdownTime time.Duration `yaml:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" env-description:"Server ShutdownTime" env-default:"10s"`
	} `yaml:"server"`
	DataBase struct {
		ConnStr string `env:"DB_CONNECTION_STRING" env-description:"db string"`

		Host     string `yaml:"host" env:"HOST_DB" env-description:"db host"`
		User     string `yaml:"username" env:"POSTGRES_USER" env-description:"db username"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-description:"db password"`
		Url      string `yaml:"db-url" env:"URL_DB" env-description:"db url"`
		Name     string `yaml:"db-name" env:"POSTGRES_DB" env-description:"db name"`
		Port     string `yaml:"port" env:"PORT_DB" env-description:"db port"`

		PoolMax      int           `yaml:"pool-max" env:"PG_POOL_MAX" env-description:"db PoolMax"`
		ConnAttempts int           `yaml:"conn-attempts" env:"PG_CONN_ATTEMPTS" env-description:"db ConnAttempts"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"PG_TIMEOUT" env-description:"db ConnTimeout"`
	} `yaml:"database"`
	RedisDB struct {
		Addr        string        `yaml:"addr" env:"REDIS_ADDR_PORT" env-description:"redis Addr" env-default:"localhost:6379"`
		Password    string        `yaml:"password" env:"REDIS_PASSWORD" env-description:"redis password" env-default:""`
		DB          int           `yaml:"db" env:"REDIS_DB" env-description:"redis db" env-default:"0"`
		Protocol    int           `yaml:"protocol" env:"REDIS_PROTOCOL" env-description:"redis protocol" env-default:"2"`
		PoolSize    int           `yaml:"pool-size" env:"REDIS_POOL_SIZE" env-description:"redis PoolSize" env-default:"10"`
		PoolTimeout time.Duration `yaml:"pool-timeout" env:"REDIS_POOL_TIMEOUT" env-description:"redis PoolTimeout" env-default:"5s"`

		ConnAttempts int           `yaml:"redis-attempts" env:"REDIS_CONN_ATTEMPTS" env-description:"redis ConnAttempts"`
		ConnTimeout  time.Duration `yaml:"redis-timeout" env:"REDIS_TIMEOUT" env-description:"redis ConnTimeout"`
		Expiration   time.Duration `yaml:"redis-exp" env:"REDIS_EXP" env-description:"redis Expiration" env-default:"1m"`
	} `yaml:"redis"`
	Nats struct {
		NatsURL        string        `yaml:"nats-natsURL" env:"NATS_CONNECT_URL" env-description:"reponats NatsURL"`
		AllowReconnect bool          `yaml:"nats-reconnect" env:"NATS_ALLOW_RECONNECT" env-description:"reponats AllowReconnect" env-default:"true"`
		MaxReconnect   int           `yaml:"nats-max-reconnect" env:"NATS_MAX_RECONNECT" env-description:"reponats MaxReconnect" env-default:"10"`
		ReconnectWait  time.Duration `yaml:"nats-wait-reconnect" env:"NATS_WAIT_RECONNECT" env-description:"reponats ReconnectWait" env-default:"1s"`
		Timeout        time.Duration `yaml:"nats-timeout" env:"NATS_TIMEOUT" env-description:"reponats timeout" env-default:"1s"`
	} `yaml:"nats"`
	Clickhouse struct {
		DSN string `env:"DSN_CONNECTION_STRING" env-description:"DSN string"`

		Host       string `yaml:"host" env:"HOST_CH" env-description:"ch host"`
		Port       string `yaml:"port" env:"PORT_CH" env-description:"ch port"`
		UserCh     string `yaml:"user-ch" env:"CLICKHOUSE_USER" env-description:"UserCh"`
		PasswordCh string `yaml:"password-ch" env:"CLICKHOUSE_PASSWORD" env-description:"PasswordCh"`
		NameCh     string `yaml:"name-ch" env:"CLICKHOUSE_DB" env-description:"NameCh"`
		AccessCh   int    `yaml:"access-ch" env:"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT" env-description:"NameCh" env-default:"1"`

		ConnAttempts int           `yaml:"conn-attempts" env:"CH_CONN_ATTEMPTS" env-description:"CH ConnAttempts"`
		ConnTimeout  time.Duration `yaml:"conn-timeout" env:"CH_TIMEOUT" env-description:"CH ConnTimeout"`
	} `yaml:"clickhouse"`
}

func NewConfig(l logger.LoggersInterface) (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		l.Warn("системе не удается найти указанный файл .env: - %v", err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		l.Error("ошибка .env ", err)
		return nil, err
	}
	if err := cleanenv.ReadConfig("./configs/configs.yml", &cfg); err != nil {
		return nil, err
	}

	cfg.DataBase.ConnStr = initDB(cfg)
	cfg.Clickhouse.DSN = initCH(cfg)

	return &cfg, nil
}

func initDB(cfg Config) string {
	if cfg.DataBase.ConnStr != "" {
		return cfg.DataBase.ConnStr
	}
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DataBase.Host,
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Url,
		cfg.DataBase.Port,
		cfg.DataBase.Name,
	)
}

func initCH(cfg Config) string {
	if cfg.Clickhouse.DSN != "" {
		return cfg.Clickhouse.DSN
	}
	return fmt.Sprintf("http://%s:%s?database=%s&username=%s&password=%s",
		cfg.Clickhouse.Host,
		cfg.Clickhouse.Port,
		cfg.Clickhouse.NameCh,
		cfg.Clickhouse.UserCh,
		cfg.Clickhouse.PasswordCh,
	)
}

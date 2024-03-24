package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"zatrasz75/go_test/models"
	"zatrasz75/go_test/pkg/logger"
	"zatrasz75/go_test/pkg/redisdb"
)

type Store struct {
	*redisdb.Redis
	l logger.LoggersInterface
}

func New(rds *redisdb.Redis, l logger.LoggersInterface) *Store {
	return &Store{rds, l}
}

// GetList Получаем сохраненные записей
func (s *Store) GetList(key string) ([]models.Goods, error) {
	val, err := s.Rds.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Key не найден в Redis.
			return nil, nil
		}
		return nil, err
	}

	var list []models.Goods
	err = json.Unmarshal([]byte(val), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AddList Сохраняет все записи на определенное время
func (s *Store) AddList(key string, list []models.Goods) error {
	goods, err := json.Marshal(list)
	if err != nil {
		s.l.Error("ошибка при кодировании данных в JSON: ", err)
		return nil
	}
	statusCmd := s.Rds.Set(context.Background(), key, goods, s.Expiration)
	result, err := statusCmd.Result()
	if err != nil {
		return err
	}
	s.l.Info("Redis: %s", result)

	return nil
}
